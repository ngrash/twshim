# twshim

twshim is a transparent Go wrapper around the standalone Tailwind CSS CLI. 
The goal of this project is to unify how developers add Tailwind CSS to their Go project by taking care of downloading the executable for the current architecture.

## Example

```
twshim -release v3.2.4 -downloads $HOME/.twshim -- -i in.css -o out.css --minify
```

## Configuration

twshim can be configured through environment variables or command line arguments:

```
  -downloads string
        Target directory for executables (override TWSHIM_DOWNLOADS)
  -release string
        Tag of the desired release (overrides TWSHIM_RELEASE)
```

A double dash (--) is required before the Tailwind CSS parameters to distinguish twshim configuration from Tailwind arguments.

twshim uses `runtime.GOOS` and `runtime.GOARCH` to decide which executable to download.

See https://github.com/tailwindlabs/tailwindcss/releases for a list of Tailwind releases.

## Usage

You can use `go run` to invoke twshim if you want to quickly execute a specific version of Tailwind CSS CLI.

```shell
go run github.com/ngrash/twshim/cmd/twshim@v0.3.0 -release v3.2.4 -downloads $HOME/.twshim/downloads -- --help
```

You can also use `go get` to add twshim to your application and use the `twshim` package from your code.

```shell
go get github.com/ngrash/twshim@v0.3.0
```