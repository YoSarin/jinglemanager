var scripts = document.getElementsByTagName("script");
var __FILE__ = scripts[scripts.length-1].src;
var pointOrder = {
    "match_start" : 0,
    "match_end"   : 1,
    "match_none"  : 2,
}


$(document).ready(function() {
    hook();
    load();
    showHideJingleMatchDetails();
});

$(window).resize(function (event) {
    metroResize();
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

function hook() {
    $('a').click(clicker);

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
}

function metroResize() {
    $(".metro").each(function (k, v) {
        var count = $(this).find("li").length;
        var width = Math.floor($(this).width() - 0.5);
        var maxRowItems = Math.floor(width/250);
        var counter = 0;
        $(v).find("li").each(function (_, item) {
            var rowItems = Math.min(maxRowItems, count - (Math.floor(counter / maxRowItems) * maxRowItems));
            var w =  Math.floor(width / Math.min(rowItems, maxRowItems));
            $(item).width(w + "px");
            counter++;
        });
    });

    if ($(window).width() < 400) {
        $("#sideColumn").css("float", "none");
    } else {
        $("#sideColumn").css("float", "right");
    }
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

function clicker(event) {
    var callback;
    if (event == null || event.data == null || typeof event.data.callback == 'undefined') {
        callback = defaultCallback;
    } else {
        callback = event.data.callback;
    }
    var m = $(this).attr("method");
    var href = $(this).attr("href")
    if (m == "download") {
        $('iframe#downloader').attr("src", href);
    } else if (m == "visit") {
        return true;
    } else {
        $.ajax(href, {
            method: (m ? m : "POST"),
            success: function(data, status) { callback(data); },
            complete: function() { console.log("complete"); },
            error: function() { console.log("error"); }
        });
    }
    return false;
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

function defaultCallback(data) {
    console.log("DefaultCallback: ", data);
}
