
var parse = require('csv-parse/lib/sync'),
	csvstr = require('csv-stringify'),
	fs = require('fs'),
	promise = require('bluebird'),
	_ = require('lodash'),
	write_csv = (rows, fname) => new Promise((acc, rej) => {
		console.log('serialising csv ', rows.length);
		csvstr(rows, (err, output) => {
			if (!err) { 
				console.log("writing to ", output);
				fs.writeFileSync(fname, output);
			 	return acc(output);
			}
			console.error('error ', err);
			rej(err);
		});
	});

var tally_freq = (in_fname, out_fname) => {
	var app2libs = JSON.parse(fs.readFileSync(in_fname));
	console.log('ok! ', _.keys(app2libs).length, ' apps ');
	return _.reduce(app2libs, (dict, libs, app) => {
		libs.map((lib) => {
			dict[lib] = dict[lib] ? dict[lib]+1 : 1;
		});
		return dict;
	}, {});
};

if (require.main === module) { 
	var input = process.argv[2],
		freqout = process.argv[3];
	if (input) { 
		console.log('reading from ', input);
		var freq = _.toPairs(tally_freq(input));
		freq.sort((x,y) => y[1] - x[1]);
		write_csv(freq,freqout).then((rows) => console.log('written ', rows.length));
	}
}