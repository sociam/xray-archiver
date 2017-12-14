var ref = require('ref');
var ffi = require('ffi');

// typedef
var bool = ref.refType(ref.types.bool);
var stringPtr = ref.refType(ref.types.CString);

// binding to a few "libsqlite3" functions...
var basicAnalyse = ffi.Library('basicAnalyze', {
  'simpleAnalyze': [ stringPtr, [ 'string' ] ],
  'checkReflect': [ 'bool', [ 'string' ] ]
});

var result = basicAnalyse.checkReflect(process.argv[0])
console.log(result)
