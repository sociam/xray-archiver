const config = require('/etc/xray/config.json');
//Analyzer: checks the database for unanalyzed apps and invokes the analyzers on them

//Takes: information uniqely identifying an apk
//Returns: path to that apk
function resolveApp(store, region, id, version) {
	return config.datadir + '/apps/' + id + '/' + store + '/' + region + '/' +  version;
}

//Takes: the path to an application APK
//Does: invokes each analyzer in the 'modules' dir on the path to the app
function analyze(store, region, id, version) {
	const { exec } = require('child_process');
	const path = require('path');
	const fs = require("fs");

	app = resolveApp(store, region, id, version);
	console.log('unpacking app: ' + app);
	exec('apktool d -s ' + app + '/' + id + '.apk' + " -o " + app + ' -f', (error, stdout, stderr) => {
  		if (error) {
    		console.error(`exec error: ${error}`);
    		return;
  		}
  		console.log(analyzer + ': ' + `stdout: ${stdout}`);
  		console.log(analyzer + ': ' + `stderr: ${stderr}`);
	});

	
	console.log('Analyzing app: ' + app);
	analyzers = config.analyzers;

	analyzers.forEach(function callback(analyzer) {
		console.log('Running analysis module: ' + analyzer);
		exec('node modules/' + analyzer + ' ' + app, (error, stdout, stderr) => {
  			if (error) {
    			console.error(`exec error: ${error}`);
    			return;
  			}
  			console.log(analyzer + ': ' + `stdout: ${stdout}`);
  			console.log(analyzer + ': ' + `stderr: ${stderr}`);
		});
	});
}

function main() {
	//while(true){
	//}
	// get unanylsed apps
	analyze('play', 'uk', 'com.amazon.kindle', '2017-12-06')
}

main()
