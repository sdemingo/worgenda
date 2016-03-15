
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
    if (!datestr){
	return
    }
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


$.fn.serializeObject = function()
{
    var o = {};
    var a = this.serializeArray();
    $.each(a, function() {
        if (o[this.name] !== undefined) {
            if (!o[this.name].push) {
                o[this.name] = [o[this.name]];
            }
            o[this.name].push(this.value || '');
        } else {
            o[this.name] = this.value || '';
        }
    });
    return o;
};


// from: http://stackoverflow.com/questions/1064089/inserting-a-text-where-cursor-is-using-javascript-jquery

function insertAtCaret(areaId,text) {
    var txtarea = document.getElementById(areaId);
    var scrollPos = txtarea.scrollTop;
    var strPos = 0;
    var br = ((txtarea.selectionStart || txtarea.selectionStart == '0') ? 
              "ff" : (document.selection ? "ie" : false ) );
    if (br == "ie") { 
        txtarea.focus();
        var range = document.selection.createRange();
        range.moveStart ('character', -txtarea.value.length);
        strPos = range.text.length;
    }
    else if (br == "ff") strPos = txtarea.selectionStart;

    var front = (txtarea.value).substring(0,strPos);  
    var back = (txtarea.value).substring(strPos,txtarea.value.length); 
    txtarea.value=front+text+back;
    strPos = strPos + text.length;
    if (br == "ie") { 
        txtarea.focus();
        var range = document.selection.createRange();
        range.moveStart ('character', -txtarea.value.length);
        range.moveStart ('character', strPos);
        range.moveEnd ('character', 0);
        range.select();
    }
    else if (br == "ff") {
        txtarea.selectionStart = strPos;
        txtarea.selectionEnd = strPos;
        txtarea.focus();
    }
    txtarea.scrollTop = scrollPos;
}
