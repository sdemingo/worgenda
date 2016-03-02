

$(document).ready(function(){
    var session=localStorage.getItem("sessionKey")
    if (session){
	localStorage.removeItem("sessionKey")
	setCookie("sessionKey",session,15)
	window.location.href="/"
    }
})

