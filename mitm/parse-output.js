
var parse = require('csv-parse/lib/sync'),
	fs = require('fs'),
	_ = require('lodash'),
	headers,
	hosts,
	counts;

data = parse(fs.readFileSync('/tmp/out.csv').toString());
headers = data[0];
data = data.slice(1);
data = data.map((x) => _.zipObject(headers,x));
hosts = _(data).map((x)=>x.host).uniq().value();

console.log('hosts : ', hosts);

// sum up counts
counts = _(data).reduce((y,x) => { y[x.host] = y[x.host] ? y[x.host] + 1 : 1; return y; },{})
console.log('counts: ', counts);

