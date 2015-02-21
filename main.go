package main

import (
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gitbao/gitbao/builder"
	"github.com/gitbao/gitbao/github"
	"github.com/gitbao/gitbao/model"
	"github.com/gitbao/gitbao/router"
	"github.com/sqs/mux"
)

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/{username}/{gist-id}", DownloadHandler).Methods("GET")
	r.HandleFunc("/poll/{id}/{line-count}/", PollHandler).Methods("GET")
	http.Handle("/", r)
	go http.ListenAndServe(":8000", nil)

	http.ListenAndServe(":8001", &router.Router{})
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("index"))
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
	_ = model.DB.Create(&bao)

	stringId := strconv.Itoa(int(bao.Id))
	bao.Location = model.Location{
		Destination: "localhost:8000",
		Subdomain:   bao.GistId + "-" + stringId,
	}

	bao.Console = "Welcome to GitBao!!\n" +
		"Getting ready to wrap up a tasty new Bao.\n" +
		"Found some files:\n"

	var isGo bool
	for _, value := range bao.Files {
		bao.Console += "    " + value.Filename + "\n"
		if value.Language == "Go" {
			isGo = true
		}
	}

	if isGo != true {
		bao.Console += "Whoops!\n" +
			"GitBao only supports Go programs at the moment.\n" +
			"Quitting...."
		bao.IsComplete = true
		model.DB.Save(&bao)
	} else {
		go func() {
			err := builder.StartBuild(&bao)
			if err != nil {
				panic(err)
			}
		}()
	}

	tmpl, err := template.New("body").Parse(siteBody)
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

func RouterHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hello"))
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

longpoll("/poll/{{.Id}}/0/", writeToBody)
</script>
	</body>
</html>

`
