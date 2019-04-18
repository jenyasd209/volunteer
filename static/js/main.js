var main = function(){
  "use strict";
  // $('.parallax').parallax();
  // $('.sidenav').sidenav();
  $(".dropdown-trigger").dropdown();
}

$(document).ready(main);

var nav = document.getElementById('nav');
var divNav = document.getElementById('div-nav');
var navSourceBottom = divNav.getBoundingClientRect().bottom + window.pageYOffset;

window.onscroll = function() {
  if (divNav.classList.contains('navbar-fixed') && window.pageYOffset < navSourceBottom) {
    divNav.classList.add('navbar');
    divNav.classList.remove('navbar-fixed');
    nav.classList.remove('small-nav');
  } else if (window.pageYOffset > navSourceBottom) {
    divNav.classList.remove('navbar');
    divNav.classList.add('navbar-fixed');
    nav.classList.add('small-nav');
  }
};

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
