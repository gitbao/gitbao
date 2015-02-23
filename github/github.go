package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gitbao/gitbao/model"
)

type githubApiGistResponse struct {
	Html_url     string
	Git_pull_url string
	Files        map[string]githubApiResponseFile
}
type githubApiResponseFile struct {
	Filename string
	Language string
	Raw_url  string
}

type gistJson struct {
	Files []string
}

var github_access_key string

func init() {
	github_access_key = os.Getenv("GITHUB_GIST_ACCESS_KEY")
	if github_access_key == "" {
		panic("Github access key required")
	}
}

func GetGistData(bao *model.Bao, path string, useAlternate bool) (err error) {
	if useAlternate != true {
		err = GetData(bao)
	} else {
		err = GetDataAlternate(bao, path)
	}
	// fmt.Printf("%#v", bao)
	return
}

func GetData(b *model.Bao) error {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/gists/"+b.GistId,
		nil,
	)
	if err != nil {
		return err
	}
	req.SetBasicAuth(
		github_access_key,
		"",
	)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error code %d: %s", resp.StatusCode, string(contents))
	}
	if err != nil {
		return err
	}
	var responseStruct githubApiGistResponse
	err = json.Unmarshal(contents, &responseStruct)
	if err != nil {
		return err
	}
	for key, value := range responseStruct.Files {
		b.Files = append(b.Files, model.File{
			Filename: key,
			Language: value.Language,
			RawUrl:   value.Raw_url,
		})
	}

	b.Url = responseStruct.Html_url
	b.GitPullUrl = responseStruct.Git_pull_url
	return nil
}

func GetDataAlternate(b *model.Bao, path string) error {
	files := make(map[string]string)
	gistUrl := "https://gist.github.com" + path
	rawPath := "https://gist.githubusercontent.com" + path

	resp, err := http.Get(gistUrl + ".json")
	if err != nil {
		return err
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", string(contents))
	}
	if err != nil {
		return err
	}
	if resp.StatusCode == 404 {
		err = fmt.Errorf("Gist not found")
		return err
	}
	var gistJson gistJson
	err = json.Unmarshal(contents, &gistJson)
	if err != nil {
		return err
	}

	for _, value := range gistJson.Files {
		var resp *http.Response
		resp, err = http.Get(rawPath + "/raw/" + value)
		if err != nil {
			return err
		}
		var contents []byte
		contents, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		files[value] = string(contents)
	}
	b.Url = gistUrl
	b.GitPullUrl = "https://gist.github.com/" + b.GistId + ".git"
	return nil
}
