//go:build mage
// +build mage

package main

import (
	"os"
	"path"
	"path/filepath"

	"github.com/aserto-dev/mage-loot/buf"
	"github.com/aserto-dev/mage-loot/common"
	"github.com/aserto-dev/mage-loot/deps"
	"github.com/aserto-dev/mage-loot/fsutil"
)

func init() {
}

func Deps() {
	deps.GetAllDeps()
}

// Lint runs linting for the entire project.
func Lint() error {
	return common.Lint()
}

func BuildAll() error {
	if err := BufBuild(); err != nil {
		return err
	}
	if err := BufGenerate(); err != nil {
		return err
	}
	return Build()
}

// Build binaries.
func Build() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	cfg := path.Join(cwd, ".goreleaser.yml")
	return common.BuildReleaser("--config", cfg, "--clean")
}

func BufBuild() error {
	if err := bufModUpdate("proto"); err != nil {
		return err
	}

	if err := bufBuild("bin/dsload.bin"); err != nil {
		return err
	}

	return nil
}

func BufGenerate() error {
	bufImage := "bin/dsload.bin"

	os.Setenv("BUF_BETA_SUPPRESS_WARNINGS", "1")

	if err := bufGenerate(bufImage); err != nil {
		return err
	}

	return nil
}

func bufBuild(bin string) error {
	fsutil.EnsureDir(path.Dir(bin))

	return buf.Run(
		buf.AddArg("build"),
		buf.AddArg("--output"),
		buf.AddArg(bin),
	)
}

func bufModUpdate(dir string) error {
	return buf.Run(
		buf.AddArg("mod"),
		buf.AddArg("update"),
		buf.AddArg(dir),
	)
}

func bufGenerate(image string) error {
	if err := bufModUpdate("."); err != nil {
		return err
	}

	oldPath := os.Getenv("PATH")

	pathSeparator := string(os.PathListSeparator)
	path := oldPath +
		pathSeparator +
		filepath.Dir(deps.GoBinPath("protoc-gen-go-grpc")) +
		pathSeparator +
		filepath.Dir(deps.GoBinPath("protoc-gen-go"))

	return buf.RunWithEnv(map[string]string{
		"PATH": path,
	},
		buf.AddArg("generate"),
		buf.AddArg(image),
	)
}
