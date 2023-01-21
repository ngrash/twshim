package main

import (
	"github.com/ngrash/twshim"
	"log"
	"os"
)

func main() {
	l := log.New(os.Stderr, "twshim: ", 0)
	twshim.Log = l

	// Fail early if architecture is not supported.
	assetName, err := twshim.RuntimeAssetName()
	if err != nil {
		l.Fatal(err)
	}

	tag, found := os.LookupEnv("TWTAG")
	if !found {
		l.Fatal("Set TWTAG to select a release.\n\nSee https://github.com/tailwindlabs/tailwindcss/releases for available releases.")
	}

	root, found := os.LookupEnv("TWROOT")
	if !found {
		l.Fatal("Set TWROOT to configure a download directory.")
	}

	cmd, err := twshim.Command(root, tag, assetName, os.Args[1:]...)
	if err != nil {
		l.Fatal(err)
	}

	if err := cmd.Run(); err != nil {
		l.Fatal(err)
	}
}
