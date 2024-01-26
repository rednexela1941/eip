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
	helpFlag := flag.Bool("h", false, "print help")
	edsFlag := flag.Bool("eds", false, "print generated eds file to stdout")
	ifaceNameFlag := flag.String("iface", MyNetworkInterfaceName, "name of the network interface to listen on")

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	iface, err := linux.CreateInterface(*ifaceNameFlag)
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
