# bazooka

[![License](https://img.shields.io/badge/license-MIT%2FApache--2.0-blue)](https://github.com/lightclient/fast-evm)

An attack orchestration tool targeting Ethereum clients.

### Overview

Like a rocket launcher, `bazooka` is designed to overwhelm its target. It
carries out pre-defined attack strategies in a deterministic fashion.
By executing atop `devp2p` it is portable across all proper Ethereum clients.
Unlike other tools [[1]](https://github.com/ethereum/hive)
[[2]](https://github.com/ethereum/retesteth) that focus on consistency &
compliance, `bazooka` focuses on creating maximally adverse operating
conditions for honest nodes. These conditions can be used to detect performance
regressions and DoS attack vectors in candidate modifications to clients.

### Usage

A patched version of `geth` is used to bypass PoW verification. To download and
build it from source, run the following:

```console
$ git clone git@github.com:lightclient/go-ethereum.git && cd go-ethereum && make
$ mv build/bin/geth build/bin/bgeth
$ export PATH="build/bin":$PATH
```

To initialize & start the patched version of `geth`:

```console
$ bgeth --datadir ~/.eth/bazooka init genesis.json
$ bgeth --datadir ~/.eth/bazooka --nodiscover --fakepow --syncmode full --verbosity 5 --bootnodes "" --networkid 1337
```

Finally, to begin an attack from this repository:

```console
$ bazooka run {attack yaml}
```

### Specifying an Attack

TODO

### Contributions

TODO
