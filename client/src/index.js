'use strict'
const Fingerprint2 = require('fingerprintjs2');

new Fingerprint2().get(function(result, components){
  console.log(result); //a hash, representing your device fingerprint
  console.log(components); // an array of FP components
  var fpResult = {}
  components.forEach(element => {
    fpResult[element.key] = element.value;
  });
  console.log(JSON.stringify({fp: fpResult}))
  fetch('/un-scrapy/result', {credentials: 'same-origin', method: 'POST', body: JSON.stringify({fp: fpResult})})
    .then(result => {
      if (!result.ok) {
        console.error('Error while validating browser', err)
        return document.write('<h1>Unexpected error, please reload the page</h1>');
      }
      result.json().then(data => {
        if (data.status === 'Success') {
          return window.location.reload(true);
        }
        return document.write('<h1>Access Denied</h1>');
      })
    })
    .catch(err => {
      console.error('Error while validating browser', err)
      return document.write('<h1>Unexpected error, please reload the page</h1>');
    });
});
