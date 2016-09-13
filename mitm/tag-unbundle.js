

// textutil -convert txt *

var parse = require('csv-parse/lib/sync'),
	csvstr = require('csv-stringify'),
	fs = require('fs'),
	promise = require('bluebird'),
	_ = require('lodash'),
	qs = require('querystring'),
	config = JSON.parse(fs.readFileSync('./munge-config.json')),
	cutils = require('./client-utils'),
	round_re = /ROUND\W*(\d+)/,
	headers,
	exit_re = /exit\W*\>|exit interview|exit questions|exit survey/i;

var loadCSV = (fname) => {
	var text = 	fs.readFileSync(fname).toString();
	console.log("Parsing file ", fname, "(", text.length, ")");
	var data = parse(text, {max_limit_on_data_read:9999999999});
	headers = data[0];
	data = data.slice(1);
	data = data.map((x) => _.zipObject(headers,x));
	return data;
}, s = (str) => (str||'').split(' ').map((x) => x.trim().toUpperCase().trim(',')).filter((x) => x.length),
	uniquetags = (fname, rater_a, rater_b) => {
		var rows = loadCSV(fname),
			cols = [rater_a, rater_b];
		return _.uniq(_.flatten(rows.map((r) => _.flatten(cols.map((c) => s(r[c]))))));
	},
  main = (fname, cols, choose_tags) => {
	console.info('got cols');
	var rows = loadCSV(fname),
		people = cols,
		rater_a = cols[0], rater_b = cols[1],
		tags = choose_tags || _.uniq(_.flatten(rows.map((r) => _.flatten(cols.map((c) => s(r[c])))))),
		header = ['id', rater_a, rater_b ],
		rowos = [header].concat(
			rows.filter((r) => {
				return r[rater_a].length && r[rater_b].length;
			}).reduce((cat, r) => {
				// [[T], [T], [T], [T]]
				return cat.concat(tags.map((T) => {
					return [ r.id + '-' + T, s(r[rater_a]).indexOf(T) >= 0 ? 'YES' : 'NO', s(r[rater_b]).indexOf(T) >= 0 ? 'YES' : 'NO' ];
				}));
			}, []));

	console.log('rowos ', rowos);

	return new Promise((acc, rej) => {
		csvstr(rowos, (err, output) => {
			if (!err) { return acc(output); }
			console.error("Error ", err);
			rej(err);
		});
	});	
}, krips = (fname, raters, choose_tags) => {
	console.info('got cols');
	var rows = loadCSV(fname),
		rater_a = raters[0], rater_b = raters[1],
		tags = choose_tags || _.uniq(_.flatten(rows.map((r) => _.flatten(raters.map((c) => s(r[c])))))).sort(),
		header = ['id'].concat(tags.reduce((cat, T) => cat.concat(raters.map((r) => T + '' + (parseInt(raters.indexOf(r)) + 1))), [])),
		rowos = [header].concat(
			rows.filter((r) => {
				return r[rater_a].length || r[rater_b].length;
			}).map((r) => {
				return [r.id].concat(tags.reduce((cat, T) => {
					return cat.concat(raters.map((rater) => s(r[rater]).indexOf(T) >= 0 ? 1 : 0));
				}, []));
			}));

	console.log('tags ', rowos);
	return new Promise((acc, rej) => {
		csvstr(rowos, (err, output) => {
			if (!err) { return acc(output); }
			console.error("Error ", err);
			rej(err);
		});
	});	
}, consolidate = (fname, raters, choose_tags) => {
	// console.info('got cols');
	var rows = loadCSV(fname),
		rater_a = raters[0], rater_b = raters[1],
		tags = choose_tags || _.uniq(_.flatten(rows.map((r) => _.flatten(raters.map((c) => s(r[c])))))).sort(),
		header = ['id'].concat(tags.reduce((cat, T) => cat.concat([T]), [])),
		rowos = [header].concat(
			rows.map((r) => {
				return [r.id].concat(tags.reduce((cat, T) => {
					return cat.concat( _.flatten(raters.map((rater) => s(r[rater]))).indexOf(T) >= 0 ? 1 : 0 );
				}, []));
			}));

	// console.log('tags ', rowos);
	return new Promise((acc, rej) => {
		csvstr(rowos, (err, output) => {
			if (!err) { return acc(output); }
			console.error("Error ", err);
			rej(err);
		});
	});	
};

if (require.main === module) { 
	// node command fileout pers1 pers2 krips/separate
	if (process.argv[5] === 'krips') {
		console.info('krips!');
		krips(process.argv[2], process.argv.slice(3,5)).then((data) => {
			// write it !
			console.info("Writing output ", '../mitm_out/tag-krips.csv', data.length);
			fs.writeFileSync('../mitm_out/tag-krips.csv', data);
			console.log('done');
		});
		return;
	}
	if (process.argv[5] === 'consolidate') {
		console.info('consolidate!');
		consolidate(process.argv[2], process.argv.slice(3,5)).then((data) => {
			console.info("Writing output ", '../mitm_out/tag-c.csv', data.length);
			fs.writeFileSync('../mitm_out/tag-c.csv', data);
			console.log('done');
		});
		return;
	}
	if (process.argv[5] !== 'separate') { 
		main(process.argv[2], process.argv.slice(3)).then((data) => {
			// write it !
			console.info("Writing output ", '../mitm_out/tag-kappa.csv', data.length);
			fs.writeFileSync('../mitm_out/tag-kappa.csv', data);
			console.log('done');
		});
	} else {
		console.info('separate ');
		var tags = uniquetags(process.argv[2], process.argv[3], process.argv[4]);
		console.info('unique tags ', tags);		
		tags.map((t) => {
			main(process.argv[2], process.argv.slice(3,5), [t]).then((data) => {
				// write it !
				var fname = '../mitm_out/qual-'+t+'.csv';
				console.info("Writing output ", t, ' ---> ', fname, data.length);
				fs.writeFileSync(fname, data);
				console.log('done');
			});			
		});
	}

}