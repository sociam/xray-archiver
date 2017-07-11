
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



//Iterate over all the gplay store files... 
//Can do this through a method of og g-play-scraper via category. 
//Could do it through searching for title a* ... there might be more than MAX query a* though


//Comms with python download

var PythonShell = require('python-shell');

//$PYTHONPATH check for gplaycli or call the exe


var storeLoci = "APK_ARCHIVE"
if(!fs.existsSync(dir)){
  fs.mkdirSync(dir)
}

var options = {
  scriptPath: '../../scrapeTools/gplaycli/gplaycli', 
  args: ["-d", "com.dxco.pandavszombies",
         "-f", storeLoci,
         "-p",]
};

//Pretty nice tool can just use this and then organise apps by section...
PythonShell.run('gplaycli.py', options, function (err, results) {
  if (err) throw err;
  // results is an array consisting of messages collected during execution
  console.log('results: %j', results);
});

