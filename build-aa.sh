#/bin/bash

if [ "$#" -ne 3 ]; then
    echo "usage: build-aa.sh [loops] [success] [gas_price]"
    exit
fi

printf "0x%064d%064d%064d\n" $1 $2 $3
