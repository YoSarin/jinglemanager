$(document).ready(function() {
    metroResize();
});

$(window).resize(function (event) {
    metroResize();
});

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
