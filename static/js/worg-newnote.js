
$(document).ready(function(){
    $('.clockpicker').clockpicker({
	donetext: "Hecho",
	placement: "bottom",
	align:"right"
    });

    date=$("#StringDate").html()
    date=localStringDate(date)
    $("#StringDate").html(date)


    var selectedText

    // Forms buttons

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

    $("#ToDo").click(function(e){
	var title = $("#Title").val().trim()
	if ($(this).prop("checked")){
	    if (!title.startsWith("TODO")){
		title = "TODO "+title	
	    }
	}else{
	    title = title.replace(/^TODO/,"")
	}
	$("#Title").val(title)
    })

    $("#Bold").click(function(e){
	insertAtCaret("Body"," * *")
    })

    $("#Italic").click(function(e){
	insertAtCaret("Body"," / /")
    })

    $("#Picture").click(function(e){
	insertAtCaret("Body"," [[url]]")
    })

    $("#Link").click(function(e){
	insertAtCaret("Body"," [[url][texto]]")
    })

    $("#List").click(function(e){
	insertAtCaret("Body"," \n-\n-\n-")
    })


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
		$("#Title").val("")
		$("#Body").val("")
		window.location.href="/"
	    },
    	    error: function(error){
		showError("Hubo problemas. La nota puede no haberse guardado")
	    }
	})    	
    })

    $("#cancel-send-note").click(function(e){
	e.preventDefault()
	$("#Title").val("")
	$("#Body").val("")
	window.location.href="/"
    })
})
