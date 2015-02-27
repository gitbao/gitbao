package github

import (
	"fmt"
	"testing"

	"github.com/gitbao/gitbao/model"
)

func TestGetGistData(t *testing.T) {
	// /maxmcd/ba67234b79784c75cfd9
	// useAlternate := true
	bao := model.Bao{
		GistId: "ba67234b79784c75cfd9",
	}
	// for {
	err := GetGistData(&bao)
	if err != nil {
		t.Error(err)
	}

	if bao.GitPullUrl != "https://gist.github.com/ba67234b79784c75cfd9.git" {
		t.Error(fmt.Errorf("Wrong git pull url"))
	}
	if bao.GistId != "ba67234b79784c75cfd9" {
		t.Error(fmt.Errorf("Wrong gist id"))
	}
	var url string
	// if useAlternate == true {
	// url = "https://gist.github.com/maxmcd/ba67234b79784c75cfd9"
	// } else {
	url = "https://gist.github.com/ba67234b79784c75cfd9"
	// }
	if bao.Url != url {
		t.Error(fmt.Errorf("Wrong url"))
	}
	// 	if useAlternate == true {
	// 		useAlternate = false
	// 	} else {
	// 		break
	// 	}
	// }
}
