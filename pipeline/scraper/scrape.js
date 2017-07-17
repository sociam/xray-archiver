/* 

Scraper script searching store for apps and archiving app data said apps.

NOTE: 120 is max number of apps to receive at once through list method...
NOTE: this has a permission section could inspect when throw this into folder
NOTE: cronjob...

TODO: to avoid the situation of askign for a captcha use a throttle keyword
, all methods now support a throttle property, which defines an upper bound to 
the amount of requests that will be attempted per second. Once that limit is reached, 
further requests will be held until the second passes.

*/
var config = require("/etc/xray/config.json"); //See example_config.json
var gplay = require("google-play-scraper");
var fs = require("fs");
var _ = require("lodash");

//Logging mechanisim for script
const EMERG = 0,
    ALERT = 1,
    CRIT = 2,
    ERR = 3,
    WARN = 4,
    NOTICE = 5,
    INFO = 6,
    DEBUG = 7;

var prefixes = ['<0>', '<1>', '<2>', '<3>', '<4>', '<5>', '<6>', '<7>'];

function log(level, txt) {
    console.log(prefixes[level] + txt);
}
//TODO: Use Logger eg - log(0, "test log");


//TODO: move region to config or section to iterate over
var region = "us";
var appStore = "play";

async function downloadAppApk(appData) {
    // TODO: Refactor 
    // TODO: Remove Async, use promises.
    // Set up directory for archiving of app apks
    var p = require("path");
    // archive location *heavily* dependent on config file settings.
    console.log("ad:", appData);
    if (!appData.appId) { return Promise.reject(appData.appId); }

    console.log("appdir:", config.appdir, "\nappId", appData.appId, "\nappStore", appStore, "\nregion", region, "\nversion", appData.version);
    
    //NOTE: If app version is undefined set at date
    if (!appData.version) {
        appData.version = appData.updated;
    }

    let appSaveDir = p.join(config.appdir, appData.appId, appStore, region, appData.version);

    /* check that the dir created from config exists. */
    if (!fs.existsSync(appSaveDir)) {
        var shell = require("shelljs");
        shell.mkdir("-p", appSaveDir);
    }

    console.log("App save directory ", appSaveDir);
    var args = ["-pd", appData.appId, "-f", appSaveDir, "-c", config.credDownload]; /* Command line args for gplay cli */
    console.log("Python downloader playstore starting");

    // TODO: (from dean) this could be used in a promise??
    const spw = require("child_process").spawn;

    const apk_downloader = await spw("gplaycli", args);
    //console.log("Apk downloader", apk_downloader);

    apk_downloader.stdout.on("data", data => {
        console.log(`stdout: ${data}`);
    });

    apk_downloader.stderr.on("data", data => {
        console.log(`stderr: ${data}`);
    });

    /* Waiting for the process to finish before handling anything */
    apk_downloader.on("close", async code => {
        /* Check for errors in downloading first. */
        if (code != 0 || !fs.existsSync(appSaveDir)){
            console.log("err could not download");
            console.log(`child process exited with code ${code}`);
            console.log('process could not save,')
            return;
        }

        

        console.log("Download process complete for ", appData.appId);

        // TODO: DB Comms... this can be factorised.
        var db = require('./db');
        var dbId = await db.insertPlayApp(appData, region);

        var client = unix.createSocket('unix_dgram');
        var unix = require('unix-dgram');

        // TODO: Check that '-' won't mess things up on the DB side... eg if region was something like 'en-gb'
        var message = Buffer(dbId + "-" + appData.appId + "-" + config.appStore + "-" + region + "-" + appData.version);

        client.on('error', console.error);
        client.send(message, 0, message.length, config.sockpath);
        client.close(); /* The end of one single app download and added to the DB */
    });
}

//Base scrapes array apps based on google-play-scraper app json format - PROMISE FORMAT
function scrape(appsData) {
    return appsData.map((val) => {
        return gplay.app({ appId: val.appId }).then(function(some_other_val) {
            return downloadAppApk(val).then(() => {
                console.log('finished downloading', val.appId);
                return val; // whatever you return here will get passed on to the next val in the promise chain..
            }).catch((e) => {
                console.error('error downloading ', val.appId, e.toString());
                throw e;
            });
        });
    });
}
function scrapeWords(wordList) {
    /* Map an array of gplay search results */
    return _.map(wordList, word => {
        console.log("Word defintion", word);

        let scraped = gplay.search({
            term: word,
            num: 120,
            region: region,
            fullDetail: true,
            throttle: 0.1
        });

        console.log(scraped);
        /*  */
        scraped.then(appsScraped => {
            appsScraped.map(app => {
                console.log("search chunk", app.appId);
                downloadAppApk(app);
            });
        });

    });

}

/* Get an array of gplay Search results */
function scrapeWord(word) {
    return gplay.search({
        term: word,
        num: 120,
        region: region,
        fullDetail: true,
        throttle: 1
    });
}

var wordStash = config.wordStashDir;

//Reading from folder of csv files
var fs = require('fs');
var fs_promise = require('fs-readdir-promise');
var readline = require('readline');

function reader (filepath) {
    return readline.createInterface({
        input: fs.createReadStream(filepath)
    });
}


var wordStashFiles = fs_promise(wordStash);

wordStashFiles.then(files => {
    var q = Promise.resolve();

    files.map(file => {
        q = q.then(() => {
            return new Promise((resolve, reject) => {
                var filepath = require("path").join(wordStash,file);
                
                var rd = reader(filepath);
                
                var p = Promise.resolve();

                rd.on('line', (word) => {
                    p = p.then(() => {
                        console.log("searching on word:", word);

                        return scrapeWord(word).then(function(appsData){

                            console.log("Appdata",appsData.appId);
                          
                            var r = Promise.resolve();

                            appsData.forEach(app => {

                                r = r.then( () => {
                                    console.log(app);
                                    return downloadAppApk(app);  
                                });
                                
                            });

                        }, (err) => { console.log("scrapeword failed:", err)});
                    });
                });

                rd.on('end', () => {
                    p.then(() => { resolve(); }, (err) => { console.log("last dl failed:", err); });
                });
            });
        }, (err) => { console.log("q failed:", err); });
    }, (err) => { console.log("stashfiles failed:", err); });
    
    // _.map(files, file => {




    //     rd.on('line', function(word) {
    //         p = p.then(() => {            
    //             console.log("searching on word:",word);
    //             var scrapedAppData = scrapeWord(word);

    //             return scrapedAppData.then(appData => {
    //                 console.log(appData.appId);
    //                 downloadAppApk(appData);
    //             });
    //         });
    //     });

    // });

}).catch(function(err) {
  console.log("Err with word stash",err.message);
});




//var parse = require('csv-parse');
//var async = require('async');
//TODO: Promise that you'll change this to promises.
/* parse 'Top Words' and download apps based on the search results. */
// NOTE: word csv's are not comma seperated, actually '\n'...
// var parser = parse({ delimiter: ',' }, function(err, data) {
//     async.eachSeries(data, function(currWord, callback) {
//         scrapeWord(currWord).then(appsScraped => {

//             appsScraped.map(app => {
//                 console.log("Downloading app: ", app.appId);
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
//         console.error("Could not list the directory.", err);
//         process.exit(1);
//     }
    
//     files.forEach(file => {
//         var p = require("path");
//         fs.createReadStream(p.join(wordStash,file)).pipe(parser);
//     });
// });

