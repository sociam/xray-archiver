const logger = require('../../util/logger');

const DB = require('../../db/db');
const db = new DB('downloader');

const config = require('/etc/xray/config.json');

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


async function updateAppVersionAPKDBPath(versionID, apkPath) {
    await db.updateAppVersionAPKLocation(versionID, apkPath);
}

async function updateAppversionHasAPKFlag(versionID, hasAPK) {
    await db.updateAppVersionHasAPKFlag(versionID, hasAPK);
}

async function updateAppVersionAPKServerLocation(versionID, serverName) {
    await db.updateServerLocation(versionID, serverName);
}


async function main() {
    let appPackageVersions = await db.selectAllAppPackageNameVersionNumbers();

    for(const appVer of appPackageVersions) {
        for(const versionID of appVer.versions) {
            const versionDetails = await db.selectAppVersion(versionID);

            const appVersionDir = resolveAppVersionDir(versionDetails);
            const dirExists = fs.existsSync(appVersionDir);
            await updateAppVersionAPKDBPath(versionID, dirExists ? appVersionDir : '');

            const hasAPK = fs.existsSync(`${appVersionDir}/${versionDetails.app}.apk`);
            await updateAppversionHasAPKFlag(versionID, hasAPK);

            await updateAppVersionAPKServerLocation(versionID, dirExists ? config.vmname : '');


            // Now the DB has been updated... do any extra stuff.

            // IF THERE IS AN APK...

                // if rsync flag is set. do something like...
                // rsync -a --relative /var/xray/apps/zen.basketball/play/uk/1.1.3/ hapax:/
                // rsync -a --relative ${appVersionDir} ${newServer}:${rootPath}

                // if update DB with rsync flag is set...
                // updateAppVersionApkDBPath ( versionID, `${rootPath}/${appVersionDir}`);
                // updateAppVersionAPKServerLocation( versionID, `${newServer}`);

        }
    }

}

main();