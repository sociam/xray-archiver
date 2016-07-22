var parser = require('parse-whois');
var whois = require('node-whois');
var openCorporates = require('opencorporates')()

host = 'xiaomi.com'

company = ''
 
whois.lookup(host, function(err, data){
	if (err) throw err;
	parsed = parser.parseWhoIsData(data);
	company = parsed[14].value
	console.log(company)
	openCorporates.companies.search('xiamoi Inc', function(err, res){
	console.log(JSON.stringify(res, null, 2))
})
});