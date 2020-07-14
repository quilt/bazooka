import sys
import os
import yaml
import peer_dos
import block_dos
from eth_utils import keccak, to_checksum_address, to_bytes
from account import Account
from transaction import Transaction, AATransaction
from helpers import create2, make_routine, make_fixture, SEND_BLOCK, SEND_TXS, SLEEP


def main():
    prefix = ""
    if len(sys.argv) == 2:
        prefix = sys.argv[1]

    #  save(block_dos.make(400000), os.path.join(prefix, "block-dos-400k.yaml"))
    #  save(block_dos.make(300000), os.path.join(prefix, "block-dos-300k.yaml"))
    #  save(block_dos.make(200000), os.path.join(prefix, "block-dos-200k.yaml"))
    #  save(block_dos.make(100000), os.path.join(prefix, "block-dos-100k.yaml"))

    save(peer_dos.make(1000), os.path.join(prefix, "peer-dos-1k.yaml"))
    save(peer_dos.make(2000), os.path.join(prefix, "peer-dos-2k.yaml"))
    save(peer_dos.make(3000), os.path.join(prefix, "peer-dos-3k.yaml"))
    save(peer_dos.make(4000), os.path.join(prefix, "peer-dos-4k.yaml"))

    print("Done.")


def save(obj, path):
    f = open(path, "w")
    f.write(yaml.dump(obj))
    f.close()


if __name__ == "__main__":
    main()
