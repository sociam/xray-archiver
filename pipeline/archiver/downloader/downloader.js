/*
Download process spawner
*/

const config = require('/etc/xray/config.json');
const fs = require('fs-extra');
const path = require('path');
const logger = require('../../util/logger');
const util = require('util');
const bashExec = util.promisify(require('child_process').exec);
const db = new (require('../../db/db'))('downloader');


function mkdirp(dir) {
    dir.split(path.sep).reduce((parentDir, childDir) => {
        const curDir = path.join(parentDir, childDir);
        if (!fs.existsSync(curDir)) {
            fs.mkdirSync(curDir);
        }
        return curDir;
    }, path.isAbsolute(dir) ? path.sep : '');
}

function parseDFOutputForFirstFS(fsString) {
    const fsLines = fsString.split('\n').filter((part) => part != '');
    return fsLines[1].split(' ')[0];
}

function parseDFOutputToJSON(fsString) {
    const fsLines = fsString.split('\n').filter((part) => part != '');
    const header = fsLines[0].split(' ').filter((part) => part != '');
    const fileSystems = fsLines.slice(1);
    let parsedOutput = {};
    for(const fs of fileSystems) {
        const lineParts = fs.split(' ').filter((part) => part != '');
        parsedOutput[lineParts[0]] = {};
        for(let i=1 ; i < header.length; i++) {
            if(lineParts[i]) {
                parsedOutput[lineParts[0]][header[i].toLowerCase()] = lineParts[i];
            }
        }
    }
    return parsedOutput;
}

async function getAvailableDiskSpace(path) {
    logger.debug(`Using 'df' to get the disk space for the drive containing '${path}'`);
    try {
        const {stdout, stderr} = await bashExec(`df ${path} -BG`);

        if(stderr) {
            logger.err(`getDiskSpace: df wrote to stderr. throwing err.`);
            throw stderr;
        }

        const dfJSON = parseDFOutputToJSON(stdout);
        const fs = parseDFOutputForFirstFS(stdout);
        return parseInt(dfJSON[fs]['available'].replace('G',''));
    }
    catch(err){
        logger.err(`Error getting the disk space for disk containing ${path}. Error: ${err}`);
    }
}


async function resolveAPKDir(appData) {
    const appsSaveDir = getLocationWithLeastSpace();
    logger.debug(`appdir: ${appsSaveDir}`, `\nappId ${appData.app}`, `\nappStore ${appData.store}`,
        `\nregion ${appData.region}`, `\nversion ${appData.version}`);

    const appSavePath = path.join(appsSaveDir, appData.app,
        appData.store, appData.region, appData.version);
    logger.info(`App desired save dir ${appSavePath}`);

    await mkdirp(appSavePath);

    return appSavePath;
}

function downloadApp(appData, appSavePath) {
    // Command line args for gplay cli
    const args = ['-pd', appData.app, '-f', appSavePath, '-c', config.system_config.downloader_credentials];
    const spw = require('child-process-promise').spawn;
    logger.debug(`Passing args to downloader${args}`);
    const apkDownloader = spw('gplaycli', args);

    const downloadProcess = apkDownloader.childProcess;

    logger.info('DL process %d for %s-%s started.', downloadProcess.pid,
        appData.app, appData.version);

    downloadProcess.stdout.on('data', (data) => {
        const prefix = `'DL process ${downloadProcess.pid} stdout: `;
        logger.debug(prefix + data.toString().replace(/\n/g, `\n${prefix}`));
    });

    downloadProcess.stderr.on('data', (data) => {
        const prefix = `'DL process ${downloadProcess.pid} stdout: `;
        logger.warning(prefix + data.toString().replace(/\n/g, `\n${prefix}`));
    });

    return apkDownloader;
}

async function download(app) {
    logger.info('Starting download attempt for:', app.app);
    // Could be move to the call to DL app. but this is where the whole DL process starts.
    db.updatedDlAttempt(app);
    let appSavePath;
    try {
        appSavePath = await resolveAPKDir(app);
    } catch (err) {
        await new Promise((resolve) => setTimeout(resolve, 6000));
        return Promise.reject(`Did not have access to resolve dir: ${err.message}`);
    }

    try {
        await downloadApp(app, appSavePath);
    } catch (err) {
        logger.debug('Attempting to remove created dir');
        await fs.rmdir(appSavePath).catch(logger.warning);
        return Promise.reject(`Downloading failed with err: ${err}`);
    }

    try {
        const apkPath = path.join(appSavePath, `${app.app}.apk`);

        if (fs.existsSync(apkPath)) {
            // Perform a check on apk size
            await fs.stat(apkPath, async(err, stats) => {
                if (stats.size == 0 || stats.size == undefined) {
                    await fs.rmdir(appSavePath).catch(logger.warning);
                    return Promise.reject('File did not successfully download and is a empty size');
                }

                await db.updateDownloadedApp(app, appSavePath, config.system_config.vm_name);
                return undefined;
            });
        }
    } catch (err) {
        // TODO: Maybe do something else? Destroying process as we have apks that
        // don't exist in db...
        return Promise.reject('Err when updated the downloaded app', err);
    }
    return undefined;
}

async function ensureDirectoriesExist(directories) {
    for(const location of directories) {
        try{
            await fs.ensureDir(path.join(location, 'apps'));
        }
        catch(err) {
            logger.err(`Error ensuring directory '${location}' exisits. Error: ${err}`);
        }
    }
}

async function main() {

    // Ensure that directory structures exist.
    await ensureDirectoriesExist(config.storage_config.apk_download_directories);

    for (;;) {
        let apps;
        try {
            apps = await db.queryAppsToDownload(4000);
        } catch (err) {
            await new Promise((resolve) => setTimeout(resolve, 1000));
            continue;
        }

        await Promise.each(apps, download).catch(logger.err);
    }
}

main();
