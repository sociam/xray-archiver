# xray
X-ray project - information controller transparency tools.

## write mitmproxy logs for apps

1. make sure your mitmproxy is in your python path before running
2. create your own `mitm-config.json` to customise your environment
3. run the script like `mitmdump -s mitm-save.py -p 8081`

## process the logs



<!-- The playstorepermissions.py script returns the permissions list of an app. Because the permissions are listed as dynamic content, we have to use selenium to simulate the browser. Unfortunately I can only acheive this by launching a browser window. It is possible to do this silently using a headless browser, but then the simulation fails to open the permissions pop-up so it doesn't return anything.
 -->