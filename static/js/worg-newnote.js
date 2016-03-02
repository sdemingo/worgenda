
$(document).ready(function(){
    $('.clockpicker').clockpicker({
	donetext: "Hecho",
	placement: "bottom",
	align:"right"
    });

    date=$("#Stamp").val()
    date=localStringDate(date)
    $("#Stamp").val(date)
})
