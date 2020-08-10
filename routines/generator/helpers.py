from eth_utils import keccak, to_checksum_address, to_bytes
from constants import SEND_BLOCK, SEND_TXS, SLEEP, STOP


def make_routine(ty, txs=[], duration=0):
    if ty == SEND_BLOCK or ty == SEND_TXS:
        return {
            "ty": ty,
            "transactions": txs
        }
    if ty == SLEEP:
        return {
            "ty": ty,
            "sleep-duration": "{}s".format(duration)
        }
    if ty == STOP:
        return {
            "ty": ty,
        }


def make_fixture(accounts, routines, height=100, wait_time=20):
    f = {
            "initialization": {
                "height": height,
                "accounts": {}
            },
            "routines": [
                make_routine(SLEEP, duration=2),
                make_routine(SEND_BLOCK, []),
                make_routine(SLEEP, duration=2),
                *routines,
                make_routine(SLEEP, duration=wait_time),
                make_routine(STOP)
            ]
    }

    for a in accounts:
        f["initialization"]["accounts"][a.addr] = a.as_obj()

    return f


def create2(deployer, salt, code):
    salt = to_bytes(hexstr=salt)
    sender_bytes = to_bytes(hexstr=deployer)
    contract_hash = keccak(to_bytes(hexstr=code))
    raw = b"".join([b"\xFF", sender_bytes, salt, contract_hash])
    h = keccak(raw)
    address_bytes = h[12:]

    return to_checksum_address(address_bytes)


def pad_left(raw, padder=b"\x00", chunk_size=32):
    return padder * (chunk_size - len(raw)) + raw
