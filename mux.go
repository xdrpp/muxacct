
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
	fmt.Fprintf(out, "       %s dump <muxed account ID>\n", progname)
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
	args := os.Args
	if len(args) < 2 {
		usage(1)
	}
	switch arg := args[1]; arg {
	case "help", "-h", "-help", "--help":
		usage(0)
	case "mux":
		if len(args) != 4 {
			usage(1)
		}
		var m MuxedAccount
		m.Type = KEY_TYPE_MUXED_ED25519
		if _, err := fmt.Sscan(args[2], &m.Med25519().Id); err != nil {
			fmt.Fprintf(os.Stderr, "%s: can't parse %q as integer: %s\n",
				progname, args[2], err)
			os.Exit(1)
		}
		var pk PublicKey
		if _, err := fmt.Sscan(args[3], &pk); err != nil ||
			pk.Type != PUBLIC_KEY_TYPE_ED25519 {
			fmt.Fprintf(os.Stderr,
				"%s: can't parse %q as ed25519 public key\n",
				progname, args[3])
			os.Exit(1)
		}
		copy(m.Med25519().Ed25519[:], pk.Ed25519()[:])
		fmt.Println(m)
	case "demux", "dump":
		if len(args) != 3 {
			usage(1)
		}
		var m MuxedAccount
		if _, err := fmt.Sscan(args[2], &m); err != nil ||
			m.Type != KEY_TYPE_MUXED_ED25519 {
			fmt.Fprintf(os.Stderr,
				"%s: can't parse %q as muxed ed25519 account ID\n",
				progname, args[2])
			os.Exit(1)
		}
		pk := PublicKey{ Type: PUBLIC_KEY_TYPE_ED25519 }
		copy(pk.Ed25519()[:], m.Med25519().Ed25519[:])
		fmt.Println(m.Med25519().Id, pk)
		if arg == "dump" {
			dumpXdr(&m)
		}
	default:
		usage(1)
	}
}
