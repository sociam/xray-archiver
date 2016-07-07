
var parse = require('csv-parse/lib/sync'),
	fs = require('fs'),
	_ = require('lodash'),
	headers,
	hosts,
	counts,
	qs = require('querystring'),
	config = JSON.parse(fs.readFileSync('./config.json'));

var loadFile = (fname) => {
	text = 	fs.readFileSync(fname).toString()
	console.log("Parsing file ", fname, "(", text.length, ")");
	data = parse(text, {max_limit_on_data_read:9999999999});
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
	return fs.readdirSync(srcdir)
		.filter((fname) => fname.indexOf('.csv') >= 0)
		.reduce((arr,fname) => arr.concat(loadFile([srcdir,fname].join('/'))), []);
}, decode_urls = (datas) =>  {
	return datas.map((x) => {
		var url = decodeURIComponent(x.url);	 
		if (url.indexOf('?') >= 0) { 
			// chop off the querystring for urls that have it
			return qs.parse(url.slice(url.indexOf('?')+1));
		}
		// dont return anything for those that don't
	}).filter((x)=>x);
}, count_hosts = (data, app) => {
	hosts = _(data).filter((x) => x.app == app).map((x)=>x.host).uniq().value();	
	counts = _(data).reduce((y,x) => { y[x.host] = y[x.host] ? y[x.host] + 1 : 1; return y; },{})
};

exports.decode_urls = decode_urls;
exports.load = loadDir;

