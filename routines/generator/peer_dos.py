from constants import AA_CODE, BASE_GAS_NEW, DEPLOYER, LOOP_GAS
from account import Account, ContractAccount
from transaction import AATransaction, Transaction
from helpers import make_fixture, make_routine, pad_left, SEND_BLOCK, SEND_TXS, SLEEP
from random import randrange


GAS_LIMIT = 400000
LOOPS = (GAS_LIMIT - BASE_GAS_NEW) // LOOP_GAS


def make(tx_count):
    accounts = []

    for i in range(0, tx_count):
        accounts.append(ContractAccount(DEPLOYER, i, AA_CODE, 400000))

    # make tx package
    txs = []
    for a in accounts:
        tx = AATransaction(a.addr, 0, LOOPS, False, 1).as_tx(GAS_LIMIT)
        txs.append(tx)

    tx_pkg = make_routine(SEND_TXS, list(map(lambda x: x.as_obj(), txs)))

    routines = [
            tx_pkg,
    ]

    return make_fixture(accounts, routines, height=10000)


def make_normal(tx_count):
    accounts = []

    for _ in range(0, tx_count):
        a = Account(400000)

        # make tx signatures invalid
        a.pk = '0x' + pad_left(str(randrange(0, 2**160)), padder='0', chunk_size=64)

        accounts.append(a)


    # make tx package
    txs = []
    for a in accounts:
        tx = Transaction(a.addr, "0xDEADBEEF00000000000000000000000000000000", "", 0, 0, 1, 400000)
        txs.append(tx)

    tx_pkg = make_routine(SEND_TXS, list(map(lambda x: x.as_obj(), txs)))

    routines = [
            tx_pkg,
    ]

    return make_fixture(accounts, routines, height=10000)
