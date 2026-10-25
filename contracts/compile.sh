#!/bin/bash
solc --abi --bin ./contracts/CollectionNFT.sol -o ./contracts/bin --base-path $1 --include-path $1/node_modules/ --overwrite
