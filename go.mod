module github.com/lightclient/bazooka

go 1.13

require (
	github.com/ethereum/go-ethereum v1.9.13
	github.com/mattn/go-colorable v0.1.0
	github.com/mattn/go-isatty v0.0.5-0.20180830101745-3fb116b82035
	github.com/spf13/cobra v1.0.0
)

replace github.com/ethereum/go-ethereum => github.com/lightclient/go-ethereum v1.9.13-lc1
