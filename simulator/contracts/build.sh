#!/bin/sh

abigen --abi build/deployer.abi --bin build/deployer.bin --pkg contract --type Deployer --out deployer.go
