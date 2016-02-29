

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
})

function getCookie(name) {
    var value = "; " + document.cookie;
    var parts = value.split("; " + name + "=");
    if (parts.length == 2) return parts.pop().split(";").shift();
}


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

    $(".new-event").on("click",function(e){
	e.preventDefault()
	$.ajax({
    	    url:"/notes/new",
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
	title=title.replace("Monday","Lunes")
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
	$(this).html(title)
    })
	}
