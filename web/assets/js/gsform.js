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
var formElementsTypes = ["input", "input", "input", "button"];
var formElementsIds = ["material", "qty", "comment"];
var formElementsClassNames = [
    "form-control form-control-sm",
    "form-control form-control-sm",
    "form-control form-control-sm",
    "deleteForm btn",
];

var formAllowedIdsArr = [];
var rowsAmount = 50;
var finalData = [];
var labels = ["Material", "Qty", "Comment"];

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
        if (i != 3) {
            element.placeholder = labels[i];
        }
        form.appendChild(element);

        if (i == 3) {
            var icon = document.createElement("i");
            var deleteButton = document.getElementById(element.id);
            deleteButton.appendChild(icon);
            icon.className = "fas fa-times";
        }
    }

    //add eventlisteners to newly created elements

    document
        .getElementById(`${formElementsIds[3]}-${form.id}`)
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
        var validMaterial = /\b31\d{5}\b/;

        currentDiv = mainDiv.childNodes[i];

        if (
            !validMaterial.test(resultObj["material"]) ||
            parseInt(resultObj["qty"]) < 0
        ) {
            invalidResult = true;
            //addValidationClass(" is-invalid", currentDiv);
            //currentDiv.className += " is-invalid";
            //console.log(currentDiv.childNodes);
            drawErrorMessage(mainDiv.childNodes[i]);

        } else {
            resultObj["qty"] = parseInt(resultObj["qty"]);
            finalData.push(resultObj);
            //addValidationClass(" is-valid", currentDiv);
            drawSuccessMessage(mainDiv.childNodes[i]);
        }
    }
    console.log(finalData);
    return JSON.stringify(finalData);
}

/*function addValidationClass(addedClass, div) {
    for (i = 0; i < div.childNodes - 1; i++) {
        console.log(i)
        //div.childNodes[i].className += addedClass
    }
}*/

function drawErrorMessage(div) {
    var errorMessage = document.createElement("div");
    errorMessage.id = "error-message";
    errorMessage.innerHTML = "Ошибка при передаче данных";
    div.appendChild(errorMessage);
}

function drawSuccessMessage(div) {
    var successMessage = document.createElement("div");
    successMessage.id = "success-message";
    successMessage.innerHTML = "Данные переданы успешно";
    div.appendChild(successMessage);
}

function sendData() {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "http://10.1.20.110:3001/shipmentbysap", true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(getData());

    xhr.onreadystatechange = function () {
        this.readyState == 4 && this.status == 200
    };

    //resetForms();
    //addForm();
}

//let's go
init();