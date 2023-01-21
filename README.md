# twshim

twshim is a transparent Go wrapper around the standalone TailwindCSS CLI. 
The goal of this project is to unify how developers add TailwindCSS to their Go project by taking care of downloading the executable for the current architecture.

## Configuration

Because all parameters are passed to the TailwindCSS executable, twshim itself is configured through environment variables.

* `TWTAG` is the tag of the desired TailwindCSS release, e.g. `v3.2.4`.
* `TWROOT` is the directory for downloaded TailwindCSS executables, e.g. `$HOME/.twshim/downloads`.
 
twshim uses `runtime.GOOS` and `runtime.GOARCH` to decide which executable to download.

## Usage

You can use `go run` to invoke twshim if you want to quickly execute a specific version of TailwindCSS CLI.

```shell
TWTAG=v3.2.4 TWROOT=$HOME/.twshim/downloads go run github.com/ngrash/twshim/cmd/twshim@v0.2.0
```

You can also use `go get` to add twshim to your application and use the `twshim` package from your code.

```shell
go get github.com/ngrash/twshim@v0.2.0
```