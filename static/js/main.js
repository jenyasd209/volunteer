var main = function(){
  "use strict";
  // $('.parallax').parallax();
  // $('.sidenav').sidenav();
  $(".dropdown-trigger").dropdown();
}

$(document).ready(main);

document.addEventListener('DOMContentLoaded', function() {
    var mainBanner = document.querySelectorAll('.parallax');
    var instances = M.Parallax.init(mainBanner, {});
    var sidebar = document.querySelectorAll('.sidenav');
    var instances_sidebar = M.Sidenav.init(sidebar, {});
});

function showMobileMenu() {
  var x = document.getElementById("user-menu");
  if (x.style.display === "block") {
    x.style.display = "none";
  } else {
    x.style.display = "block";
  }
}

function setSelected(elem){
  var a = document.getElementsByTagName('a');
  for (i = 0; i < a.length; i++) {
      a[i].classList.remove('selected-item')
  }
  elem.classList.add('selected-item');
  console.log("output");
}
