function init() {
    resetForms();
    createForm();

    document.addEventListener(
        "DOMContentLoaded",
        function () {
            document.getElementById("add-condition").onclick = addFormButton; //addForm
            document.getElementById("clear-form").onclick = init;
            document.getElementById("apply").onclick = sendData;
        },
        false
    );
}
//global variables
//
//

var mainDiv = document.getElementById("form"); //ref to main div with forms

//attributes for form fields, they will be added during the creation of forms
var formElementsTypes = ["input", "button"];//https://www.jetbrains.com/idea/features/editions_comparison_matrix.html
var formElementsIds = ["scanid", "button"];
var formElementsClassNames = [
    "form-control form-control-sm col-lg-4",
    "deleteForm btn",
];

var formAllowedIdsArr = [];
var rowsAmount = 50;
var finalData = [];
var labels = ["Сканируем ЕО катушки"];
const regexp = /\bP\d{7}LK\d{9}R\d{10}Q\d{5}D\d{8}\b/;
const regexp2 = /\bP\d{7}L\d{10}R\d{10}Q\d{5}D\d{8}\b/;

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
        if (i != 1) {
            element.placeholder = labels[i];
        }
        form.appendChild(element);

        if (i == 1) {
            var icon = document.createElement("i");
            var deleteButton = document.getElementById(element.id);
            deleteButton.appendChild(icon);
            icon.className = "fa fa-times";
        }
    }
    form.childNodes[0].focus();

    //add eventlisteners to newly created elements

    var input = document.getElementById(`${formElementsIds[0]}-${form.id}`);
    var deleteCross = document.getElementById(`${formElementsIds[1]}-${form.id}`);

    // create error if incorrect data were entered
    function createError(errorText) {
        var errorMessage = document.createElement("div");
        errorMessage.className = "col-sm-3";
        var message = document.createElement("small");
        message.id = "error-help";
        message.className = "text-danger";
        message.innerHTML = errorText;
        errorMessage.appendChild(message);
        return errorMessage;
    }

    // check for an errors if user change activity to another element
    input.onblur = function () {
        errorMsg = createError("Введён некорректный номер");
        addRemoveError(this, errorMsg);
    };

    // check for an errors if user push enter and add new input
    input.addEventListener("keyup", function (event) {
        if (event.keyCode === 13) {
            event.preventDefault();
            errorMsg = createError("Введён некорректный номер");
            addRemoveError(this, errorMsg);
            if (checkCanAddValueOnContext(this)) {
                addForm();
            }
        }
    });

    deleteCross.addEventListener("click", deleteForm);

    canShowAddButton(); //check if we can show add button
    canShowDeleteButton(); //check if we can show delete button
}

//check if we can add form and add it
function addForm() {
    if (formAllowedIdsArr.length > 0) {
        createForm();
    }
}

// count similar strings entered in all input fields
function checkValue(inputValue) {
    var forms = document.getElementsByClassName("form-control");
    var counter = 0;
    for (i = 0; i < forms.length; i++) {
        if (inputValue == forms[i].value) {
            counter += 1;
        }
    }

    return counter;
}

// check if we have any invalid string
function checkAllValues() {
    var result = true;
    var forms = document.getElementsByClassName("form-control");
    for (i = 0; i < forms.length; i++) {
        if (forms[i].className.includes("is-invalid")) {
            result = false;
            break;
        }
    }
    return result;
}

// add new input if we didn't have duplicates
function addFormButton() {
    var check = checkAllValues();
    if (check) {
        addForm();
    }
}

// check input value on context
function checkCanAddValueOnContext(e) {
    var numberToCheck = checkValue(e.value);
    var result = true;
    if (!regexp.test(e.value) || numberToCheck > 1) {
        result = false;
    } else {
        result = true;
    }

    return result;
}

// add or remove error message
function addRemoveError(context, errorMsg) {
    if (!checkCanAddValueOnContext(context)) {
        if (!context.className.includes("is-invalid")) {
            context.className = context.className.concat(" is-invalid");
        }
        if (context.parentNode.children.length == 2) {
            context.parentNode.appendChild(errorMsg);
        }
    } else {
        if (context.parentNode.lastChild.className == "col-sm-3") {
            context.parentNode.removeChild(context.parentNode.lastChild);
        }
        context.className = context.className.replace(" is-invalid", "");
    }
}

function deleteForm() {
    formAllowedIdsArr.push(this.parentNode.id);
    this.parentNode.parentNode.removeChild(this.parentNode);
    canShowAddButton();
    canShowDeleteButton();
}

// generate array of acceptable form ids
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

function getData() {
    var resultData = [];
    var forms = document.getElementsByClassName("form-control");
    for (i = 0; i < forms.length; i++) {
        if (forms[i].value) {
            resultData.push({ scanid: forms[i].value });
        }
    }
    return JSON.stringify(resultData);
}

function drawErrorMessage(div) {
    var errorMessage = document.createElement("div");
    errorMessage.id = "error-message";
    errorMessage.innerHTML = "Ошибка при передаче данных";
    //mainDiv.appendChild(errorMessage);
    //div.appendChild(errorMessage);
    return errorMessage;
}

function drawSuccessMessage(div) {
    var successMessage = document.createElement("div");
    successMessage.id = "success-message";
    successMessage.innerHTML = "Данные переданы успешно";
    //div.appendChild(successMessage);
    return successMessage;
}

function sendData() {
    dataToSend = getData();
    resetForms();
    createForm();
    if (dataToSend != "[]") {
        var xhr = new XMLHttpRequest();
        //   xhr.open("POST", "http://10.1.20.110:3001/ininspection", true);
        xhr.open("POST", "http://localhost:3001/ininspection", true);
        xhr.setRequestHeader("Content-Type", "application/json");

        xhr.send(dataToSend);
        xhr.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                //mainDiv.innerHTML = drawSuccessMessage();
                console.log("success");
            } else {
                //mainDiv.innerHTML = drawErrorMessage();
                console.log("error");
            }
        };
    } else {
        console.log("There is no data to send");
    }
}

//let's go
init();
