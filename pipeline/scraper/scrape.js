/* 

Scraoper script searching google playstore for apps and archving said apps.

NOTE: 120 is max number of apps to receive at once through list method...
NOTE: this has a permission section could inspect when throw this into folder
NOTE: cronjob...

TODO: to avoid the situation of askign for a captcha use a throttle keyword
, all methods now support a throttle property, which defines an upper bound to 
the amount of requests that will be attempted per second. Once that limit is reached, 
further requests will be held until the second passes.

*/

var gplay = require("google-play-scraper");
var config = require("/etc/xray/config.json"); //See example_config.json
var fs = require("fs");
var _ = require("lodash");

//Logging mechanisim for script
const EMERG = 0, ALERT = 1, CRIT = 2, ERR = 3, WARN = 4, NOTICE = 5, INFO = 6, DEBUG = 7;

var prefixes = ['<0>', '<1>', '<2>', '<3>', '<4>', '<5>', '<6>', '<7>'];
function log(level, txt){
  console.log(prefixes[level] + txt);
}
//log(0, "test log");


// process.on('unhandledRejection', (reason, p) => {
//   console.log('Unhandled Rejection at: Promise', p, 'reason:', reason);
//   // application specific logging, throwing an error, or other logic here
// });

async function downloadAppApk(appData) {
      var p = require("path");

      var saveDir = p.join(config.appdir,appData.appId,appStore,region,
        appData.version
      )

      if (!fs.existsSync(saveDir)) {
        var shell = require("shelljs");
        shell.mkdir("-p", saveDir);
      }

      console.log("App save directory ", saveDir);
      var args = ["-pd", appData.appId, "-f", saveDir, "-c", config.credDownload];
      console.log("Python downloader playstore starting");

      const spw = require("child_process").spawn;

      const apk_downloader = await spw("gplaycli", args);
      apk_downloader.catch(err => {
        console.log(err);
      })


      apk_downloader.stdout.on("data", data => {
        console.log(`stdout: ${data}`);
      });

      apk_downloader.stderr.on("data", data => {
        console.log(`stderr: ${data}`);
      });

      apk_downloader.on("close", async code => {
        if (code != 0) {
          console.log("err");
          console.log(`child process exited with code ${code}`);
          if (!fs.existsSync(saveDir + appData.appId + ".apk")) {
            fs.rmdirSync(saveDir)
          } 
          return;
        }
        
        console.log("Download process complete for ", appData.appId);
                    console.log("Download process complete for ", appData.appId);

        var db = require('./db');
        var dbId = await db.insertPlayApp(appData, region);
        var client = unix.createSocket('unix_dgram');
        var unix = require('unix-dgram');
        
        var message = Buffer(dbId + "-"+ appData.appId + "-" + "play" +"-"+region + "-" + appData.version);
      
        client.on('error', console.error);
        client.send(message, 0, message.length, config.sockpath);
        client.close();

      });

}

async function rescrapeAppId(scrapeBase) {
  return await _.map(scrapeBase, async val => {
      var id = val.appId;
      //console.log("Scraping details on: ",id);

      var appData = gplay.app({
        appId: id
      }).then(function(val){
        
        downloadAppApk(val).then(console.log,console.log).catch(console.log);

      }).catch(function(e){
        console.log("Catch promise err",e);
      
      });   
    });
}

function scrape(scrapeBase) {
    return scrapeBase.map((val) => {
        return gplay.app({appId: val.appId}).then(function(some_other_val){ 
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


    

//TODO: move region to config or section to iterate over
var region = "us";
var appStore = "play";

// var gplay = require('google-play-scraper');
// gplay.app({appId: 'com.dxco.pandavszombies'})
//   .then(console.log, console.log);

// let searchResult = gplay.list({
//         collection: gplay.collection.NEW_FREE,
//         category: gplay.category.BUSINESS,
//         num: 20,
//         region: region,
//         fullDetail: true,
//         throttle: 10
//       })
      
// searchResult.then(res => {
//    console.log("Results found: ",res.length);
// })
      
  
//scrape(searchResult);


//console.log(res);





async function gatherResults() {
      await Promise.all(_.flatMap(gplay.category, async catg => {  
        _.map(await Promise.all(_.chunk(_.map(gplay.collection, async coll => {

          var res = await gplay.list({
            collection: coll,
            category: catg,
            num: 120,
            region: region,
            fullDetail: true,
            throttle: 10
          });
          downloadAppApk(res).then(console.log,console.log).catch(console.log);

        })), 10), async (collChunk) => {
          return await Promise.all(_.map(collChunk, (e) => {
            //TODO: check if already matches before download
            downloadAppApk(e);
          }));
        });
    }), e => { return e.appId; } ).catch( err => {
      console.log(err);
    })
}

(async () => {
  var scrapeResults = gatherResults();
})();


//Promise version

// function gatherScapingResults() {
//   return _.map(gplay.category, (catg) => {
//     console.log("Began gathering on category",catg);
//     return _.chunk( _.map(gplay.collection, coll => {
//       console.log("Began gathering on coll",coll);
//       var res = gplay.list({
//             collection: gplay.collection.NEW_FREE,
//             category: gplay.category.BUSINESS,
//             num: 10,
//             region: region,
//             fullDetail: true,
//             throttle: 10
//       }).then( (chunkData) => {
//           if(chunkData !== undefined){ //Does not function on paid apps
//             console.log("Initalising downloading: ", chunkData.appId);
//             // downloadAppApk(chunkData).then(() => { 
//             // console.log('finished downloading:', chunkData.appId); 
//             // return val; // whatever you return here will get passed on to the next val in the promise chain..
//             // }).catch((e) => {
//             //     console.error('error downloading ', chunkData.appId, e.toString());
//             //     throw e;
//             // });
//           }
//       }).catch((e) => {
//         console.error('error downloading ', chunkData.appId, e.toString());
//         throw e;
//       });;
//     },10));
//   });
// }

// let results = gatherScapingResults();
