// Готує демо-проєкт для перевірки розгортання контракту: шари, згенерована колекція,
// whitelist. Власник - гаманець з імітації MetaMask (тестовий акаунт №1 Anvil), тож
// проєкт видно в браузері з цією імітацією після входу.
//
//   node scripts/seed-demo.js
//   PRIVATE_KEY=0x... PROJECT_NAME="..." NFT_COUNT=10 BASE=... node scripts/seed-demo.js
//   ART=data/demo-art PROJECT_NAME="Cosmic Critters" NFT_COUNT=100 node scripts/seed-demo.js
//
// ART - каталог процедурного набору шарів з вагами (scripts/demo-art/render.js); без нього - три прості шари.
//
// Якщо проєкт з такою назвою вже є, скрипт лише виводить його ідентифікатор.

const { ethers, Platform } = require('./lib/platform');

// Публічно відомі тестові ключі Anvil (лише для локальної розробки).
const ANVIL_KEYS = [
    '0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d', // №1
    '0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a', // №2
    '0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6', // №3
];

const owner = new ethers.Wallet(process.env.PRIVATE_KEY || ANVIL_KEYS[0]);
const projectName = process.env.PROJECT_NAME || 'Demo Collection';
const nftCount = Number(process.env.NFT_COUNT || 10);

(async () => {
    const platform = new Platform(process.env.BASE);
    console.log(`Платформа: ${platform.base}`);
    await platform.login(owner);

    const existing = await platform.findProject(projectName);
    if (existing) {
        console.log(`\nПроєкт "${projectName}" уже є: ${existing.id}`);
        return;
    }

    const project = await platform.createProject(projectName, nftCount);
    if (process.env.ART) await platform.uploadArt(project.id, process.env.ART);
    else await platform.uploadLayers(project.id);
    await platform.generate(project.id, nftCount);
    // Власник проєкту завжди у whitelist, далі - тестові акаунти Anvil (разом три адреси).
    const whitelist = [{ name: 'Owner', walletAddress: owner.address }];
    for (const [i, key] of ANVIL_KEYS.entries()) {
        const walletAddress = new ethers.Wallet(key).address;
        if (whitelist.length < 3 && walletAddress !== owner.address) whitelist.push({ name: `Anvil #${i + 1}`, walletAddress });
    }
    await platform.addToWhitelist(project.id, whitelist);

    console.log(`\nГОТОВО. Проєкт "${projectName}" ${project.id}, власник ${owner.address}`);
    console.log(`Сторінка контракту: ${platform.base}/create-contract/${project.id}`);
})().catch((err) => {
    console.error(`ПОМИЛКА: ${err.message}`);
    process.exit(1);
});
