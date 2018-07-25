const logger = require('../../util/logger');
const DB = require('../../db/db');
const db = new DB('downloader');


// Get list of all app package names.

// check if the app package has:
    // Directory...
        // APK
        // ICON

// Then update DB accordingly

// if move APK Flag is set
    // Move file from one location to the specified location
    // update DB to reflect relocation.

// if move ICON flag is set
    // move file from one location to the specified location


// if delete old apk flag is set
    // delete old apk after move

// if delete old icon flag is set
    // delete old icon after move


async function main() {
    
}

main();