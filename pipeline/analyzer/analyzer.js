//Analyzer: checks the database for unanalyzed apps and invokes the analyzers on them

//Takes: the path to an application APK
//Does: invokes each analyzer in the 'modules' dir on the path to the app
function analyze(app) {
const path = require('path');
const fs = require("fs");
	console.log('Analyzing app: ' + app);

	modules = fs.readdirSync('./modules/');
	console.log('modules:' + modules)

	analyzers = modules.filter(function(file) {
    	return path.extname(file).toLowerCase() === '.js';
	});
	
	analyzers.forEach(function callback(analyzer) {
		console.log('Running analysis module: ' + analyzer);
		const { exec } = require('child_process');
		exec('node modules/' + analyzer + ' ' + app, (error, stdout, stderr) => {
  			if (error) {
    			console.error(`exec error: ${error}`);
    			return;
  			}
  			console.log(`stdout: ${stdout}`);
  			console.log(`stderr: ${stderr}`);
		});
	});
}

function main() {
	//while(true){
	//}
	// get unanylsed apps
	analyze('test')
}

main()
