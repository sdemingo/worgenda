

var W={
    datesMarked:[],
    datepickerFormatLayout : "dd M yy"
}

$(document).ready(function(){

    W.Events=[]
    
    $("#datepicker").datepicker({
	dayNames: [ "Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado" ],
	dayNamesMin: [ "Do", "Lu", "Ma", "Mi", "Ju", "Vi", "Sa" ],
	dayNamesShort: [ "Dom", "Lun", "Mar", "Mier", "Jue", "Vie", "Sab" ],
	monthNames: [ "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", 
		      "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre" ],
	firstDay: 1,
	dateFormat:W.datepickerFormatLayout,
	beforeShowDay: markDates,
	onSelect: getEventsForADate
    })

    getMarkDates()
    var today=$.datepicker.formatDate(W.datepickerFormatLayout, new Date()) 
    getEventsForADate(today)
})


function getMarkDates(){

    $.ajax({
    	url:"/notes/dates",
    	type: 'get',
	dataType: 'json',
    	success: function (dates){
	    for (var i=0;i<dates.length;i++){
		W.datesMarked.push(new Date(dates[i]))
	    }

	    console.log(W.datesMarked)
	    $("#datepicker").datepicker( "refresh" )
	},
    	error: function(error){
	    console.log(error)
	}
    })
}

function getEventsForADate(sdate){
  
    $.ajax({
    	url:"/notes/events",
    	type: 'post',
	data: sdate,
    	success: function (html){
	    loadHTML(html)
	},
    	error: function(error){
	    console.log(error)
	}
    })    
}


function markDates(date) {
    var dates=W.datesMarked
    for (var i = 0; i < dates.length; i++) {
	if (dates[i].getTime() == date.getTime()){
	    return [true, "dp-marked-date"];
        }
    }
    return [true, ""];
} 


function loadHTML(html){
    $("#content").html(html)
}
