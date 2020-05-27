#!/bin/sh


GETH=$(pwd)/geth
SIM=$(pwd)/.eth-sim


usage() {
	echo "usage: run.sh { build | [run ROUTINE] }"
}

init() {
	echo "Cloning quilt/geth"
	git clone -b account-abstraction https://github.com/quilt/go-ethereum geth
	make -C $GETH
}

build() {
	make -C $GETH
}

run() {
	rm -r $SIM

	$GETH/build/bin/geth --datadir $SIM init $(pwd)/genesis.json
	echo "INIT COMPLETE -- STARTING GETH"

	$GETH/build/bin/geth --datadir $SIM --nodiscover --fakepow --syncmode full --verbosity 5 --networkid 1337 &
	P1=$!

	go run main.go run $1 &
	P2=$!

	wait $P2
	kill $P1

	rm -r $SIM

	echo ""
	echo "Done."
}

rm -r $SIM 2> /dev/null


# check if quilt/geth already exists
if [ ! -d $GETH ]
then
	init
fi

# run script
case $1 in 
	"build")
		build
		;;
	"run")
		run $2
		;;
	*)
		usage
		;;
esac
