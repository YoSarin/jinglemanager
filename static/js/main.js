$(document).ready(function() {
    load();

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
});

function load() {
    $.ajax("/track/list", {
        success: function(data, status) { listTracks(data); }
    });
}

function listTracks(data) {
    try {
        var data = $.parseJSON(data);
        $(data).each(function (k, v) {
            if($("#song-" + v.ID).length == 0) {
                $("#songs").append('<li id="song-' + v.ID + '" class="song"><strong>' + v.File + ' <a class="control" href="/track/play/' + v.ID + '">play</a> | <a class="control" href="/track/pause/' + v.ID + '">pause</a> | <a class="control" href="/track/stop/' + v.ID + '">stop</a></strong></li>');
            }
            $("#song-" + v.ID).find("a").each(function (k, v) {
                $(v).prop('onclick', null).off('click');
                $(v).click(controlClicker);
            });
        });
    } catch (e) {
        console.log(e);
    }
}

function controlClicker() {
    var href = $(this).attr("href");
    $.ajax(href, {
        method: "POST",
        success: function(data, status) { listTracks(data); },
        complete: function() { console.log("complete"); },
        error: function() { console.log("error"); }
    })
    return false;
}
