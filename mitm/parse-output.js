
var parse = require('csv-parse/lib/sync'),
	fs = require('fs'),
	_ = require('lodash'),
	headers,
	qs = require('querystring'),
	company_domains,
	COMPANY_DOMAINS = 'curated/company-domains.csv', 
	config = JSON.parse(fs.readFileSync('./config.json'));

var loadFile = (fname) => {
	var text = 	fs.readFileSync(fname).toString();
	console.log("Parsing file ", fname, "(", text.length, ")");
	var data = parse(text, {max_limit_on_data_read:9999999999});
	headers = data[0];
	data = data.slice(1);
	data = data.map((x) => _.zipObject(headers,x));
	return data;
}, loadDir = () => {
	// loads all of the data in the specified directory
	var srcdir = config.inputdir;
	return fs.readdirSync(srcdir)
		.filter((fname) => fname.indexOf('.csv') >= 0)
		.reduce((arr,fname) => arr.concat(loadFile([srcdir,fname].join('/'))), []);
}, getCompanyDomains = () => {
	if (company_domains === undefined) { 
		company_domains = loadFile(COMPANY_DOMAINS).map((x) => { x.domains = x.domains.split(' '); return x; });
	}
	return company_domains;
}, only_third_parties = (data) => {		
	return data.filter((r) => {
		// first attempt: try to filter out for hosts that have substrings with company or app name
		var appname = r.app.toLowerCase().split(' '),
			host = r.host.toLowerCase(),
			company = r.company && r.company.toLowerCase().split(' ').filter((x) => x.length > 2);

		if (company && _.some(company, (cfrag) => host.indexOf(cfrag) >= 0)) {
			console.info('detected company name in hostname ', r.host, r.company);
			return false;
		}
		if (_.some(appname, (namefrag) => host.indexOf(namefrag) >= 0)) {
			console.info('detected name frag in appname ', appname, r.host);
			return false;
		}
		if (_.some(getCompanyDomains()[r.company] || [], (serverTLD) => host.indexOf(serverTLD) >= 0)) {
			console.info('detected company TLD ', appname, r.host);			
			return false;
		}
		if (!r.platform) { console.warn(" no platform for ", r.app); }
		if (_.some(getCompanyDomains()[r.platform], (plat) => r.host.indexOf(plat) >= 0)) { 
			console.info('detected platform domain ', r.app, r.host);			
			return false;
		}
		return true;
	});

}, decode_url = (url) => {
	url = decodeURIComponent(url);	 
	if (url.indexOf('?') >= 0) { 
		// chop off the querystring for urls that have it
		return qs.parse(url.slice(url.indexOf('?')+1));
	}
}, decode_all = (datas) => { 
	return datas.map((x) => decode_url(x.url)).filter((x)=>x);
}, count_hosts = (data, app) => {
	if (app !== undefined) { data = _(data).filter((x) => x.app === app); }
	return _(data).reduce((y,x) => { y[x.host] = y[x.host] ? y[x.host] + 1 : 1; return y; },{});
};


exports.decode_url = decode_url;
exports.decode_all = decode_all;
exports.count_hosts = count_hosts;
exports.only_third_parties = only_third_parties;
exports.load = loadDir;

var main = (app) => { 
	var data = loadDir();
	console.log('decoded urls', decode_all(data)); 
	console.log('count hosts ', count_hosts(only_third_parties(data), app));
};

if (require.main === module) { 
	// console.info(process.argv.length);
	// process.argv.forEach(function (val, index, array) {
	//   console.log(index + ': ' + val);
	// });	
	if (process.argv.length === 3) { 
		// console.log('fo ', process.argv.length);
		return main(process.argv[2]);
	}  else {
		main(); 
	}
}