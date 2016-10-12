
var fs = require('fs'),
	promise = require('bluebird'),
	_ = require('lodash'),
	config = JSON.parse(fs.readFileSync('./config.json')),
	spawn = require('child_process');

var trie_mkchild = (name) => ({ name:name, children:{}, subtree:[] })
	trie_root = trie_mkchild('.'),
	trie_add = (path) => {
		// path is an array
		var cur = trie_root;
		path.map((p) => {
			cur.subtree.push(path.join('.'))
			if (cur.children[p] === undefined) {
				cur.children[p] = trie_mkchild(p)
			}
			cur = cur.children[p];
		});
		cur.subtree.push(path.join('.'))
	};

console.log(JSON.stringify(trie_root, null, 4));

// now let's try doing some magic

var apktoolpath = config.apktoolpath,
	tmpdir = config.tmpdir,
	toplevel = () => {
	fs.readdirSync(config.appsdir).map((appname) => {
		var apppath = [config.appsdir,appname].join('/'),
			cmd = `java -jar ${apktoolpath} d ${apppath}`;
		console.log('executing ', cmd);
		spawn.execSync(cmd, { cwd:tmpdir });
	});
};

toplevel();
