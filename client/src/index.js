'use strict'
const Fingerprint2 = require('fingerprintjs2');
console.log("Entry");

new Fingerprint2().get(function(result, components){
  console.log(result); //a hash, representing your device fingerprint
  console.log(components); // an array of FP components

});
