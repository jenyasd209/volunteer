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

function renderStars() {
  let result = "";
  for (let i = 1.0; i <= 5.0; i++) {
    if (rating >= i) {
      result += `<span class="fa fa-star checked"></span>`;
    }else if (rating < i && rating > i - 1.0) {
      result += `<span class="fa fa-star-half-empty checked"></span>`;
    }else if (rating < i) {
      result += `<span class="fa fa-star"></span>`;
    }
  }
  return result + `<span>(${rating})</span>`
}

document.getElementById("rating").innerHTML = renderStars(4);