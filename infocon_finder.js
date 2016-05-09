// First install unirest
// npm install unirest

var unirest = require('unirest');

unirest.get('https://jsonwhois.com/api/v1/whois')
 .headers({
    'Accept': 'application/json',
    'Authorization': 'Token token=<Api Key>'
 })

   .query({
      "domain": "google.com"
       })

   .end(function (response) {
        console.log(response.body);
});

record.WhoisRecord.registrant.organization
