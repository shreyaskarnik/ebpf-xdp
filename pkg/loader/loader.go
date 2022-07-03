package loader

import (
	"log"
	"net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
	"github.com/shreyaskarnik/ebpf-xdp/pkg/printer"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS bpf ../../ebpf/xdp.c -- -I../../ebpf/headers
func SetupEBPF(ifaceName string) {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Fatal(err)
	}
	// load
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatal(err)
	}
	defer objs.Close()
	l, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.XdpProgFunc,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	log.Printf("XDP attached to %q (index %d)", iface.Name, iface.Index)
	log.Printf("Press Ctrl-C to exit and remove the program")
	p := tea.NewProgram(printer.NewModel(objs.XdpStatsMap))
	err = p.Start()
	if err != nil {
		log.Fatal(err)
	}
}
