$(document).ready(function() {
    $("a").click(function () {
        var href = $(this).attr("href");
        var anchor = $(this);
        if (href.match(/^\/track\/add/)) {
            $.ajax(href, {
                method: "POST",
                success: function(data, status) { data = $.parseJSON(data); anchor.attr("href", "/track/play/" + data.ID); anchor.text("Play") },
                complete: function() { console.log("complete"); },
                error: function() { console.log("error"); }
            })
            return false;
        }
        if (href.match(/^\/track\/play/)) {
            $.ajax(href, {
                method: "POST",
                success: function(data, status) { data = $.parseJSON(data); anchor.attr("href", "/track/pause/" + data.ID); anchor.text("Stop") },
                complete: function() { console.log("complete"); },
                error: function() { console.log("error"); }
            })
            return false;
        }
        if (href.match(/^\/track\/(stop|pause)/)) {
            $.ajax(href, {
                method: "POST",
                success: function(data, status) { data = $.parseJSON(data); anchor.attr("href", "/track/play/" + data.ID); anchor.text("Play") },
                complete: function() { console.log("complete"); },
                error: function() { console.log("error"); }
            })
            return false;
        }
    })
});
