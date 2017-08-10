// Other's packages
const request = require('request');
const cheerio = require('cheerio');
const _ = require('lodash');

// Our packages
const logger = require('../../util/logger');
const Database = require('../../db/db');
var db = new Database('suggester');


/**
 * Fetch fetchAlternatives
 * 
 * For a specific App's Alternative page. Collect Titles and AltTo URLs
 * that are lsited for that app.
 * 
 * The urls collected here will point towards another alternativeTo Page
 * that will have a store link (hopefully)
 */
function scrapePageForAlts(URLString) {

    //TODO: Check that page is an app page.
    // Request the page from the Site.
    request(URLString, (err, res, html) => {

        // if there wasn't an err with the request.
        if (err) {
            logger.err('Failed to fetch HTML: ' + err);
            return err;
        }

        // Initialising Variables and Loading HTML into Cheerio
        var $ = cheerio.load(html);
        var altIDs = [];
        var altApps = [];


        // Selecting All Alternative App ID's from the DOM
        $('ul#alternativeList').find('li').each((i, elem) => {
            // Each alt app has an ID for the app title.
            var val = $(elem).attr('id');
            // if the element exists on the page.
            if (val) {
                altIDs.push(val);
            }
        });

        // Bail if there were no alternative IDs.
        if (!altIDs.length) { 
            logger.debug('No Alternative Apps suggested.');
            return 'No Alternative Apps suggested';
        }

        // for each ID find the Title and URL to AltTo page.
        altApps = _.map(altIDs, (altID) => {
            //logger.debug(altID);
            var aTag = $('#' + altID).find('h3').find('a').first();
            return {
                'title': aTag.text(),
                'altToURL': 'http://alternativeto.net' + aTag.attr('href'),
                'altAppIconURL': '',
                'gPlayURL': '',
                'gPlayAppID': '',
                'officalSiteURL': ''
            };
        });

        // if no titles or URLs wer found.
        if(!altApps) {
            return 'No Alternative App links/Titles could be found. Bailing!';
        }

        _.forEach(altApps, (altApp) => {
            scrapeAltAppPage(altApp);
        });

    });
}

/**
 * Given an alt app and a URL this functino adds the URL
 * to the gPlayURL field in the JSON.
 */
function addGPlayURL(altApp, URL) {
    if (!URL) {
        logger.debug('No Google Play Store Link');
        return altApp;
    }
    altApp.gPlayURL = URL;
    altApp.gPlayAppID = URL.split('id=')[1]; // TODO: cut App ID from string.
    
    return altApp;
}

/**
 * Give an AltApp and a URL this function adds the URL to
 * the official site field in the json.
 */
function addOfficialSiteURL(altApp, URL) {
    if (!URL) {
        logger.debug('No Official site link');
        return altApp;
    }

    altApp.officalSiteURL = URL;
    return altApp;

}

/**
 * Give a URL to a specific suggested Alternative App page's URL
 * this function will scrape the Playstore URL or failing that
 * the Developers Website page. Failing that it will just return
 * Alternative To URL for the page.
 * 
 * Scrapes for:
 *  GPlayAppStore
 *  OfficialAppSite
 *  App's Icon URL
 *  App's GPlayStore ID.
 */
function scrapeAltAppPage(altApp) {
    request(altApp.altToURL, (err, res, html) => {

        // if there wasn't an err with the request.
        if (err) {
            logger.err('Failed to fetch HTML: ' + err);
            return err;
        }

        // Initialising Variables and Loading HTML into Cheerio
        var $ = cheerio.load(html);
        altApp = addGPlayURL(
            altApp,
            $('a[data-link-action="AppStores Link"]:contains("Google")').attr('href')
        );

        altApp = addOfficialSiteURL(
            altApp,
            $('a[data-link-action="Official Website Button"]').attr('href')
        );

        altApp.altAppIconURL = $('#appHeader').find('img').first().attr('data-src-retina');

        logger.debug(altApp);
    });
}

db.getAppsToFindAltsForThatHaventYetHadThemFound()



scrapePageForAlts('http://alternativeto.net/browse/search/?license=free&platform=android&q=facebook');
