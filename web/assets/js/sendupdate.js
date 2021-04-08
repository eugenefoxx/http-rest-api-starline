//const btn = document.querySelector('button');

function sendupdateInspectionComponent() {
/*const form = document.querySelector('form[name="valform"]');
const status = form.elements['status'].value;
const note = form.elements['note'].value;

alert(status);
console.log(note);
alert(note);*/

var resultData = [];
    var forms = document.getElementsByClassName("form-control");
    for (i = 0; i < forms.length; i++) {
        if (forms[i].value) {
            resultData.push({ updateInspecionEO: forms[i].value });
        }
    }     
console.log(resultData);
return resultData;
}

async function postData(url = '', data = {}) {
    // Default options are marked with *
    const response = await fetch(url, {
      method: 'POST', // GET, POST, PUT, DELETE, etc.
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
    return await response.json(); // parses JSON response into native JavaScript objects
  }

  dataToSend = sendupdateInspectionComponent();

  postData('http://localhost:3001/updateinspection', dataToSend)
  .then((data) => {
    console.log(data); // JSON data parsed by `response.json()` call
  });
/*
btn.addEventListener( 'click', function() {

    sendupdate()

// let elements = document
//sendData( {test:'ok'} );
})
*/
