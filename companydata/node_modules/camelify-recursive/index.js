var changeCase = require('change-case');
require('es6-shim');

var log = console.log.bind(log);

var kind = function(item) {
	var getPrototype = function(item) {
		return Object.prototype.toString.call(item).slice(8, -1);
	};
	var kind, Undefined;
	if (item === null ) {
		kind = 'null';
	} else {
		if ( item === Undefined ) {
			kind = 'undefined';
		} else {
			var prototype = getPrototype(item);
			if ( ( prototype === 'Number' ) && isNaN(item) ) {
				kind = 'NaN';
			} else {
				kind = prototype;
			}
		}
	}
	return kind;
};

// Inspired from https://gist.github.com/Sneagan/8366247
var camelifyObject = function(obj) {
	var keys = Object.keys(obj)
	var convertKey = function(oldKeyName) {
		if ( oldKeyName.includes('_') ) {
			var oldValue = obj[oldKeyName]
			var newKeyName = changeCase.camelCase(oldKeyName);
			obj[newKeyName] = oldValue;
			delete obj[oldKeyName];
		}
	}
	keys.forEach(function(key, index){
		var value = obj[key]

		if ( kind(value) === 'Object' ) {
			obj[key] = camelifyObject(value);
		}
		if ( kind(value) === 'Array' ) {
			value.forEach(function(item, index){
				value = camelifyObject(obj[key][index]);
			});
		}
		// Now we've fixed any potential subkeys, convert the actual key
		convertKey(key);
	});
	return obj;
};

module.exports = camelifyObject
