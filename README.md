# xray / data controller indicators for smartphone apps

## capture logs 

0. check out the repo
1. create a virtualenv, enter the virtualenv; pip install mitmproxy
2. create your own `mitm-config.json` based on mitm-config.sample.json to customise your environment
3. run the script like `mitmdump -s mitm-save.py -p 8081`

## process the logs

1. create a copy of `config.sample.json` called `config.json`
2. you need bower, make sure you have it (or `npm install -g bower`)
3. at the top level, run bower install
4. go to the mitm directory
5. run `npm install`
6. make a 'mitm_out' directory (or whatever destination you specified in `config.json`)
7. run `node parse-output.js` to generate the data files 
8. stand at the base directory, run `http-server`
9. point a browser to to localhost:8080/prototypes/p1.html

