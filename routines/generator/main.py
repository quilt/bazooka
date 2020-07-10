import sys
import os
import yaml
import block_dos
from eth_utils import keccak, to_checksum_address, to_bytes
from account import Account
from transaction import Transaction, AATransaction
from helpers import create2, make_routine, make_fixture, SEND_BLOCK, SEND_TXS, SLEEP

DEPLOYER = "0xD2192C7F2EAEb1f05279c45D19828118e3D6f46C"
AA_CODE = "0x6055600C60003960556000F33373ffffffffffffffffffffffffffffffffffffffff1460245736601f57005b600080fd5b60016020355b81900380602a576000358060005414604157600080fd5b600101600055604035604f57005b606035aa00"
BLOCK_LIMIT = 8000000

def main():
    prefix = ""
    if len(sys.argv) == 2:
        prefix = sys.argv[1]

    save(block_dos.make(400000), os.path.join(prefix, "block-dos-400k.yaml"))
    save(block_dos.make(300000), os.path.join(prefix, "block-dos-300k.yaml"))
    save(block_dos.make(200000), os.path.join(prefix, "block-dos-200k.yaml"))
    save(block_dos.make(100000), os.path.join(prefix, "block-dos-100k.yaml"))

    print("Done.")


def save(obj, path):
    f = open(path, "w")
    f.write(yaml.dump(obj))
    f.close()


if __name__ == "__main__":
    main()
