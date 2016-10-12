
var fs = require('fs'),
	promise = require('bluebird'),
	_ = require('lodash'),
	config = JSON.parse(fs.readFileSync('./config.json')),
	spawn = require('child_process');

var trie_mkchild = (name, fullname) => ({ name:name, fullname: fullname, children:{}, subtree:[] }),
	trie_root = trie_mkchild('', ''),
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

				console.log('checking dupe ', so_far, x.name, 
					'unique ? ', 
					so_far.filter((sf) => x.name.indexOf(sf) >= 0),
					so_far.filter((sf) => x.name.indexOf(sf) >= 0).length === 0
				);

				// only add us if we are the shortest
				if (so_far.filter((sf) => x.name.indexOf(sf) >= 0).length === 0) { 
					so_far.push(x.name);
				}
			}
		});
		var trees = so_far.map((n) => ft.filter((x) => x.name === n)[0]);
		console.log('!~ names ', so_far, ' trees ', trees);
		return trees;
	};

// trie_add(['a','b','c']);
// trie_add(['a','b','d']);
// trie_add(['a','b','c','e']);
// trie_add(['a','b','c','f']);
// console.log(JSON.stringify(flattened_trie(trie_root), null, 4));

var walkDir = (dirname, appname, basedir) => {
	fs.readdirSync(dirname).map((filename) => {
		var fullpath = [dirname,filename].join('/'),
			stat = fs.statSync(fullpath),
			pathsplits = fullpath.slice(basedir.length+appname.length+'smali'.length+3).split('/');
		console.log('fullpath > ', fullpath, pathsplits);
		if (stat && stat.isDirectory()) { 
			trie_add(pathsplits); 
			walkDir(fullpath, appname, basedir);
		}
	});
};


// now let's try doing some magic

var apktoolpath = config.apktoolpath,
	tmpdir = config.tmpdir,
	toplevel = () => {
	fs.readdirSync(config.appsdir).map((appname) => {
		var apk = appname.indexOf('.apk') >= 0,
			apppath = [config.appsdir,appname].join('/'),
			cmd = `java -jar ${apktoolpath} d ${apppath} -f`,
			unpackdirname = [config.tmpdir, appname.slice(0,-4), 'smali'].join('/');

		console.log('appname ', appname, ' apk ', apk);
		if (!apk) { console.log('skipping ', appname); return; }
		console.log('executing ', cmd);
		spawn.execSync(cmd, { cwd:tmpdir });
		console.log('walking ', unpackdirname);
		walkDir(unpackdirname, appname.slice(0,-4), tmpdir);
	});
	console.log(' ----------------------> full trie -------> ');
	console.log(JSON.stringify(trie_root, null, 4));

	console.log(' ----------------------> flat trie -------> ');
	console.log(flattened_trie(trie_root));

	console.log(' ----------------------> find packages -------> ');
	console.log(find_packages(flattened_trie(trie_root)));
};

toplevel();
