package builder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/gitbao/gitbao/model"
	"github.com/jinzhu/gorm"
)

func StartBuild(b *model.Bao, server model.Server) error {
	var docker model.Docker
	docker.ServerId = server.Id

	err := configDocker(&docker)
	if err != nil {
		writeToBao(b, "Error configuring Docker: "+err.Error(), true)
		return nil
	}

	writeToBao(b, "Cloning gist files", false)
	directory, err := DownloadFromRepo(b)
	if err != nil {
		writeToBao(b, "Error cloning files", true)
		return err
	} else {
		writeToBao(b, "Files cloned successfully", false)
	}
	writeToBao(b, "Creating dockerfile", false)
	err = CreateDockerfile(directory)
	if err != nil {
		writeToBao(b, "Error creating dockerfile", true)
		return err
	}
	writeToBao(b, "Building dockerfile (this could take a while)", false)
	dockerId, err := BuildDockerfile(b, directory, docker)
	if err != nil {
		writeToBao(b, err.Error()+"\nquitting...", true)
		return err
	}

	docker.DockerId = dockerId
	docker.BaoId = b.Id

	b.Location.BaoId = b.Id
	b.Location.Subdomain = fmt.Sprintf("%s-%d", b.GistId, b.Id)
	b.Location.Destination = fmt.Sprintf("%s:%d", server.Ip, docker.Port)

	writeToBao(b, "Site hosted on: "+b.Location.Subdomain+".gitbao.com", true)
	model.DB.Save(b)
	model.DB.Create(&docker)

	err = os.RemoveAll(directory)

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
	contents := `FROM golang:onbuild
EXPOSE 8080
`

	err := ioutil.WriteFile(path+"/Dockerfile", []byte(contents), 0644)
	return err
}

func BuildDockerfile(b *model.Bao, path string, docker model.Docker) (
	dockerId string, err error) {
	cmd := exec.Command("sudo", "docker", "build", "-t", "outyet", path)
	var stdobuild bytes.Buffer
	// var stdebuild bytes.Buffer
	cmd.Stdout = &stdobuild
	cmd.Stderr = &stdobuild
	err = cmd.Run()

	buildError := stdobuild.Bytes()
	if err != nil {
		err = fmt.Errorf("Error building application: \n%s", string(buildError))
		return
	}

	writeToBao(b, "Application built successfully\nStarting application:", false)

	// writeToBao(b, string(stdobuild.Bytes()))
	cmd = exec.Command("sudo", "docker", "run",
		"--publish", fmt.Sprintf("%d:8080", docker.Port),
		"--name", path,
		"--detach",
		"outyet",
	)
	output, err := cmd.Output()
	dockerId = string(output)
	fmt.Println(dockerId)

	if err != nil {
		err = fmt.Errorf("Error running application: %s\n", err)
		return
	}

	// writeToBao(b, string(stderun.Bytes()))
	return
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

func writeToBao(b *model.Bao, text string, isComplete bool) error {
	b.IsComplete = isComplete
	b.Console = b.Console + text + "\n"
	model.DB.Save(b)
	return nil
}

func configDocker(d *model.Docker) error {
	var lastDocker model.Docker
	query := model.DB.Order("port desc").Not("port = ?", 0).Where("server_id = ?", 5).First(&lastDocker)
	if query.Error != nil && query.Error != gorm.RecordNotFound {
		return query.Error
	}
	if lastDocker.Port < 9000 {
		d.Port = 9000
	} else {
		d.Port = lastDocker.Port + 1
	}
	fmt.Printf("%d %d", lastDocker.Port, d.Port)
	return nil
}

// func BuildDockerfile(path string) error {
//     cmd := exec.Command("name", ...)
// }
