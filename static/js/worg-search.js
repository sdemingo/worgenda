
$(document).ready(function(){

    $("p.full-body").hide()

    $("a.list-group-item").click(function(e){
	e.preventDefault()

	if ($(this).hasClass("expand")){
	    $(this).removeClass("expand")
	    $(this).children(".full-body").hide()
	}else{
	    $(this).addClass("expand")
	    $(this).children(".full-body").show()
	}
    })

})
