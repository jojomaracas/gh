package main

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Git struct {
	Cmd string
}

func (git *Git) Version() (string, error) {
	output, err := git.execGitCmd([]string{"version"})
	if err != nil {
		return "", err
	}

	return output[0], err
}

func (git *Git) Dir() (string, error) {
	output, err := git.execGitCmd([]string{"rev-parse", "-q", "--git-dir"})
	if err != nil {
		return "", err
	}

	gitDir := output[0]
	gitDir, err = filepath.Abs(gitDir)
	if err != nil {
		return "", err
	}

	return gitDir, nil
}

func (git *Git) PullReqMsgFile() (string, error) {
	gitDir, err := git.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(gitDir, "PULLREQ_EDITMSG"), nil
}

func (git *Git) Editor() (string, error) {
	output, err := git.execGitCmd([]string{"var", "GIT_EDITOR"})
	if err != nil {
		return "", err
	}

	return output[0], nil
}

func (git *Git) EditorPath() (string, error) {
	gitEditor, err := git.Editor()
	if err != nil {
		return "", err
	}

	editorPath, err := exec.LookPath(gitEditor)
	if err != nil {
		return "", err
	}

	return editorPath, nil
}

func (git *Git) Owner() (string, error) {
	remote, err := git.Remote()
	if err != nil {
		return "", err
	}

	return mustMatchGitUrl(remote)[1], nil
}

func (git *Git) Project() (string, error) {
	remote, err := git.Remote()
	if err != nil {
		return "", err
	}

	return mustMatchGitUrl(remote)[2], nil
}

func (git *Git) Head() (string, error) {
	output, err := git.execGitCmd([]string{"symbolic-ref", "-q", "--short", "HEAD"})
	if err != nil {
		return "master", err
	}

	return output[0], nil
}

// FIXME: only care about origin push remote now
func (git *Git) Remote() (string, error) {
	r := regexp.MustCompile("origin\t(.+github.com.+) \\(push\\)")
	output, err := git.execGitCmd([]string{"remote", "-v"})
	if err != nil {
		return "", err
	}

	for _, o := range output {
		if r.MatchString(o) {
			return r.FindStringSubmatch(o)[1], nil
		}
	}

	return "", errors.New("Can't find remote")
}

func (git *Git) Log(sha1, sha2 string) (string, error) {
	execCmd := NewExecCmd("git")
	execCmd.WithArg("log").WithArg("--no-color")
	execCmd.WithArg("--format=%h (%aN, %ar)%n%w(78,3,3)%s%n%+b")
	execCmd.WithArg("--cherry")
	shaRange := fmt.Sprintf("%s...%s", sha1, sha2)
	execCmd.WithArg(shaRange)

	outputs, err := execCmd.ExecOutput()
	if err != nil {
		return "", err
	}

	return outputs, nil
}

func (git *Git) execGitCmd(input []string) (outputs []string, err error) {
	cmd := NewExecCmd(git.Cmd)
	for _, i := range input {
		cmd.WithArg(i)
	}

	out, err := cmd.ExecOutput()
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(out, "\n") {
		outputs = append(outputs, string(line))
	}

	return outputs, nil
}
