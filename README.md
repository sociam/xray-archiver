# X-Ray Archiver

SOCIAM Project for the archiving of Mobile applications. Stores Application metadata

## App Data Retriever
Node JS script for the retrieval of Data relating to mobile applications found on the Google playstore. The script relies on the initialisation of a postgres database and the existance of search terms in the database.

This script utilises the [Google Play Scraper](https://github.com/facundoolano/google-play-scraper) GitHub projet for fetching information from the google play store.

## Search Term Explorer
The search Term explorer generates search terms that can be used by the retriever for fetching app meta data from the google play store. it uses the same [Google Play Scraper](https://github.com/facundoolano/google-play-scraper) for this task.

The search terms generated are auto-completion suggestions made by the google play store. These mean that search terms used by the retreiver will lead to popular apps searched by users.

The explorer depends on the initialisation of a Postgres DB being set up according to the init_db.sql file.

## APK Downloader
The APK downloader fetches Android app APK's using the app data that has been stored in the Postgres DB. The script utilises the [GPlayCli](https://github.com/matlink/gplaycli) for connecting to the google play store for downloading application APKs.

## Database
A Postgres database contains a series of tables required by all elements of the project. Tables for search term data and app meta data are required for each script to function correctly.

An init_db.sql file located in the db folder of this project can be used to initial a postgres database.

# Installation

## install.sh
*info about the install bash script*

## init_db.sql
*info about the init db sql file.*

## Dependencies:
* Node 
    * See - package.json
* [Gplaycli] (https://github.com/matlink/gplaycli)
    * Expecting gplaycli installed to path
* PostgreSQL 
    * db setup to handle the app data scraped. See db/init_db.sql. 
