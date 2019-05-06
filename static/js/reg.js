window.onload = function() {
  password.oninput = function(){
    checkPassword()
  }

  repeatPassword.oninput = function(){
    checkPassword()
  }
  // displayBlock();
}

function checkPassword(){
  pswd = document.forms["form-reg"]["password"]
  rptPswd = document.forms["form-reg"]["repeatPassword"]
  if (pswd.value !== rptPswd.value){
      document.forms["form-reg"]["registration-btn"].disabled = true;
      document.forms["form-reg"]["registration-btn"].className = "btn-disabled";

      document.getElementById("alert_text").classList.remove("hide");
    }else{
      document.forms["form-reg"]["registration-btn"].disabled = false;
      document.forms["form-reg"]["registration-btn"].className = "light-blue";

      document.getElementById("alert_text").classList.add("hide");
    }
}

function selectGroupCheck(){
  let group = document.getElementsByName('group');
  let freelancer = document.getElementById('freelancer-block');
  let customer = document.getElementById('customer-block');
  if (group[0].checked) {
    freelancer.classList.remove("hide");
    customer.classList.add("hide");
    document.getElementById("specializations").required = true;
    document.getElementById("organization-name").required = false;
  }if (group[1].checked){
    freelancer.classList.add("hide");
    customer.classList.remove("hide");
    document.getElementById("organization-name").required = true;
    document.getElementById("specializations").required = false;
  }
  // elem1.classList.add('selected-group');
  // elem2.classList.remove('selected-group');
  // displayBlock();
}

// function displayBlock() {
//   var inp = document.getElementsByName('group');
//   var freelancer = document.getElementById('freelancer-block');
//   var customer = document.getElementById('customer-block');
//   if (inp[0].checked) {
//     freelancer.style.display = "block";
//     customer.style.display = "none";
//   }if (inp[1].checked){
//     freelancer.style.display = "none";
//     customer.style.display = "block";
//   }
// }
