// Збирає init-скрипт з імітацією MetaMask для наскрізних тестів і перевірок у браузері без розширення.
// Імітація підписує повідомлення й транзакції тестовим ключем, а решту RPC-запитів
// пересилає на вузол (за замовчуванням - Anvil з `make chain`).
//
//   node scripts/build-metamask-mock.js [вихідний файл]
//   MOCK_RPC_URL=... MOCK_PRIVATE_KEY=... MOCK_CHAIN_ID=0x7a69 node scripts/build-metamask-mock.js
//
// Ключ за замовчуванням - тестовий акаунт №1 Anvil (публічно відомий, лише для розробки).

const fs = require('fs');
const path = require('path');

const out = path.resolve(process.argv[2] || path.join(__dirname, '../tests/e2e/.mock/metamask-mock.js'));
const settings = {
    rpcUrl: process.env.MOCK_RPC_URL || 'http://localhost:8545',
    privateKey: process.env.MOCK_PRIVATE_KEY || '0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d',
    chainId: process.env.MOCK_CHAIN_ID || '0x7a69',
};

const ethersPath = path.join(__dirname, '../web/app/node_modules/ethers/dist/ethers.umd.min.js');
if (!fs.existsSync(ethersPath)) {
    console.error('Не знайдено ethers: виконайте cd web/app && npm ci --legacy-peer-deps');
    process.exit(1);
}

function mock(settings) {
    if (window.ethereum || window.top !== window) return;

    const { ethers } = window;
    const wallet = new ethers.Wallet(settings.privateKey);
    const listeners = {};
    let requestId = 0;

    async function rpc(method, params) {
        const res = await fetch(settings.rpcUrl, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ jsonrpc: '2.0', id: ++requestId, method, params: params || [] }),
        });
        const body = await res.json();
        if (body.error) throw Object.assign(new Error(body.error.message), { code: body.error.code, data: body.error.data });
        return body.result;
    }

    function emit(event, value) {
        (listeners[event] || []).forEach((listener) => listener(value));
    }

    function toText(hex) {
        if (!/^0x[0-9a-f]*$/i.test(hex)) return hex;
        return ethers.utils.toUtf8String(hex);
    }

    const ethereum = {
        isMetaMask: true,
        chainId: settings.chainId,
        networkVersion: String(parseInt(settings.chainId, 16)),
        selectedAddress: wallet.address,
        isConnected: () => true,
        on(event, listener) {
            (listeners[event] = listeners[event] || []).push(listener);
            return ethereum;
        },
        removeListener(event, listener) {
            listeners[event] = (listeners[event] || []).filter((l) => l !== listener);
            return ethereum;
        },
        async request({ method, params = [] }) {
            console.log(`[MockMetaMask] ${method}`);
            switch (method) {
                case 'eth_requestAccounts':
                case 'eth_accounts':
                    return [wallet.address];
                case 'eth_chainId':
                    return ethereum.chainId;
                case 'net_version':
                    return ethereum.networkVersion;
                case 'personal_sign':
                    return wallet.signMessage(toText(params[0]));
                case 'eth_sign':
                    return wallet.signMessage(ethers.utils.arrayify(params[1]));
                case 'wallet_switchEthereumChain': {
                    const target = params[0].chainId.toLowerCase();
                    if (target !== settings.chainId.toLowerCase()) {
                        // Як MetaMask для мережі, якої немає в гаманці.
                        throw Object.assign(new Error(`Unrecognized chain ID "${target}" (mock is bound to ${settings.chainId})`), { code: 4902 });
                    }
                    return null;
                }
                case 'wallet_addEthereumChain':
                    return null;
                case 'eth_sendTransaction': {
                    const tx = params[0];
                    // Надбавку валідатору бере з мережі, як MetaMask; ethers 5 інакше ставить 2.5 gwei.
                    const [priority, block] = await Promise.all([
                        rpc('eth_maxPriorityFeePerGas'),
                        rpc('eth_getBlockByNumber', ['latest', false]),
                    ]);
                    const maxPriorityFeePerGas = ethers.BigNumber.from(tx.maxPriorityFeePerGas || priority);
                    const maxFeePerGas = tx.maxFeePerGas
                        ? ethers.BigNumber.from(tx.maxFeePerGas)
                        : ethers.BigNumber.from(block.baseFeePerGas).mul(2).add(maxPriorityFeePerGas);
                    const signer = wallet.connect(new ethers.providers.JsonRpcProvider(settings.rpcUrl));
                    const sent = await signer.sendTransaction({
                        to: tx.to || undefined,
                        data: tx.data,
                        value: tx.value || undefined,
                        gasLimit: tx.gas || tx.gasLimit || undefined,
                        maxPriorityFeePerGas,
                        maxFeePerGas,
                    });
                    return sent.hash;
                }
                default:
                    return rpc(method, params);
            }
        },
    };
    // Застарілі методи, якими користуються старі бібліотеки.
    ethereum.enable = () => ethereum.request({ method: 'eth_requestAccounts' });
    ethereum.send = (method, params) => ethereum.request({ method, params });

    window.ethereum = ethereum;
    window.dispatchEvent(new Event('ethereum#initialized'));
    console.log(`[MockMetaMask] ${wallet.address} on chain ${settings.chainId} via ${settings.rpcUrl}`);
}

const script = [
    '// Згенеровано scripts/build-metamask-mock.js - не редагувати вручну.',
    '(function () {',
    'if (window.top !== window || window.ethereum) return;',
    fs.readFileSync(ethersPath, 'utf-8'),
    '}).call(window);',
    `(${mock.toString()})(${JSON.stringify(settings)});`,
    '',
].join('\n');

fs.mkdirSync(path.dirname(out), { recursive: true });
fs.writeFileSync(out, script);
console.log(`Записано ${out} (${(script.length / 1024).toFixed(0)} KB): ${settings.chainId} через ${settings.rpcUrl}`);
