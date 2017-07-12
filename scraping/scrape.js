
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

TODO: we can a call a update to a folder through python periodically... does node.js forever support that kind of shenanigans...

*/ 
var gplay = require('google-play-scraper'); //google play store query scrapper. 
var PythonShell = require('python-shell'); //Shell based python spawn - cleaner + get errs
var fs = require('fs');    
var _ = require('lodash');


// var http = require('http'), unixSocket = require("unix-socket");
// /*Socket for pushing out data found */
// var server = http.Server();

// var option = { 
//     path: "/var/apkArchive",
//     mode: 0666 
// };
 
// unixSocket.listen(server, option, function(result) {
//     if (result) {
//         console.log('Server started on ' + result);
//     } else {
//         console.error('Error');
//         process.exit(0);
//     }
// });

//space for socket

//Iterate over all the gplay store files... 
//Can do this through a method of og g-play-scraper via category. 
//Could do it through searching for title a* ... there might be more than MAX query a* though

function requeryOnAppId(scrapeBase) {
  //console.log(scrapeBase);
  _.forEach(scrapeBase, function(index){
      var id = index.appId;
      //console.log(id);
      var scrapeResults = gplay.app({
        appId: id
          }).then(function(data) {
              //console.log(data);
              return data;
          });
  });
}



//Iterate over json top apps  
function scrapeCollectionResult(collectionType) {
    return gplay.list({
      collection: gplay.collection.TRENDING,
      num: 2
    }).then(function(result) {
      //console.log(result);
      return result;
    });
}


var someResult = scrapeCollectionResult(gplay.collection.TRENDING).then(requeryOnAppId);
console.log(someResult.then(console.log));


//Comms with python download
//await scrapeResults;

scrapeResults.then(console.log,console.log);
var PythonShell = require('python-shell');

//$PYTHONPATH check for gplaycli or call the exe

var saveDir = "APK_ARCHIVE"
if(!fs.existsSync(saveDir)){
  fs.mkdirSync(saveDir);
} 

//iterate results over gplay list
scrapeResults.then(function(result) { 
  console.log(result);
  _.forEach(result, function(element) {
    console.log(element.appId);
    var options = {
      scriptPath: '../../scrapeTools/gplaycli/gplaycli', 
      args: ["-d", element.appId,
            "-f", saveDir,
            "-p",]
    };

    //TODO: Split version after 
    //todo: net socket datagram

    //Pretty nice tool can just use this and then organise apps by section...
    console.log("Python downloader playstore starting");
    PythonShell.run('gplaycli.py', options, function (err, results) {
      if (err) throw err;
      // results is an array consisting of messages collected during execution
      console.log('results: %j', results);
    });
  });
});
