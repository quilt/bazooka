#/bin/bash

if [ "$#" -ne 4 ]; then
    echo "usage: build-aa.sh [nonce] [loops] [success] [gas_price]"
    exit
fi

printf "0x%064d%064d%064d%064d\n" $1 $2 $3 $4
