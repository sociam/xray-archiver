'use strict';
/*
Scraper script searching store for apps and archiving app data said apps.

NOTE: 120 is max number of apps to receive at once through list method...
NOTE: this has a permission section could inspect when throw this into folder
NOTE: cronjob...

TODO: to avoid the situation of askign for a captcha use a throttle keyword
, all methods now support a throttle property, which defines an upper bound to
the amount of requests that will be attempted per second. Once that limit is reached,
further requests will be held  until the second passes.
*/
const config = require('/etc/xray/config.json'); //See example_config.json
const gplay = require('google-play-scraper');
//Reading from folder of csv files
const fs = require('fs');
const fs_promise = require('fs-readdir-promise');
const readline = require('readline');

//Unix pipeline setup
const unix = require('unix-dgram');
const client = unix.createSocket('unix_dgram');

const logger = require('./logger.js');

let appsSaveDir = require('path').join(config.datadir, 'apps');


if (!require('fs').existsSync(appsSaveDir)) {
    logger.info('New apps folder needed ' + appsSaveDir);
    require('shelljs').mkdir('-p', appsSaveDir);
}

//TODO: move region to config or section to iterate over
let region = 'us';
let appStore = 'play';


function resolveAPKDir(appData) {
    let path = require('path');
    //console.log('appdir:'+ config.datadir, '\nappId'+ appData.appId, '\nappStore'+ appStore, '\nregion'+ region, '\nversion'+ appData.version);

    //log('appdir:'+ config.appdir, '\nappId'+ appData.appId, '\nappStore'+ appStore, '\nregion'+ region, '\nversion'+ appData.version);
    //NOTE: If app version is undefined setting to  date
    if (!appData.version || appData.version === 'Varies with device') {
        logger.debug('Version not found defaulting too', appData.updated);
        let formatDate = appData.updated.replace(/\s+/g, '').replace(',', '/');
        appData.version = formatDate;
    }

    let appSavePath = path.join(appsSaveDir, appData.appId, appStore, region, appData.version);
    logger.info('App desired save dir ' + appSavePath);

    /* check that the dir created from config exists. */
    const fsEx = require('fs-extra');

    return Promise.all([fsEx.pathExists(appSavePath), Promise.resolve(appSavePath)]);
}

function spawnGplayDownloader(args) {

    const spw = require('child-process-promise').spawn;
    logger.info('Passing args to downloader' + args);
    const apkDownloader = spw('gplaycli', args);

    let downloadProcess = apkDownloader.childProcess;

    logger.info('APK downloader created childProcess.pid: ' + downloadProcess.pid);

    downloadProcess.stdout.on('data', data => {
        logger.info(`The downloader process produce the following stdout: ${data}`);
    });

    downloadProcess.stderr.on('data', data => {
        logger.warning(`The downloader process produce the following stderr: ${data}`);
    });

    return apkDownloader;
}


//TODO: check dir setup before attempting to search on that word

function extractAppData(appData) {
    //Check appData state
    if (!appData.appId) { return Promise.reject('Invalid appdata', appData.appId); }

    let resolveApk = resolveAPKDir(appData).then(appPathInfo => {
        let exists = appPathInfo[0],
            appSavePath = appPathInfo[1];

        logger.info('Does app save exist already? : ' + exists);
        if (exists) {
            logger.info('App version already exists');
            //logger.err('Could not save apps ', err.message);
            return Promise.reject('App version already exists');
        } else {
            logger.info('New app version ');
            require('shelljs').mkdir('-p', appSavePath);
            appSavePath;
        }
    }).catch(function(err) {
        logger.err('Could not create a app save dir ' + err.message);
        return Promise.reject(appData.appId);
    });

    resolveApk.then(appSaveDir => {
        let args = ['-pd', appData.appId, '-f', appSaveDir, '-c', config.credDownload]; /* Command line args for gplay cli */
        logger.info('Python downloader playstore starting');
        return spawnGplayDownloader(args).catch((err) => {
            logger.warning('Downloading failed with error: ' + err.message);
            return err;
        });
    }).then( () => {
        logger.info('Download process complete for ' + appData.appId);

        //TODO: append previous and write seperately
        appData.push({isDownloaded:true});

        // TODO: DB Comms... this can be factorised.
        let db = require('./db');
        return db.insertPlayApp(appData, region).catch((err) => {
            logger.err('Inserting play app failed'+ err +  region);
            return err;
        });
    }).then((dbId) => {
        // TODO: if unix fails keep trying the socket
        if (!require('fs').existsSync(config.sockpath)) {
            logger.err('Could not bind to socket... try again later  ');
            return Promise.resolve();
        }

        // TODO: Check that '-' won't mess things up on the DB side... eg if region was something like 'en-gb'
        let message = Buffer(dbId + '-' + appData.appId + '-' + config.appStore + '-' + region + '-' + appData.version);

        //client.on('error', logger.err);
        client.send(message, 0, message.length, config.sockpath);
        client.close(); /* The end of one single app download and added to the DB */
    }).catch(function(err) {
        logger.err('Could not write to db ' + err.message);
        return;
    });
}

//Base scrapes array apps based on google-play-scraper app json format - PROMISE FORMAT
// function scrape(appsData) {
//     return appsData.map((val) => {
//         return gplay.app({ appId: val.appId }).then(function(some_other_val) {
//             return extractAppData(val).then(() => {
//                 logger('finished downloading', val.appId);
//                 return val; // whatever you return here will get passed on to the next val in the promise chain..
//             }).catch((e) => {
//                 logger.err('error downloading ', val.appId, e.toString());
//                 throw e;
//             });
//         }).catch(function(err) {
//             logger.err('Could not save app ', err);
//             return Promise.reject(appData.appId);
//         });
//     });
// }


/* Get an array of gplay Search results */
function scrapeWord(word) {
    return gplay.search({
        term: word,
        num: 4,
        region: region,
        price: 'free',
        fullDetail: true,
        throttle: 0.01,
    });
}

let wordStash = config.wordStashDir;


function reader(filepath) {
    return readline.createInterface({
        input: fs.createReadStream(filepath)
    });
}


// //Do processing syncrounously do prevent gplay having a moan
// function processAppData(appsData, processFn) {
//     let index = 0;

//     function next() {
//         if (index < appsData.length) {
//             logger.logger('Processing ', index);
//             processFn(appsData[index++])
//                 .then(next)
//                 .catch((err) => { logger.warning('downloading app failure:', err) });
//         }
//     }
//     next();
// }

function wipe_scraped_word() {
    fs.writeFile(config.datadir + '/scraped_words.txt', '', function(err) {
        if (err) {
            logger.err('Unable wipe the scraped word file');
        }
    });
}

function write_latest_word(word) {
    fs.appendFile(config.datadir + '/scraped_words.txt', word + '\n', function(err) {
        if (err) {
            logger.err('Unable to log to the scraped word file' + err.message);
        }
    });
}

// function topApps() {
//     gplay.category.forEach( cat => {
//         gplay.collection.forEach( coll => {
//             //TODO: this might all happen at once... review owrdStashFiles
//             //finish below off then scrape data + download
//             gplay.list({
//                 collection
//                 category
//                 num: 12,
//                 region: region,
//                 fullDetail: true,
//                 throttle: 0.01
//             }).then( app => {

//                 //TODO: later before begin scraping you do a similar search here

//                 return extractAppData(app);
//             }).catch( err => {
//                 log('Err with word stash',err.message);
//             });

//         });
//     });
// }

wipe_scraped_word();

let wordStashFiles = fs_promise(wordStash);

// TODO: Refactor this....
wordStashFiles.then(files => {
    let q = Promise.resolve();

    files.map(file => {
        q = q.then(() => {
            return new Promise((resolve) => {
                logger.info('Resolving word stash' + wordStash);
                let filepath = require('path').join(wordStash, file);

                let rd = reader(filepath);

                let p = Promise.resolve();

                rd.on('line', (word) => {
                    p = p.then(() => {
                        logger.info('searching on word:' + word);
                        write_latest_word(word);
                        return scrapeWord(word).catch((err) => logger.err('scraping app on word failed:' + err));
                    }).then(function(appsData) {
                        logger.info('Search apps total: ' + appsData.length);
                        let r = Promise.resolve();

                        appsData.forEach(app => {
                            r = r.then(() => {
                                logger.info('Attempting to download:' + app.appId);
                                return extractAppData(app);
                            }, (err) => logger.warning('downloading app failed:' + err));
                        });
                        //processAppData(appsData,extractAppData);

                    }, (err) => logger.err('going through word list failed:' + err));
                });

                rd.on('end', () => {
                    p.then(() => resolve());
                    p.catch((err) => logger.err('last data word failed:' + err));
                });
            });
        }, (err) => logger.err('could no iterate through words in file:' + err));
    }, (err) => logger.err('iterating through dir word list failed::' + err));
}).catch(function(err) {
    logger.err('Err with word stash' + err.message);
});




//let parse = require('csv-parse');
//let async = require('async');
//TODO: Promise that you'll change this to promises.
/* parse 'Top Words' and download apps based on the search results. */
// NOTE: word csv's are not comma seperated, actually '\n'...
// let parser = parse({ delimiter: ',' }, function(err, data) {
//     async.eachSeries(data, function(currWord, callback) {
//         scrapeWord(currWord).then(appsScraped => {

//             appsScraped.map(app => {
//                 log('Downloading app: '+ app.appId);
//                 downloadAppApk(app);
//             });

//             // when processing finishes invoke the callback to move to the next one
//             callback();
//         });
//     })
// });

// // Loop through all the files in the word stash
// fs.readdir(wordStash, function(err, files) {
//     if (err) {
//         logger.err('Could not list the directory.'+ err);
//         process.exit(1);
//     }

//     files.forEach(file => {
//         let p = require('path');
//         fs.createReadStream(p.join(wordStash,file)).pipe(parser);
//     });
// });
