package datastoretest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func dockerHaveImage(name string) (ok bool, err error) {
	out, err := exec.Command("docker", "images", "--no-trunc").Output()
	if err != nil {
		return
	}
	return bytes.Contains(out, []byte(name)), nil
}

func dockerRun(args ...string) (containerID string, err error) {
	runOut, err := exec.Command("docker", append([]string{"run"}, args...)...).Output()
	if err != nil {
		return
	}
	containerID = strings.TrimSpace(string(runOut))
	if containerID == "" {
		return "", errors.New("unexpected empty output from `docker run`")
	}
	return
}

func dockerKillContainer(container string) error {
	return exec.Command("docker", "kill", container).Run()
}

func dockerPull(name string) error {
	out, err := exec.Command("docker", "pull", name).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%v: %s", err, out)
	}
	return err
}

func dockerIP(containerID string) (string, error) {
	out, err := exec.Command("docker", "inspect", containerID).Output()
	if err != nil {
		return "", err
	}
	type networkSettings struct {
		IPAddress string
	}
	type container struct {
		NetworkSettings networkSettings
	}
	var c []container
	if err := json.NewDecoder(bytes.NewReader(out)).Decode(&c); err != nil {
		return "", err
	}
	if len(c) == 0 {
		return "", errors.New("no output from docker inspect")
	}
	if ip := c[0].NetworkSettings.IPAddress; ip != "" {
		return ip, nil
	}
	return "", errors.New("no IP. Not running?")
}
