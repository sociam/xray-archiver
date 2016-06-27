
var parse = require('csv-parse/lib/sync'),
	fs = require('fs'),
	_ = require('lodash'),
	headers,
	hosts,
	counts,
	qs = require('querystring'),
	config = JSON.parse(fs.readFileSync('./config.json'));

console.log("Parsing file ", config.input);
data = parse(fs.readFileSync(config.input).toString());
headers = data[0];
data = data.slice(1);
data = data.map((x) => _.zipObject(headers,x));
hosts = _(data).map((x)=>x.host).uniq().value();
console.log('hosts : ', hosts);

// sum up counts
counts = _(data).reduce((y,x) => { y[x.host] = y[x.host] ? y[x.host] + 1 : 1; return y; },{})
console.log('counts: ', counts);

_(data).map((x) => console.log(x.params))

// decodes the urlstrings of the gets
var decode_urls = (datas) =>  {
	return datas.map((x) => {
		var url = decodeURIComponent(x.url);	 
		if (url.indexOf('?') >= 0) { 
			// chop off the querystring for urls that have it
			return qs.parse(url.slice(url.indexOf('?')+1));
		}
		// dont return anything for those that don't
	}).filter((x)=>x);
};

exports.data = data;
exports.counts = counts;
exports.decode_urls = decode_urls;

