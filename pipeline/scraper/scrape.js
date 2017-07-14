/* 

Scraoer script searching google playstore for apps and archving said apps.



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


//Logger
const EMERG = 0, ALERT = 1, CRIT = 2, ERR = 3, WARN = 4, NOTICE = 5, INFO = 6, DEBUG = 7;

var prefixes = ['<0>', '<1>', '<2>', '<3>', '<4>', '<5>', '<6>', '<7>'];
function log(level, txt){
  console.log(prefixes[level] + txt);
}
//log(0, "test log");

async function downloadAppApk(element) {
      var p = require("path");
      // console.log(config.appdir,element.appId,appStore,region,
      //   element.version
      // );

      var saveDir = p.join(config.appdir,element.appId,appStore,region,
        element.version
      );

      if (!fs.existsSync(saveDir)) {
        var shell = require("shelljs");
        shell.mkdir("-p", saveDir);
      }

      console.log("App save directory ", saveDir);
      var args = ["-pd", element.appId, "-f", saveDir, "-c", config.credDownload];
      console.log("Python downloader playstore starting");

      const spw = require("child_process").spawn;

      const apk_downloader = await spw("gplaycli", args);

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

          //TODO: clean up directory on 
          return;
        }
        console.log("Download process complete for ", element.appId);
 
      });
}

async function rescrapeAppId(scrapeBase) {
  return await _.map(scrapeBase, async val => {
      var id = val.appId;
      //console.log("Scraping details on: ",id);

      var appData = gplay.app({
        appId: id
      }).then(function(val){
        downloadAppApk(val);
      }).catch(function(e){
        console.log(e);
      });   
    });
}

//TODO: move region to config or section to iterate over
var region = "us";
var appStore = "play";


async function scrapeColl(collectionType) {
  var res = await gplay.list({
    collection: collectionType,
    num: 1,
    region: region
  });
  return await rescrapeAppId(res);
}

async function gatherResults() {
      await Promise.all(_.flatMap(gplay.category, async catg => {  
        _.map(await Promise.all(_.chunk(_.map(gplay.collection, async coll => {

          var res = await gplay.list({
            collection: coll,
            category: catg,
            num: 1,
            region: region
          });

          await rescrapeAppId(res);
        })), 10), async (collChunk) => {
          return await Promise.all(_.map(collChunk, (e) => {
            //TODO: check if already matches before download
            downloadAppApk(e);



          }));
        });
    }), e => { return e.appId; } );
}


(async () => {
  // let scrapeResults = await scrapeColl(gplay.collection.TRENDING);
  // gplay.collection
  var scrapeResults = gatherResults();

  //console.log("Number of apps to download: ",Object.keys(scrapeResults).length);

  //Staggering results to prevent blowing the stack
  // //_.chunk(scrapeResults, 2).forEach(arr => {
  // async.map(scrapeResults, arr => {  
  //   _.map(arr, function(element) {
  //     //console.log(element);
  //     var p = require("path");
  //     // console.log(config.appdir,element.appId,appStore,region,
  //     //   element.version
  //     // );

  //     var saveDir = p.join(config.appdir,element.appId,appStore,region,
  //       element.version
  //     );

  //     if (!fs.existsSync(saveDir)) {
  //       var shell = require("shelljs");
  //       shell.mkdir("-p", saveDir);
  //     }

  //     console.log("App save directory ", saveDir);
  //     var args = [
  //       "-d",
  //       element.appId,
  //       "-f",
  //       saveDir,
  //       "-c",
  //       config.credDownload,
  //       "-v"
  //     ];

  //     console.log("Python downloader playstore starting");

  //     const spw = require("child_process").spawn;

  //     const apk_downloader = spw("gplaycli", args);

  //     apk_downloader.stdout.on("data", data => {
  //       console.log(`stdout: ${data}`);
  //     });

  //     apk_downloader.stderr.on("data", data => {
  //       console.log(`stderr: ${data}`);
  //     });

  //     apk_downloader.on("close", async code => {
  //       if (code != 0) {
  //         console.log("err");
  //         console.log(`child process exited with code ${code}`);
  //         return;
  //       }

        
  //       console.log("Download process complete for ", element.appId);
  //       //TODO: check the actual success of the download + see if the app is done? or just write tests right?
  //       // var db = require('./db');

  //       // var dbId = await db.insertPlayApp(element, region);

  //       // // Send a single message to the server.

  //       var client = unix.createSocket('unix_dgram');
  //       var unix = require('unix-dgram');
        
  //       var message = Buffer(dbId + "-"+ element.appId + "-" + "play" +"-"+region + "-" + element.version);
        
  //       client.on('error', console.error);
  //       client.send(message, 0, message.length, config.sockpath);
  //       client.close();

  //     });
  //   });
  // });
})();
