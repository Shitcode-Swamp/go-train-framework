![gopher-train](./gopher-train.png)

# go-train-framework

[![Latest release](https://img.shields.io/github/v/release/Shitcode-Swamp/go-train-framework)](https://github.com/Shitcode-Swamp/go-train-framework/releases/latest)

Custom framework made in edutcational puposes.

## Overall description

Framework automaticaly builds AppCtx with needed dependencies
To use the framework build tool you might run binary via
`chmod +x ./build/go-train && ./build/go-train` or `run.sh`

You must describe what directory to watch via flag `--root=path`

How the framework is built is descrived in the `structure.excalidraw`

## Using it as a Go tool in other repos

The command lives at `cmd/go-train` so it can be installed like any other Go tool.

**Global install** (puts `go-train` on your `$GOBIN`/`$GOPATH/bin`):

```sh
go install github.com/Shitcode-Swamp/go-train/cmd/go-train@latest
go-train -root=./path/to/watch -ctxPath=./path/to/output
```

**As a per-module tool dependency** (Go 1.24+, pins the version in the consuming
repo's `go.mod` via a `tool` directive):

```sh
go get -tool github.com/Shitcode-Swamp/go-train/cmd/go-train@latest
go tool go-train -root=./path/to/watch -ctxPath=./path/to/output
```

While developing locally without a published module, point a consuming repo's
`go.mod` at your local checkout instead:

```go
replace github.com/Shitcode-Swamp/go-train => /path/to/go-train-framework
```
