package main

import (
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gitbao/gitbao/builder"
	"github.com/gitbao/gitbao/github"
	"github.com/gitbao/gitbao/model"
	"github.com/sqs/mux"
)

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/{username}/{gist-id}", DownloadHandler).Methods("GET")
	r.HandleFunc("/poll/{id}/{line-count}/", PollHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}

func DownloadHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	gistId := vars["gist-id"]
	username := vars["username"]

	path := "/" + username + "/" + gistId

	bao, err := github.GetGistData(gistId, path, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	model.DB.Create(&bao)

	bao.Console = "Welcome to GitBao!!\n" +
		"Getting ready to wrap up a tasty new Bao.\n"

	go func() {
		err := builder.StartBuild(&bao)
		if err != nil {
			panic(err)
		}
	}()

	type Inventory struct {
		Material string
		Count    uint
	}
	tmpl, err := template.New("bdoy").Parse(siteBody)
	if err != nil {
		// panic(err)
	}
	err = tmpl.Execute(w, bao)
	if err != nil {
		// panic(err)
	}

	// w.Write([]byte(siteBody))

}

func PollHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Second * 1)

	vars := mux.Vars(req)

	id := vars["id"]
	intid, err := strconv.Atoi(id)
	int64id := int64(intid)

	if err != nil || id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for {
		var bao model.Bao
		model.DB.Find(&bao, int64id)
		if bao.IsComplete == true {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(bao.Console))
			return
		} else {
			w.Write([]byte(bao.Console))
			return
		}
		time.Sleep(time.Second * 2)
	}
}

const siteBody = `

<html>
	<head>
		<script type="text/javascript" src="https://code.jquery.com/jquery-2.1.3.min.js"></script>
	</head>
	<body>
		<pre>{{.Console}}</pre>

<script type="text/javascript">
function longpoll(url, callback) {

    var req = new XMLHttpRequest();
    req.open('GET', url, true);

    req.onreadystatechange = function(aEvt) {
        if (req.readyState == 4) {
            if (req.status == 200) {

            	if (req.responseText != "done") {
            		longpoll(url, callback);
            	}
            } else {
                console.log("long-poll connection lost");
            }
            callback(req.responseText);

        }
    };

    req.send(null);
}
function writeToBody(text) {
	$('pre').text(text)
}

longpoll("/poll/1/0/", writeToBody)
</script>
	</body>
</html>

`
