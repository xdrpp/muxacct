
package main

import (
	"fmt"
	"github.com/xdrpp/goxdr/xdr"
	"os"
	"path"
)

//go:generate goxdr -o muxedaccount.go muxedaccount.x

var progname string

func usage(stat int) {
	out := os.Stderr
	if stat == 0 {
		out = os.Stdout
	}
	fmt.Fprintf(out, "usage: %s demux <muxed account ID>\n", progname)
	fmt.Fprintf(out, "       %s mux <uint64> <pubkey>\n", progname)
	os.Exit(stat)
}

func dumpXdr(t xdr.XdrType) {
	bin := XdrToBytes(t)
	fmt.Print("{")
	for i := range bin {
		if i % 8 == 0 {
			fmt.Print("\n    ")
		} else {
			fmt.Print(" ")
		}
		fmt.Printf("0x%02x,", bin[i])
	}
	fmt.Print("\n}\n")
}

func main() {
	if len(os.Args) > 0 {
		progname = path.Base(os.Args[0])
	}
	if len(os.Args) < 2 {
		usage(1)
	}
	switch os.Args[1] {
	case "help", "-h", "-help", "--help":
		usage(0)
	case "mux":
		if len(os.Args) != 4 {
			usage(1)
		}
		var m MuxedAccount
		m.Type = KEY_TYPE_MUXED_ED25519
		if _, err := fmt.Sscan(os.Args[2], &m.Med25519().Id); err != nil {
			fmt.Fprintf(os.Stderr, "%s: can't parse %q as integer: %s\n",
				progname, os.Args[2], err)
			os.Exit(1)
		}
		var pk PublicKey
		if _, err := fmt.Sscan(os.Args[3], &pk); err != nil ||
			pk.Type != PUBLIC_KEY_TYPE_ED25519 {
			fmt.Fprintf(os.Stderr,
				"%s: can't parse %q as ed25519 public key\n",
				progname, os.Args[3])
			os.Exit(1)
		}
		copy(m.Med25519().Ed25519[:], pk.Ed25519()[:])
		fmt.Println(m)
		dumpXdr(&m)
	case "demux":
		if len(os.Args) != 3 {
			usage(1)
		}
		var m MuxedAccount
		if _, err := fmt.Sscan(os.Args[2], &m); err != nil ||
			m.Type != KEY_TYPE_MUXED_ED25519 {
			fmt.Fprintf(os.Stderr,
				"%s: can't parse %q as muxed ed25519 account ID\n",
				progname, os.Args[2])
			os.Exit(1)
		}
		pk := PublicKey{ Type: PUBLIC_KEY_TYPE_ED25519 }
		copy(pk.Ed25519()[:], m.Med25519().Ed25519[:])
		fmt.Println(m.Med25519().Id, pk)
		dumpXdr(&m)
	case "dump":
		if len(os.Args) != 3 {
			usage(1)
		}
		var m MuxedAccount
		if _, err := fmt.Sscan(os.Args[2], &m); err != nil {
			fmt.Fprintf(os.Stderr,
				"%s: can't parse %q as muxed account ID\n",
				progname, os.Args[2])
			os.Exit(1)
		}
		dumpXdr(&m)
	default:
		usage(1)
	}
}
