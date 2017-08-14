'use strict';
/*global require, module */

const config = require('/etc/xray/config.json');
const pg = require('pg');
const logger = require('../util/logger');

class DB {
    constructor(module) {
        //WISHLIST: initialise pool for desired db option from the config here.
        var db_cfg = config.db;
        db_cfg.user = config[module].db.user;
        db_cfg.password = config[module].db.password;
        db_cfg.max = 10;
        db_cfg.idleTimeoutMillis = 30000;

        //this initializes a connection pool
        //it will keep idle connections open for 30 seconds
        //and set a limit of maximum 10 idle clients
        this.pool = new pg.Pool(db_cfg);

        this.pool.on('error', function(err) {
            logger.err('idle client error', err.message, err.stack);
        });
    }


    //export the query method for passing queries to the pool
    query(text, values) {
        if (values) logger.debug('query:', text, values);
        else logger.debug('query:', text, values);
        return this.pool.query(text, values);
    }

    async connect() {
        logger.debug('connecting to db pool');
        let ret = await this.pool.connect();
        ret.lquery = (text, values) => {
            if (values) logger.debug('lquery:', text, values);
            else logger.debug('lquery:', text);
            return ret.query(text, values);
        };

        return ret;
    }

    async insertDev(dev) {
        try {
            var res = await this.query('SELECT id FROM developers WHERE $1 = ANY(email)', [dev.email]);
            if (res.rowCount > 0) {
                logger.debug('developer with email %s exists', dev.email);
                return res.rows[0].id;
            }

            logger.debug('inserting developer with email %s', dev.email);

            // maybe dev id needs to be URL encoded?
            let store_site = 'https://play.google.com/store/apps/developer?id=' + dev.id;
            res = await this.query('INSERT INTO developers(email,name,store_site,site) VALUES ($1, $2, $3, $4) RETURNING id', [
                [dev.email], dev.name, store_site, dev.site
            ]);
        } catch (err) { logger.err(err); }
        return res.rows[0].id;
    }

    async updatedLastAltChecked(appID) {
        var client = await this.connect();

        try {
            await client.lquery('UPDATE app_versions SET last_alt_checked = CURRENT_TIMESTAMP where app = $1', [appID]);

        } catch (err) {
            logger.err(err);
        } finally {
            client.release();
        }
    }

    async insertAltApp(altApp) {
        // check if the current alt app exists in the db and has been analysed.
        var isCollected = false;

        if (altApp.gPlayID != '') {
            var analysedRes = await this.query('SELECT analyzed from app_versions where app = $1 ', [altApp.gPlayAppID]);
            if (analysedRes.rowCount > 0) {
                isCollected = true;
            }
        }

        var client = await this.connect();
        logger.debug('Connected');

        var checkRes = await client.lquery(
            'SELECT * FROM alt_apps WHERE alt_app_title = $1 and app_id = $2', [altApp.title, altApp.appID]
        );

        logger.debug(checkRes.rowCount);

        if (checkRes.rowCount == 0) {
            try {
                await client.lquery('BEGIN');
                await client.lquery(
                    'INSERT INTO alt_apps( \
                        app_id, alt_app_title, alt_to_url, g_play_url, g_play_id, icon_url, official_site_url, is_scraped) \
                            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)', [
                        altApp.appID,
                        altApp.title,
                        altApp.altToURL,
                        altApp.gPlayURL,
                        altApp.gPlayAppID,
                        altApp.altAppIconURL,
                        altApp.officialSiteURL,
                        isCollected
                    ]);
                logger.debug(altApp.title + ' added to DB');
                await client.lquery('COMMIT');
            } catch (err) {
                logger.err('Error with previous query:', err);
                await client.lquery('ROLLBACK');
            } finally {
                client.release();
            }
        } else {
            logger.debug('%s already exists, skipping.', altApp.title);
            client.release();
        }
    }

    async getAppsToFindAltsForThatHaventYetHadThemFound(limit) {
        var res = await this.query(
            'SELECT a.title, v.app FROM app_versions v FULL OUTER JOIN playstore_apps a \
                ON (a.id = v.id) \
                    WHERE v.app NOT IN (SELECT app_id AS app FROM alt_apps) \
                         ORDER BY last_alt_checked NULLS FIRST LIMIT $1', [limit]);

        if (res.rowCount <= 0) {
            logger.debug('No apps need alternatives to be searched for. Or something has screwed up');
        }

        logger.info(res.rowCount, 'Apps need alternatives to be scraped.');
        return res.rows;
    }

    async queryAppsToDownload(batch) {
        var res = await this.query('SELECT * FROM (SELECT * FROM app_versions ORDER BY last_dl_attempt) AS apps WHERE downloaded = False LIMIT $1 ', [batch]);

        if (res.rowCount <= 0) {
            return Promise.reject('No downloads found. Consider slowing down downloader or speeding up scraper');
        }
        logger.info('Found apps to download:', res.rowCount);
        return res.rows;
    }

    async updateDownloadedApp(app) {
        await this.query('UPDATE app_versions SET downloaded=True WHERE app = $1', [app.app]);
    }

    async updatedDlAttempt(app) {
        await this.query('UPDATE app_versions SET last_dl_attempt=CURRENT_TIMESTAMP WHERE app = $1', [app.app]);
    }

    async getAppData() {
        logger.debug('Fetching Search Terms');
        var res = await this.query('SELECT search_term FROM search_terms WHERE age(last_searched) > interval \'1 month\'');
        logger.debug(res.rows.length + ' terms fetched');
        return res.rows;
    }

    /**
     *  Query the search_terms table to get a list of terms that are stale
     */
    async getStaleSearchTerms() {
        logger.debug('Fetching Search Terms');
        var res = await this.query('SELECT search_term FROM search_terms WHERE age(last_searched) > interval \'1 month\'');
        logger.debug(res.rows.length + ' terms fetched');
        return res.rows;
    }

    /**
     *  Insert A Manually Curated App and ID into the Database
     */
    async insertManualSuggestion(altPair) {
        logger.debug('Inserting - Source: ' + altPair.source + '  Alt:' + altPair.alt);

        try {

            logger.debug('checking if ' + altPair.source + ' and ' + altPair.alt + ' exists in db.');
            var check_res = await this.query('SELECT source_id, alt_id FROM manual_alts WHERE source_id = $1 and alt_id = $2', [altPair.source, altPair.alt]);
            logger.debug(check_res.rowCount + ' rows found for ' + altPair.source + ' and ' + altPair.alt);

            if (check_res.rowCount) {
                logger.debug(altPair.source + ' and ' + altPair.alt + ' in the manual alt table');
                await this.query('INSERT INTO manual_alts VALUES ($1, $2)', [altPair.source, altPair.alt]);
                logger.debug(altPair.source + ' and ' + altPair.alt + ' Added to the DB.');
            }
        } catch (err) {
            logger.err('ERROR!!!!!' + err);
        }
    }

    /**
     * Sets the last_searched date of a specified search term to be the current date.
     * Used to track 'stale' search terms.
     */
    async updateLastSearchedDate(searchTerm) {
        logger.debug('Setting last searched date for ' + searchTerm + ' to current date');
        var client = await this.connect();
        logger.debug('connected');

        logger.debug('checking if ' + searchTerm + ' exists in db.');
        var check_res = await client.lquery('SELECT search_term FROM search_terms WHERE search_term = $1', [searchTerm]);
        logger.debug(check_res.rowCount + ' rows found for ' + searchTerm);
        if (check_res.rowCount > 0) {
            logger.debug(searchTerm + ' exists, updating last searched date.');
            var update_res = await client.lquery('UPDATE search_terms SET last_searched = CURRENT_DATE WHERE search_term = $1', [searchTerm]);
        }
        client.release();
        return update_res;
    }

    /**
     *  Add a search term to the table if it doesn't already exist.
     */
    async insertSearchTerm(searchTerm) {
        var client = await this.connect();
        logger.debug('Connected');

        logger.debug('Checking if ' + searchTerm + ' exists before adding to search_terms');
        var checkRes = await client.lquery('SELECT search_term FROM search_terms WHERE search_term = $1', [searchTerm]);

        if (checkRes.rowCount == 0) {
            try {
                await client.lquery('BEGIN');
                await client.lquery('INSERT INTO search_terms VALUES ($1, \'epoch\')', [searchTerm]);
                logger.debug(searchTerm + ' added to DB');
                await client.lquery('COMMIT');
            } catch (err) {
                logger.err('Error with previous query:', err);
                await client.lquery('ROLLBACK');
            } finally {
                client.release();
            }
        } else {
            logger.debug('%s already exists, skipping.', searchTerm);
            client.release();
        }
    }

    /**
     *  Check if app already exists in the apps DB. used before attempting to log again.
     */
    async doesAppExist(app) {
        var res = await this.query('SELECT * FROM apps WHERE id = $1', [app.appId]);
        logger.debug('app query res count: ' + res.rowCount);
        return (res.rowCount > 0);
    }

    /**
     *  Inserts App Data scraped from the google play store into the DB.
     */
    async insertPlayApp(app, region) {

        logger.debug('Inserting Dev Information. ');
        var devId = await this.insertDev({
            name: app.developer,
            id: app.developerId,
            email: app.developerEmail,
            site: app.developerWebsite
        });


        var appExists = false,
            verExists = false;
        var verId;
        var res = await this.query('SELECT * FROM apps WHERE id = $1', [app.gPlayAppID]);

        if (res.rowCount != 0) {
            appExists = true;
            // app exists in database, check if version does as well
            var res1 = await this.query(
                'SELECT id FROM app_versions WHERE app = $1 AND store = $2 AND region = $3 AND version = $4', [app.appId, 'play', region, app.version]);

            if (res1.rowCount > 0) {
                // app version is also in database
                verExists = true;
                verId = res1.rows[0].id;
            }
        }

        var client = await this.connect();
        logger.debug('Connected');

        await client.lquery('BEGIN'); // maybe this should be inside the try?
        try {
            if (!verExists) {
                if (!appExists) {
                    await client.lquery('INSERT INTO apps VALUES ($1, $2)', [app.appId, []]);
                }
                let res = await client.lquery(
                    'INSERT INTO app_versions(app, store, region, version, downloaded, last_dl_attempt, analyzed )' +
                    'VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id', [app.appId, 'play', region, app.version, 0, 'epoch', 0]
                );
                verId = res.rows[0].id;

                await client.lquery('UPDATE apps SET versions=versions || $1 WHERE id = $2', [
                    [verId], app.appId
                ]);

            }

            await client.lquery(
                'INSERT INTO playstore_apps VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, current_date)', [
                    verId,
                    app.title,
                    app.summary,
                    app.description,
                    app.url,
                    app.price,
                    app.free,
                    app.score,
                    app.reviews,
                    app.genreId,
                    app.familyGenreId,
                    app.minInstalls,
                    app.maxInstalls,
                    devId,
                    app.updated,
                    app.androidVersion,
                    app.contentRating,
                    app.screenshots,
                    app.video,
                    app.recentChanges
                ]);
            await client.lquery('COMMIT');
        } catch (e) {
            logger.debug(e);
            await client.lquery('ROLLBACK');
            throw e;
        } finally {
            client.release();
        }
        return verId;
    }

}

module.exports = DB;