function longpoll(url, callback) {
    var req = new XMLHttpRequest();
    req.open('GET', url, true);
    req.responseType = "text"

    req.onreadystatechange = function(aEvt) {
        if (req.readyState == 4) {
            var response = JSON.parse(req.responseText)
            if (req.status == 200) {
                if (response.IsComplete != true) {
                    longpoll(url, callback);
                }
            } else {
                console.log("long-poll connection lost");
            }
            callback(response);

        }
    };

    req.send(null);
}
function writeToBody(responseJson) {
    $('pre.console').text(responseJson.Console)
}