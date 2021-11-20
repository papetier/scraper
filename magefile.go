//go:build mage
// +build mage

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/magefile/mage/mg"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type Build mg.Namespace

const (
	binaryPath      = buildPath + "/waiter"
	binaryPathLocal = buildPath + "/waiter_local"
	buildPath       = "build"
	cmdPath         = "./cmd/server"
	packageName     = "github.com/papetier/scraper"
	prodArch        = "amd64"
	prodOs          = "linux"
)

// Build the production binary (linux/amd64)
func (Build) Prod() error {
	err := prepareEnv(prodArch, prodOs)
	if err != nil {
		return err
	}
	ldFlags := getXLdflags(prodArch, prodOs)
	fullLdFlags := "-w -s " + ldFlags
	runAndStreamOutput("go", "build", "-v", "-ldflags", fullLdFlags, "-o", binaryPath, cmdPath)
	return nil
}

// Build a local binary to run on your current device
func (Build) Local() error {
	err := prepareEnv(runtime.GOARCH, runtime.GOOS)
	if err != nil {
		return err
	}
	ldFlags := getXLdflags(runtime.GOARCH, runtime.GOOS)
	runAndStreamOutput("go", "build", "-v", "-ldflags", ldFlags, "-o", binaryPathLocal, cmdPath)
	return nil
}

func getCommitShortHash() string {
	commitShortHash, err := runCmdWithOutput("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		fmt.Printf("Error getting commit short hash: %s\n", err)
		os.Exit(1)
	}
	return strings.Trim(string(commitShortHash), "\n")
}

func getXLdflags(buildArch, buildOs string) string {
	currentTime := time.Now().Format(time.RFC3339)

	return `-X '` + packageName + `/pkg/version.Version=` + getVersion() + `' -X '` + packageName + `/pkg/version.CommitShortHash=` + getCommitShortHash() + `' -X '` + packageName + `/pkg/version.Arch=` + buildArch + `' -X '` + packageName + `/pkg/version.Os=` + buildOs + `' -X '` + packageName + `/pkg/version.Time=` + currentTime + `'`
}

func getVersion() string {
	version, err := runCmdWithOutput("git", "describe", "--tags", "--always", "--abbrev=10")
	if err != nil {
		fmt.Printf("Error getting version: %s\n", err)
		os.Exit(1)
	}
	return strings.Trim(string(version), "\n")
}

func prepareEnv(buildArch, buildOs string) error {
	// common flags
	err := os.Setenv("CGO_ENABLED", "0")
	if err != nil {
		return fmt.Errorf("setting CGO_ENABLED environment: %w", err)
	}

	// target specific flags
	err = os.Setenv("GOARCH", buildArch)
	if err != nil {
		return fmt.Errorf("setting GOARCH environment: %w", err)
	}
	err = os.Setenv("GOOS", buildOs)
	if err != nil {
		return fmt.Errorf("setting GOOS environment: %w", err)
	}
	return nil
}

func runCmdWithOutput(name string, arg ...string) (output []byte, err error) {
	cmd := exec.Command(name, arg...)
	output, err = cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running command: %s, %w", cmd, err)
	}

	return output, nil
}

func runAndStreamOutput(cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.Env = os.Environ()

	fmt.Printf("%s\n\n", c.String())

	stdout, _ := c.StdoutPipe()
	errbuf := bytes.Buffer{}
	c.Stderr = &errbuf
	c.Start()

	reader := bufio.NewReader(stdout)
	line, err := reader.ReadString('\n')
	for err == nil {
		fmt.Print(line)
		line, err = reader.ReadString('\n')
	}

	if err := c.Wait(); err != nil {
		fmt.Printf(errbuf.String())
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
