#!/bin/sh

abigen --abi build/deployer.abi --bin build/deployer.bin --pkg contracts --type Deployer --out deployer.go
