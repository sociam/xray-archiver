
var fs = require('fs'),
	parse = require('csv-parse/lib/sync'),
	csvstr = require('csv-stringify'),
	_ = require('lodash'),	
    loadCSV = (fname) => {
		var text = 	fs.readFileSync(fname).toString();
		console.error("Parsing file ", fname, "(", text.length, ")");
		var data = parse(text, {max_limit_on_data_read:9999999999});
		var headers = data[0];
		data = data.slice(1);
		data = data.map((x) => _.zipObject(headers,x));
		return data;
	},	
	write_csv = (rows, fname) => new Promise((acc, rej) => {
		console.info('serialising csv ', rows.length);
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

exports.loadCSV = loadCSV;
exports.saveCSV = write_csv;