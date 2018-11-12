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

let downloadLocations = [];

async function ensureDirectoriesExist(directories) {
    const validDirectories = [];
    for (const dir of directories) {
        try {
            await fs.ensureDir(path.join(dir.path, 'apps'));
            validDirectories.push(dir);
        } catch (err) {
            logger.err(
                `
                Error ensuring directory '${dir.path}' exists.
                Error: ${err}.
                Exluding ${dir.name} at ${dir.path} from the download location options.
                `
            );
        }
    }
    return validDirectories;
}

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
    const parsedOutput = {};
    for (const fs of fileSystems) {
        const lineParts = fs.split(' ').filter((part) => part != '');
        parsedOutput[lineParts[0]] = {};
        for (let i=1; i < header.length; i++) {
            if (lineParts[i]) {
                parsedOutput[lineParts[0]][header[i].toLowerCase()] = lineParts[i];
            }
        }
    }
    return parsedOutput;
}

async function getUUID(devicePath) {
    // Can use df to get the filesystem and mounted on path.
    //      df -BG /dev/disk/by-uuid/5D3E-D824
    //      Filesystem     1G-blocks  Used Available Use% Mounted on
    //      /dev/sdc1            58G   14G       44G  24% /mnt/sanDiskUSB
    //
    try {
        const { stdout, stderr } = await bashExec(`sudo blkid -s UUID -o value ${devicePath}`);
        if (stderr) {
            logger.err(`blkid wrote to STDERR: ${stderr}`);
            throw stderr;
        }
        return stdout;
    } catch (err) {
        logger.err(`Error getting UUID for ${devicePath} using 'blkid'. Error: ${err}`);
        return err;
    }
}

async function df(path='') {
    const { stdout, stderr } = await bashExec(`df ${path} -BG`);
    if (stderr) {
        logger.err('getDiskSpace: df wrote to stderr. throwing err.');
        throw stderr;
    }
    return { stdout, stderr };
}

async function getPathFileSystem(path) {
    try {
        const { stdout, stderr } = await df(path);
        if (stderr) {
            throw stderr;
        }
        return parseDFOutputForFirstFS(stdout);
    } catch (err) {
        logger.err(`Error getting the Filesystem containing ${path}. Error: ${err}`);
        return err;
    }
}
async function getAvailableDiskSpace(path) {
    logger.debug(`Using 'df' to get the filesystem space for the filesystem containing '${path}'`);
    try {
        const { stdout, stderr } = await df(path);
        const dfJSON = parseDFOutputToJSON(stdout);
        const fs = parseDFOutputForFirstFS(stdout);
        if (stderr) {
            throw stderr;
        }
        return parseInt(dfJSON[fs]['available'].replace('G', ''));
    } catch (err) {
        logger.err(`Error getting the filesystem space for filesystem containing ${path}.
        Error: ${err}`);
        return err;
    }
}

async function getLocationWithLeastSpace() {
    logger.debug('Getting Save Location with the lowest amount of Space.');
    downloadLocations = await ensureDirectoriesExist(
        config.storage_config.apk_download_directories
    );
    const dirSpaces = [];
    for (const dir of downloadLocations) {
        dirSpaces.push({
            name: dir.name,
            path: dir.path,
            available: await getAvailableDiskSpace(dir.path),
        });
    }

    const dirsWithSomeSpace = dirSpaces.filter((dir) => {
        return dir.available >= config.storage_config.minimum_gb_required;
    });

    if(dirsWithSomeSpace.length == 0) {
        var err = new Error(
            'NoDiskspaceError: No disks exist that have space available.',
            '\nMinimum Space Required:',
            config.storage_config.minimum_gb_required,
            'Disks Checked:',
            dirSpaces
        )
        logger.err(err);
        process.exit(-1);
    }

    return dirsWithSomeSpace.sort((a, b) => {
        return a.available > b.available ? 1 : -1;
    })[0];
}

async function resolveAPKSaveInfo(appData) {
    const appsSaveDir = await getLocationWithLeastSpace();
    const filesystem = await getPathFileSystem(appsSaveDir.path);
    const UUID = await getUUID(filesystem);

    logger.debug(`appdir: ${appsSaveDir.path} - space remaining: ${appsSaveDir.available}`,
        `\nappId ${appData.app}`, `\nappStore ${appData.store}`,
        `\nregion ${appData.region}`, `\nversion ${appData.version}`);

    const appSavePath = path.join(appsSaveDir.path, 'apps', appData.app,
        appData.store, appData.region, appData.version);
    logger.info(`App Save Directory formed: '${appSavePath}'`);

    await mkdirp(appSavePath);

    return {
        appSavePath: appSavePath,
        appSaveFS: filesystem,
        appSaveFSName: appsSaveDir.name,
        appSavePathRoot: appsSaveDir.path,
        appSaveUUID: UUID,
    };
}

function downloadApp(appData, appSavePath) {
    // Command line args for gplay cli
    const args = [
        '-pd', appData.app,
        '-f', appSavePath,
        '-c', config.system_config.downloader_credentials,
    ];
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
    let appSaveInfo;
    try {
        appSaveInfo = await resolveAPKSaveInfo(app);
    } catch (err) {
        await new Promise((resolve) => setTimeout(resolve, 6000));
        return Promise.reject(`Did not have access to resolve dir: ${err.message}`);
    }

    try {
        await downloadApp(app, appSaveInfo.appSavePath);
    } catch (err) {
        logger.debug('Attempting to remove created dir');
        await fs.rmdir(appSaveInfo.appSavePath).catch(logger.warning);
        return Promise.reject(`Downloading failed with err: ${err}`);
    }

    try {
        const apkPath = path.join(appSaveInfo.appSavePath, `${app.app}.apk`);

        if (fs.existsSync(apkPath)) {
            // Perform a check on apk size
            await fs.stat(apkPath, async(err, stats) => {
                if (stats.size == 0 || stats.size == undefined) {
                    await fs.rmdir(appSaveInfo.appSavePath).catch(logger.warning);
                    return Promise.reject('File did not successfully download and is a empty size');
                }

                await db.updateDownloadedApp(app, appSaveInfo, config.system_config.vm_name);
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

async function main() {
    // Ensure that directory structures exist.
    downloadLocations = await ensureDirectoriesExist(
        config.storage_config.apk_download_directories
    );

    for (;;) {
        let apps;
        try {
            apps = await db.queryAppsToDownload(4000);
        } catch (err) {
            await new Promise((resolve) => setTimeout(resolve, 1000));
            continue;
        }

        for (const app of apps ) {
            try {
                await download(app).catch((err) => {
                    throw err;
                });
            } catch (err) {
                logger.err(
                    `Error Downloading application with package name: ${app.app}.`,
                    `Error: ${err}`
                );
            }
        }

        // await Promise.each(apps, download).catch(logger.err);
    }
}

main();
