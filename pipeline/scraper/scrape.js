
/* 

Background ndoe script scraping apks from playstore.
Categorising into folders on genre. 
Does a apk download step by step and just spawning a child 
python process. That downloader was cool so hey, why not use it!

https://www.npmjs.com/package/google-play-scraper

NOTE: 120 is max number of apps to receive at once through list method...
NOTE: this has a permission section could inspect when throw this into folder
NOTE: cronjob...
TODO: to avoid the situation of askign for a captcha use a throttle keyword
, all methods now support a throttle property, which defines an upper bound to the amount of requests that will be attempted per second. Once that limit is reached, further requests will be held until the second passes.

*/ 
var gplay = require('google-play-scraper'); //google play store query scrapper. 
var fs = require('fs');    
var _ = require('lodash');
var config = require('/etc/xray/config.json')


async function rescrapeAppId(scrapeBase) {

  return data = await Promise.all(_.map(scrapeBase, async (val) => {
      var id = val.appId;
      console.log(id);
      var data = await gplay.app({
        appId: id
      });
      //console.log(data);
      return data;
  }));
}

/*await versions */ 

var region = 'us';

async function scrapeCollection(collectionType) {
  var res = await gplay.list({
      collection: collectionType,
      num: 2,
      region: region
    });
  return await rescrapeAppId(res);
}

(async () => {
  let scrapeResults = await scrapeCollection(gplay.collection.TRENDING);
  //scrapeResults = scrapeCollectionResult(gplay.collection.TRENDING).then(requeryOnAppId);

  //  console.log(scrapeResults);

function sleep(ms) {
      return new Promise(resolve => setTimeout(resolve, ms));
}

_.chunk(scrapeResults, 10).forEach((arr) => {
//iterate results over gplay list
  _.forEach(arr, function(element) {
    
    var p = require('path');
    console.log(element.version);      
    var saveDir = p.join(config.appdir, element.appId,'play',region,element.version);
    console.log(saveDir);
    if(!fs.existsSync(saveDir)) {
      //fs.mkdirSync(saveDir);
      //Needed tp create intermediate dirs if neccassary...did not know if path had this optition
      var shell = require('shelljs');
      shell.mkdir('-p',saveDir);
      console.log(saveDir);
    } 

    
      
      
    var args =  ["-d", element.appId,
                "-f", saveDir,
		            "-c", config.credDownload,
                "-p"]
    //console.log(args);
    
    //Pretty nice tool can just use this and then organise apps by section...
    console.log("Python downloader playstore starting");
    
    const spw = require('child_process').spawn;
     
      
    const apk_downloader = spw('gplaycli',args);

    apk_downloader.stdout.on('data', (data) => {
      console.log(`stdout: ${data}`);
      //Was a sucess pipe the output
    });

    apk_downloader.stderr.on('data', (data) => {
      console.log(`stderr: ${data}`);
    });

    apk_downloader.on('close', (code) => {
      if(code != 0){
        console.log('err')
        console.log(`child process exited with code ${code}`);
        return;
      }

      //The expectance: apk -store -region -version over pipe

      // Send a single message to the server.

      var unix = require('unix-dgram');
      var message = Buffer(element.appId + "-" + "play" +"-"+region + "-" + element.version);
      var client = unix.createSocket('unix_dgram');
      client.on('error', console.error);
      client.send(message, 0, message.length, config.sockpath);
      client.close();
    });
  });
});
})();
