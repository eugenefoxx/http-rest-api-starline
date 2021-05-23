function myFunction() {
    var input, filter, table, tr, td, cell, i, j;
    input = document.getElementById("myInput");
    filter = input.value.toUpperCase();
    table = document.getElementById("myTable");
    tr = table.getElementsByTagName("tr");
    for (i = 1; i < tr.length; i++) {
        // Hide the row initially.
        tr[i].style.display = "none";

        td = tr[i].getElementsByTagName("td");
        for (var j = 0; j < td.length; j++) {
            cell = tr[i].getElementsByTagName("td")[j];
            if (cell) {
                if (cell.innerHTML.toUpperCase().indexOf(filter) > -1) {
                    tr[i].style.display = "";
                    break;
                }
            }
        }
    }
}

function searchDebitorModalWindow() {
    var input, filter, table, tr, td, cell, i, j;
    input = document.getElementById("debitorModalInput");
    filter = input.value.toUpperCase();
    table = document.getElementById("debitorModalTable");
    tr = table.getElementsByTagName("tr");
    for (i = 1; i < tr.length; i++) {
        // Hide the row initially.
        tr[i].style.display = "none";

        td = tr[i].getElementsByTagName("td");
        for (var j = 0; j < td.length; j++) {
            cell = tr[i].getElementsByTagName("td")[j];
            if (cell) {
                if (cell.innerHTML.toUpperCase().indexOf(filter) > -1) {
                    tr[i].style.display = "";
                    break;
                }
            }
        }
    }
}

function sortTable(n) {
    var table, rows, switching, i, x, y, shouldSwitch, dir, switchcount = 0;
    table = document.getElementById("myTable");
    switching = true;
    //Set the sorting direction to ascending:
    dir = "asc";
    /*Make a loop that will continue until
    no switching has been done:*/
    while (switching) {
        //start by saying: no switching is done:
        switching = false;
        rows = table.rows;
        /*Loop through all table rows (except the
        first, which contains table headers):*/
        for (i = 1; i < (rows.length - 1); i++) {
            //start by saying there should be no switching:
            shouldSwitch = false;
            /*Get the two elements you want to compare,
            one from current row and one from the next:*/
            x = rows[i].getElementsByTagName("TD")[n];
            y = rows[i + 1].getElementsByTagName("TD")[n];
            /*check if the two rows should switch place,
            based on the direction, asc or desc:*/
            if (dir == "asc") {
                if (x.innerHTML.toLowerCase() > y.innerHTML.toLowerCase()) {
                    //if so, mark as a switch and break the loop:
                    shouldSwitch = true;
                    break;
                }
            } else if (dir == "desc") {
                if (x.innerHTML.toLowerCase() < y.innerHTML.toLowerCase()) {
                    //if so, mark as a switch and break the loop:
                    shouldSwitch = true;
                    break;
                }
            } else if (Number(x.innerHTML) > Number(y.innerHTML)) {
                //if so, mark as a switch and break the loop:
                shouldSwitch = true;
                break;
            }
        }
        if (shouldSwitch) {
            /*If a switch has been marked, make the switch
            and mark that a switch has been done:*/
            rows[i].parentNode.insertBefore(rows[i + 1], rows[i]);
            switching = true;
            //Each time a switch is done, increase this count by 1:
            switchcount++;
        } else {
            /*If no switching has been done AND the direction is "asc",
            set the direction to "desc" and run the while loop again.*/
            if (switchcount == 0 && dir == "asc") {
                dir = "desc";
                switching = true;
            }
        }
    }
}

function sortDebitorModalTable(n) {
    var table, rows, switching, i, x, y, shouldSwitch, dir, switchcount = 0;
    table = document.getElementById("debitorModalTable");
    switching = true;
    //Set the sorting direction to ascending:
    dir = "asc";
    /*Make a loop that will continue until
    no switching has been done:*/
    while (switching) {
        //start by saying: no switching is done:
        switching = false;
        rows = table.rows;
        /*Loop through all table rows (except the
        first, which contains table headers):*/
        for (i = 1; i < (rows.length - 1); i++) {
            //start by saying there should be no switching:
            shouldSwitch = false;
            /*Get the two elements you want to compare,
            one from current row and one from the next:*/
            x = rows[i].getElementsByTagName("TD")[n];
            y = rows[i + 1].getElementsByTagName("TD")[n];
            /*check if the two rows should switch place,
            based on the direction, asc or desc:*/
            if (dir == "asc") {
                if (x.innerHTML.toLowerCase() > y.innerHTML.toLowerCase()) {
                    //if so, mark as a switch and break the loop:
                    shouldSwitch = true;
                    break;
                }
            } else if (dir == "desc") {
                if (x.innerHTML.toLowerCase() < y.innerHTML.toLowerCase()) {
                    //if so, mark as a switch and break the loop:
                    shouldSwitch = true;
                    break;
                }
            } else if (Number(x.innerHTML) > Number(y.innerHTML)) {
                //if so, mark as a switch and break the loop:
                shouldSwitch = true;
                break;
            }
        }
        if (shouldSwitch) {
            /*If a switch has been marked, make the switch
            and mark that a switch has been done:*/
            rows[i].parentNode.insertBefore(rows[i + 1], rows[i]);
            switching = true;
            //Each time a switch is done, increase this count by 1:
            switchcount++;
        } else {
            /*If no switching has been done AND the direction is "asc",
            set the direction to "desc" and run the while loop again.*/
            if (switchcount == 0 && dir == "asc") {
                dir = "desc";
                switching = true;
            }
        }
    }
}

function alternate(id) {

    if (document.getElementsByTagName) {

        var table = document.getElementById(id);

        var rows = table.getElementsByTagName("tr");

        for (i = 0; i < rows.length; i++) {

            //manipulate rows 

            if (i % 2 == 0) {

                rows[i].className = "even";

            } else {

                rows[i].className = "odd";

            }

        }

    }

}

function editVendor(id) {
    axios.get('/operation/updatevendor/' + id)
        .then(function (response) {
            document
                .getElementById('modal-2')
                .getElementsByClassName('modal__dialog-body')[0].innerHTML = response.data;

        });
}

function editInspection(id) {
    axios.get('/operation/updateinspection/' + id)
        .then(function (response) {
            document
                .getElementById('modal-2')
                .getElementsByClassName('modal__dialog-body')[0].innerHTML = response.data;
            //    debugger;
                let sendForm = document.querySelector('#send-form');
                console.log(sendForm);
                //    let form = document.querySelector('form');
                
                if (sendForm !== null) {
                //  alert("Hi");
                  
                sendForm.onclick = function (event) {
                    event.preventDefault();
                    // Get the form
                    let form = document.querySelector('#update-component-form');
                    // Get all field data from the form
                    // returns a FormData object
                    let data = new FormData(form);
                    console.log('work');
                    //  console.log(serialize(form));
                    // Convert to an object
                    let formObj = serialize(data);
                    console.log(formObj);
                 /*   let resultData = [];
                    resultData = serialize(data);
                    console.log(resultData);*/
                    //  resultData = formObj;
                    //arrRezalt = serializeArray(form);
                   // console.log(arrRezalt);
                    // action="/operation/updateinspection/{{.GET.ID}}"
                    fetch('http://localhost:3001/operation/updateinspection',{
  
		                method: 'POST',
		                body: JSON.stringify(serialize(data)),
		                headers: {
			                'Content-type': 'application/json; charset=UTF-8'
		                }
	                }).then(function (response) {
		                if (response.ok) {
                          	return response.json();
		                }
		                return Promise.reject(response);
	                }).then(function (data) {
                      /*  const updateid = document.querySelector('td[updateid-id="'+data.id+'"]')
                        if (!updateid || !updateid.innerHTML){
                            return false;
                        }*/
                       //  debugger;
                        console.log(data.id)
                        console.log(data.status)
                       // console.log("data.id", tt)
                      //  document.querySelector('#updatestatus').innerHTML = data.status;
                        document.querySelector('td[updateid-status="'+data.id+'"]').innerHTML = data.status;
                        console.log(data.note)
                     //   if (data.note = undefined){
                     //       return false;
                     //   }
                      //  document.querySelector('#updatenote').innerHTML = data.note;
                      //  if {
                        document.querySelector('td[updateid-note="'+data.id+'"]').innerHTML = data.note;
		                console.log(data.message)
                //  document.querySelector('label[response="'+data.message+'"]').innerText = data.message;
                        document.querySelector('#responseInspection').innerText = data.message;
                        //}
                        console.log(data);
                    //}
	                }).catch(function (error) {
		                console.warn(error);
	                });
                }
              }  
        });
        
}

function editInspectionMix(id) {
    axios.get('/operation/updateinspectionmix/' + id)
        .then(function (response) {
            document
                .getElementById('modal-2')
                .getElementsByClassName('modal__dialog-body')[0].innerHTML = response.data;
            //    debugger;
                let sendForm = document.querySelector('#send-form');
                console.log(sendForm);
                //    let form = document.querySelector('form');
                
                if (sendForm !== null) {
                //  alert("Hi");
                  
                sendForm.onclick = function (event) {
                    event.preventDefault();
                    // Get the form
                    let form = document.querySelector('#update-component-form');
                    // Get all field data from the form
                    // returns a FormData object
                    let data = new FormData(form);
                    console.log('work');
                    //  console.log(serialize(form));
                    // Convert to an object
                    let formObj = serialize(data);
                    console.log(formObj);
                 /*   let resultData = [];
                    resultData = serialize(data);
                    console.log(resultData);*/
                    //  resultData = formObj;
                    //arrRezalt = serializeArray(form);
                   // console.log(arrRezalt);
                    // action="/operation/updateinspection/{{.GET.ID}}"
                    fetch('http://localhost:3001/operation/updateinspectionmix',{
  
		                method: 'POST',
		                body: JSON.stringify(serialize(data)),
		                headers: {
			                'Content-type': 'application/json; charset=UTF-8'
		                }
	                }).then(function (response) {
		                if (response.ok) {
                          	return response.json();
		                }
		                return Promise.reject(response);
	                }).then(function (data) {
                      /*  const updateid = document.querySelector('td[updateid-id="'+data.id+'"]')
                        if (!updateid || !updateid.innerHTML){
                            return false;
                        }*/
                        
                        console.log(data.id)
                        console.log(data.status)
                       // console.log("data.id", tt)
                      //  document.querySelector('#updatestatus').innerHTML = data.status;
                        document.querySelector('td[updateid-status="'+data.id+'"]').innerHTML = data.status;
                        console.log(data.note)
                     //   if (data.note = undefined){
                     //       return false;
                     //   }
                      //  document.querySelector('#updatenote').innerHTML = data.note;
                      //  if {
                        document.querySelector('td[updateid-note="'+data.id+'"]').innerHTML = data.note;
		                console.log(data.message)
                //  document.querySelector('label[response="'+data.message+'"]').innerText = data.message;
                        document.querySelector('#responseInspection').innerText = data.message;
                        //}
                        console.log(data);
                    //}
	                }).catch(function (error) {
		                console.warn(error);
	                });
                }
              }  
        });
        
}

function serializeArray (form) {
    // Create a new FormData object
    const formData = new FormData(form);
  
    // Create an array to hold the name/value pairs
    const pairs = [];
  
    // Add each name/value pair to the array
    for (const [name, value] of formData) {
      pairs.push({ name, value });
    }
  
    // Return the array
    return pairs;
  }

function serialize (data) {
    
	let obj = {};
	for (let [key, value] of data) {
		if (obj[key] !== undefined) {
			if (!Array.isArray(obj[key])) {
				obj[key] = [obj[key]];
			}
			obj[key].push(value);
		} else {
			obj[key] = value;
		}
	}
    return [obj];
  
 /*
  var array = Object.keys(obj)
    .map(function(key) {
        return obj[key];
    });
	console.log(array);
    return array;
    */
}

/*
function editInspection(id) {
    axios.get('/operation/updateinspection/' + id)
        .then(function (response) {
            document
                .getElementById('modal-2')
                .getElementsByClassName('modal__dialog-body')[0].innerHTML = response.data;

        });
}
*/
function editUserQuality(id) {
    axios.get('/operation/updateuserquality/' + id)
        .then(function (response) {
            document
                .getElementById('modal-2')
                .getElementsByClassName('modal__dialog-body')[0].innerHTML = response.data;

        });
  }

  function acceptInspection(id) {
    axios.get('/operation/acceptinspectiontowh/' + id)
        .then(function (response) {
            document
                .getElementById('modal-4')
                .getElementsByClassName('modal__dialog-body')[0].innerHTML = response.data;

        });
  }

  function editUserWarehouse(id) {
    axios.get('/operation/updateuserwarehouse/' + id)
        .then(function (response) {
            document
                .getElementById('modal-6')
                .getElementsByClassName('modal__dialog-body')[0].innerHTML = response.data;

        });
  }

  function redirectTostatusinspection() {
    window.location.replace("http://localhost:3001/operation/statusinspection");
}