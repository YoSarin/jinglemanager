$(document).ready(function() {
    $("#tournamentList a").bind('click', {callback: function () {
        window.location = '/'
    }}, clicker)
});
