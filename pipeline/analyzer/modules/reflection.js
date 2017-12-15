function analyse() {
	apk = process.argv[2];

	const { exec } = require('child_process');
	console.log("running grep -Paqh \\x00\\x00\\x00.Ljava/lang/reflect[/a-zA-Z]*;\\x00\\x00\\x00 -- " + apk + "classes.dex")

	exec("grep -Paqh \\x00\\x00\\x00.Ljava/lang/reflect[/a-zA-Z]*;\\x00\\x00\\x00 " + apk + "/classes.dex", (error, stdout, stderr) => {
		if (error) {
			console.error(`exec error: ${error}`);
			return;
		}
		console.log(reflect + ': ' + `stdout: ${stdout}`);
		console.log(reflect + ': ' + `stderr: ${stderr}`);
	});
}

analyse();

//  6     out, err := cmd.Output()
//  7     if err != nil && strings.TrimSpace(string(out)) != "" {
//  8         fmt.Printf("Error checking for reflection: output below\n%s\n\n", string(out))
//  9         return false
// 10     }
// 11
// 12     return (err == nil)

