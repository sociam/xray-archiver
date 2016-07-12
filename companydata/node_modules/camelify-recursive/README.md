# Recursively convert keys to camelCase

Have an API response from someone that doesn't use JavaScript?

Want to be able to use the keys without having to fix the case all the time?

Let camelify do the work.

It will **recursively** convert all the keys in the object to camelCase, so you can use those keys directly in your JavaScript (or JSON database) and not have to worry about converting things manually all the time.

## Usage

Just:

	var camelify = require('camelify-recursive');

Then, to fix the object `someObject`

	someObject = camelify(someObject);

## Yes, this really does what you think it does

For example, this a camelified result from  the OpenCorporates API, which was once filled with underscores:

	{
		"results": {
			"company": {
				"name": "GITHUB, INC.",
				"inactive": false,
				"source": {
					"publisher": "California Secretary of State",
					"url": "http://kepler.sos.ca.gov/",
					"retrievedAt": "2012-04-03T07:19:16+00:00"
				},
				"data": null,
				"filings": [],
				"officers": [],
				"companyNumber": "C3268102",
				"jurisdictionCode": "us_ca",
				"incorporationDate": "2009-12-31",
				"dissolutionDate": null,
				"companyType": "Domestic Stock",
				"registryUrl": "https://businessfilings.sos.ca.gov/frmDetail.asp?CorpID=03268102",
				"branchStatus": null,
				"currentStatus": "Active",
				"createdAt": "2011-09-17T15:33:31+00:00",
				"updatedAt": "2013-10-27T06:27:24+00:00",
				"retrievedAt": "2012-04-03T07:19:16+00:00",
				"opencorporatesUrl": "https://opencorporates.com/companies/us_ca/C3268102",
				"previousNames": [],
				"agentName": "BRANDEE MURPHY",
				"agentAddress": "582 MARKET STREET STE 1700, SAN FRANCISCO, CA 94104",
				"registeredAddressInFull": "548 4TH STREET, SAN FRANCISCO, CA 94107",
				"registeredAddress": {
					"locality": null,
					"region": null,
					"country": "United States",
					"streetAddress": "548 4TH STREET, SAN FRANCISCO, CA 94107",
					"postalCode": null
				},
				"corporateGroupings": [],
				"industryCodes": [],
				"financialSummary": null,
				"homeCompany": null,
				"controllingEntity": null
			}
		},
		"apiVersion": "0.3.2"
	}

## Can I have emoji camels?

🐫🐪🐫🐪🐫🐪🐫🐪🐫🐪


