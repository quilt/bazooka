import yaml
from eth_utils import keccak, to_checksum_address, to_bytes
from account import Account
from transaction import Transaction, AATransaction
from helpers import make_routine, make_fixture, SEND_BLOCK, SEND_TXS, SLEEP

LOOPS = 12750

DEPLOYER = "0xD2192C7F2EAEb1f05279c45D19828118e3D6f46C"
AA_CODE = "0x6055600C60003960556000F33373ffffffffffffffffffffffffffffffffffffffff1460245736601f57005b600080fd5b60016020355b81900380602a576000358060005414604157600080fd5b600101600055604035604f57005b606035aa00"


def main():
    accounts = []
    for i in range(0, 292):
        accounts.append(Account(i, AA_CODE, 999999999))

    # make blocks to set nonce to non-zero value
    txs = []
    for a in accounts:
        tx = AATransaction(a.addr, 0, 1, True, 1).as_tx(42500)
        txs.append(tx)

    init1 = make_routine(SEND_BLOCK, list(map(lambda x: x.as_obj(), txs[0:188])))
    init2 = make_routine(SEND_BLOCK, list(map(lambda x: x.as_obj(), txs[188:])))

    # make block to invalidate mempool
    txs = []
    for a in accounts:
        tx = AATransaction(a.addr, 1, 1, True, 1).as_tx(27500)
        txs.append(tx)

    invalidator = make_routine(SEND_BLOCK, list(map(lambda x: x.as_obj(), txs)))

    # make tx package
    txs = []
    for a in accounts:
        tx = AATransaction(a.addr, 1, LOOPS, True, 1).as_tx(400000)
        txs.append(tx)

    tx_pkg = make_routine(SEND_TXS, list(map(lambda x: x.as_obj(), txs)))

    routines = [
            init1,
            init2,
            make_routine(SLEEP, duration=2),
            tx_pkg,
            invalidator
    ]

    print(yaml.dump(make_fixture(accounts, routines)))


if __name__ == "__main__":
    main()
