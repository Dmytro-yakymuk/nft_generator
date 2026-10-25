"""Writes <Name>.abi and <Name>.bin from `solc --combined-json abi,bin` output.

Usage: python3 scripts/split-solc-output.py <combined.json> <ContractName> <out dir>
"""
import json
import os
import sys

combined, name, out_dir = sys.argv[1:4]
contracts = json.load(open(combined))["contracts"]
key = next(k for k in contracts if k.endswith(":" + name))
abi = contracts[key]["abi"]
with open(os.path.join(out_dir, name + ".abi"), "w") as f:
    f.write(abi if isinstance(abi, str) else json.dumps(abi))
with open(os.path.join(out_dir, name + ".bin"), "w") as f:
    f.write(contracts[key]["bin"])
print(f"{name}: abi {len(abi)} entries, bin {len(contracts[key]['bin']) // 2} bytes")
