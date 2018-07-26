const logger = require('../../util/logger');

const DB = require('../../db/db');
const db = new DB('downloader');

const config = require('/etc/xray/config.json');

const fs = require('fs');

const argv = require('minimist')(process.argv.slice(2));

const { execSync } = require('child_process');

let argumentsInvalid = false;

const ROOT_FILE_PATH = argv.appRoot ? argv.appRoot : '/var/xray/apps';

// Check Flags are set correctly.
if (argv.deleteAPK && !argv.updateDB) {
    logger.info('Invalid argument configuration: UpdateDB Flag must be set to delete APKs.');
    argumentsInvalid = true;
}

// if any of the rsync related flags are set, but they aren't all set.
if ( (argv.rsync || argv.server || argv.root) && !(argv.rsync && argv.root && argv.server)) {
    logger.info(`Invalid argument configuration: When wanting to rsync files, all three rsync
                flags must be set: rsync, server, and root.`);
    argumentsInvalid = true;
}

if (argumentsInvalid) {
    process.exit();
}

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

function resolveAppVersionDir(appInfo, rootPath) {
    const packageName = appInfo.app;
    const storeName = appInfo.store;
    const region = appInfo.region;
    const version = appInfo.version;

    rootPath = rootPath ? rootPath : ROOT_FILE_PATH;

    return `${rootPath}/${packageName}/${storeName}/${region}/${version}`;
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
    const appPackageVersions = await db.selectAllAppPackageNameVersionNumbers();

    for (const appVer of appPackageVersions) {
        for (const versionID of appVer.versions) {
            const versionDetails = await db.selectAppVersion(versionID);
            // If the APK is still on this VM...
            if (versionDetails.apk_server_location == config.vmname) {
                // Check and update APK status and location.
                const appVersionDir = resolveAppVersionDir(versionDetails);

                const dirExists = fs.existsSync(appVersionDir);
                await updateAppVersionAPKDBPath(versionID, dirExists ? appVersionDir : '');
                const apkPath = `${appVersionDir}/${versionDetails.app}.apk`;
                const hasAPK = fs.existsSync(apkPath);
                await updateAppversionHasAPKFlag(versionID, hasAPK);

                await updateAppVersionAPKServerLocation(versionID, dirExists ? config.vmname : '');

                // If the flags for rsyncing are set....
                if ( hasAPK && argv.rsync && argv.server && argv.root) {
                    const rsyncString =
                        `rsync -a --relative ${appVersionDir} ${argv.server}:${argv.root}`;
                    logger.info(`Performing an RSync for app version: ${versionID}.`);
                    logger.info(`RSync String: ${rsyncString}`);

                    execSync(rsyncString);

                    if (argv.deleteAPK && argv.updateDB) {
                        logger.info(`Deleting old APK now it is moved to ${argv.server}.
                                    Deleting: ${apkPath}`);
                        fs.unlinkSync(apkPath);
                    }

                    if (argv.updateDB) {
                        logger.info(`Updating DB with new APK storage details: ${argv.server}`);

                        const pathString =
                            argv.root == '/' ? appVersionDir : `${argv.root}/${appVersionDir}`;

                        await updateAppVersionAPKDBPath(versionID, pathString);
                        await updateAppVersionAPKServerLocation(versionID, argv.server);
                    }
                }
            }
        }
    }
}

main();
