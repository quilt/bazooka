from account import ContractAccount
from transaction import AATransaction, Transaction
from helpers import make_fixture, make_routine, SEND_BLOCK, SEND_TXS, SLEEP
from constants import AA_CODE, BLOCK_LIMIT, DEPLOYER
import time

BASE_GAS = 21470
LOOP_GAS = 26

def make(gas_limit, block_limit=BLOCK_LIMIT):
    loops = (gas_limit - BASE_GAS) // LOOP_GAS
    max_txs = block_limit // (BASE_GAS + LOOP_GAS)

    accounts = []
    for i in range(0, max_txs):
        accounts.append(ContractAccount(DEPLOYER, i, AA_CODE, 999999999))

    # make blocks to set nonce to non-zero value
    txs = []
    init_blocks = []
    for a in accounts:
        tx = AATransaction(a.addr, 0, 1, True, 1).as_tx(36500)
        txs.append(tx)

    for i in range(0, len(txs), 200):
        init_blocks.append(make_routine(SEND_BLOCK, list(map(lambda x: x.as_obj(), txs[i:i+200]))))

    # make block to invalidate mempool
    txs = []
    for a in accounts:
        tx = AATransaction(a.addr, 1, 1, True, 1).as_tx(21500)
        txs.append(tx)

    invalidator = make_routine(SEND_BLOCK, list(map(lambda x: x.as_obj(), txs)))

    # make tx package
    txs = []
    for a in accounts:
        tx = AATransaction(a.addr, 1, loops, True, 1).as_tx(gas_limit)
        txs.append(tx)

    tx_pkg = make_routine(SEND_TXS, list(map(lambda x: x.as_obj(), txs)))

    routines = init_blocks
    routines.append(make_routine(SLEEP, duration=2))
    routines.append(tx_pkg)
    routines.append(make_routine(SLEEP, duration=2))
    routines.append(invalidator)

    return make_fixture(accounts, routines)
