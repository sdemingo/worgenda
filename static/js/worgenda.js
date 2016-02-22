

var W={}

$(document).ready(function(){

    W.Events=[]
    
    $("#datepicker").datepicker({
	dayNames: [ "Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado" ],
	dayNamesMin: [ "Do", "Lu", "Ma", "Mi", "Ju", "Vi", "Sa" ],
	dayNamesShort: [ "Dom", "Lun", "Mar", "Mier", "Jue", "Vie", "Sab" ],
	monthNames: [ "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", 
		      "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre" ],
	firstDay: 1
    })

    //getEvents()
})


function getEvents(){

    $.ajax({
    	url:"/events",
    	type: 'get',
	dataType: 'json',
    	success: function (events){
	    console.log(events)
	},
    	error: function(error){
	    console.log(error)
	}
    })
}
