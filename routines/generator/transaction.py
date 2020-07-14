from eth_utils import to_bytes
from helpers import pad_left
from constants import SEND_BLOCK, SEND_TXS, SLEEP, STOP

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
        return pad_left(to_bytes(self.nonce)) + pad_left(to_bytes(self.loops)) + pad_left(to_bytes(self.success)) + pad_left(to_bytes(self.gas_price))

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
