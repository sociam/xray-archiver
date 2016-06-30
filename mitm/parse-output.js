
var parse = require('csv-parse/lib/sync'),
	fs = require('fs'),
	_ = require('lodash'),
	headers,
	hosts,
	counts,
	qs = require('querystring'),
	config = JSON.parse(fs.readFileSync('./config.json'));

var loadFile = (fname) => {
	console.log("Parsing file ", fname);
	data = parse(fs.readFileSync(fname).toString());
	headers = data[0];
	data = data.slice(1);
	data = data.map((x) => _.zipObject(headers,x));
	// for giggles
	// hosts = _(data).map((x)=>x.host).uniq().value();	
	// counts = _(data).reduce((y,x) => { y[x.host] = y[x.host] ? y[x.host] + 1 : 1; return y; },{})
	// console.log('counts: ', counts);
	// _(data).map((x) => console.log(x.params))
	return data;
}, loadDir = () => {
	// loads all of the data in the specified directory
	var srcdir = config.inputdir;
	return fs.readdirSync(srcdir).reduce((arr,fname) => arr.concat(loadFile([srcdir,fname].join('/'))), []);
}, decode_urls = (datas) =>  {
	return datas.map((x) => {
		var url = decodeURIComponent(x.url);	 
		if (url.indexOf('?') >= 0) { 
			// chop off the querystring for urls that have it
			return qs.parse(url.slice(url.indexOf('?')+1));
		}
		// dont return anything for those that don't
	}).filter((x)=>x);
};

exports.decode_urls = decode_urls;
exports.load = loadDir;

