$(document).ready(function() {
  var url = window.location.href;
  if (url.indexOf('?showmodal=1') != -1) {
    $("#modal-1").modal('show');
  }
  if (url.indexOf('?showmodal=2') != -1) {
    $("#modal-2").modal('show');
  }
});