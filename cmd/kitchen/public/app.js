$(function() {
    $('a.deploy').click(function() {
        $.post($(this).attr("href"), function( data ) {
            console.log(data)
        });
        $(this).attr("disable", true)
        return false
    })
})
function longpoll(url, callback) {
    var req = new XMLHttpRequest();
    req.open('GET', url, true);
    req.responseType = "text"
    req.onerror = function(aEvt) {
        console.log("error")
    };
    req.onreadystatechange = function(aEvt) {
        if (req.readyState == 4) {
            if (req.status == 200) {
                var response = JSON.parse(req.responseText)
                if (response.IsComplete != true) {
                    longpoll(url, callback);
                }
                if (response.IsReady == true) {
                    $('a.deploy').attr('disabled', false)
                }
                callback(response);
            } else {
                $('pre.console').append("\nError connecting to deployment server.")
            }
        }
    };

    req.send(null);
}
function writeToBody(responseJson) {
    $('pre.console').text(responseJson.Console)
}