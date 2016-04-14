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
                        $(v).click(controlClicker);
                    });
                }
        });
    } catch (e) {
        console.log(e);
    }
}

function controlClicker() {
    var href = $(this).attr("href");
    var m = $(this).attr("method");
    $.ajax(href, {
        method: (m ? m : "POST"),
        success: function(data, status) { listTracks(data); },
        complete: function() { console.log("complete"); },
        error: function() { console.log("error"); }
    })
    return false;
}
