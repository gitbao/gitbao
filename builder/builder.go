package builder

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/gitbao/gitbao/model"
)

func StartBuild(b *model.Bao) error {
	writeToBao(b, "Cloning gist files")
	directory, err := DownloadFromRepo(b)
	if err != nil {
		writeToBao(b, "Error cloning files")
		return err
	} else {
		writeToBao(b, "Files cloned successfully")
	}
	writeToBao(b, "Creating dockerfile")
	err = CreateDockerfile(directory)
	if err != nil {
		writeToBao(b, "Error creating dockerfile")
		return err
	}

	b.IsComplete = true
	model.DB.Save(b)

	return nil
}

func DownloadFromRepo(b *model.Bao) (directory string, err error) {
	path := "."
	directory, err = ioutil.TempDir(path, "forBuild")
	if err != nil {
		return
	}
	// fmt.Printf("%#v", b)
	err = runCommand(b, "git", "clone", b.GitPullUrl, path+"/"+directory)
	if err != nil {
		return
	}
	return
}

func CreateDockerfile(path string) error {
	contents := "FROM golang:onbuild\nEXPOSE 8080"
	err := ioutil.WriteFile(path+"/Dockerfile", []byte(contents), 0644)
	return err
}

func runCommand(b *model.Bao, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	// cmd.Stderr = os.Stderr

	// stderr, err := cmd.StderrPipe()
	// if err != nil {
	// 	return err
	// }

	// stdout, err := cmd.StdoutPipe()
	// if err != nil {
	// 	return err
	// }

	// http: //golang.org/pkg/bufio/#Scanner
	_ = cmd.Start()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// // b.Console = ""
	// scanner := bufio.NewScanner(stderr)
	// for scanner.Scan() {
	// 	writeToBao(b, scanner.Text())
	// }
	// scanner.Split(bufio.ScanWords)
	// for scanner.Scan() {
	// 	writeToBao(b, scanner.Text())
	// }

	// r := bufio.NewReader(stdout)
	// for line, _ := r.ReadLine() {
	// 	writeToBao(b, line)
	// }

	cmd.Wait()
	return nil
}

func writeToBao(b *model.Bao, text string) error {
	b.Console = b.Console + text + "\n"
	fmt.Println(b.Console)
	model.DB.Save(b)
	return nil
}

// func BuildDockerfile(path string) error {
//     cmd := exec.Command("name", ...)
// }
