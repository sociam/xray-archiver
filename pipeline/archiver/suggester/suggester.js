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
function scrapePageForAlts(URLString, appID) {

    // Request the page from the Site.
    return request(URLString, (err, res, html) => {

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
                'officialSiteURL': '',
                'appID': appID
            };
        });

        // if no titles or URLs wer found.
        if(!altApps) {
            return 'No Alternative App links/Titles could be found. Bailing!';
        }

        _.forEach(altApps, (altApp) => {
            scrapeAltAppPage(altApp);
        });
        return 'Apps Scraped';
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

    altApp.officialSiteURL = URL;
    return altApp;

}

/**
 * Give a URL to a specific suggested Alternative App page's URL
 * this function will scrape the Playstore URL or failing that
 * the Developers Website page. Failing th
at it will just return
 * Alternative To URL for the page.
 * 
 * Scrapes for:
 *  GPlayAppStore
 *  OfficialAppSite
 *  App's Icon URL
 *  App's GPlayStore ID.
 */
function scrapeAltAppPage(altApp) {
    return request(altApp.altToURL, (err, res, html) => {

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

        return db.insertAltApp(altApp);
    });
}


function findAppAltPage(URLString, appID) {
    // Request the page from the Site.
    return request(URLString, (err, res, html) => {

        // if there wasn't an err with the request.
        if (err) {
            logger.err('Failed to fetch HTML: ' + err);
            return err;
        }

        // Initialising Variables and Loading HTML into Cheerio
        var $ = cheerio.load(html);
        // check if direct link to page was found or it gave searhc results.

        var titleString = $('title').text().split(' - ')[1];
        logger.debug(titleString);
        if (titleString == 'Search on AlternativeTo.net') {
            var firstRes = $('.app-list').first().find('h3').first().find('a').first();
            var resHref = firstRes.attr('href');
            URLString = 'http://alternativeto.net' + resHref; 
        }

        return scrapePageForAlts(URLString, appID);
    });
}

(async() => {
    var rows = await db.getAppsToFindAltsForThatHaventYetHadThemFound(10);

    await _.forEach(rows, async (row) => {
        var encodedURI = encodeURIComponent(row.title);
        encodedURI = encodedURI.replace(new RegExp('%20', 'g'), '+');
        findAppAltPage(
            'http://alternativeto.net/browse/search/?license=free&platform=android&q=' + encodedURI,
            row.app
        );
    });

})();
