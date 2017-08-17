// External Packages
const gplay = require('google-play-scraper');

// Our own packages.
const logger = require('../../util/logger');
const Database = require('../../db/db');
const db = new Database('retriever');

/**
 * Process Command line arguments
 * Checks that there is exactly one argument passed to this script.
 * Falls over with a Donezo if there isnt.
 */
function processArgs() {
    // Process command line arguments. Expecing only one. an app ID.
    const arg = process.argv.slice(2);

    // Fail if there is more than one argument passed.
    if (arg.length != 1) {
        throw Error(`DONEZO:\n\tExpected eactly one argument. Actually got: ${arg.length}.` +
                    `\n\tArguments passed are: ${arg}`);
    }

    return arg[0];
}

function validAppID(appID) {
    return /^[a-zA-Z0-9.]*$/.test(appID);
}

function insertAppData(appData) {
    const formatDate = new Date().toISOString().substring(0, 10);

    // Checking version data - correct version to update date
    if (!appData.updated) {
        logger.debug('No updated value found. Setting to todays date.');
        appData.updated = formatDate;
    }

    if (!appData.version || appData.version === 'Varies with device') {
        logger.debug('Version not found defaulting too', appData.updated);

        // let formatDate = appData.updated.replace(/\s+/g, '').replace(',', '/');
        const formatDate = new Date(appData.updated).toISOString().substring(0, 10);
        appData.version = formatDate;
    }

    // push the app data to the DB
    return db.insertPlayApp(appData, 'us');
}

function fetchAppData(appID) {
    return gplay.app({
        appId: appID,
        fullDetails: true,
    }).then((res) => {
        // logger.debug(res);
        return insertAppData(res)
            .then((res) => {
                return `App inserted into the DB. Response: ${res}`;
            })
            .catch((err) => {
                return err;
            });
    }).catch((err) => {
        logger.err(err);
        return err;
    });
}

function main() {
    const arg = processArgs();
    logger.debug(arg);

    const app = {
        appId: arg,
    };

    return db.doesAppExist(app)
        .then((res) => {
            if (res) {
                logger.debug('App Already Exists. just going to stop it riiiight here.');
                process.exit(0);
            }

            if (!validAppID(arg)) {
                logger.debug('App ID is actually invalid. you\'re donzo!');
                process.exit(0);
            }

            logger.debug(`App ID is fine. lets get cooking! fetching: ${arg}`);

            return fetchAppData(arg)
                .then(() => {
                    logger.debug('AppData fetched and inserted to the DB. We\'re Donezo');
                    process.exit(0);
                })
                .catch((err) => {
                    logger.err(`Error Occured. failed to insert App Data.${err}`);
                });
        });
}
main();
