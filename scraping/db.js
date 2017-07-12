
var config = require('/etc/xray/config.json');
const pg = require('pg');

var db_cfg = config.scraper.db;
db_cfg.max = 10;
db_cfg.idleTimeoutMillis = 30000;
   
//this initializes a connection pool 
//it will keep idle connections open for 30 seconds 
//and set a limit of maximum 10 idle clients 
const pool = new pg.Pool(config);

pool.on('error', function (err, client) {
   console.error('idle client error', err.message, err.stack);
});
 
//export the query method for passing queries to the pool 
module.exports.query = function (text, values, callback) {
  console.log('query:', text, values);
  return pool.query(text, values, callback);
};
 
// the pool also supports checking out a client for 
// multiple operations, such as a transaction 
module.exports.connect = function (callback) {
  return pool.connect(callback);
};

module.exports.connect((thing) => {
	console.log(thing);
    console.log("Connected");
    module.exports.query("SELECT * FROM apps");
});

	
