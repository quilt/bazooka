from constants import AA_CODE, BASE_GAS_NEW, DEPLOYER, LOOP_GAS
from account import Account
from transaction import AATransaction, Transaction
from helpers import make_fixture, make_routine, SEND_BLOCK, SEND_TXS, SLEEP


GAS_LIMIT = 400000
LOOPS = (GAS_LIMIT - BASE_GAS_NEW) // LOOP_GAS


def make(tx_count):
    accounts = []
    for i in range(0, tx_count):
        accounts.append(Account(DEPLOYER, i, AA_CODE, 400000))

    # make tx package
    txs = []
    for a in accounts:
        tx = AATransaction(a.addr, 0, LOOPS, True, 1).as_tx(GAS_LIMIT)
        txs.append(tx)

    tx_pkg = make_routine(SEND_TXS, list(map(lambda x: x.as_obj(), txs)))

    routines = [
            tx_pkg,
    ]

    return make_fixture(accounts, routines, height=10000)
