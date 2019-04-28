function showMobileMenu() {
  var x = document.getElementById("user-menu");
  if (x.style.display === "block") {
    x.style.display = "none";
  } else {
    x.style.display = "block";
  }
}

// var selectedItem;

// function redirect(elem){
//   // setSelected(elem);
//   document.addEventListener("DOMContentLoaded", setSelected(elem));
//   document.location.href = elem.href;
// }

function setSelected(elem){
  selectedItem = elem;
  // document.addEventListener("DOMContentLoaded", setSelected(elem));
  console.log(document.location.href);
  var a = document.getElementsByTagName('a');
  for (i = 0; i < a.length; i++) {
      a[i].classList.remove('selected-item')
  }
  elem.classList.add('selected-item');
}
