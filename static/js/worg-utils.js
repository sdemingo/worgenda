
function getCookie(name) {
    var value = "; " + document.cookie;
    var parts = value.split("; " + name + "=");
    if (parts.length == 2) return parts.pop().split(";").shift();
}

function setCookie(cname, cvalue, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + (exdays*24*60*60*1000));
    var expires = "expires="+d.toUTCString();
    document.cookie = cname + "=" + cvalue + "; " + expires;
} 

function storeSession(){
    var cookie=getCookie("sessionKey")
    localStorage.setItem("sessionKey",cookie)
}

function localStringDate(datestr){
    datestr=datestr.replace("Monday","Lunes")
	.replace("Tuesday","Martes")
	.replace("Wednesday","Miércoles")
	.replace("Thursday","Jueves")
	.replace("Friday","Viernes")
	.replace("Saturday","Sábado")
	.replace("Sunday","Domingo")

	.replace("January","Enero")
	.replace("February","Febrero")
	.replace("March","Marzo")
	.replace("April","Abril")
	.replace("May","Mayo")
	.replace("June","Junio")
	.replace("July","Julio")
	.replace("Agoust","Agosto")
	.replace("September","Septiembre")
	.replace("October","Octubre")
	.replace("November","Noviembre")
	.replace("December","Diciembre")

    return datestr
}


function showError(msg){
    $("body").prepend('<div class="alert alert-danger alert-dismissible error-msg" role="alert"><button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button><strong>Error:</strong> '+msg+'</div>')
    window.setTimeout(function() { $(".alert").alert('close'); }, 2000);

		     }
