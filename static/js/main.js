$(document).ready(function() {
    hook();
    load();
});

var socket = new WebSocket("ws://localhost:8080/socket")
socket.onmessage = function(evt) {
    console.log(evt.data);
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
                $("#songs")
                    .append(
                        $('<div id="song-' + v.ID + '" class="song"></div>')
                        .append($('<a class="control" href="/track/play/' + v.ID + '">play</a>')).append(' | ')
                        .append($('<a class="control" href="/track/stop/' + v.ID + '">stop</a>')).append(' | ')
                        .append($('<a class="control" href="/track/pause/' + v.ID + '">pause</a>')).append(' | ')
                        .append($('<a class="control" method="delete" href="/track/delete/' + v.ID + '">delete</a>')).append(' | ')
                        .append($("<strong>").text(v.File))

                    ).find("a").each(function (k, v) {
                        $(v).prop('onclick', null).off('click');
                        $(v).click({callback: listTracks}, clicker);
                    });
                }
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
                $("#apps")
                    .append(
                        $('<div id="app-' + v.ID + '" class="app"></div>')
                        .append($('<a class="control" href="/app/delete/' + v.ID + '">delete</a>')).append(' | ')
                        .append($("<strong>").text(v.Name)).append(" ")
                        .append($("<small>").text((100 * v.Volume) + "%"))
                    ).find("a").each(function (k, v) {
                        $(v).prop('onclick', null).off('click');
                        $(v).click({callback: listApps}, clicker);
                    });
                }
        });
    } catch (e) {
        console.log(e);
    }
}

function clicker(event) {
    var callback;
    try {
        callback = event.data.callback;
        if (typeof callback == 'undefined') {
            callback = defaultCallback;
        }
    } catch (e) {
        console.log(e);
        callback = defaultCallback;
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
