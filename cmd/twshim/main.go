package main

import (
	"flag"
	"fmt"
	"github.com/ngrash/twshim"
	"log"
	"os"
)

var (
	releaseFlag   = flag.String("release", "", "Tag of the desired release (overrides TWSHIM_RELEASE)")
	downloadsFlag = flag.String("downloads", "", "Target directory for executables (override TWSHIM_DOWNLOADS)")
)

func valueOrEnv(value, env string) (string, bool) {
	if value != "" {
		return value, true
	}
	if v, ok := os.LookupEnv(env); ok {
		return v, true
	}
	return "", false
}

func main() {
	l := log.New(os.Stderr, "twshim: ", 0)
	twshim.Log = l

	// Fail early if architecture is not supported.
	assetName, err := twshim.RuntimeAssetName()
	if err != nil {
		l.Fatal(err)
	}

	u := flag.Usage
	//goland:noinspection GoUnhandledErrorResult
	flag.Usage = func() {
		u()
		w := flag.CommandLine.Output()
		fmt.Fprintln(w, "See https://github.com/tailwindlabs/tailwindcss/releases for a list of releases.")
		fmt.Fprintln(w, "Use a double dash (--) to pass parameters to Tailwind CSS CLI.")
		fmt.Fprintln(w, "Example:")
		fmt.Fprintln(w, "  twshim -release v3.2.4 -downloads /tmp/tw -- -i in.css -o out.css --minify")
		fmt.Fprintln(w)
	}
	flag.Parse()

	release, ok := valueOrEnv(*releaseFlag, "TWSHIM_RELEASE")
	if !ok {
		flag.CommandLine.SetOutput(os.Stdout) // write usage to stdout to highlight error
		flag.Usage()
		l.Fatal("Insufficient parameters: Please specify release.")
	}

	downloads, ok := valueOrEnv(*downloadsFlag, "TWSHIM_DOWNLOADS")
	if !ok {
		flag.CommandLine.SetOutput(os.Stdout) // write usage to stdout to highlight error
		flag.Usage()
		l.Fatal("Insufficient parameters: Please specify downloads directory.")
	}

	cmd, err := twshim.Command(downloads, release, assetName, flag.Args()...)
	if err != nil {
		l.Fatal(err)
	}
	if err := cmd.Run(); err != nil {
		l.Fatal(err)
	}
}
