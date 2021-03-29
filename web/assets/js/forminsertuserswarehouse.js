function init() {
    resetForms();
    addForm();

    document.addEventListener(
        "DOMContentLoaded",
        function () {
            document.getElementById("add-condition").onclick = addForm;
            document.getElementById("clear-form").onclick = init;
            document.getElementById("apply").onclick = sendData;
        },
        false
    );
}

//global variables
var mainDiv = document.getElementById("form"); //ref to main div with forms

//attributes for form fields, they will be added during the creation of forms
var formElementsTypes = ["input", "input", "input", "input", "select", "input", "button"];
var formElementsIds = ["email", "password", "firstName", "lastName", "role", "tabel"];
var formElementsClassNames = [
    "form-control form-control-sm",
    "form-control form-control-sm",
    "form-control form-control-sm",
    "form-control form-control-sm",
    "form-control form-control-sm",
    "form-control form-control-sm",
    "deleteForm btn",
];

var formAllowedIdsArr = [];
var rowsAmount = 50;
var finalData = [];
var labels = ["email", "Пароль", "Имя", "Фамилия", "Роль", "Табель"];
var roles = [
/*{
        
        id: 'Выберите позицию',
        title: 'Выберите позицию'
    },*/
    {
        id: 'старший кладовщик склада',
        title: 'старший кладовщик склада'
    },
    {
        id: 'кладовщик склада',
        title: 'кладовщик склада'
    }
];

//creating form and adding attributes
function createForm() {
    var form = document.createElement("div");
    form.className = "form-inline";
    form.id = formAllowedIdsArr[0];
    formAllowedIdsArr = formAllowedIdsArr.slice(1);
    mainDiv.appendChild(form);
    for (var i = 0; i < formElementsTypes.length; i++) {
        var element = document.createElement(formElementsTypes[i]);
        element.id = `${formElementsIds[i]}-${form.id}`;
        element.className = formElementsClassNames[i].concat(" ", "fields-style");
        if (i != 6) {
            element.placeholder = labels[i];
        }

        if (formElementsTypes[i] === 'select' && formElementsIds[i] === 'role') { // && formElementsTypes[i] === 'required'
         //   if (formElementsTypes[i] === 'label') {
         //     roles.forEach((label) => {
         //           element.innerHTML += `<option label="TTTTTT"</option>`;
                    roles.forEach((role) => {
                
                        element.innerHTML += `<option value="${role.id}">${role.title}</option>`;
                    });
         //       });
            
        //}
        }

        form.appendChild(element);

        if (i == 6) {
            var icon = document.createElement("i");
            var deleteButton = document.getElementById(element.id);
            deleteButton.appendChild(icon);
            icon.className = "fa fa-times";
        }
    }

    //add eventlisteners to newly created elements

    document
        .getElementById(`${formElementsIds[6]}-${form.id}`)
        .addEventListener("click", deleteForm);

    canShowAddButton(); //check if we can show add button
    canShowDeleteButton(); //check if we can show delete button
}

function createSelectMenuOptions(array, menu) {
    for (var i = 0; i < array.length; i++) {
        var option = document.createElement("option");
        option.value = array[i];
        option.text = array[i];
        menu.appendChild(option);
    }
}

//check if we can add form and add it
function addForm() {
    if (formAllowedIdsArr.length > 0) {
        createForm();
    }
}

function deleteForm() {
    formAllowedIdsArr.push(this.parentNode.id);
    this.parentNode.parentNode.removeChild(this.parentNode);
    canShowAddButton();
    canShowDeleteButton();
}

function formIdGenerator(num) {
    var formIdArray = [];
    for (i = 0; i < num; i++) {
        formIdArray.push(`form-${i + 1}`);
    }
    return formIdArray;
}

//reset all forms, fields, console and generate new allowed IDs array
function resetForms() {
    formAllowedIdsArr = formIdGenerator(rowsAmount);
    while (mainDiv.firstChild) {
        mainDiv.removeChild(mainDiv.firstChild);
    }
    finalData = [];
    //console.clear();
}

function canShowAddButton() {
    if (formAllowedIdsArr.length == 0) {
        document.getElementById("add-condition").disabled = true;
    } else {
        document.getElementById("add-condition").disabled = false;
    }
}

function canShowDeleteButton() {
    if (formAllowedIdsArr.length < rowsAmount - 1) {
        document.getElementsByClassName("deleteForm")[0].disabled = false;
    } else {
        document.getElementsByClassName("deleteForm")[0].disabled = true;
    }
}

//function checkData() {}

function getData() {
    var usedRows = rowsAmount - formAllowedIdsArr.length;
    var invalidResult = false;
    for (i = 0; i < usedRows; i++) {
        var resultObj = {};
        for (j = 0; j < mainDiv.childNodes[i].childNodes.length - 1; j++) {
            var resultValue = mainDiv.childNodes[i].childNodes[j].value;
            var prop = formElementsIds[j];
            /*if (prop == "qty") {
              resultObj[`${prop}`] = parseInt(resultValue);
            } else {
              resultObj[`${prop}`] = resultValue;
            }*/
            resultObj[`${prop}`] = resultValue;
        }
    //    var validCodeDebitor = /\bK\d{5}\b/;


        currentDiv = mainDiv.childNodes[i];

    /*    if (
       //     !validCodeDebitor.test(resultObj["code_debitor"]) ||
            resultObj["code_debitor"] != " "
         //   parseInt(resultObj["qty"]) < 1
        ) {
            invalidResult = true;
            //  var validMaterial = /\b31\d{5}\b/;
            drawErrorMessage(mainDiv.childNodes[i]);
        } else { */
          //  resultObj["qty"] = parseInt(resultObj["qty"]);
            finalData.push(resultObj);
            drawSuccessMessage(mainDiv.childNodes[i]);
      //  }
    }
    console.log(finalData);
    return JSON.stringify(finalData);
}

function drawErrorMessage(div) {
    var errorMessage = document.createElement("div");
    errorMessage.id = "error-message";
    errorMessage.innerHTML = "Ошибка при передаче данных";
    //mainDiv.appendChild(errorMessage);
    div.appendChild(errorMessage);
}

function drawSuccessMessage(div) {
    var successMessage = document.createElement("div");
    successMessage.id = "success-message";
    successMessage.innerHTML = "Данные переданы успешно";
    //mainDiv.appendChild(successMessage);
    div.appendChild(successMessage);
}

function sendData() {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "http://localhost:3001/createuserswarehouse", true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(getData());

    xhr.onreadystatechange = function () {
      //  this.readyState == 4 && this.status == 200
            if (this.readyState == 4 && this.status == 200) {
                drawSuccessMessage();
            } else {
                drawErrorMessage();
           }
    };
}

//let's go
init();