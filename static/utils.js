
var fs = require('fs'),
	parse = require('csv-parse/lib/sync'),
	csvstr = require('csv-stringify'),
	_ = require('lodash'),	
    loadCSV = (fname, columns) => {
		var text = 	fs.readFileSync(fname).toString();
		console.error("Parsing file ", fname, "(", text.length, ")");
		var data = parse(text, {max_limit_on_data_read:9999999999});
		var headers = columns || data[0];
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
	}),
	getSLDs = () => {
		var ccsld = fs.readFileSync('../mitm/curated/ccsld.txt').toString();
		return ccsld.split('\n').filter((x) => (x && x.trim().length > 0 && x.indexOf('.') >= 0 && x.indexOf('//') < 0 &&  x.indexOf('!') < 0 && x.indexOf('*') < 0));
	}, 	
	ipv4re = /((^\s*((([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))\s*$)|(^\s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?\s*$))/,
	shorten_2ld = (host) => {
		host = host.trim().toLowerCase();
		// raw ip addresses don't get shortened
		if (ipv4re.test(host)) { 
			console.log('detected an ipv4 address, skipping shortening', host);
			return host;	
		}
		var match = host.match(/([^\.]*)\.([^\.]*)$/);
		if (match) { 
			var short = match[0];
			if (exports.ccslds.indexOf(short) >= 0) { 
				var onemore = host.slice(0,host.length - short.length - 1).match(/([^\.]*)$/);
				if (onemore) { 
					return [onemore[0], short].join('.');
				} else {
					return host; 
				}
			}
			return short;
		}
		return host;
	};


exports.ccslds = getSLDs();
exports.loadCSV = loadCSV;
exports.saveCSV = write_csv;