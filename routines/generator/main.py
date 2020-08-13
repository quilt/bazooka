import sys
import os
import yaml
import peer_dos
import block_dos
from eth_utils import keccak, to_checksum_address, to_bytes
from account import Account
from transaction import Transaction, AATransaction
from helpers import create2, make_routine, make_fixture, SEND_BLOCK, SEND_TXS, SLEEP
from constants import BASE_GAS_NEW


def main():
    prefix = ""
    if len(sys.argv) == 2:
        prefix = sys.argv[1]

    #  save(block_dos.make(400000), os.path.join(prefix, "block-dos-400k.yaml"))
    #  save(block_dos.make(300000), os.path.join(prefix, "block-dos-300k.yaml"))
    #  save(block_dos.make(200000), os.path.join(prefix, "block-dos-200k.yaml"))
    #  save(block_dos.make(100000), os.path.join(prefix, "block-dos-100k.yaml"))

    save(peer_dos.make(1000, 100000), os.path.join(prefix, "peer-dos-1k-100k.yaml"))
    save(peer_dos.make(2000, 100000), os.path.join(prefix, "peer-dos-2k-100k.yaml"))
    save(peer_dos.make(3000, 100000), os.path.join(prefix, "peer-dos-3k-100k.yaml"))
    save(peer_dos.make(4000, 100000), os.path.join(prefix, "peer-dos-4k-100k.yaml"))
    save(peer_dos.make(8000, 100000), os.path.join(prefix, "peer-dos-8k-100k.yaml"))
    save(peer_dos.make(12000, 100000), os.path.join(prefix, "peer-dos-12k-100k.yaml"))
    save(peer_dos.make(16000, 100000), os.path.join(prefix, "peer-dos-16k-100k.yaml"))

    save(peer_dos.make(1000, 200000), os.path.join(prefix, "peer-dos-1k-200k.yaml"))
    save(peer_dos.make(2000, 200000), os.path.join(prefix, "peer-dos-2k-200k.yaml"))
    save(peer_dos.make(3000, 200000), os.path.join(prefix, "peer-dos-3k-200k.yaml"))
    save(peer_dos.make(4000, 200000), os.path.join(prefix, "peer-dos-4k-200k.yaml"))
    save(peer_dos.make(8000, 200000), os.path.join(prefix, "peer-dos-8k-200k.yaml"))
    save(peer_dos.make(12000, 200000), os.path.join(prefix, "peer-dos-12k-200k.yaml"))
    save(peer_dos.make(16000, 200000), os.path.join(prefix, "peer-dos-16k-200k.yaml"))

    save(peer_dos.make(1000, 300000), os.path.join(prefix, "peer-dos-1k-300k.yaml"))
    save(peer_dos.make(2000, 300000), os.path.join(prefix, "peer-dos-2k-300k.yaml"))
    save(peer_dos.make(3000, 300000), os.path.join(prefix, "peer-dos-3k-300k.yaml"))
    save(peer_dos.make(4000, 300000), os.path.join(prefix, "peer-dos-4k-300k.yaml"))
    save(peer_dos.make(8000, 300000), os.path.join(prefix, "peer-dos-8k-300k.yaml"))
    save(peer_dos.make(12000, 300000), os.path.join(prefix, "peer-dos-12k-300k.yaml"))
    save(peer_dos.make(16000, 300000), os.path.join(prefix, "peer-dos-16k-300k.yaml"))

    save(peer_dos.make(1000, 400000), os.path.join(prefix, "peer-dos-1k-400k.yaml"))
    save(peer_dos.make(2000, 400000), os.path.join(prefix, "peer-dos-2k-400k.yaml"))
    save(peer_dos.make(3000, 400000), os.path.join(prefix, "peer-dos-3k-400k.yaml"))
    save(peer_dos.make(4000, 400000), os.path.join(prefix, "peer-dos-4k-400k.yaml"))
    save(peer_dos.make(8000, 400000), os.path.join(prefix, "peer-dos-8k-400k.yaml"))
    save(peer_dos.make(12000, 400000), os.path.join(prefix, "peer-dos-12k-400k.yaml"))
    save(peer_dos.make(16000, 400000), os.path.join(prefix, "peer-dos-16k-400k.yaml"))

    #  save(peer_dos.make(1000, gl=BASE_GAS_NEW), os.path.join(prefix, "peer-dos-fast-1k.yaml"))
    #  save(peer_dos.make(2000, gl=BASE_GAS_NEW), os.path.join(prefix, "peer-dos-fast-2k.yaml"))
    #  save(peer_dos.make(3000, gl=BASE_GAS_NEW), os.path.join(prefix, "peer-dos-fast-3k.yaml"))
    #  save(peer_dos.make(4000, gl=BASE_GAS_NEW), os.path.join(prefix, "peer-dos-fast-4k.yaml"))
    #  save(peer_dos.make(8000, gl=BASE_GAS_NEW), os.path.join(prefix, "peer-dos-fast-8k.yaml"))
    #  save(peer_dos.make(12000, gl=BASE_GAS_NEW), os.path.join(prefix, "peer-dos-fast-12k.yaml"))
    #  save(peer_dos.make(16000, gl=BASE_GAS_NEW), os.path.join(prefix, "peer-dos-fast-16k.yaml"))

    #  save(peer_dos.make_normal(1000), os.path.join(prefix, "peer-dos-normal-1k.yaml"))
    #  save(peer_dos.make_normal(2000), os.path.join(prefix, "peer-dos-normal-2k.yaml"))
    #  save(peer_dos.make_normal(3000), os.path.join(prefix, "peer-dos-normal-3k.yaml"))
    #  save(peer_dos.make_normal(4000), os.path.join(prefix, "peer-dos-normal-4k.yaml"))
    #  save(peer_dos.make_normal(8000), os.path.join(prefix, "peer-dos-normal-8k.yaml"))
    #  save(peer_dos.make_normal(12000), os.path.join(prefix, "peer-dos-normal-12k.yaml"))

    #  save(peer_dos.make_valid_normal(2000), os.path.join(prefix, "peer-dos-valid-normal-2k.yaml"))
    #  save(peer_dos.make_valid_normal(3000), os.path.join(prefix, "peer-dos-valid-normal-3k.yaml"))
    #  save(peer_dos.make_valid_normal(4000), os.path.join(prefix, "peer-dos-valid-normal-4k.yaml"))
    #  save(peer_dos.make_valid_normal(8000), os.path.join(prefix, "peer-dos-valid-normal-8k.yaml"))
    #  save(peer_dos.make_valid_normal(12000), os.path.join(prefix, "peer-dos-valid-normal-12k.yaml"))
    #  save(peer_dos.make_valid_normal(16000), os.path.join(prefix, "peer-dos-valid-normal-16k.yaml"))

    #  save(peer_dos.make_multiple_normal(10, 16000), os.path.join(prefix, "peer-dos-normal-10x16k.yaml"))

    print("Done.")


def save(obj, path):
    f = open(path, "w")
    f.write(yaml.dump(obj))
    f.close()


if __name__ == "__main__":
    main()
