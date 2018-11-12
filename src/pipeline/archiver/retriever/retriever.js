
const gplay = require('google-play-scraper');

const logger = require('../../util/logger');
const db = new (require('../../db/db'))('retriever');

const region = 'uk';

/*
 * Inserts app data into the db using db.js
 */
function insertAppData(appData) {
    // Checking version data - correct version to update date
    if (!appData.version || appData.version === 'Varies with device') {
        logger.debug('Version not found defaulting too', appData.updated);
        // let formatDate = appData.updated.replace(/\s+/g, '').replace(',', '/');
        const formatDate = new Date(appData.updated).toISOString().substring(0, 10);
        appData.version = formatDate;
    }

    // push the app data to the DB
    return db.insertPlayApp(appData, region);
}

// TODO Add Permission list to app Data JSON
async function fetchAppData(searchTerm, numberOfApps, perSecond) {
    const appSearchResults = await gplay.search({
        term: searchTerm,
        num: numberOfApps,
        throttle: perSecond,
        country: region,
        // fullDetail: true,
    });

    for (const result of appSearchResults) {
        logger.debug(`Requsting App Data for: ${result.appId}`);
        await gplay.app(
            { appId: result.appId }
        )
            .then(
                async(appData) => {
                    logger.debug(`inserting ${appData.title} to the DB`);
                    await insertAppData(appData).catch((err) => logger.err(err));
                },
                (err) => logger.err(
                    `Error Requesting appData for App: ${result.appId}. Error: ${err}`
                )
            );
    }
}

(async() => {
    const dbRows = await db.getStaleSearchTerms();
    for (const dbRow of dbRows) {
        logger.info(`searching for: ${dbRow.search_term}`);
        await fetchAppData(dbRow.search_term, 250, 10);
        await db.updateLastSearchedDate(dbRow.search_term);
    }
})();
