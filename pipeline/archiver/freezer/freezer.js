
const gplay = require('google-play-scraper'),
      fs = require('fs'),
      _ = require('lodash'),
      Promise = require('bluebird'),
      path = require('path'),
      logger = require('../../util/logger'),
      db = new (require('../../db/db'))('retriever');

const region = 'uk';

// source
const source = '/mnt/a/apps', dest = '/mnt/cold';

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


function containsAPK(dirname) {
    // returns the path to the apk if apk exists
    walk(dirname, (f, src_stats) => {
	if (f.toLowerCase().indexOf('.apk') === f.length - '.apk'.length) {
	    // console.log(`got a candidate! ${f.slice(source.length+1)}`)
	    const apk_path = f.slice(source.length+1);
	    try {
		const dstat =  fs.statSync(path.join(dest, apk_path));
		if (dstat && dstat.isFile() && dstat.size === src_stats.size) {
		    console.log('exists on cold ', path.join(dest, apk_path), ' size ', dstat.size, ' === ', src_stats.size);
		} else {
		    console.log('exists on cold but mismatch ', path.join(dest, apk_path), ' size ', dstat.size, ' <> ', src_stats.size);
		}
	    } catch(e) {
		if (e && e.errno === -2) { 
		    console.log('no such file on dest ', path.join(dest, apk_path));
		} else {
    		    console.log(e);
		}
	    }
	}
    });	
};



(()=> { 
    fs.readdir(source, (err, items) => {
	console.log(`got ${items.length} items`);
	items.map((dirname) => {
	    return containsAPK(path.join(source,dirname));
	});		  
    });
})();

