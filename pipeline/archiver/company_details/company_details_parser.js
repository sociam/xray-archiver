'use-strict';
const fs = require('fs');
const logger = require('../../util/logger');

const DB = require('../../db/db');
const db = new DB('suggester');

fs.readFile('./company_details.json', (err, data) => {
    if (err) {
        logger.debug(err);
        process.abort();
    }
    const companies = JSON.parse(data);
    logger.debug(`${Object.keys(companies).length} companies to log`);
    Object.keys(companies).map((key) =>
        companies[key].domains.map((domain) =>
            db.insertCompanyDomain(key, domain)));
});
