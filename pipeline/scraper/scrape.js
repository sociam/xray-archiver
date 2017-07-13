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


async function rescrapeAppId(scrapeBase) {
  return await Promise.all(_.map(scrapeBase, async val => {
      var id = val.appId;
      console.log("Scraping details on: ",id);
      return await gplay.app({
        appId: id
      });
    }));
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


(async () => {
  // let scrapeResults = await scrapeColl(gplay.collection.TRENDING);
  // gplay.collection
  let scrapeResults = 
  _.uniqBy(
    await Promise.all(_.flatMap(gplay.category, async catg => {  
      return await Promise.all(_.map(gplay.collection, async coll => {
        
        // console.log(catg,coll);

        var res = await gplay.list({
          collection: coll,
          category: catg,
          num: 1,
          region: region
        });
        
        return await rescrapeAppId(res);
      }));
  }), e => { return e.appId; } ));
  

  console.log("Number of apps to download: ",Object.keys(scrapeResults).length);

  //Staggering results to prevent blowing the stack
  _.chunk(scrapeResults, 2).forEach(arr => {
    _.map(arr, function(element) {
      //console.log(element);
      var p = require("path");
      console.log(config.appdir,element.appId,appStore,region,
        element.version
      );

      var saveDir = p.join(config.appdir,element.appId,appStore,region,
        element.version
      );

      if (!fs.existsSync(saveDir)) {
        var shell = require("shelljs");
        shell.mkdir("-p", saveDir);
      }

      console.log("App save directory ", saveDir);
      var args = [
        "-d",
        element.appId,
        "-f",
        saveDir,
        "-c",
        config.credDownload,
        "-v"
      ];

      console.log("Python downloader playstore starting");

      const spw = require("child_process").spawn;

      const apk_downloader = spw("gplaycli", args);

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
          return;
        }

        
        console.log("Download process complete for ", element.appId);
        //TODO: check the actual success of the download + see if the app is done? or just write tests right?
        // var db = require('./db');

        // var dbId = await db.insertPlayApp(element, region);

        // // Send a single message to the server.

        var client = unix.createSocket('unix_dgram');
        var unix = require('unix-dgram');
        
        var message = Buffer(dbId + "-"+ element.appId + "-" + "play" +"-"+region + "-" + element.version);
        
        client.on('error', console.error);
        client.send(message, 0, message.length, config.sockpath);
        client.close();

      });
    });
  });
})();
