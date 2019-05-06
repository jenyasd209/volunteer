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
