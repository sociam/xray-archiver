
var fs = require('fs'),
	promise = require('bluebird'),
	_ = require('lodash'),
	config = JSON.parse(fs.readFileSync('./config.json')),
	spawn = require('child_process'),
	deleteFolderRecursive = function(path) {
	  if( fs.existsSync(path) ) {
	      fs.readdirSync(path).forEach(function(file) {
	        var curPath = path + "/" + file;
	          if(fs.statSync(curPath).isDirectory()) { // recurse
	              deleteFolderRecursive(curPath);
	          } else { // delete file
	              fs.unlinkSync(curPath);
	          }
	      });
	      fs.rmdirSync(path);
	    }
	};

var trie_mkchild = (name, fullname) => ({ name:name, fullname: fullname, children:{}, subtree:[] }),
	trie_root = trie_mkchild('', ''),
	reset_root = () => { trie_root = trie_mkchild('',''); },
	trie_add = (path) => {
		// path is an array
		var cur = trie_root;
		path.map((p) => {
			var idx = path.indexOf(p),
				partial_path = path.slice(0,idx+1).join('.');
			cur.subtree.push(path.join('.'));
			if (cur.children[p] === undefined) {
				cur.children[p] = trie_mkchild(partial_path);
			}
			cur = cur.children[p];
		});
		cur.subtree.push(path.join('.'));
	},
	flattened_trie = (t) => {
		return _.values(t.children).map(flattened_trie).reduce((lst, a) => lst.concat(a), [t]);
	},
	find_packages = (ft) => {
		// takes flattened tree in
		const MIN_LENGTH = 2;
		var so_far = [];

		ft.map((x) => {
			if (x.name.split('.').length > MIN_LENGTH && _.values(x.children).length > 1) {
				// replace longer matches with this shorter one

				// get rid of all that are longer
				so_far = so_far.filter((sfm) => sfm.indexOf(x.name) < 0);

				// console.log('checking dupe ', so_far, x.name, 
				// 	'unique ? ', 
				// 	so_far.filter((sf) => x.name.indexOf(sf) >= 0),
				// 	so_far.filter((sf) => x.name.indexOf(sf) >= 0).length === 0
				// );

				// only add us if we are the shortest
				if (so_far.filter((sf) => x.name.indexOf(sf) >= 0).length === 0) { 
					so_far.push(x.name);
				}
			}
		});
		var trees = so_far.map((n) => ft.filter((x) => x.name === n)[0]);
		// console.log('!~ names ', so_far, ' trees ', trees);
		return trees;
	};

// trie_add(['a','b','c']);
// trie_add(['a','b','d']);
// trie_add(['a','b','c','e']);
// trie_add(['a','b','c','f']);
// console.log(JSON.stringify(flattened_trie(trie_root), null, 4));

var walkDir = (dirname, appname, subdirname, basedir) => {
	fs.readdirSync(dirname).map((filename) => {
		var fullpath = [dirname,filename].join('/'),
			stat = fs.statSync(fullpath),
			relevant_part = fullpath.slice(basedir.length+appname.length+subdirname.length+3),
			pathsplits = relevant_part.split('/');

//		console.log('fullpath > ', relevant_part);

		if (stat && stat.isDirectory()) { 
			trie_add(pathsplits); 
			walkDir(fullpath, appname, subdirname, basedir);
		}
	});
};


// now let's try doing some magic

var apktoolpath = config.apktoolpath,
	tmpdir = config.tmpdir,
	by_app = {},
	toplevel = () => {
	fs.readdirSync(config.appsdir).map((apkname) => {
		var apk = apkname.indexOf('.apk') >= 0,
			appname = apkname.slice(0,-4),
			apkpath = [config.appsdir,apkname].join('/'),
			cmd = `java -jar ${apktoolpath} d ${apkpath} -f`,
			unpackroot = [config.tmpdir, appname].join('/');

		if (!apk) { console.error('skipping ', apkname); return; }
		// console.error('executing ', cmd);
		try {
			spawn.execSync(cmd, { cwd:tmpdir });
			// console.log('walking ', unpackdirname);
			fs.readdirSync(unpackroot).map((sdname) => {
				var fullpath = [unpackroot,sdname].join('/');
				if (sdname.indexOf('smali') >= 0 || sdname.indexOf('unknown') >= 0 && fs.statSync(fullpath).isDirectory()) {
					console.error('walking ', fullpath);
					walkDir(fullpath, appname, sdname, tmpdir);
				}
			});
			by_app[appname] = find_packages(flattened_trie(trie_root)).map((p) => p.name);
			reset_root();
		} catch(e) {
			console.error('skipping ', appname);
		}
		try {
			// delete the app 
			console.error('cleaning up ', unpackroot);			
			deleteFolderRecursive(unpackroot);
		} catch(e) { 
			console.error("error cleaning up ", e);
		}
	});
	// console.log(' ----------------------> full trie -------> ');
	// console.log(JSON.stringify(trie_root, null, 4));

	// console.log(' ----------------------> flat trie -------> ');
	// console.log(flattened_trie(trie_root));

	// console.log(' ----------------------> find packages -------> ');
	// console.log(find_packages(flattened_trie(trie_root)));

	console.log(JSON.stringify(by_app, null, 4));

};

toplevel();
