// Get the form
let form = document.querySelector('#update-component-form');

// Get all field data from the form
// returns a FormData object
let data = new FormData(form);

// Convert to a query string
//let queryString = new URLSearchParams(data).toString();

// Convert to an object
let formObj = serialize(data);
console.log(formObj);

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
	return obj;
}

