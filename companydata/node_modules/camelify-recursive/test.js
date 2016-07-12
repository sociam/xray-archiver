var camelify = require('./index.js')

var log = console.log.bind(console)

var original = {
  "api_version": "0.3.2",
  "results": {
    "page": 1,
    "per_page": 30,
    "total_pages": 1,
    "total_count": 0,
    "filings": []
  }
}

log('ORIGINAL', original)

log('FINAL', camelify(original))
