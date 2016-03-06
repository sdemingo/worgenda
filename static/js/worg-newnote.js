
$(document).ready(function(){
    $('.clockpicker').clockpicker({
	donetext: "Hecho",
	placement: "bottom",
	align:"right"
    });

    date=$("#StringDate").val()
    date=localStringDate(date)
    $("#StringDate").val(date)


    $("#send-note").click(function(e){
	e.preventDefault()
	var event = $("#new-note").serializeObject()
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
})
