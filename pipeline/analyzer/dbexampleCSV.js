var config = require('/etc/xray/scraping/config.json');
const pg = require('pg');
var db_cfg = config.scraper.db;
db_cfg.max = 10;
db_cfg.idleTimeoutMillis = 30000;
const pool = new pg.Pool(db_cfg);
pool.on('error', function (err, client) {
	console.error('idle client error', err.message, err.stack);
});
function connect() {
	return pool.connect();
};
function query(text, values) {
	console.log('query:', text, values);
	return pool.query(text, values);
};

//const text = 'INSERT INTO companies(id,company_ch,name,company_old,hosts,founded,acquired,type,typetag,jurisdiction,jurisdiction_code,parent,capital,equity,min_size,max_size,data_sources,description) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18),'
//const values = ['google', 'Google',	google.com google.com googleapis.com 1e100.com appspot.com ggpht.com googlevideo.com gstatic.com googleusercontent.com dmtry.com gvt1.com i.ytimg.com google.co.uk google.co.uk google.co.in blogspot.com	2002		analytics profiling advertising login	app	us					whois	Google is a multinational corporation that is specialized in internet-related services and products. The company’s product portfolio includes Google Search, which provides users with access to information online	
var columns = ['company_ch', 'id', 'company_old', 'name', 'hosts', 'founded', 'acquired', 'type', 'typetag', 'jurisdiction', 'jurisdiction_code', 'parent', 'capital', 'equity', 'min_size', 'max_size', 'data_sources', 'description'];
var ws = fs.createWriteStream('company-metadata-v4.csv');

pool.query = client.query('SELECT '+columns.join(', ')+' FROM companies');
query.on('row', function(row) {
    var values = [];
    // process column values; if you need to do special formatting (i.e. dates) don't loop and instead handle each one specially
    columns.forEach(function(col) {
        values = row[col];
    });
    ws.write(values.join('| '));
});

query.on('end', function(result) {
    ws.close();
});

/*async function insertCom(com) {
	var res = await pool.query("SELECT id FROM companies WHERE $1 = ANY(name);", [com.name]);
	if (res.length > 0) {
		return res.rows[0].id;
	}
	// maybe dev id needs to be URL encoded?
	let store_site = 'https://play.google.com/store/apps/developer?id='+com.id;
	res = await query("INSERT INTO companies(company_ch,name,company_old,hosts,founded,acquired,type,typetag,jurisdiction,jurisdiction_code,parent,capital,equity,min_size,max_size,data_sources,description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING id;",
	                  [com_ch, com.name, com_old, [com.hosts], com.founded, com.acquired, [com.type], [com.typetag], com.jurisdiction, com.jurisdiction_code, com.parent, com.capital, com.equity, com.min_size, com.max_size, [com.data_sources], com.description]);
	return res.rows[0].id
};*/