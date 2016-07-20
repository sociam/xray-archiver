
var parse = require('csv-parse/lib/sync'),
	fs = require('fs'),
	_ = require('lodash'),
	headers,
	qs = require('querystring'),
	company_domains,
	COMPANY_DOMAINS = 'curated/company-domains.csv', 
	platform_companies,
	PLATFORM_COMPANIES = 'curated/platform-company.csv',
	config = JSON.parse(fs.readFileSync('./config.json')),
	detectors = require('./detect-pitypes').detectors;

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
		company_domains = loadFile(COMPANY_DOMAINS)
			.map((x) => { x.domains = x.domains.split(' '); return x; })
			.reduce((obj, x) => { obj[x.company] = x.domains; return obj; }, {});
	}
	return company_domains;
}, getPlatformCompanies = () => {
	if (platform_companies === undefined) { 
		platform_companies = loadFile(PLATFORM_COMPANIES)
			.reduce((obj, x) => { obj[x.platform] = x.company; return obj; }, {});
	}
	return platform_companies;
}, only_third_parties = (data) => {		
	return data.filter((r) => {
		// first attempt: try to filter out for hosts that have substrings with company or app name
		var appname = r.app.toLowerCase().split(' '),
			host = r.host.toLowerCase(),
			company = r.company && r.company.toLowerCase().split(' ').filter((x) => x.length > 2);

		if (company && _.some(company, (cfrag) => host.indexOf(cfrag) >= 0)) {
			// console.info('detected company name in hostname ', r.host, r.company);
			return false;
		}
		if (_.some(appname, (namefrag) => host.indexOf(namefrag) >= 0)) {
			// console.info('detected name frag in appname ', appname, r.host);
			return false;
		}
		if (_.some(getCompanyDomains()[r.company] || [], (serverTLD) => host.indexOf(serverTLD) >= 0)) {
			// console.info('detected company TLD ', appname, r.host);			
			return false;
		}
		// if (!r.platform) { console.warn(" no platform for ", r.app); }
		var pc = getPlatformCompanies()[r.platform.toLowerCase().trim()];
		// console.info('platform ', r.platform.toLowerCase(), pc, getCompanyDomains()[pc]);
		if (pc && _.some(getCompanyDomains()[pc], (plat) => r.host.indexOf(plat) >= 0)) { 
			// console.info('detected platform domain ', r.app, r.host);			
			return false;
		}
		return true;
	});
}, decodeURL = (url) => {
	url = decodeURIComponent(url);	 
	if (url.indexOf('?') >= 0) { 
		// chop off the querystring for urls that have it
		return qs.parse(url.slice(url.indexOf('?')+1));
	}
}, decode_all = (datas) => { 
	return datas.map((x) => decodeURL(x.url)).filter((x)=>x);
}, count_hosts = (data, app) => {
	if (app !== undefined) { data = _(data).filter((x) => x.app === app); }
	return _(data).reduce((y,x) => { y[x.host] = y[x.host] ? y[x.host] + 1 : 1; return y; },{});
}, getParty = (data, party) => {
	party = party.toLowerCase().trim();
	var hosts = getCompanyDomains()[party] || [];
	return data.filter((x) => {
		var host = x.host.toLowerCase();
		return _.some([ host.indexOf(party) >= 0 ].concat(hosts.map((h) => host.indexOf(h) >= 0)));
	});
}, getSLDs = () => {
	var ccsld = fs.readFileSync('curated/ccsld.txt').toString();
	return ccsld.split('\n').filter((x) => (x && x.trim().length > 0 && x.indexOf('.') >= 0 && x.indexOf('//') < 0 &&  x.indexOf('!') < 0 && x.indexOf('*') < 0));
}, decodeHeaders = (record) => record && record.headers && JSON.parse(decodeURIComponent(record.headers)
), decodeBody = (record) => record && record.body && decodeURIComponent(record.body),
decode = (record) => _.extend({}, decodeURL(record.url), decodeHeaders(record) || {}, decodeBody() || {}),
detect = (data) => {
	// returns an array one per data element
	// [ { record: {<data item>}, types: [ 'DEVICE_SOFT' ] } ... { } { } ]
	return data.map((x) => ({ record: x, decode: decode(x) })).map((pair) => {
		var d = pair.decode,
			types = _.keys(d).map((k) => detectors.map((detector) => { 
			// console.log(detector.type, ' testing ', k, d[k], detector.kv(k,d[k]), detector.kv(k,d[k]) ? true : false);
			return detector.kv(k,d[k]) ? detector.type : undefined; 
		}).filter((x) => x)).reduce((a,x) => a.concat(x), []);
		// console.info('types detected for ', d, _.uniq(types));
		return { record: pair.record, types: _.uniq(types) };
	});
}, detect_by_host = (detected, hostkey) => detected.reduce((dict, x) => {
	var host = hostkey && x.record[hostkey] || x.record.host;
	dict[host] = _.uniq((dict[host] || []).concat(x.types))
	return dict;
}, {}), hosts_by_app = (data, hostkey) => {
	return _(data).reduce((y,x) => { 
		var host = hostkey && x[hostkey] || x.host;
		y[x.app] = y[x.app] || {};
		y[x.app][host] = y[x.app][host] ? y[x.app][host] + 1 : 1; 
		return y;
	},{});
}, shorten_2ld = (host) => {
	var match = host.match(/([^\.]*)\.([^\.]*)$/);
	if (match) { 
		var short = match[0];
		if (exports.ccslds.indexOf(short) >= 0) { 
			var onemore = host.slice(0,host.length - short.length - 1).match(/([^\.]*)$/);
			if (onemore) { 
				return [onemore[0], short].join('.');
			} else {
				// fallback
				return host; 
			}
		}
		return short;
	}
	return host;
}, fold_into_2ld = (data) => {
	data.map((x) => { x.host_2ld = shorten_2ld(x.host);	});
	return data;
};

exports.decode_all = decode_all;
exports.count_hosts = count_hosts;
exports.only_third_parties = only_third_parties;
exports.getCompanyDomains = getCompanyDomains;
exports.getPlatformCompanies = getPlatformCompanies;
exports.load = loadDir;
exports.getParty = getParty;
exports.decodeURL = decodeURL;
exports.decodeHeaders = decodeHeaders;
exports.decodeBody = decodeBody;
exports.decode = decode;
exports.ccslds = getSLDs();
exports.data = fold_into_2ld(loadDir());
exports.detect = detect;
exports.detected = detect_by_host(detect(exports.data));
exports.hosts_by_app = hosts_by_app(exports.data); // hosts_by_app(exports.data, 'host_2ld');
exports.detectors = detectors;

var main = (app) => { 
	// var data = loadDir();
	var data = exports.data;
	console.log('decoded urls', decode_all(data)); 
	console.log('count hosts ', count_hosts(only_third_parties(data), app));
	console.log('hba', exports.hosts_by_app);

	if (config.out_hosts_by_app) { 
		console.info("writing Hosts By App table to:", config.out_hosts_by_app, _.keys(exports.hosts_by_app).length, ' apps');
		fs.writeFileSync(config.out_hosts_by_app, JSON.stringify(exports.hosts_by_app));
	}
	if (config.out_pi_by_host) { 
		console.info("writing PI types by host to:", config.out_pi_by_host, _.keys(exports.detected).length, ' hosts');
		fs.writeFileSync(config.out_pi_by_host, JSON.stringify(exports.detected));
	}
	if (config.out_data) { 
		console.info("writing all data records to:", config.out_data, _.keys(exports.data).length, ' records');
		fs.writeFileSync(config.out_data, JSON.stringify(exports.data));
	}
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