//go:build mage

package main

import (
	"errors"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
)

// type Run mg.Namespace

// // Runs the site using hugo.
// func (Run) Client() error {
// 	return goRun("cmd/pmss/main.go", "scan")

// }

// // Runs the pdf docs.
// func (Run) Server() error {
// 	return goRun("cmd/pmssd/main.go", "server")
// }

type Build mg.Namespace

func (Build) Client() error {
	return goBuild("cmd/pmss/main.go", "bin/pmss")
}

func (Build) Server() error {
	return goBuild("cmd/pmssd/main.go", "bin/pmssd")
}

func goBuild(file, outFile string) error {
	return goCmd("build", "-o", outFile, file)
}

func Bootstrap() error {
	var err error
	err = errors.Join(err, goCmd("mod", "download"))
	err = errors.Join(err, Generate())
	return err
}

func Generate() error {
	return goCmd("generate", "./...")
}

func Test() error {
	return goCmd("test", "./...")
}

func goCmd(args ...string) error {
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
