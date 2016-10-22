
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
		isAd = _.reduce(libtotype, ((d, lib, libname) => {
			var short = chop_2ld(libname);
			if (lib.type.indexOf('ads') >= 0) { 
				d[short] = true;
			}
			return d;
		}), {}),
		outputfname = process.argv[4];

	console.log('libtotype ', JSON.stringify(libtotype, 0, 4));
	var libsbyapp = _.reduce(trie5k, (byapp, libs, app_pkg) => {
			libs.map((libname) => {
				var short = chop_2ld(libname);
				if (!libtotype[libname]) { 
					console.warn('NO lib in type ', libname); 
					return;
				}
				console.info('got! type ', libname, libtotype[libname].type, app_pkg, short);
				// if (libtotype[libname].type.indexOf('ads') >= 0) {
				if (isAd[short]) { 
					byapp[app_pkg] = _.union(byapp[app_pkg] || {}, [short]);
				}
			});
			return byapp;
		}, {}),
		rowcounts = _.map(libsbyapp, ((libs, app) => [app, libs.length]));

	console.log('libsbyapp ', JSON.stringify(libsbyapp, 0, 4));

	console.error('Writing output ', rowcounts.length, ' to ', outputfname);
	utils.saveCSV(rowcounts, outputfname).then((x) => {
		console.log("Done!");
	});
}