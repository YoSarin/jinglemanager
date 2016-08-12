var scripts = document.getElementsByTagName("script");
var __FILE__ = scripts[scripts.length-1].src;
var pointOrder = {
    "match_start" : 0,
    "match_end"   : 1,
    "match_none"  : 2,
}
var pointerTimeout = null;


$(document).ready(function() {
    healthcheck();
    hook();
    load();
    showHideJingleMatchDetails();
});

$(window).resize(function() {
    slotDisplay();
});


var Handler = {
    "song_changed"   : songChange,
    "app_added"      : appAdd,
    "app_removed"    : appRemove,
    "volume_changed" : volumeChange,
    "cleanup"        : load,
    "jingle_added"   : jingleAdd,
    "jingle_changed" : jingleChange,
    "jingle_removed" : jingleRemove,
    "log"            : log,
}

var slots = new Timescale();
var resolution = 1; // pixels per minute

connectSocket("changes");
connectSocket("logs");

function connectSocket(name) {
    var socket;
    socket = new WebSocket("ws://" + window.location.hostname + ":8080/" + name);
    socket.onopen = function(evt) {
        log({
            Severity: "info",
            Time: new Date(),
            Message: "WebSocket " + name + " connected",
            File: __FILE__
        })
    }
    socket.onmessage = function(evt) {
        try {
            var data = $.parseJSON(evt.data);
            if (typeof Handler[data.Type] == 'undefined') {
                console.warn(data.Type, "Not implemented");
                return;
            }
            Handler[data.Type](data.Data);
        } catch (e) {
            console.error(e);
            console.info(evt.data);
        }
    }
    socket.onclose = function(evt) {
        setTimeout(function () {
            connectSocket(name);
        }, 5000);
    }
}

function healthcheck() {
    $.ajax("/alive?", {
        error: function() {
            alert("Máchale, spadlo ti to...\n:(\nSpusť znovu jinglemanager.exe");
            setTimeout(healthcheck, 5000);
        },
        success: function() {
            setTimeout(healthcheck, 5000);
        }
    });
}

function hook() {
    $("form.ajax").submit(function () {
        var f = $(this);
        $.ajax(f.attr("action"), {
            method: f.attr("method"),
            data: f.serialize(),
            success: function(data, status) {  },
            complete: function() {  },
            error: function() { console.log("error"); }
        })
        return false;
    });

    $("#jingles").sortable({
        stop: function(event, ui) {
            // event.toElement is the element that was responsible
            // for triggering this event. The handle, in case of a draggable.
            $( event.originalEvent.target ).one('click', function(e){ e.stopImmediatePropagation(); } );
        }
    }).disableSelection()

    $("#addJingle select[name=play]").change(function() {
        showHideJingleMatchDetails();
    });

    $(".refreshSlots").unbind("click", clicker);
    $(".refreshSlots").click(
        {
        callback: function(data) {
                clearSlots();
                listSlots(data);
            }
        }, clicker
    );
    $(".advanced").addClass("hidden");
    $(".toggleAdvanced").unbind("click", clicker);
    $(".toggleAdvanced").click({callback: function() { $(".advanced").toggleClass("hidden"); }}, clicker)
}

function load() {
    $.ajax("/jingle/list", {
        success: function(data, status) { listJingles(data); }
    });
    $.ajax("/track/list", {
        success: function(data, status) { listTracks(data); }
    });
    $.ajax("/app/list", {
        success: function(data, status) { listApps(data); }
    });
    $.ajax("/slot/list", {
        success: function(data, status) { listSlots(data); }
    });
}

function showHideJingleMatchDetails() {
    var v = $("#addJingle select[name=play]").val()
    if (v == "match_related") {
        $(".matchOnly").removeClass("hidden");
    } else {
        $(".matchOnly").addClass("hidden");
    }
}

function listTracks(data) {
    try {
        var data = $.parseJSON(data);
        $("#songs").empty();
        $.each(data, function (k, v) {
            songChange(v);
        });
    } catch (e) {
        console.log(e);
    }
}

function listJingles(data) {
    try {
        var data = $.parseJSON(data);
        $("#jingles").empty();
        $.each(data, function (k, v) {
            if($("#jingle-" + v.Song.ID).length == 0) {
                jingleAdd(v);
            }
            jingleChange(v);
        });
    } catch (e) {
        console.log(e);
    }
}

function listApps(data) {
    try {
        var data = $.parseJSON(data);
        $("#apps").empty();
        $.each(data, function (k, v) {
            if($("#app-" + v.ID).length == 0) {
                appAdd(v)
            }
            appChange();
        });
    } catch (e) {
        console.log(e);
    }
}

function clearSlots() {
    slots = new Timescale();
}

function listSlots(data) {
    try {
        var data = $.parseJSON(data);
        $.each(data, function (k, v) {
            slotAdd(v)
        });
    } catch (e) {
        console.log(e);
    }
}

function slotAdd(slot) {
    slots.add(new Slot(slot.StartsAt, slot.Duration));
    slotDisplay();
}

function slotDisplay() {
    $("#slots").empty();
    var totalDuration = slots.duration();
    $("#slots").append($("<div>").addClass("pointer").append('&dtrif;'));
    resolution = parseInt($("#slots").width()) / slots.duration();
    // console.log(resolution, $("#slots").width(), slots.duration());
    movePointer();
    $.each(slots.displayList(), function(k, v) {
        var item = $('<div class="slot upcoming">')
            .css("width", (resolution * v.duration) + "px")
            .css("border-left-width", (resolution * v.gapBefore) + "px");

        if ((v.duration * resolution) >= 60) {
            item.append($("<div>").addClass("inner")
                .append("&#8614;" + formatSlotDate(v.start) + "<br />&rarrb;" + formatSlotDate(v.end))
                .append($("<br />"))
                .append($("<small>")
                    .append(v.duration + " minut")
                )
            );
        }

        window.setTimeout(function () {
            item.addClass("current");
            item.removeClass("upcoming");
        }, v.start - Date.now());

        window.setTimeout(function () {
            item.removeClass("current");
            item.addClass("done");
        }, v.end - Date.now());

        $("#slots").append(item);
    });

    metroResize();
}

function movePointer() {
    // just to make sure it won't run multiple times at once
    window.clearTimeout(pointerTimeout)
    var elapsed = (Date.now() - slots.start()) / 1000 / 60;
    var missed = (Date.now() - slots.end()) / 1000 / 60;
    var position = 0;
    if (elapsed < 0) {
        var h = lpad(Math.floor(Math.abs(elapsed / 60)), 2, "0");
        var m = lpad(Math.floor(Math.abs(elapsed % 60)), 2, "0");
        var s = lpad(Math.floor(Math.abs((elapsed * 60) % 60)), 2, "0")
        $("#slots .pointer").html('&#9666; ' + h + ':' + m + ':' + s);
        position = 0;
    } else if (elapsed > (slots.end() - slots.start())/1000/60) {
        var h = lpad(Math.floor(Math.abs(missed / 60)), 2, "0");
        var m = lpad(Math.floor(Math.abs(missed % 60)), 2, "0");
        var s = lpad(Math.floor(Math.abs((missed * 60) % 60)), 2, "0")
        $("#slots .pointer").html(h + ':' + m + ':' + s + ' &#9656;');
        position = $(document).width() - ($("#slots .pointer").width() + 20);
    } else {
        $("#slots .pointer").html('&#9662;');
        position = elapsed * resolution - Math.ceil($("#slots .pointer").width()/2.0);
    }
    $("#slots .pointer").css("left", position + "px")

    pointerTimeout = window.setTimeout(movePointer, 100 * 1);
}

function formatSlotDate(v) {
    return lpad(v.getHours(), 2, "0") + ":" + lpad(v.getMinutes(), 2, "0");
}

function appAdd(app) {
    if($("#app-" + app.ID).length > 0) {
        return;
    }
    $("#apps")
        .append(
            $('<li id="app-' + app.ID + '" class="app">')
            .append($('<div class="outer">')
                .append($('<div class="content">')
                    .append($("<strong>").text(app.Name)).append(" ")
                    .append($('<small class="state">').text("[" + Math.round(100 * app.Volume) + "%]")).append("<br />")
                    .append($('<a class="control" method="delete" href="/app/delete/' + app.ID + '">delete</a>'))
                )
            )
        ).find("a").each(function (k, v) {
            $(v).prop('onclick', null).off('click');
            $(v).click(clicker);
        });
    metroResize();
}

function appRemove(app) {
    $("#app-" + app.ID).remove()
    metroResize();
}

function appChange(app) {

}

function volumeChange(app) {
    $("#app-" + app.ID).each(function () {
        $(this).find(".state").text( Math.round(app.Volume * 100) + "%");
    });
}

function jingleAdd(jingle) {
    if($("#jingle-" + jingle.Song.ID).length > 0) {
        return;
    }
    var newJingle = $('<li id="jingle-' + jingle.Song.ID + '" class="song" time="' + jingle.TimeBeforePoint + '" point="' + jingle.Point + '">')
        .append($('<div class="outer">')
            .append($('<div class="progress">').css("width", Math.round(jingle.Song.Position * 100) + "%"))
            .append($('<div class="content">')
                .append($("<strong>").text(jingle.Name)).append('<br />')
                .append($('<small class="songtitle">').text(jingle.Song.File))
            )
        ).each(function (k, v) {
            $(v).prop('onclick', null).off('click');
            $(v).multi_click(jinglePlayPause, jingleStop, jingleDelete, 500);
        });

    $("#jingles").append(newJingle);

    $("#jingles li").each(function (k, v) {
        if (compareJingles(newJingle, $(v)) == -1) {
            newJingle.insertBefore(v);
            return false;
        }
    });

    metroResize();
}

function compareJingles(a, b) {
    if (pointOrder[a.attr("point")] > pointOrder[b.attr("point")]) {
        return 1;
    } else if (pointOrder[a.attr("point")] < pointOrder[b.attr("point")]) {
        return -1;
    } else if (-1*parseInt(a.attr("time")) > -1*parseInt(b.attr("time"))) {
        return 1;
    } else if (-1*parseInt(a.attr("time")) < -1*parseInt(b.attr("time"))) {
        return -1;
    }
    return 0;
}

function jinglePlayPause(event) {
    var s = $(this);
    var action = s.hasClass("playing") ? "pause" : "play";
    var id = s.attr("id").replace("jingle-", "");
    var url = "/track/" + action + "/" + id;
    $.ajax(url, {
        method: "POST"
    });
}

function jingleStop(event) {
    var s = $(this);
    var id = s.attr("id").replace("jingle-", "");
    var url = "/track/stop/" + id;
    $.ajax(url, {
        method: "POST"
    });
}

function jingleDelete(event) {
    var s = $(this);
    console.log(s);
    var id = s.attr("id").replace("jingle-", "");
    var url = "/track/delete/" + id;

    if (confirm("Fakt smazat jingle '" + s.find("strong").text() + "'?")) {
        $.ajax(url, {
            method: "DELETE"
        });
    }
}

function jingleChange(jingle) {
}

function jingleRemove(jingle) {
    $("#jingle-" + jingle.ID).remove()
    metroResize();
}

function songChange(song) {
    $("#jingle-" + song.ID).each(function () {
        $(this).find('.progress').css("width", Math.round(song.Position * 100) + "%");
        $(this).removeClass("playing");
        $(this).removeClass("paused");
        $(this).addClass(song.IsPlaying ? "playing" : song.Position > 0 ? "paused" : "");
    });
}

function log(log) {
    $('#logs')
        .prepend(
            $('<div>').attr('class', log.Severity)
            .append($('<small>').text(new Date(log.Time).toLocaleString())).append(' | ')
            .append($('<strong>').text(log.Severity)).append(' | ')
            .append(log.Message).append(' | ')
            .append($('<small>').text(log.File))
        )
}

function lpad (string, lenght, char) {
    char = char || " ";
    var o = string.toString()
	while (o.length < lenght) {
		o = char + o;
	}
	return o;
}
