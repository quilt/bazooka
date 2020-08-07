#!/bin/sh

DATADIR=eth-sim

usage() {
	echo "usage: run.sh { build | [run ROUTINE] }"
}

init() {
	echo "Cloning quilt/geth"
	git clone -b aa-data-collection https://github.com/quilt/go-ethereum geth
	git --git-dir geth/.git --work-tree geth cherry-pick "34566fbe5d71d689cfda691c2163e31d19142542"
	git --git-dir geth/.git --work-tree geth cherry-pick "9f0d98d6da24a0a693c3c99876f8930f74d87314"
	git --git-dir geth/.git --work-tree geth cherry-pick "d24cda71dc4d3a9424058714e83c61f190b1e716"
	make -C geth
}

build() {
	make -C geth
}

run() {
	rm -rf $DATADIR

	# update the shared genesis to be something convincing and pass this check:
	# https://github.com/ethereum/go-ethereum/blob/56a319b9daa5228a6b22ecb1d07f8183ebd98106/eth/sync.go#L327
	date=$(date +%s -d '12 hours ago')
	sed -i 's/^  "timestamp": "[0-9]*",/  "timestamp": "'${date}'",/' genesis.json

	geth/build/bin/geth --datadir $DATADIR init genesis.json
	echo "INIT COMPLETE -- STARTING GETH"

	geth/build/bin/geth --datadir $DATADIR --nodiscover --fakepow --syncmode full --verbosity 5 --networkid 1337 &
	P1=$!

	go run main.go run $1 --target-data-dir=$DATADIR &
	P2=$!

	wait $P2
	kill $P1

	rm -rf $DATADIR

	echo ""
	echo "Collecting event logs..."
	./data-collection/db_populate.py

	echo ""
	echo "Done."
}

if [ ! -d geth ]
then
	init
fi

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
