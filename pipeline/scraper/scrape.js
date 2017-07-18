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


function resolveAPKDir(appData){

    let path = require("path");
    //console.log("appdir:", config.appdir, "\nappId", appData.appId, "\nappStore", appStore, "\nregion", region, "\nversion", appData.version);
    //NOTE: If app version is undefined setting to  date
    if (!appData.version) {
        appData.version = appData.updated;
    }

    let appSavePath = path.join(config.appdir, appData.appId, appStore, region, appData.version,appData.appId+".apk");
    console.log("App desired save dir ", appSavePath);

    /* check that the dir created from config exists. */
    const fsEx = require('fs-extra');
   
    return fsEx.pathExists(appSavePath).then(exists => {
            console.log("Does app save exist already? : ", exists);
            if(exists) {
                console.log("App version already exists", appSavePath);
                return Promise.reject(appData.appId); 
            } else {
                console.log("New app version", appSavePath);
                require("shelljs").mkdir("-p", appSavePath);
                return Promise.resolve(appSavePath);
            }
        }).catch(function (err) {
            console.error('Could not create a app save dir ', err);
            return Promise.reject(appData.appId); 
        });
}


function spawnGplayDownloader(args) {

    const spw = require('child-process-promise').spawn;
    const apkDownloader = spw("gplaycli", args);

    var downloadProcess = apkDownloader.childProcess;

    console.log('[spawn] APK downloader childProcess.pid: ', downloadProcess.pid);

    downloadProcess.stdout.on("data", data => {
        console.log(`stdout: ${data}`);
    });

    downloadProcess.stderr.on("data", data => {
        console.log(`stderr: ${data}`);
    });

    return apkDownloader;
}

function extractAppData(appData) {
    //Check appData state
    if (!appData.appId) { return Promise.reject("Invalid appdata",appData.appId); }

    var resolveApk = resolveAPKDir(appData);
    //console.log("Resolve apk",resolveApk).then(() => { resolve(); }, (err) => { console.log("last dl failed:", err); });
    
    resolveApk.then(appSaveDir => {
        
        let args = ["-pd", appData.appId, "-f", appSaveDir, "-c", config.credDownload]; /* Command line args for gplay cli */
        
        console.log("Python downloader playstore starting");
    
        let spawnGplay = spawnGplayDownloader(args);
        //console.log("Gplay spwaner",spawnGplay);

        spawnGplay.then(pipeCode => {

            console.log("Download process complete for ", appData.appId);

            // TODO: DB Comms... this can be factorised.
            var db = require('./db');
            var dbId = db.insertPlayApp(appData, region);

            dbId.then(() => {
                var client = unix.createSocket('unix_dgram');
                var unix = require('unix-dgram');

                // TODO: Check that '-' won't mess things up on the DB side... eg if region was something like 'en-gb'
                var message = Buffer(dbId + "-" + appData.appId + "-" + config.appStore + "-" + region + "-" + appData.version);

                client.on('error', console.error);
                client.send(message, 0, message.length, config.sockpath);
                client.close(); /* The end of one single app download and added to the DB */
            }).catch(function(err) {
                console.error('Could not write to db ', err.message);
                return Promise.reject(appData.appId);
            });       
        }).catch(function (err) {
            console.error('[spawn] download ERROR: ', err.message);
            
            return Promise.reject(appData.appId); 
        });   
    }).catch(function (err) {
        console.error('Could not save app ', err.message);
        return Promise.reject(appData.appId); 
    });   
}

//Base scrapes array apps based on google-play-scraper app json format - PROMISE FORMAT
function scrape(appsData) {
    return appsData.map((val) => {
        return gplay.app({ appId: val.appId }).then(function(some_other_val) {
            return extractAppData(val).then(() => {
                console.log('finished downloading', val.appId);
                return val; // whatever you return here will get passed on to the next val in the promise chain..
            }).catch((e) => {
                console.error('error downloading ', val.appId, e.toString());
                throw e;
            });
        }).catch(function (err) {
            console.error('Could not save app ', err);
            return Promise.reject(appData.appId); 
        });   
    });
}


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


//Do processing syncrounously do prevent gplay having a moan
function processAppData(appsData,processFn) {
    var index = 0;

    function next() {
        if(index < appsData.length) {
            console.log("Processing ", index);
             processFn(appsData[index++])
             .then(next)
             .catch((err) => { console.log("downloading app failed:", err)});
        }     
    }
    next();
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
//                 console.log("Err with word stash",err.message);
//             });

//         });
//     });
// }

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

                            console.log("Search apps total: ",appsData.length);
                          
                            var r = Promise.resolve();

                            appsData.forEach(app => {

                                r = r.then( () => {
                                    console.log("Attempting to download:",app.appId);
                                    return extractAppData(app);  
                                }, (err) => { console.log("downloading app failed:", err)});
                            });
                            //processAppData(appsData,extractAppData);

                        }, (err) => { console.log("scraping app on word failed:", err)});
                    }), (err) => { console.log("scraping apps cailes :", err)};
                });

                rd.on('end', () => {
                    p.then(() => { resolve(); }, (err) => { console.log("last dl failed:", err); });
                });
            });
        }, (err) => { console.log("q failed:", err); });
    }, (err) => { console.log("stashfiles failed:", err); });
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

