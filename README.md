# xray
X-ray project - information controller transparency tools.

## Scraping playstore permissions

The playstorepermissions.py script returns the permissions list of an app. Because the permissions are listed as dynamic content, we have to use selenium to simulate the browser. Unfortunately I can only acheive this by launching a browser window. It is possible to do this silently using a headless browser, but then the simulation fails to open the permissions pop-up so it doesn't return anything.
