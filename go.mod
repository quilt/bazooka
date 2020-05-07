module github.com/lightclient/bazooka

go 1.13

require (
	github.com/ethereum/go-ethereum v1.9.13
	github.com/ledgerwatch/turbo-geth v0.0.0-20200401160441-9bb7f3056d4a
	github.com/mattn/go-colorable v0.1.2
	github.com/mattn/go-isatty v0.0.9
	github.com/spf13/cobra v1.0.0
	gopkg.in/yaml.v2 v2.2.4
)

replace github.com/ethereum/go-ethereum => github.com/lightclient/go-ethereum v1.9.13-lc1
