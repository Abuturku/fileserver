/**
  * @author Andreas Schick (2792119), Linda Latreider (7743782), Niklas Nikisch (9364290)
  */

 /**
 * handle error message if a user tries to set her/his password to an empty string
 */
function emptyRegPw(){
	var p1 = document.getElementById("r_password").value;
	var p2 = document.getElementById("r_password2").value;
	if(p1 === "" || p2 === ""){
		alert("A password must have at least 1 character.\nPlease fill both \"password\" fields.");
		return false;
	}
}

.addEventListener("DOMContentLoaded", function(){
	 /**
	 * reading parameters that are added to the url by the server
	 */
	function getUrlParameter(paramName){
		var result = "-1",
			tmp = [];
		location
			.search.substr(1)
			.split("&")
			.forEach(function (item) {
				tmp = item.split("=");
				if (tmp[0] === paramName){
					result = decodeURIComponent(tmp[1]);
				}
		});
		return result;
	}
	
	 /**
	 * if the registration or the login attempt failed, the server adds a parameter to the url.
	 * -> the respective problem can be evaluated, so that an error message can be displayed
	 * if there is no problem, the server redirects the user to the landrive.html page.
	 */
	window.onload = function validate(){
		var message = "none";
		if(getUrlParameter("register")==="userfalse"){
			var message="Registration failed. Username already in use.";
			document.getElementById("r_username").style.backgroundColor="lightpink";
		}
		if(getUrlParameter("register")==="pwfalse"){
			var message="Registration failed. Passwords do not match.";
			document.getElementById("r_password").style.backgroundColor="lightpink";
			document.getElementById("r_password2").style.backgroundColor="lightpink";
		}
		if(getUrlParameter("login")==="false"){
			var message="Login failed. Username and/or password are/is incorrect.";
			document.getElementById("l_username").style.backgroundColor="lightpink";
			document.getElementById("l_password").style.backgroundColor="lightpink";
		}
		
		var messagefield = document.getElementById("errormessage");
		if(message!=="none"){
			messagefield.style.display='style';
			messagefield.innerHTML = message;
		} else {
			messagefield.style.display='none';
		}
	}
})