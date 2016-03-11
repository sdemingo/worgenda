
$(document).ready(function(){
    $('.clockpicker').clockpicker({
	donetext: "Hecho",
	placement: "bottom",
	align:"right"
    });

    date=$("#StringDate").html()
    date=localStringDate(date)
    $("#StringDate").html(date)


    $("#send-note").click(function(e){
	e.preventDefault()
	var event = $("#new-note").serializeObject()
	if (event.Title.trim() == ""){
	    showError("La nota debe tener un título")
	    return
	}
	$.ajax({
    	    url:"/notes/add",
    	    type: "post",
	    dataType:"json",
	    data: JSON.stringify(event),
    	    success: function (html){
		window.location.href="/"
	    },
    	    error: function(error){
		showError("Hubo problemas. La nota puede no haberse guardado")
	    }
	})    	
    })

    $("#cancel-send-note").click(function(e){
	e.preventDefault()
	window.location.href="/"
    })

    $("#AllDay").click(function(e){
    	if ($(this).prop("checked")){
    	    $("#Hour").prop("readonly",true)
	    $("#Hour").val("00:00")
	    $(".clockpicker").addClass("no-events")
    	}else{
    	    $("#Hour").prop("readonly",false)
	    $("#Hour").val("09:00")
	    $(".clockpicker").removeClass("no-events")
    	}
    })
})
