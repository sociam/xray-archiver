/* eslint no-console: 0 */
/*
To log error messages to systemd.
*/

const util = require('util');
const EMERG = 0,
    ALERT = 1,
    CRIT = 2,
    ERR = 3,
    WARNING = 4,
    NOTICE = 5,
    INFO = 6,
    DEBUG = 7;

const prefixes = ['<0>', '<1>', '<2>', '<3>', '<4>', '<5>', '<6>', '<7>'];

function log(level, args) {
    const format = util.format.apply(null, args).replace(/\n/g, `\n${prefixes[level]}`);
    process.stdout.write(`${prefixes[level] + format}\n`);
}

module.exports = {
    emerg: function(...args) {
        log(EMERG, args);
    },
    alert: function(...args) {
        log(ALERT, args);
    },
    crit: function(...args) {
        log(CRIT, args);
    },
    err: function(...args) {
        log(ERR, args);
    },
    warning: function(...args) {
        log(WARNING, args);
    },
    notice: function(...args) {
        log(NOTICE, args);
    },
    info: function(...args) {
        log(INFO, args);
    },
    debug: function(...args) {
        log(DEBUG, args);
    },
};
