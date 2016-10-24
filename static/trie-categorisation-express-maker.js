
var utils = require('./utils'),
	fs = require('fs'),
	promise = require('bluebird'),
	_ = require('lodash'),
	chop_2ld = (name) => {
		var match = name.match(/^[^\.]*\.[^\.]*/);
		if (match) { 
			return match[0];
		}
		return name;
	};


if (require.main === module) { 
	console.error('input : < trie-5k json >  < libtotype csv > < output csv >');
	var trie5k = JSON.parse(fs.readFileSync(process.argv[2])),
		libtotype = utils.loadCSV(process.argv[3]).reduce((d, lib) => {
			d[lib.package] = lib;
			return d;
		}, {}),	
		short2app = _.reduce(trie5k, (d, libs, app) => {
			libs.map((lib) => {
				var short = chop_2ld(lib);
				d[short] = _.union(d[short] || [], [app]);
			});
			return d;
		}, {}),
		isAd = _.reduce(libtotype, ((d, lib, libname) => {
			var short = chop_2ld(libname);
			if (lib.type.indexOf('ads') >= 0 || libname.indexOf('ads') >= 0) { 
				d[short] = true;
			}
			return d;
		}), {}),
		hasVisited = _.reduce(libtotype, ((d, lib, libname) => {
			var short = chop_2ld(libname);
			if (isAd[short] || (lib.type && lib.type.length > 0)) { 
				d[short] = (d[short] || []).concat([libname]);
			}
			return d;
		}), {}),		
		outputfname = process.argv[4];

	var missing_shorts = _.reduce(trie5k, (missing, libs, app_pkg) => {
			libs.map((libname) => {
				var short = chop_2ld(libname);
				if (!hasVisited[short]) { missing.push(short); }
				hasVisited[short] = (hasVisited[short] || []).concat([libname]);				
			});
			return missing;
		}, []),
		rows = missing_shorts.map((short) => [short, short2app[short].length, hasVisited[short].join(';')]);

	console.error('Writing output ', rows.length, ' to ', outputfname);

	utils.saveCSV(rows, outputfname).then((x) => console.log("Done!"));
}