
function initAccept() {
    resetFormsAccept();
    createFormAccept();
    
    document.addEventListener(
        "DOMContentLoaded",
        function () {
            document.getElementById("add-conditionAccept").onclick = addFormButtonAccept; //addForm
            document.getElementById("clear-formAccept").onclick = initAccept;
            document.getElementById("applyAccept").onclick = sendDataAccept;
        },
        false
    );
}
//global variables
//
//

var mainDivAccept = document.getElementById("formAccept"); //ref to main div with forms

//attributes for form fields, they will be added during the creation of forms
var formElementsTypesAccept = ["input", "button"];//https://www.jetbrains.com/idea/features/editions_comparison_matrix.html
var formElementsIdsAccept = ["scanidAccept", "buttonaccept"];
var formElementsClassNamesAccept = [
    "form-control form-control-sm col-lg-4",
    "deleteForm btn",
];

var formAllowedIdsArrAccept = [];
var rowsAmountAccept = 50;
var finalDataAccept = [];
var labelsAccept = ["Сканируем QR-code катушки"];
const regexpAccept = /\bP\d{7}LK\d{9}R\d{10}Q\d{5}D\d{8}\b/;
const regexp2Accept = /\bP\d{7}L\d{10}R\d{10}Q\d{5}D\d{8}\b/;
// P1016624L1000037226R1000317938Q00550D20200311

//creating form and adding attributes
function createFormAccept() {
    var form = document.createElement("div");
    form.className = "form-inline";
    form.id = formAllowedIdsArrAccept[0];
    formAllowedIdsArrAccept = formAllowedIdsArrAccept.slice(1);
    mainDivAccept.appendChild(form);
    for (var i = 0; i < formElementsTypesAccept.length; i++) {
        
        var element = document.createElement(formElementsTypesAccept[i]);
        element.id = `${formElementsIdsAccept[i]}-${form.id}`;
        element.className = formElementsClassNamesAccept[i].concat(" ", "fields-style");
        if (i != 1) {
            element.placeholder = labelsAccept[i];
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

    var input = document.getElementById(`${formElementsIdsAccept[0]}-${form.id}`);
    var deleteCross = document.getElementById(`${formElementsIdsAccept[1]}-${form.id}`);

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
        addRemoveErrorAccept(this, errorMsg);

        toggleSubmitBtnAccept();
    };

    // check for an errors if user push enter and add new input
    input.addEventListener("keyup", function (event) {
        if (event.keyCode === 13 || 40) {
            event.preventDefault();
            errorMsg = createError("Введён некорректный номер");
            
            addRemoveErrorAccept(this, errorMsg);
            if (validateQrCode(this.value)) {
                addFormAccept();
            }
        }

        toggleSubmitBtnAccept();
    });

    deleteCross.addEventListener("click", deleteFormAccept);

    canShowAddButtonAccept(); //check if we can show add button
    canShowDeleteButtonAccept(); //check if we can show delete button
}

function toggleSubmitBtnAccept() {
    // кнопка не активна, если 
    document.getElementById("applyAccept").disabled = !checkAllValuesAccept();  
}

//check if we can add form and add it
function addFormAccept() {
    if (formAllowedIdsArrAccept.length > 0) {
        createFormAccept();
    }
}

// count similar strings entered in all input fields
function checkValueAccept(inputValue) {
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
function checkAllValuesAccept() {
    var result = true;
    var forms = document.getElementsByClassName("form-control");
    for (i = 0; i < forms.length; i++) {
        if (forms[i].className.includes("is-invalid") && forms[i].value != '') {
            result = false;
            break;
        }
    }
    return result;
}

// add new input if we didn't have duplicates
function addFormButtonAccept() {
    var check = checkAllValuesAccept();
    if (check) {
        addFormAccept();
    }
}

// check input value on context
// function checkCanAddValueOnContextAccept(e) {
//     var numberToCheck = checkValueAccept(e.value);
//     var result = true;

//     var table = document.getElementById('thetable');
//     let rows = table.querySelectorAll('tr');
//     let check = e.value;
//     let check2 = check.toUpperCase();

//     for (var i = 0; i < rows.length; i++) {
//         var fullname = rows[i].querySelectorAll('td');
//         fullname = fullname[0].innerHTML.toUpperCase();
//         let status = rows[i].querySelectorAll('td');
//         status = status[9].innerHTML.toUpperCase();
    
//             if (
//                 !(regexpAccept.test(e.value) || regexp2Accept.test(e.value) ) || check2 != fullname || numberToCheck > 1
//             ) {
//                     console.log(check2)
//                     console.log(status)
//                     if (status !== 'OK' || status !== 'NG') {
        
//                     result = false;
//                     } else {
//                         console.log("result true ", check2)
//                         console.log("result true status " ,status)
//                         result = true;
//                     }
//             }

//         return result;
//     }
// }

function validateQrCode(code) {
 
    const regexpContractEO = /\bP\d{7}LK\d{9}R\d{10}Q\d{5}D\d{8}\b/; 
    const regexpStarLineEO = /\bP\d{7}L\d{10}R\d{10}Q\d{5}D\d{8}\b/;
    const regexpCart = /\bP\d{7}LR\d{10}Q\d{5}D\d{8}\b/;
    const column = document.querySelector('td[data-material-id="'+code+'"]');
    let result = true;

    if (!column || !column.innerHTML) {
        return false;
    }

    const row = document.querySelector('td[data-material-id="'+code+'"]').parentElement; 
    const stateIsValid = row.classList.contains('NG') || row.classList.contains('OK');

  //  if (rowid.includes(code)) {
    if (
        (
            regexpStarLineEO.test(code) ||
            regexpContractEO.test(code) ||
            regexpCart.test(code)
        ) && (stateIsValid)
    ) {  
        result = true;
    } else {
        result = false;
    }  

    return result;
}

// add or remove error message
function addRemoveErrorAccept(context, errorMsg) {
    if (!validateQrCode(context.value)) {
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

function deleteFormAccept() {
    formAllowedIdsArrAccept.push(this.parentNode.id);
    this.parentNode.parentNode.removeChild(this.parentNode);
    canShowAddButtonAccept();
    canShowDeleteButtonAccept();
    toggleSubmitBtnAccept();
}

// generate array of acceptable form ids
function formIdGeneratorAccept(num) {
    var formIdArray = [];
    for (i = 0; i < num; i++) {
        formIdArray.push(`form-${i + 1}`);
    }
    return formIdArray;
}

//reset all forms, fields, console and generate new allowed IDs array
function resetFormsAccept() {
    formAllowedIdsArrAccept = formIdGeneratorAccept(rowsAmountAccept);
    while (mainDivAccept.firstChild) {
        mainDivAccept.removeChild(mainDivAccept.firstChild);
    }
    finalDataAccept = [];
    //console.clear();
}

function canShowAddButtonAccept() {
    if (formAllowedIdsArrAccept.length == 0) {
        document.getElementById("add-conditionAccept").disabled = true;
    } else {
        document.getElementById("add-conditionAccept").disabled = false;
    }
}

function canShowDeleteButtonAccept() {
    if (formAllowedIdsArrAccept.length < rowsAmountAccept - 1) {
        document.getElementsByClassName("deleteForm")[0].disabled = false;
    } else {
        document.getElementsByClassName("deleteForm")[0].disabled = true;
    }
}

function getDataAccept() {
    var resultData = [];
    var forms = document.getElementsByClassName("form-control");
    for (i = 0; i < forms.length; i++) {
        if (forms[i].value) {
            resultData.push({ scanidAccept: forms[i].value });
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

function sendDataAccept() {
    for (el of document.getElementsByClassName('is-invalid')) {
        el.parentNode.remove();
    }

 //   if (checkCanAddValueOnContext(context) == true) {
    dataToSend = getDataAccept();
    resetFormsAccept();
    createFormAccept();
    if (dataToSend != "[]") {
        var xhr = new XMLHttpRequest();
        //   xhr.open("POST", "http://10.1.20.110:3001/ininspection", true);
        xhr.open("POST", "http://localhost:3001/operation/acceptgroupsinspectiontowh", true);
      //  xhr.open("POST", "http://", true);
        xhr.setRequestHeader("Content-Type", "application/json");

        xhr.send(dataToSend);
        console.log(dataToSend);
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
//}

//let's go
initAccept();

