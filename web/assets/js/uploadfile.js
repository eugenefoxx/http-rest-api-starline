
var windowOnloadAdd = function (event) {
    if ( window.onload ){
       window.onload = window.onload + event;
    } else {
       window.onload = event;
    };
 };
 windowOnloadAdd(function() {
     // Select your input type file and store it in a variable
const input = document.getElementById('fileUpload');

// This will upload the file after having read it
const upload = (file) => {
  fetch('http://localhost:3001/operation/uploadfile', { // Your POST endpoint
    method: 'POST',
    headers: {
      // Content-Type may need to be completely **omitted**
      // or you may need something
      'Content-Type': 'multipart/form-data;'
    },
    body: file // This is your file object
  }).then(
    response => response.json() // if the response is a JSON object
  ).then(data => {
    console.log(data.message)
    document.querySelector('#responseUpdateFile').innerText = data.message;
  console.log(data)
    }).catch(
    error => console.log(error) // Handle the error response object
  );
};

// Event handler executed when a file is selected
const onSelectFile = () => upload(input.files[0]);

// Add a listener on your input
// It will be triggered when a file will be selected
input.addEventListener('change', onSelectFile, false);

/*document.querySelector('#fileUpload').addEventListener('change', event => {
    handleImageUpload(event)
  })

  const handleImageUpload = event => {
    const files = event.target.files
    const formData = new FormData()
    formData.append('myFile', files[0])
  
    fetch('http://localhost:3001/operation/uploadfile', {
      method: 'POST',
      body: formData
    })
    .then(response => response.json())
    .then(data => {
        console.log(data.message)
        document.querySelector('#responseUpdateFile').innerText = data.message;
      console.log(data)
    })
    .catch(error => {
      console.error(error)
    })
  }*/
});

