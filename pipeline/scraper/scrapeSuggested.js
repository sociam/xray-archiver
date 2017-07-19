/**
 * Get the suggested search Autocompetes and find the apps associated with them.
 */


var gplay = require('google-play-scraper');
var alphabet = require('alphabet');
var _ = require('lodash');
var config = require('/etc/xray/config.json');
var fs = require('fs-extra');

_.forEach(alphabet.lower, (letter) => {
    gplay.suggest({ term: letter })
        .then(
            (suggestion) => {
                _.forEach(suggestion, (word) => {
                    console.log(word);
                    gplay.suggest({ term: word })
                        .then((suggestion) => {
                            _.forEach(suggestion, (suggestion) => {
                                // TODO: Output to a file for words.
                                console.log(suggestion);
                            })
                        })
                })
            },
            (err) => console.log(err)
        );
});