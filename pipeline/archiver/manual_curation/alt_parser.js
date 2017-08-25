const fs = require('fs');
const logger = require('../../util/logger');

const DB = require('../../db/db');
const db = new DB('suggester');

const childProcess = require('child_process');

function readCSV(path) {
    const lines = fs.readFileSync(path).toString().split('\n');
    const pairs = lines.reduce((json, line) => {
        const split = line.replace('\r', '').split(',');
        return json.concat([
            { source: split[0], alt: split[1] },
            { source: split[1], alt: split[0] },
        ]);
    }, []);
    logger.debug(pairs);
    return pairs;
}

function flatten(arr) {
    return arr.reduce((a, b) => a.concat(b), []);
}

function cartesianProductAltArray(...args) {
    return args.reduce((prods, arr) =>
        flatten(prods.map((prod) => arr.map((v) => prod.concat(v)))), [
            []
        ]);
}

function nwayAlts(alts) {
    let curr = alts[0].source;
    let group = [];
    let nway = [];
    alts.forEach((alt) => {
        if (alt.source != curr) {
            const res = cartesianProductAltArray(group, group);
            nway = nway.concat(res);
            group = [];
            curr = alt.source;
        }
        group.push(alt.alt);
    });
    const res = cartesianProductAltArray(group, group);
    nway = nway.concat(res);
    return nway;
}

function parseAltCSVToJSON(path) {
    const array = readCSV(path);
    return array.concat(nwayAlts(array).map((alt) => {
        return { source: alt[0], alt: alt[1] };
    }));
}

function parseAltCSVtoArray(path) {
    const lines = fs.readFileSync(path).toString().split('\n');
    const arr = flatten(lines.map((line) => line.split(',')));
    return arr.reduce((res, e) => res.includes(e) ? res : res.concat([e]), []);
}

function scrapeAppID(appID) {
    logger.debug(`Attempting to Scrape ${appID}`);
    return childProcess.execSync(`node ../retriever/idFetch.js ${appID}`,
        (error, stdout, stderr) => {
            logger.debug(`stdout: ${stdout}`);
            logger.debug(`stderr: ${stderr}`);
            if (error !== null) {
                logger.debug(`exec error: ${error}`);
            }
        });
}

function main() {
    const alts = parseAltCSVToJSON('alt_apps.csv');
    const apps = parseAltCSVtoArray('alt_apps.csv');

    logger.debug(`Apps Parsed. Line Count: ${apps.length}`);
    logger.debug(`App-Alt Pairs Parsed: ${alts.length}`);
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