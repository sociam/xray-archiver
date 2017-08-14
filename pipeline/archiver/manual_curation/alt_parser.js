'use strict';
const fs = require('fs');
const logger = require('../../util/logger');
const _ = require('lodash');

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


function main() {
    var array = readCSV('alt_apps.csv');
    logger.debug(array.length);
    var subset = array.slice(0, 5);
    logger.debug(subset[0].source);
    array = array.concat(nwayAlts(array).map((alt) => {
        if (alt[0] != alt[1]) {
            return {
                'source': alt[0],
                'alt': alt[1]
            };
        }
    }));
    logger.debug(array);
}

main();