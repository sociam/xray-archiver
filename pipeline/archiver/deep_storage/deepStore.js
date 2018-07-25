const logger = require('../../util/logger');
const DB = require('../../db/db');
const db = new DB('downloader');

const fs = require('fs');


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


function resolveAppVersionDir(appInfo) {

    const packageName = appInfo.app;
    const storeName = appInfo.store;
    const region = appInfo.region;
    const version = appInfo.version;

    return `/var/xray/apps/${packageName}/${storeName}/${region}/${version}`;
}

async function main() {
    let appPackageVersions = await db.selectAllAppPackageNameVersionNumbers();

    for(const appVer of appPackageVersions) {
        for(const versionID of appVer.versions) {
            const versionDetails = await db.selectAppVersion(versionID);

            const appVersionDir = resolveAppVersionDir(versionDetails);

            const dirExists = fs.existsSync(appVersionDir);
            console.log(`App Version has Directory Structure: ${dirExists ? 'TRUE' : 'FALSE'} - ${appVersionDir}`);
        }
    }

}

main();