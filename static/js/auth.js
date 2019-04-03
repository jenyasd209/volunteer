password.oninput = function(){
  checkPassword()
}

repeatPassword.oninput = function(){
  checkPassword()
}

function checkPassword(){
  pswd = document.forms["form-reg"]["password"]
  rptPswd = document.forms["form-reg"]["repeatpassword"]
  if (pswd.value !== rptPswd.value){
      document.forms["form-reg"]["registration-btn"].disabled = true;
      document.forms["form-reg"]["registration-btn"].className = "btn-disabled";

      document.getElementById("alert_text").style.visibility = "visible";
    }else{
      document.forms["form-reg"]["registration-btn"].disabled = false;
      document.forms["form-reg"]["registration-btn"].className = "light-blue";

      document.getElementById("alert_text").style.visibility = "hidden";
    }
}
