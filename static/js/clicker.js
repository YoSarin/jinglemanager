$(document).ready(function() {
    $('a').click(clicker);
});

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
    } else if (m == "none") {
        callback();
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

function defaultCallback(data) {
    console.log("DefaultCallback: ", data);
}
