
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
	console.error('arguments : < trie-5k json > < non-collapsed-labels.csv > < collapsed-labels.csv > < top5k-apps.csv > < category-mapping.json > < output csv >');
	var trie5k = JSON.parse(fs.readFileSync(process.argv[2])),
		libtotype = utils.loadCSV(process.argv[3]).reduce((d, lib) => {
			d[lib.package] = lib;
			return d;
		}, {}),
		// shortened base -> type
		short2type = utils.loadCSV(process.argv[4]).reduce((d, lib) => {
			d[lib.base] = lib.type;
			return d;
		}, {}), 
		// category mapping
		catmap = JSON.parse(fs.readFileSync(process.argv[6])),
		// app 2 category
		app2cat = utils.loadCSV(process.argv[5]).reduce((d,apk, lib) => {
			if (apk.category && catmap[apk.category]) { 
				d[apk.id] = catmap[apk.category];
			}
			return d;
		}, {}),
		isAd = _.reduce(libtotype, ((d, lib, libname) => {
			var short = chop_2ld(libname);
			if (lib.type.indexOf('ads') >= 0 || short2type[short] && short2type[short].indexOf('ads') >= 0) { 
				d[short] = true;
			}
			return d;
		}), {}),
		outputfname = process.argv[7],
		hist = (vals) => {
			var mx = _.max(vals);
			return _.range(0,mx).map((x) => vals.filter((i) => i === x).length);
		};

	var adsbyapkbycat = _.reduce(trie5k, (bycat, libs, app_pkg) => {
			var category = app2cat[app_pkg];
			if (!category) { console.error("no category for app ", app_pkg); return bycat; }
			bycat[category] = bycat[category] || {};
			var adlibs = _.uniq(libs.filter((libname) => isAd[chop_2ld(libname)]).map((lib) => chop_2ld(lib)));
			console.info(libs.filter((libname) => isAd[chop_2ld(libname)]), ' => ', adlibs);
			bycat[category][app_pkg] = adlibs.length;
			return bycat;
		}, {}),
		rows = _.map(adsbyapkbycat, (n_adsbyapk,cat) => {
			var row = {  cat : cat };
			_.extend(row, hist( _.values(n_adsbyapk) ));
			return row;
		});
	console.error('Writing output ', rows.length, ' to ', outputfname);
	utils.saveCSV(rows, outputfname).then((x) => console.log("Done!"));
}