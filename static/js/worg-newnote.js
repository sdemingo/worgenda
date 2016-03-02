
$(document).ready(function(){
    $('.clockpicker').clockpicker({
	donetext: "Hecho",
	placement: "bottom",
	align:"right"
    });

    date=$("#Stamp").val()
    date=localStringDate(date)
    $("#Stamp").val(date)


    $("#send-note").click(function(e){
	e.preventDefault()
	showError("Esta característica aún no ha sido implementada")
    })

    
})
