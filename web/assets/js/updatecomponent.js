//const btn = document.querySelector('button');

function updateInspectionComponent() {
/*const form = document.querySelector('form[name="valform"]');
const status = form.elements['status'].value;
const note = form.elements['note'].value;

alert(form);
console.log(note);
alert(note);*/

let sendForm = document.querySelector('#send-form');
let form = document.querySelector('form');
sendForm.onclick = function (event) {
  debugger;
  event.preventDefault();
  console.log('work');
  console.log(serialize(form));
}
/*
var resultData = [];
    var forms = document.getElementsByClassName("form-control");
    for (i = 0; i < forms.length; i++) {
        if (forms[i].value) {
            resultData.push({ updateInspecionEO: forms[i].value });
        }
    }     
console.log(resultData);
return resultData;
*/
}
/*
data = updateInspectionComponent();
console.log("data - ", data);
async function postDataIn(data) {
//  debugger;
  let response = await fetch('http://localhost:3001/operation/updateinspection', {
    method: 'POST',
    body: JSON.stringify(data),
    headers: {
      'Content-Type': 'application/json;charset=utf-8'
    }
  });
}
*/
/*
data = sendupdateInspectionComponent();
let response = await fetch('http://localhost:3001/operation/updateinspection', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json;charset=utf-8'
  },
  body: JSON.stringify(data)
});

let result = await response.json();
alert(result.message);
*/
//data = sendupdateInspectionComponent();
/*
async function postData(url = '', data = {}) {
    // Default options are marked with *
    const response = await fetch(url, {
      method: 'GET', // GET, POST, PUT, DELETE, etc.
  //    mode: 'cors', // no-cors, *cors, same-origin
      cache: 'no-cache', // default, no-cache, reload, force-cache, only-if-cached
      credentials: 'same-origin', // include, *same-origin, omit
      headers: {
        'Content-Type': 'application/json'
        // 'Content-Type': 'application/x-www-form-urlencoded',
      },
      redirect: 'follow', // manual, *follow, error
      referrerPolicy: 'no-referrer', // no-referrer, *client
      body: JSON.stringify(data) // body data type must match "Content-Type" header
      
    });
//    console.log("date send", data)
    return await response.json(); // parses JSON response into native JavaScript objects
  }

  dataToSend = sendupdateInspectionComponent();
//  console.log("dataToSend", dataToSend)
  postData('http://localhost:3001/operation/updateinspection', dataToSend)
  .then((data) => {
    console.log(data); // JSON data parsed by `response.json()` call
  });
*/ 
/*
btn.addEventListener( 'click', function() {

    sendupdate()

// let elements = document
//sendData( {test:'ok'} );
})
*/
