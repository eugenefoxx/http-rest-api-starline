let url = 'http://localhost:3001/operation/updateinspection/?id';
let response = await fetch(url);

let commits = await response.json(); // читаем ответ в формате JSON

console.log(commits);