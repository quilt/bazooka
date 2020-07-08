import yaml
from eth_utils import keccak, to_checksum_address, to_bytes

LOOPS = 12750

DEPLOYER = "0xD2192C7F2EAEb1f05279c45D19828118e3D6f46C"
AA_CODE = "0x6055600C60003960556000F33373ffffffffffffffffffffffffffffffffffffffff1460245736601f57005b600080fd5b60016020355b81900380602a576000358060005414604157600080fd5b600101600055604035604f57005b606035aa00"

SEND_TXS = 0
SEND_BLOCK = 1
SLEEP = 2
STOP = 3

class Account:
    addr = ""
    code = ""
    salt = 0
    balance = 0

    def __init__(self, salt, code, balance):
        self.salt = "0x" + to_bytes(salt).hex().zfill(64)
        self.code = code
        self.balance = balance
        self.addr = create2(DEPLOYER, self.salt, code)


    def as_obj(self):
        return {
            "code": self.code,
            "salt": self.salt,
            "balance": self.balance
        }


class Transaction:
    sender = ""
    to = ""
    data = ""
    nonce = 0
    amount = 0
    gas_price = 0
    gas_limit = 0

    def __init__(self, sender, to, data, nonce, amount, gp, gl):
        self.sender = sender
        self.to = to
        self.data = data
        self.nonce = nonce
        self.amount = amount
        self.gas_price = gp
        self.gas_limit = gl

    def as_obj(self):
        return {
                "from": self.sender,
                "to": self.to,
                "data": self.data,
                "nonce": self.nonce,
                "amount": self.amount,
                "gas-price": self.gas_price,
                "gas-limit": self.gas_limit,
        }


class AATransaction:
    sender = ""
    nonce = 0
    loops = 0
    success = True
    gas_price = 0

    def __init__(self, sender, nonce, loops, success, gas_price):
        self.sender = sender
        self.nonce = nonce
        self.loops = loops
        self.success = success
        self.gas_price = gas_price

    def as_bytes(self):
        return pad(to_bytes(self.nonce)) + pad(to_bytes(self.loops)) + pad(to_bytes(self.success)) + pad(to_bytes(self.gas_price))

    def as_tx(self, gl):
        return Transaction(
            "0x0000000000000000000000000000000000000000",
            self.sender,
            "0x" + self.as_bytes().hex(),
            0,
            0,
            0,
            gl
        )

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


def make_fixture(accounts, routines):
    f = {
            "initialization": {
                "height": 100,
                "accounts": {}
            },
            "routines": [
                make_routine(SLEEP, duration=2),
                make_routine(SEND_BLOCK, []),
                make_routine(SLEEP, duration=2),
                *routines,
                make_routine(SLEEP, duration=20),
                make_routine(STOP)
            ]
    }

    for a in accounts:
        f["initialization"]["accounts"][a.addr] = a.as_obj()

    return f

def main():
    accounts = []
    for i in range(0, 188):
        accounts.append(Account(i, AA_CODE, 999999999))

    # make tx package
    txs = []
    for a in accounts:
        tx = AATransaction(a.addr, 0, LOOPS, True, 1).as_tx(400000)
        txs.append(tx)

    for i in range(0, len(txs)):
        txs[i] = txs[i].as_obj()

    txs_routine = make_routine(SEND_TXS, txs)

    # make block
    block_txs = []
    for a in accounts:
        tx = AATransaction(a.addr, 0, 1, True, 1).as_tx(42500)
        block_txs.append(tx)

    for i in range(0, len(block_txs)):
        block_txs[i] = block_txs[i].as_obj()

    block_routine = make_routine(SEND_BLOCK, block_txs)


    routines = [txs_routine, make_routine(SLEEP, duration=2), block_routine]
    #  routines = [block_routine]
    print(yaml.dump(make_fixture(accounts, routines)))


# utilities
def create2(deployer, salt, code):
    salt = to_bytes(hexstr=salt)
    sender_bytes = to_bytes(hexstr=deployer)
    contract_hash = keccak(to_bytes(hexstr=code))
    raw = b"".join([b"\xFF", sender_bytes, salt, contract_hash])
    h = keccak(raw)
    address_bytes = h[12:]

    return to_checksum_address(address_bytes)


def pad(raw, padder=b"\x00", chunk_size=32):
    return padder * (chunk_size - len(raw)) + raw


def pad_right(raw, padder="0", chunk_size=64):
    return raw + padder * (chunk_size - len(raw))


if __name__ == "__main__":
    main()
