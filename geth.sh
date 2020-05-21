#! /bin/sh

BAZOOKA=~/development/go-workspace/src/github.com/lightclient/bazooka
GETH=~/development/go-workspace/src/github.com/ethereum/go-ethereum
SIM=~/eth/aasim

make -C $GETH
rm -r $SIM
$GETH/build/bin/geth --datadir $SIM init $BAZOOKA/genesis.json
echo "INIT COMPLETE -- STARTING GETH"
$GETH/build/bin/geth --datadir $SIM --nodiscover --fakepow --syncmode full --verbosity 5 --bootnodes "" --networkid 1337

