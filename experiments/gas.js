// Експеримент E5 (розділ 3 роботи): газ whitelist - корінь дерева Меркла (CollectionNFT)
// проти переліку адрес у mapping (experiments/gas/CollectionNFTMapping.sol, решта коду та сама).
// Для whitelist з N = 10 / 100 / 1000 / 10000 адрес вимірюються газ запису whitelist у контракт
// і газ mintWhitelist одного токена (середнє для трьох покупців, однаковий порядок в обох варіантах).
// Газ детермінований, тому кожне значення вимірюється один раз. Мережа - тимчасовий Anvil на порті 8546
// (правила основної мережі Ethereum; збережений стан make chain не змінюється).
//
//   node experiments/gas.js     # записує experiments/results/gas.csv

const fs = require('fs');
const path = require('path');
const { execFileSync } = require('child_process');
const { ethers } = require('../scripts/lib/platform');

const ROOT = path.join(__dirname, '..');
const OUT = path.join(__dirname, 'results', 'gas.csv');
const SIZES = [10, 100, 1000, 10000];
const BATCH = 500; // адрес за транзакцію addToWhitelist: ~11 млн газу при ліміті блоку 30 млн
const PRICE = ethers.utils.parseEther('0.001');
// умови продажу в конструкторі: розмір колекції й ліміт за транзакцію (покупці мінтять по одному токену)
const MAX_SUPPLY = 10000;
const MAX_PER_MINT = 5;
const ANVIL = 'generator-experiment-anvil';
const RPC = 'http://127.0.0.1:8546';
// перші акаунти Anvil (мнемоніка за замовчуванням): №0 - власник, №1-3 - покупці
const MNEMONIC = 'test test test test test test test test test test test junk';

function docker(...args) {
    return execFileSync('docker', args, { encoding: 'utf8', maxBuffer: 64 * 2 ** 20 });
}

function compileMapping() {
    const output = docker('run', '--rm', '-v', `${ROOT}/contracts:/sources:ro`, '-v', `${__dirname}/gas:/exp:ro`,
        'ethereum/solc:0.8.24', '--optimize', '--optimize-runs', '200', '--base-path', '/exp',
        '--include-path', '/sources/node_modules', '--combined-json', 'abi,bin', '/exp/CollectionNFTMapping.sol');
    const contract = Object.entries(JSON.parse(output).contracts).find(([name]) => name.endsWith(':CollectionNFTMapping'))[1];
    return { abi: typeof contract.abi === 'string' ? JSON.parse(contract.abi) : contract.abi, bin: contract.bin };
}

/** Дерево Меркла, сумісне з MerkleProof OpenZeppelin: листок - keccak256(адреса), пари сортуються. */
function merkleTree(addresses) {
    const hashPair = (a, b) => ethers.utils.keccak256(a < b ? ethers.utils.concat([a, b]) : ethers.utils.concat([b, a]));
    const layers = [addresses.map((a) => ethers.utils.keccak256(a))];
    while (layers[layers.length - 1].length > 1) {
        const layer = layers[layers.length - 1];
        const next = [];
        for (let i = 0; i < layer.length; i += 2) next.push(i + 1 < layer.length ? hashPair(layer[i], layer[i + 1]) : layer[i]);
        layers.push(next);
    }
    const proof = (index) => layers.slice(0, -1).flatMap((layer, level) => {
        const sibling = (index >> level) ^ 1;
        return sibling < layer.length ? [layer[sibling]] : [];
    });
    return { root: layers[layers.length - 1][0], proof };
}

async function gasOf(txPromise) {
    const receipt = await (await txPromise).wait();
    return receipt.gasUsed.toNumber();
}

(async () => {
    docker('run', '-d', '--rm', '--name', ANVIL, '-p', '8546:8545', 'ghcr.io/foundry-rs/foundry:latest', 'anvil --host 0.0.0.0');
    try {
        const provider = new ethers.providers.JsonRpcProvider(RPC);
        for (let i = 0; i < 50; i++) {
            try { await provider.getBlockNumber(); break; } catch { await new Promise((r) => setTimeout(r, 200)); }
        }
        const accounts = [0, 1, 2, 3].map((i) => ethers.Wallet.fromMnemonic(MNEMONIC, `m/44'/60'/0'/0/${i}`).connect(provider));
        const [owner, ...buyers] = accounts;

        const merkleAbi = JSON.parse(fs.readFileSync(path.join(ROOT, 'contracts/bin/CollectionNFT.abi'), 'utf8'));
        const merkleBin = fs.readFileSync(path.join(ROOT, 'contracts/bin/CollectionNFT.bin'), 'utf8').trim();
        const mapping = compileMapping();

        const rows = ['whitelist_size,variant,setup_gas,setup_txs,mint_gas_mean,proof_length,network'];
        for (const size of SIZES) {
            const others = Array.from({ length: size - buyers.length }, () => ethers.Wallet.createRandom().address);
            const whitelist = [...buyers.map((b) => b.address), ...others];

            // Корінь Меркла: одна транзакція setMerkleRoot після розгортання з нульовим коренем.
            const merkle = await new ethers.ContractFactory(merkleAbi, merkleBin, owner)
                .deploy('http://localhost/metadata/', ethers.constants.HashZero, 'Experiment', MAX_SUPPLY, PRICE, MAX_PER_MINT);
            await merkle.deployed();
            const tree = merkleTree(whitelist);
            const merkleSetup = await gasOf(merkle.setMerkleRoot(tree.root));
            const merkleMints = [];
            for (const [i, buyer] of buyers.entries()) {
                merkleMints.push(await gasOf(merkle.connect(buyer).mintWhitelist(tree.proof(i), 1, { value: PRICE })));
            }

            // mapping: адреси пакетами по BATCH у транзакціях addToWhitelist.
            const listed = await new ethers.ContractFactory(mapping.abi, mapping.bin, owner)
                .deploy('http://localhost/metadata/', 'Experiment', MAX_SUPPLY, PRICE, MAX_PER_MINT);
            await listed.deployed();
            let mappingSetup = 0;
            let txs = 0;
            for (let i = 0; i < whitelist.length; i += BATCH) {
                mappingSetup += await gasOf(listed.addToWhitelist(whitelist.slice(i, i + BATCH), { gasLimit: 29_000_000 }));
                txs++;
            }
            const mappingMints = [];
            for (const buyer of buyers) {
                mappingMints.push(await gasOf(listed.connect(buyer).mintWhitelist(1, { value: PRICE })));
            }

            const mean = (values) => Math.round(values.reduce((a, b) => a + b, 0) / values.length);
            const proofLength = tree.proof(0).length;
            rows.push(`${size},merkle,${merkleSetup},1,${mean(merkleMints)},${proofLength},anvil`);
            rows.push(`${size},mapping,${mappingSetup},${txs},${mean(mappingMints)},0,anvil`);
            console.log(`N=${size}: Меркл - запис ${merkleSetup}, мінт ${mean(merkleMints)} (доказ ${proofLength}); `
                + `mapping - запис ${mappingSetup} (${txs} тр.), мінт ${mean(mappingMints)}`);
        }
        fs.writeFileSync(OUT, `${rows.join('\n')}\n`);
        console.log(`-> ${path.relative(ROOT, OUT)}`);
    } finally {
        docker('rm', '-f', ANVIL);
    }
})().catch((err) => {
    console.error(`ПОМИЛКА: ${err.message}`);
    process.exit(1);
});
