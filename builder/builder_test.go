package builder

import (
	"os"
	"testing"

	"github.com/gitbao/gitbao/model"
)

var bao model.Bao

func init() {
	bao = model.Bao{
		GistId:     "ba67234b79784c75cfd9",
		GitPullUrl: "https://gist.github.com/ba67234b79784c75cfd9.git",
	}
}

func TestDownloadFromRepo(t *testing.T) {

	directory, err := DownloadFromRepo(&bao)
	if err != nil {
		t.Error(err)
	}
	// fmt.Printf("%s", directory)

	err = os.RemoveAll(directory)
	if err != nil {
		t.Error(err)
	}
}

func TestWriteToBao(t *testing.T) {
	testText := "This is some text"
	model.DB.Create(&bao)
	writeToBao(&bao, testText, false)
	if bao.Console != testText+"\n" {
		t.Errorf("Not writing to bao correctly")
	}
}

func TestCreateDockerfile(t *testing.T) {
	err := CreateDockerfile(".")
	if err != nil {
		t.Error(err)
	}
}
