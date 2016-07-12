var parser = require('parse-whois');
var whois = require('node-whois');
 
whois.lookup('github.com', function(err, data){
	if (err) throw err;
 
	console.log(data);
 
	console.log(parser.parseWhoIsData(data));
});