package docker

import (
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/artifactory/utils"
	"io"
	"os/exec"
	"path"
	"strings"
)

func New(imageTag string) Image {
	return &image{tag: imageTag}
}

// Docker image
type Image interface {
	Push() error
	Id() (string, error)
	ParentId() (string, error)
	Tag() string
	Path() string
	Name() string
}

// Internal implementation of docker image
type image struct {
	tag string
}

// Push docker image
func (image *image) Push() error {
	cmd := &pushCmd{image: image}
	return utils.RunCmd(cmd)
}

// Get docker image tag
func (image *image) Tag() string {
	return image.tag
}

// Get docker image ID
func (image *image) Id() (string, error) {
	cmd := &getImageIdCmd{image: image}
	content, err := utils.RunCmdOutput(cmd)
	return strings.Trim(string(content), "\n"), err
}

// Get docker parent image ID
func (image *image) ParentId() (string, error) {
	cmd := &getParentId{image: image}
	content, err := utils.RunCmdOutput(cmd)
	return strings.Trim(string(content), "\n"), err
}

// Get docker image relative path in Artifactory
func (image *image) Path() string {
	indexOfFirstSlash := strings.Index(image.tag, "/")
	indexOfLastColon := strings.LastIndex(image.tag, ":")

	if indexOfLastColon < 0 || indexOfLastColon < indexOfFirstSlash {
		return path.Join(image.tag[indexOfFirstSlash:], "latest")
	}
	return path.Join(image.tag[indexOfFirstSlash:indexOfLastColon], image.tag[indexOfLastColon+1:])
}

// Get docker image name
func (image *image) Name() string {
	indexOfFirstSlash := strings.Index(image.tag, "/")
	indexOfLastColon := strings.LastIndex(image.tag, ":")

	if indexOfLastColon < 0 || indexOfLastColon < indexOfFirstSlash {
		return image.tag[indexOfFirstSlash+1:] + ":latest"
	}
	return image.tag[indexOfFirstSlash+1:]
}

// Image push command
type pushCmd struct {
	image *image
}

func (pushCmd *pushCmd) GetCmd() *exec.Cmd {
	var cmd []string
	cmd = append(cmd, "docker")
	cmd = append(cmd, "push")
	cmd = append(cmd, pushCmd.image.tag)
	return exec.Command(cmd[0], cmd[1:]...)
}

func (pushCmd *pushCmd) GetEnv() map[string]string {
	return map[string]string{}
}

func (pushCmd *pushCmd) GetStdWriter() io.WriteCloser {
	return nil
}
func (pushCmd *pushCmd) GetErrWriter() io.WriteCloser {
	return nil
}

// Image get image id command
type getImageIdCmd struct {
	image *image
}

func (getImageId *getImageIdCmd) GetCmd() *exec.Cmd {
	var cmd []string
	cmd = append(cmd, "docker")
	cmd = append(cmd, "images")
	cmd = append(cmd, "--format", "{{.ID}}")
	cmd = append(cmd, "--no-trunc")
	cmd = append(cmd, getImageId.image.tag)
	return exec.Command(cmd[0], cmd[1:]...)
}

func (getImageId *getImageIdCmd) GetEnv() map[string]string {
	return map[string]string{}
}

func (getImageId *getImageIdCmd) GetStdWriter() io.WriteCloser {
	return nil
}
func (getImageId *getImageIdCmd) GetErrWriter() io.WriteCloser {
	return nil
}

// Image get parent image id command
type getParentId struct {
	image *image
}

func (getImageId *getParentId) GetCmd() *exec.Cmd {
	var cmd []string
	cmd = append(cmd, "docker")
	cmd = append(cmd, "inspect")
	cmd = append(cmd, "--format", "{{.Parent}}")
	cmd = append(cmd, getImageId.image.tag)
	return exec.Command(cmd[0], cmd[1:]...)
}

func (getImageId *getParentId) GetEnv() map[string]string {
	return map[string]string{}
}

func (getImageId *getParentId) GetStdWriter() io.WriteCloser {
	return nil
}
func (getImageId *getParentId) GetErrWriter() io.WriteCloser {
	return nil
}
