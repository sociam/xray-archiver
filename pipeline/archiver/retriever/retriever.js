const gplay = require('google-play-scraper');
const logger = require('../../util/logger');
const db = new(require('../../db/db'))('retriever');

const region = 'uk';

/**
 * Inserts app data into the db using db.js
 * @param {*The app data json that is to be inserted into the databae.} appData
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
    // logger.debug(appData.appId);
}

// TODO Add Permission list to app Data JSON
async function fetchAppData(searchTerm, numberOfApps, perSecond) {
    let res = [];

    try {
        res = await gplay.search({
            term: searchTerm,
            num: numberOfApps,
            throttle: perSecond,
            region: region,
            fullDetail: true,
        });
    } catch (err) {
        logger.err(err);
    }
    return res;
}

function scrapeFromSearchTerm() {
    return db.getStaleSearchTerms(1)
        .then(
            (searchTerms) => {
                const terms = searchTerms.map((term) => term.search_term);
                logger.debug('Marking Search terms as being used:', terms);
                terms.forEach((term) => db.updateLastSearchedDate(term)
                    .then(() => logger.debug('Last Search date updated:', term))
                    .catch((err) => logger.err('ERROR SETTING SEARCH DATE FOR:', term, err))
                );
                logger.debug('Fetching App data for searchterms:', terms);
                return Promise.all(terms.map((term) => fetchAppData(term, 120, 1) /* Need to catch...*/ ));
            },
            (err) => {
                logger.err('ERROR READING SEARCH TERMS FROM DATABASE:', err);
            }
        ).then(
            (appDatas) => {
                const appData = appDatas.reduce((a, b) => a.concat(b), []);
                logger.debug('Filtering apps that already exist: ',
                    appData.map((data) => data.appId));
                const data = Promise.all(appData.forEach((appData) => {
                    db.doesAppExist(appData)
                        .then((res) => res ? logger.debug('App Not Added') : insertAppData(appData));
                })); /* Need to catch...*/
                logger.debug(data);
                return data;
            }, (err) => logger.err('ERROR FETCHING APP DATA FOR SEARCH TERMS:', err)
        ).then(
            (appDatas) => {
                logger.debug('App Data doesn\'t exist:', appDatas.map((data) => data.appId));
                appDatas.forEach((data) => insertAppData(data)
                    .catch((err) => logger.err('ERROR INSERTING APP INTO DATABASE:', err)));
            },
            (err) => logger.err('ERROR CHECKING EXISTANCE OF AN APP:', err)
        );
}

async function scrape() {
    await scrapeFromSearchTerm()
        .then(() => {
            scrape();
        })
        .catch((err) => logger.err('ERROR:', err));
    logger.debug('Moving onto the next search term');
}

scrape();