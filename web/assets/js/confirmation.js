var elems = document.getElementsByClassName('confirmation')
var confirmIt = function (e) {
    if (!confirm('Вы уверены?')) e.preventDefault();
};
for (var i = 0, l = elems.length; i < l; i++) {
    elems[i].addEventListener('click', confirmIt, false);
}