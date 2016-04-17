$(document).ready(function() {
    hook();
    load();
});

var Handler = {
    "song_change"  : songChange,
    "song_added"   : songAdd,
    "song_removed" : songRemove,
    "app_added"    : appAdd,
    "app_removed"  : appRemove
}

var socket = new WebSocket("ws://localhost:8080/socket")
socket.onmessage = function(evt) {
    try {
        var data = $.parseJSON(evt.data);
        if (typeof Handler[data.Type] == 'undefined') {
            console.warn(data.Type, "Not implemented")
            return;
        }
        Handler[data.Type](data.Data);
    } catch (e) {
        console.error(e);
        console.info(evt.data);
    }
}

function hook() {
    $('a').click(clicker)

    $("form#addSong").submit(function () {
        var f = $(this);
        console.log(f.attr("method"));
        $.ajax("/track/add?filename=" + encodeURIComponent(f.find("input[name=filename]").val()), {
            method: "POST",
            success: function(data, status) { listTracks(data); },
            complete: function() { console.log("complete"); },
            error: function() { console.log("error"); }
        })
        return false;
    });
}

function load() {
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
    $("#apps")
        .append(
            $('<div id="app-' + app.ID + '" class="app"></div>')
            .append($('<a class="control" method="delete" href="/app/delete/' + app.ID + '">delete</a>')).append(' | ')
            .append($("<strong>").text(app.Name)).append(" ")
            .append($("<small>").text((100 * app.Volume) + "%"))
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

function songAdd(song) {
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
    $("#song-" + song.ID).each(function () {
        $(this).find('.state').text((song.IsPlaying ? "hraje" : "nehraje") + " " + Math.round(song.Position * 100) + " %");
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
    $.ajax(href, {
        method: (m ? m : "POST"),
        success: function(data, status) { callback(data); },
        complete: function() { console.log("complete"); },
        error: function() { console.log("error"); }
    })
    return false;
}


function defaultCallback(data) {
    console.log(data);
}
