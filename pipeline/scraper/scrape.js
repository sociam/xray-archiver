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
  return (data = await Promise.all(
    _.map(scrapeBase, async val => {
      var id = val.appId;
      console.log("Scraping details on: ",id);
      var data = await gplay.app({
        appId: id
      });
      return data;
    })
  ));
}

//TODO: move region to config or section to iterate over
var region = "us";

async function scrapeColl(collectionType) {
  var res = await gplay.list({
    collection: collectionType,
    num: 20,
    region: region
  });
  return await rescrapeAppId(res);
}

(async () => {
  let scrapeResults = await scrapeColl(gplay.collection.TRENDING);

  //Staggering results to prevent blowing the stack
  _.chunk(scrapeResults, 2).forEach(arr => {
    _.map(arr, function(element) {
      var p = require("path");
      console.log(
        "Resolved Save Path",
        config.appdir,
        element.appId,
        "play",
        region,
        element.version
      );
      var saveDir = p.join(
        config.appdir,
        element.appId,
        "play",
        region,
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

        //The expectance: apk -store -region -version over pipe
        console.log("Download complete for ", element.appId);

        // var db = require('./db');

        // var dbId = await db.insertPlayApp(element, region);

        // // Send a single message to the server.

        // var client = unix.createSocket('unix_dgram');
        // var unix = require('unix-dgram');
        // var message = Buffer(dbId + "-"+ element.appId + "-" + "play" +"-"+region + "-" + element.version);
        // client.on('error', console.error);
        // client.send(message, 0, message.length, config.sockpath);
        // client.close();
      });
    });
  });
})();
