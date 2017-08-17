'use strict';
const fs = require('fs');
const logger = require('../../util/logger');
const _ = require('lodash');

const DB = require('../../db/db');
var db = new DB('suggester');

const childProcess = require('child_process');


function readCSV(path) {
    var lines = fs.readFileSync(path).toString().split('\n');
    var json = [];
    lines.forEach((line) => {
        var split = line.split(',');
        json.push({
            'source': split[0].replace('\r', ''),
            'alt': split[1].replace('\r', '')
        });
    });
    return json;
}

// https://stackoverflow.com/questions/12303989/cartesian-product-of-multiple-arrays-in-javascript
function cartesianProductAltArray() {
    return _.reduce(arguments, function(a, b) {
        return _.flatten(_.map(a, function(x) {
            return _.map(b, function(y) {
                return x.concat(y);
            });
        }), true);
    }, [
        []
    ]);
}

function nwayAlts(alts) {
    var curr = alts[0].source;
    var group = [];
    var nway = [];
    alts.forEach((alt) => {
        if (alt.source != curr) {
            var res = cartesianProductAltArray(group, group);
            nway = nway.concat(res);
            group = [];
            curr = alt.source;
        }
        group.push(alt.alt);
    });
    var res = cartesianProductAltArray(group, group);
    nway = nway.concat(res);
    return nway;
}

function parseAltCSVToJSON(path) {
    var array = readCSV(path);
    return array.concat(nwayAlts(array).map((alt) => {
        if (alt[0] != alt[1]) {
            return {
                'source': alt[0],
                'alt': alt[1]
            };
        }
    }));
}

function parseAltCSVtoArray(path) {
    var lines = fs.readFileSync(path).toString().split('\n');
    var arr = _.flatten(_.map(lines, (line) => {
        return line.split(',');

    }));
    return _.uniqBy(arr, (e) => {
        return e;
    });

}

function scrapeAppID(appID) {
    logger.debug('Attempting to Scrape ' + appID);
    return childProcess.execSync('node ../retriever/idFetch.js ' + appID,
        (error, stdout, stderr) => {
            logger.debug('stdout: ' + stdout);
            logger.debug('stderr: ' + stderr);
            if (error !== null) {
                logger.debug('exec error: ' + error);
            }
        });

}

function main() {
    var alts = parseAltCSVToJSON('alt_apps.csv');
    var apps = parseAltCSVtoArray('alt_apps.csv');

    logger.debug('Apps Parsed. Line Count:' + apps.length);
    logger.debug('App-Alt Pairs Parsed: ' + alts.length);
    apps.forEach((app) => {
        scrapeAppID(app.replace('\r', ''));
    });

    alts.forEach((app) => {
        db.insertManualSuggestion(app)
            .then((res) => logger.debug(res))
            .catch((err) => logger.err(err));

    });
}


main();