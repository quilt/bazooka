#!/bin/bash

DATADIR=eth-sim
CACHE=4096

usage() {
	echo "usage: run.sh { build | [run ROUTINE CACHE] }"
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

	if [[ $2 == 'false' ]]
		then
			CACHE=0
	fi

	# update the shared genesis to be something convincing and pass this check:
	# https://github.com/ethereum/go-ethereum/blob/56a319b9daa5228a6b22ecb1d07f8183ebd98106/eth/sync.go#L327
	date=$(date +%s -d '12 hours ago')
	sed -i 's/^  "timestamp": "[0-9]*",/  "timestamp": "'${date}'",/' genesis.json

	geth/build/bin/geth --datadir $DATADIR init genesis.json
	echo "INIT COMPLETE -- STARTING GETH"

	geth/build/bin/geth --datadir $DATADIR --cache=$CACHE --nodiscover --fakepow --syncmode full --verbosity 5 --networkid 1337 &
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

run_all() {
	declare -a cached_runs=(
		"block-dos-100k.yaml"
		"block-dos-200k.yaml"
		"block-dos-300k.yaml"
		"block-dos-400k.yaml"
		"peer-dos-valid-normal-8k.yaml"
	)

	declare -a nocache_runs=(
		"peer-dos-fast-1k.yaml"
		"peer-dos-fast-2k.yaml"
		"peer-dos-fast-3k.yaml"
		"peer-dos-fast-4k.yaml"
		"peer-dos-fast-8k.yaml"
		"peer-dos-fast-12k.yaml"
		"peer-dos-fast-16k.yaml"
		"peer-dos-1k-100k.yaml"
		"peer-dos-2k-100k.yaml"
		"peer-dos-3k-100k.yaml"
		"peer-dos-4k-100k.yaml"
		"peer-dos-8k-100k.yaml"
		"peer-dos-12k-100k.yaml"
		"peer-dos-16k-100k.yaml"
		"peer-dos-1k-200k.yaml"
		"peer-dos-2k-200k.yaml"
		"peer-dos-3k-200k.yaml"
		"peer-dos-4k-200k.yaml"
		"peer-dos-8k-200k.yaml"
		"peer-dos-12k-200k.yaml"
		"peer-dos-16k-200k.yaml"
		"peer-dos-1k-300k.yaml"
		"peer-dos-2k-300k.yaml"
		"peer-dos-3k-300k.yaml"
		"peer-dos-4k-300k.yaml"
		"peer-dos-8k-300k.yaml"
		"peer-dos-12k-300k.yaml"
		"peer-dos-16k-300k.yaml"
		"peer-dos-1k-400k.yaml"
		"peer-dos-2k-400k.yaml"
		"peer-dos-3k-400k.yaml"
		"peer-dos-4k-400k.yaml"
		"peer-dos-8k-400k.yaml"
		"peer-dos-12k-400k.yaml"
		"peer-dos-16k-400k.yaml"
		"peer-dos-normal-1k.yaml"
		"peer-dos-normal-2k.yaml"
		"peer-dos-normal-3k.yaml"
		"peer-dos-normal-4k.yaml"
		"peer-dos-normal-8k.yaml"
		"peer-dos-normal-12k.yaml"
		"peer-dos-normal-16k.yaml"
		"peer-dos-normal-32k.yaml"
		"peer-dos-valid-normal-8k.yaml"
		"peer-dos-normal-10x16k.yaml"
		"peer-dos-normal-10x10mb.yaml"
	)

	wget --load-cookies /tmp/cookies.txt \
		"https://docs.google.com/uc?export=download&confirm=$(wget \
		--quiet \
		--save-cookies \
		/tmp/cookies.txt \
		--keep-session-cookies \
		--no-check-certificate \
		'https://docs.google.com/uc?export=download&id=1pvepXU7p7AQ-VyQHK91-z3QtKxH_3Qrr' \
		-O- | sed \
		-rn 's/.*confirm=([0-9A-Za-z_]+).*/\1\n/p')&id=1pvepXU7p7AQ-VyQHK91-z3QtKxH_3Qrr" \
		-O ./routines/peer-dos-normal-10x10mb.yaml && rm -rf /tmp/cookies.txt

	for i in "${cached_runs[@]}"
	do
		run "./routines/$i" || break
		echo "completed run with cache: $i"
		sleep 5
	done

	for i in "${nocache_runs[@]}"
	do
		run "./routines/$i" "false" || break
		echo "completed run no cache: $i"
		sleep 5
	done

	echo ""
	echo "finished run_all"
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
		run $2 $3
		;;
	"run_all")
		run_all
		;;
	*)
		usage
		;;
esac
