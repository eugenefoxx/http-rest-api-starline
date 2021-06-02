
let form = document.getElementById('update-component-form');
function serialize (form) {
  //  let checkForm = document.getElementById('update-component-form');
  //  console.log(checkForm);
    if (!form || form.nodeName !== "FORM" /*&& !form.getElementById('update-component-form')*/ ){
        return false;
    }
/* /    let form = document.getElementById('update-component-form')
    if (document.querySelectorAll('update-component-form')){
        return false;
    }*/
    let i,j, q = []
    for (i = form.elements.length -1; i >= 0; i = i -1) {
        if (form.elements[i].name === "") {
            continue;
        }
        debugger;
        switch (form.elements[i].nodeName) {
            case 'INPUT':
                switch (form.elements[i].type) {
                    case 'text':
                    case 'tel':
                    case 'email':
                    case 'hidden':
                    case 'password':
                    case 'button':
                    case 'reset':
                    case 'submit':
                        q.push(form.elements[i].name + "=" + encodeURIComponent(form.elements[i].value));
                        break;
                    case 'checkbox':
                    case 'radio':
                            if (form.elements[i].checked) {
                               q.push(form.elements[i].name + "=" + encodeURIComponent(form.elements[i].value));
                            } 
                           break;                   
                }
                break;
            case 'file':
                   break;
            case 'TEXTAREA':
                    q.push(form.elements[i].name + "=" + encodeURIComponent(form.elements[i].value));
                    break;
            case 'SELECT':
                switch (form.elements[i].type) {
                    case 'select-one':
                        q.push(form.elements[i].name + "=" + encodeURIComponent(form.elements[i].value));
                        break;
                    case 'select-multiple':
                        for (j = form.elements[i].options.length - 1; j >= 0; j = j - 1) {
                            if (form.elements[i].options[j].selected) {
                                q.push(form.elements[i].name + "=" + encodeURIComponent(form.elements[i].value));
                                                                      }
                        } 
                        break;
                    case 'BUTTON':
                        switch (form.elements[i].type) {
                            case 'reset':
                            case 'submit':
                            case 'button':
                                q.push(form.elements[i].name + "=" + encodeURIComponent(form.elements[i].value));
                                break;
                        }
                        break;                       
                }
        }
        return q.join("&");
    }
}