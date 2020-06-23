# bazooka

[![License](https://img.shields.io/badge/license-MIT%2FApache--2.0-blue)](https://github.com/lightclient/fast-evm)

A p2p load testing tool for Ethereum clients.

### Overview

Like a rocket launcher, `bazooka` is designed to overwhelm its target. It
carries out pre-defined load testing strategies in a deterministic fashion.
By executing atop `devp2p` it is portable across all proper Ethereum clients.
Unlike other tools [[1]](https://github.com/ethereum/hive)
[[2]](https://github.com/ethereum/retesteth) that focus on consistency &
compliance, `bazooka` focuses on creating maximally adverse operating
conditions for honest nodes. These conditions can be used to detect performance
regressions and DoS attack vectors in candidate modifications to clients.

### Usage

To initialize a patched version of Geth and run the sample routine against it,
use the following commands:

```console
$ ./run.sh build && ./run.sh run fixtures/sample.yaml
```

### Specifying an Routine

A routine has two main parts: the `initialization` and `routine`. 

#### Initialization
The first section allows you to define how many (empty) blocks the chain should
start with and the accounts to generate. The `key` of the accounts map is the
account's address. This should be derived from either the private key `key` or
using the `create2` [formula](https://crates.io/crates/create2). The depoyer's
address is `0xD2192C7F2EAEb1f05279c45D19828118e3D6f46C`.

#### Routines

There are 4 types of routines:

| id | name  | description  |
|--:|---|---|
| 0 | NewTxs   | Announces transactions to the target node  |
| 1 | NewBlock | Announces a new block to the target node |
| 2 | Sleep    | Sleeps for a certain amount of time |
| 3 | Exit     | Ends the current routine |

A script `build-aa.sh` has also been included to generate input for the sample
AA contract.

#### Examples

[sample.yaml](fixtures/sample.yaml)
[aa-sample.yaml](fixtures/aa-sample.yaml)

### Contributions

TODO
