//go:build linux
// +build linux

package main

import (
	"log"
	"os"

	"github.com/shreyaskarnik/ebpf-xdp/pkg/loader"
)

// $BPF_CLANG and $BPF_CFLAGS are set by the Makefile.
func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <ifname>", os.Args[0])
	}
	// lookup interface by name
	ifaceName := os.Args[1]
	loader.SetupEBPF(ifaceName)
}
