package main

import "github.com/aserto-dev/ds-load/sdk/common/version"

func main() {
	vi := version.GetInfo()
	println("ds-sync", vi)
}
