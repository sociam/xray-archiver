# App Scraper

Setup to scrape platstore for free apks to later analysis the permission setup.


# Dependencies:
* Node 
    * See - package.json
* [Gplaycli] (https://github.com/matlink/gplaycli)
    * Expecting gplaycli installed to path
* PostgreSQL 
    * db setup to handle the app data scraped. See db/init_db.sql. 

# Configuration:

All expectance comes from a config.json file. See example_config.json.

## Words Search

In order to a scrape a wide distribution of apps from the playstore the search bar is taken advantage off. 
A one word item per line list can then be passed to the scrape.js to scrape on.

# App Updateing



*NOTE: time to download a significant number of apps can be lenghtly. App download speed is limited to prevent bot blocking*