

var W={
    datesMarked:[],
    datepickerFormatLayout : "dd M yy",
    currentDate:new Date()
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

    storeSession()
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

	    $("#datepicker").datepicker( "refresh" )
	},
    	error: function(error){
	    window.location.href="/"
	}
    })
}

function getEventsForADate(sdate){
    W.currentDate=sdate

    $.ajax({
    	url:"/notes/events",
    	type: 'post',
	data: "date="+sdate,
    	success: function (html){
	    loadHTML(html)
	},
    	error: function(error){
	    window.location.href="/"
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

    localNames()

    //default events
    $(".close-event").on("click",function(e){
	e.preventDefault()
	getEventsForADate(W.currentDate)	
    })

    $(".new-note").on("click",function(e){
	e.preventDefault()
	$.ajax({
    	    url:"/notes/new",
    	    type: 'post',
    	    success: function (html){
		loadHTML(html)
	    },
    	    error: function(error){
		window.location.href="/"
	    }
	})  	
    })

    $("#day-events .list-group a").on("click",function(e){
	e.preventDefault()
	var href=$(this).attr("href")
	$.ajax({
    	    url:href,
    	    type: 'post',
	    data: "date="+W.currentDate,
    	    success: function (html){
		loadHTML(html)
	    },
    	    error: function(error){
		window.location.href="/"
	    }
	})  	
    })
}




function localNames(){
    $(".panel-title").each(function(i){
	var title=$(this).html()
	title=localStringDate(title)
	$(this).html(title)
    })
	}
