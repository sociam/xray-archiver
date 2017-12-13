

const gplay = require('google-play-scraper'),
      fs = require('fs'),
      config = require('/etc/xray/config.json'),
      { Pool, Client } = require('pg'),
      _ = require('lodash'),
      path = require('path'),
      Promise = require('bluebird'),
      logger = require('../../util/logger'),
      { exec } = require('child_process');

const region = 'uk';

function walk(dir, callback) {
    const files = fs.readdirSync(dir);
    files.map(function(file) {
	var filepath = path.join(dir, file), 
	    stats = fs.statSync(filepath);
	if (stats.isDirectory()) {
	    walk(filepath, callback);
	} else if (stats.isFile()) {
	    callback(filepath, stats);
	}
    });
}


function getAnalysis(client, appId, appVersion) { //
    console.log('getAnalysis ', appId, appVersion);
    return client.query('select app_versions.app,app_versions.version,app_versions.analyzed,app_analyses.analyzer,app_analyses.analysis from app_versions left outer join app_analyses on app_analyses.id = app_versions.id WHERE app_versions.app = $1 AND app_versions.version = $2', [appId, appVersion]).then((res) => {
	console.log('!!!!!!!!!!!!!!!!!!!!!!! >>>>>>>>>>>>> result ', res.rows);
	// console.log(`app query res count: ${res.rowCount}`);
	return res.rows[0];
    }).catch((e) => {
	console.error('Error query ', e);
	logger.err('Error checking if app exists in the database:', err);
	throw err;
    });
}

// relocation src/dst
const source = config.freezer.sync && config.freezer.sync.src,
      dest = config.freezer.sync && config.freezer.sync.dest;

function processAPK(client, dirname) {
    // returns the path to the apk if apk exists
    const fns = [];
    walk(dirname, (f, src_stats) => {
	fns.push(() => {
	    if (f.toLowerCase().indexOf('.apk') >= 0 && f.toLowerCase().indexOf('.apk') === f.length - '.apk'.length) {
		// console.log(`got a candidate! ${f.slice(source.length+1)}`)
		const apk_path = f.slice(source.length+1),
		      split = f.split('/'),
		      version = split[split.length - 2],
		      pkg = apk_path.slice(0,apk_path.indexOf('/'));

		console.log(`pkg path ${apk_path} - pkg:${pkg}, version:${version}`);

		const winMatch = () => {
		    // matches already, so no need to sync.
		    return getAnalysis(client, pkg, version).then((has) => {
		    	console.log('getAnalysis results ', has);
			// interpret results.. and if successful then check to delete file!
			if (!has.analysis) {
			    console.info('i think we have a false');
			    // then do analysis here...
			    // then deletion			    
			} else {			    
    			    console.info('i think we have a true');
			    // just go straight to deletion
			}
			return Promise.resolve();
		    }).catch((ea) => {
		    	console.error('ERROR with getanalysis', ea);
		    });
		}, failMatch = () => {
		    // failmatch so we first sync if
		    console.log(`sync ${config.freezer.sync}, ${config.freezer.sync.enable}`);
		    if (config.freezer.sync && config.freezer.sync.enable) {
			// first rsync then ...
			return new Promise((win, fail) => {
			    logger.info(`running rsync >> /usr/bin/rsync -avz ${path.join(source, pkg)} ${dest}/`);
			    exec(`/usr/bin/rsync -avz ${path.join(source, pkg)} ${dest}/`, (error, stdout, stderr) => {
				if (error) {
				    logger.error(`rsync error: ${error}`);
				    return fail();
				}
				logger.info(`rsync output ${stdout} ${stderr}`);
				return win();
			    });
			});			
		    }
		    return Promise.reject('error: sync is not enabled');
		};
		
		try {
		    const dstat =  fs.statSync(path.join(dest, apk_path));
		    if (dstat && dstat.isFile() && dstat.size === src_stats.size) {
			return winMatch();
		    } else {
			logger.info('exists on cold but mismatch ', path.join(dest, apk_path), ' size ', dstat.size, ' <> ', src_stats.size);
			return failMatch().then(() => winMatch());
		    }
		} catch(e) {
		    if (e && e.errno === -2) { 
			logger.info('no such file on dest ', path.join(dest, apk_path));
			return failMatch().then(() => winMatch());
		    } else {
			logger.error(`unknown error ${e}, so skipping ${dest}`);
			return Promise.resolve();
		    }
		}
	    }
	});
    });
    // now we have a list of things we can reduce
    return Promise.reduce(fns, (reduction, f) => f(), {});
};


function connect() {
    const dbCfg = _.extend({}, config.db, config.freezer.db);
    
    dbCfg.max = 10;
    dbCfg.idleTimeoutMillis = 30000;

    console.info(`authenticating using u:${dbCfg.user}, p:${dbCfg.password}`);
    console.log(dbCfg);

    const client = new Client(dbCfg);
    client.connect();
    return client;
}


(()=> {
    const client = connect();
    client.query('select * from app_versions limit 1;').then(async (res) => {
	console.log('test query res ', res.rows);
	fs.readdir(source, (err, items) => {
	    console.log(`got ${items.length} items`);
	    Promise.reduce(items, (reductor, dirname) => processAPK(client, path.join(source,dirname)), {});
	});
    });

})();

