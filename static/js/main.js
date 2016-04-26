var scripts = document.getElementsByTagName("script");
var __FILE__ = scripts[scripts.length-1].src;

$(document).ready(function() {
    hook();
    load();
});

var Handler = {
    "song_changed"   : songChange,
    "song_added"     : songAdd,
    "song_removed"   : songRemove,
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

function listTracks(data) {
    try {
        var data = $.parseJSON(data);
        $("#songs").empty();
        $.each(data, function (k, v) {
            if($("#song-" + v.ID).length == 0) {
                songAdd(v);
            }
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
            $('<div id="app-' + app.ID + '" class="app"></div>')
            .append($('<a class="control" method="delete" href="/app/delete/' + app.ID + '">delete</a>')).append(' | ')
            .append($("<strong>").text(app.Name)).append(" ")
            .append($('<small class="state">').text(Math.round(100 * app.Volume) + "%"))
        ).find("a").each(function (k, v) {
            $(v).prop('onclick', null).off('click');
            $(v).click(clicker);
        });
}

function appRemove(app) {
    $("#app-" + app.ID).remove()
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
    $("#jingles")
        .append(
            $('<div id="jingle-' + jingle.Song.ID + '" class="song"></div>')
            .append($('<a class="control" href="/track/play/' + jingle.Song.ID + '">play</a>')).append(' | ')
            .append($('<a class="control" href="/track/stop/' + jingle.Song.ID + '">stop</a>')).append(' | ')
            .append($('<a class="control" href="/track/pause/' + jingle.Song.ID + '">pause</a>')).append(' | ')
            .append($('<a class="control" method="delete" href="/track/delete/' + jingle.Song.ID + '">delete</a>')).append(' | ')
            .append($('<small class="state">').text(jingle.Song.IsPlaying ? "hraje" : "nehraje")).append(' | ')
            .append($("<strong>").text(jingle.Name)).append(' | ')
            .append($('<small class="songtitle">').text(jingle.Song.File))
        ).find("a").each(function (k, v) {
            $(v).prop('onclick', null).off('click');
            $(v).click(clicker);
        });
}

function jingleChange(jingle) {
}

function jingleRemove(jingle) {
    $("#jingle-" + jingle.ID).remove()
}

function songAdd(song) {
    if($("#song-" + song.ID).length > 0) {
        return;
    }
    $("#songs")
        .append(
            $('<div id="song-' + song.ID + '" class="song"></div>')
            .append($('<a class="control" href="/track/play/' + song.ID + '">play</a>')).append(' | ')
            .append($('<a class="control" href="/track/stop/' + song.ID + '">stop</a>')).append(' | ')
            .append($('<a class="control" href="/track/pause/' + song.ID + '">pause</a>')).append(' | ')
            .append($('<a class="control" method="delete" href="/track/delete/' + song.ID + '">delete</a>')).append(' | ')
            .append($('<small class="state">').text(song.IsPlaying ? "hraje" : "nehraje")).append(' | ')
            .append($("<strong>").text(song.File))

        ).find("a").each(function (k, v) {
            $(v).prop('onclick', null).off('click');
            $(v).click(clicker);
        });
}
function songRemove(song) {
    $("#song-" + song.ID).remove()
}
function songChange(song) {
    $("#song-" + song.ID + ", #jingle-" + song.ID).each(function () {
        $(this).find('.state').text((song.IsPlaying ? "hraje" : "nehraje") + " " + Math.round(song.Position * 100) + "%");
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
