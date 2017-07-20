//Logging mechanisim for script
const EMERG = 0,
    ALERT = 1,
    CRIT = 2,
    ERR = 3,
    WARNING = 4,
    NOTICE = 5,
    INFO = 6,
    DEBUG = 7;

var prefixes = ['<0>', '<1>', '<2>', '<3>', '<4>', '<5>', '<6>', '<7>'];

modules.exports = {
    info: function(txt) {
        console.log(prefixes[INFO], txt);
    },
    err: function(txt) {
        console.log(prefixes[ERR], txt);
    },
    alert: function(txt) {
        console.log(prefixes[ALERT], txt);
    },
    crit: function(txt) {
        console.log(prefixes[CRIT], txt);
    },
    warn: function(txt) {
        console.log(prefixes[WARNING], txt);
    },
    notice: function(txt) {
        console.log(prefixes[NOTICE], txt);
    },
    debug: function(txt) {
        console.log(prefixes[DEBUG], txt);
    }
}
logger.info("Logger initialised");

