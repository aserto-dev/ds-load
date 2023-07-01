//go:build mage
// +build mage

package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"path/filepath"

	"github.com/aserto-dev/mage-loot/buf"
	"github.com/aserto-dev/mage-loot/common"
	"github.com/aserto-dev/mage-loot/deps"
	"github.com/aserto-dev/mage-loot/fsutil"
	"github.com/itchyny/gojq"
	"github.com/magefile/mage/sh"
)

func init() {
}

func Deps() {
	deps.GetAllDeps()
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

// Release releases the project.
func Release() error {
	return common.Release()
}

func GetWorkspacePaths() []string {
	var outBuf bytes.Buffer
	ran, err := sh.Exec(nil, &outBuf, os.Stderr, "go", "work", "edit", "-json")
	if err != nil || !ran {
		return []string{}
	}

	var v interface{}
	if err := json.Unmarshal(outBuf.Bytes(), &v); err != nil {
		return []string{}
	}

	q, err := gojq.Parse(".Use[].DiskPath")
	if err != nil {
		return []string{}
	}

	results := []string{}
	iter := q.Run(v)
	for {
		result, ok := iter.Next()
		if !ok {
			break
		}
		if str, ok := result.(string); ok {
			results = append(results, str)
		}
	}

	return results
}

// Test - based on ci.yaml implementation:
// go work edit -json | jq -r '.Use[].DiskPath'  | xargs -I{} .ext/gobin/gotestsum-v1.10.0/gotestsum --format short-verbose -- -count=1 -v {}/...
func Test() error {
	for _, p := range GetWorkspacePaths() {
		if err := deps.GoDep("gotestsum")([]string{"--format", "short-verbose", "--", "-count=1", "-v", "./" + filepath.Join(p, "...")}...); err != nil {
			return err
		}
	}
	return nil
}

// Lint - based on ci.yaml implementation:
// go work edit -json | jq -r '.Use[].DiskPath'  | xargs -I{} .ext/gobin/golangci-lint-v1.52.2/golangci-lint run {}/... -c .golangci.yaml
func Lint() error {
	for _, p := range GetWorkspacePaths() {
		if err := deps.GoDep("golangci-lint")([]string{"run", "./" + filepath.Join(p, "..."), "-c", ".golangci.yaml"}...); err != nil {
			return err
		}
	}
	return nil
}
