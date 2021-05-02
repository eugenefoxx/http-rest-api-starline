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
var formElementsTypes = ["input", "input", "input", "input", "button"];
var formElementsIds = ["scanid", "qtyfact", "qtysap", "qtypanacim"];
var formElementsClassNames = [
    "form-control form-control-sm",
    "form-control form-control-sm",
    "form-control form-control-sm",
    "form-control form-control-sm",
    "deleteForm btn",
];

var formAllowedIdsArr = [];
var rowsAmount = 50;
var finalData = [];
var labels = ["Scan ID", "Qty Fact", "Qty SAP", "Qty Panacim"];

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
        if (i != 4) {
            element.placeholder = labels[i];
        }
        form.appendChild(element);

        if (i == 4) {
            var icon = document.createElement("i");
            var deleteButton = document.getElementById(element.id);
            deleteButton.appendChild(icon);
            icon.className = "fa fa-times";
        }
    }

    //add eventlisteners to newly created elements

    document
        .getElementById(`${formElementsIds[4]}-${form.id}`)
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
    //  console.clear();
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

function getData() {
    var usedRows = rowsAmount - formAllowedIdsArr.length;
    for (i = 0; i < usedRows; i++) {
        var resultObj = {};
        for (j = 0; j < mainDiv.childNodes[i].childNodes.length - 1; j++) {
            var resultValue = mainDiv.childNodes[i].childNodes[j].value;
            var prop = formElementsIds[j];
            if (prop == "scanid") {
                resultObj[`${prop}`] = resultValue;
            } else {
                resultObj[`${prop}`] = parseInt(resultValue);
            }
        }
        finalData.push(resultObj);
    }
    //console.log(finalData);
    return JSON.stringify(finalData);
}

function sendData() {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "http://localhost:3001/operation/insertIDReturn", true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(getData());
    resetForms();
    addForm();
}

//let's go
init();