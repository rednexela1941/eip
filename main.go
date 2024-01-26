package main

import (
	"flag"
	"log"
	"os"

	"github.com/rednexela1941/eip/pkg/eds"
	"github.com/rednexela1941/eip/pkg/linux"
	"github.com/rednexela1941/eip/test/adaptertest"
)

const (
	MyNetworkInterfaceName = "enp0s20f0u3u3"
)

func main() {
	edsFlag := flag.Bool("eds", false, "print generated eds file to stdout")
	flag.Parse()

	iface, err := linux.CreateInterface(MyNetworkInterfaceName)
	if err != nil {
		log.Fatal(err)
	}

	adapter := adaptertest.MakeTestAdapter()
	if err := adapter.AddNetworkInterface(iface); err != nil {
		log.Fatal(err)
	}

	if *edsFlag {
		if err := eds.WriteEDS(os.Stdout, adapter, nil); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := adapter.StartNetworkLoop(); err != nil {
		log.Fatal(err)
	}
}
