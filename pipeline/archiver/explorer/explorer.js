
const gplay = require('google-play-scraper');
const Promise = require('bluebird');

const logger = require('../../util/logger');
const DB = require('../../db/db');
const db = new DB('explorer');



/*
TODO:
 */
function flatten(arr) {
    return arr.reduce((a, b) => a.concat(b), []);
}

/**
 * Creates a cartesion product of arrays of strings.
 * https://stackoverflow.com/questions/12303989/cartesian-product-of-multiple-arrays-in-javascript
 * Eg, ['a', 'b', 'c'] x2 => ['aa' ''ab' 'ac' 'ba' 'bb'] ...
 */
function cartesianProductChars(...args) {
    return args.reduce((prods, arr) =>
        flatten(prods.map((prod) => arr.map((v) => prod.concat(v)))), [[]]);
}

/**
 * Creates a file of suggestions made by Google play when passing
 * the start of strings, eg. 'a', 'b', 'aa', 'ab' ...
 *
 * @param {*The list of words used to get autocompletes} startingWords
 */
function scrapeSuggestedWords(startingWords) {
    Promise.each(startingWords, (letter) => {
        return gplay.suggest({ term: letter, throttle: 10, region: 'uk' })
            .then((suggestion) => {
                Promise.each(suggestion, (word) => {
                    logger.debug(`Inserting to DB: ${word}`);
                    return db.insertSearchTerm(word).catch((err) => logger.err(err));
                }).catch(logger.err);
            }).catch(logger.err);
    });
}

//Alphabet 
const single = 'abcdefghijklmnopqrstuvwxyz '.split('');
const double = cartesianProductChars(single, single);
const triple = cartesianProductChars(single, single, single);

const charTriples = single.concat(double).concat(triple);

scrapeSuggestedWords(charTriples);
