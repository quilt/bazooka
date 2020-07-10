from main import AA_CODE, BLOCK_LIMIT, DEPLOYER
from account import Account
from transaction import AATransaction, Transaction
from helpers import make_fixture, make_routine, SEND_BLOCK, SEND_TXS, SLEEP

BASE_GAS = 27451
LOOP_GAS = 26

def make(gas_limit, block_limit=BLOCK_LIMIT):
    loops = (gas_limit - BASE_GAS) // LOOP_GAS
    max_txs = block_limit // (BASE_GAS + LOOP_GAS)

    accounts = []
    for i in range(0, max_txs):
        accounts.append(Account(DEPLOYER, i, AA_CODE, 999999999))

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
        tx = AATransaction(a.addr, 1, loops, True, 1).as_tx(gas_limit)
        txs.append(tx)

    tx_pkg = make_routine(SEND_TXS, list(map(lambda x: x.as_obj(), txs)))

    routines = [
            init1,
            init2,
            make_routine(SLEEP, duration=2),
            tx_pkg,
            make_routine(SLEEP, duration=2),
            invalidator
    ]

    return make_fixture(accounts, routines)
