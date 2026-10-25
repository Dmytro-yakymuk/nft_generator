// Перевірка розгорнутої платформи через API: вхід підписом гаманця -> проєкт -> шари ->
// генерація -> метадані токенів за адресою, яку записано в контракт (base URI) -> публікація в IPFS.
//
//   node scripts/smoke.js                       # http://localhost:9012
//   BASE=https://<домен>.ngrok-free.app node scripts/smoke.js

const { ethers, Platform } = require('./lib/platform');

const NFT_COUNT = 5;

(async () => {
    const platform = new Platform(process.env.BASE);
    console.log(`Платформа: ${platform.base}`);
    await platform.call('головна сторінка', 'GET', '/');

    await platform.login(ethers.Wallet.createRandom());
    const project = await platform.createProject(`Smoke ${new Date().toISOString()}`, NFT_COUNT);
    await platform.uploadLayers(project.id);

    // Ваги: шар Eyes (2 варіанти) - необов'язковий; некоректні ваги відхиляються.
    const layers = await platform.call('список шарів з вагами', 'GET', `/api/v0/collections/${project.id}/layers`);
    const eyes = layers.find((layer) => layer.name === 'Eyes');
    if (eyes.noneWeight !== 0 || eyes.weights['v1.png'] !== 100) throw new Error(`ваги за замовчуванням: ${JSON.stringify(eyes)}`);
    const weightsUrl = `/api/v0/collections/${project.id}/layers/${eyes.id}/weights`;
    await platform.call('вага 0 -> 400', 'PUT', weightsUrl, { noneWeight: 0, weights: { 'v1.png': 0 } }, [400]);
    await platform.call('вага неіснуючого зображення -> 400', 'PUT', weightsUrl, { noneWeight: 0, weights: { 'none.png': 5 } }, [400]);
    await platform.call('ваги й необов\'язковий шар', 'PUT', weightsUrl, { noneWeight: 50, weights: { 'v1.png': 30, 'v2.png': 20 } });
    const updated = (await platform.call('список шарів', 'GET', `/api/v0/collections/${project.id}/layers`)).find((layer) => layer.name === 'Eyes');
    if (updated.noneWeight !== 50 || updated.weights['v2.png'] !== 20) throw new Error(`ваги не збережено: ${JSON.stringify(updated)}`);

    await platform.generate(project.id, NFT_COUNT);

    for (let tokenId = 0; tokenId < NFT_COUNT; tokenId++) {
        const metadata = await platform.call(`метадані токена ${tokenId}`, 'GET', `/metadata/${project.id}/${tokenId}`);
        const image = await fetch(metadata.image, { headers: { 'ngrok-skip-browser-warning': '1' } });
        if (image.status !== 200) throw new Error(`зображення ${metadata.image}: HTTP ${image.status}`);
        if (tokenId === 0) console.log(JSON.stringify(metadata, null, 2));
        // атрибути - у порядку шарів, комбінація - у dna.
        const traits = metadata.attributes.map((attribute) => attribute.trait_type).join(',');
        if (traits !== 'Background,Body,Eyes' || !/^[0-9]+$/.test(metadata.dna)) throw new Error(`метадані токена ${tokenId}: ${traits}, dna ${metadata.dna}`);
    }
    console.log('ok   зображення токенів доступні');
    await platform.call('неіснуючий токен -> 404', 'GET', `/metadata/${project.id}/${NFT_COUNT}`, undefined, [404]);
    await platform.call('дашборд без розгорнутого контракту -> 409', 'GET', `/api/v0/dashboard/${project.id}/stats`, undefined, [409]);

    // IPFS: публікація дає CID каталогів зображень і метаданих, повторна - ті самі CID (адресація за вмістом),
    // контракт розгортається з базовим URI ipfs://<CID метаданих>/. Шлюз перевіряється лише локально
    // або коли задано GATEWAY.
    const published = await platform.call('публікація в IPFS', 'POST', `/api/v0/collections/${project.id}/ipfs`);
    const again = await platform.call('повторна публікація в IPFS', 'POST', `/api/v0/collections/${project.id}/ipfs`);
    if (again.imagesCid !== published.imagesCid || again.metadataCid !== published.metadataCid) {
        throw new Error(`повторна публікація дала інші CID: ${JSON.stringify([published, again])}`);
    }
    console.log(`ok   CID зображень ${published.imagesCid}, метаданих ${published.metadataCid}`);
    const gateway = process.env.GATEWAY || (process.env.BASE ? '' : 'http://localhost:8080');
    if (gateway) {
        const tokenId = NFT_COUNT - 1;
        const ipfsMetadata = await (await fetch(`${gateway}/ipfs/${published.metadataCid}/${tokenId}`)).json();
        if (ipfsMetadata.image !== `ipfs://${published.imagesCid}/${tokenId}.png`) throw new Error(`image в IPFS: ${ipfsMetadata.image}`);
        const picture = await fetch(`${gateway}/ipfs/${published.imagesCid}/${tokenId}.png`);
        if (picture.status !== 200) throw new Error(`зображення через шлюз: HTTP ${picture.status}`);
        console.log(`ok   метадані й зображення токена ${tokenId} через шлюз IPFS`);
    }
    const deployData = await platform.call('дані розгортання', 'GET', `/api/v0/contracts/data/${project.id}`);
    if (!deployData.includes(Buffer.from(`ipfs://${published.metadataCid}/`).toString('hex'))) {
        throw new Error('у даних розгортання немає базового URI ipfs://');
    }
    console.log('ok   базовий URI контракту - ipfs://<CID метаданих>/');

    // Умови продажу йдуть у конструктор: максимум - кількість NFT колекції, ціна у wei, ліміт за транзакцію.
    const word = (value) => ethers.utils.hexZeroPad(ethers.BigNumber.from(value).toHexString(), 32).slice(2);
    const termsData = await platform.call('дані розгортання з умовами продажу', 'GET',
        `/api/v0/contracts/data/${project.id}?price=0.02&maxPerMint=2`);
    if (!termsData.includes(word(NFT_COUNT) + word(ethers.utils.parseEther('0.02')) + word(2))) {
        throw new Error('у даних розгортання немає умов продажу (кількість, ціна, ліміт)');
    }
    console.log(`ok   умови продажу в конструкторі: ${NFT_COUNT} токенів, 0.02 ETH, до 2 за транзакцію`);
    await platform.call('ціна не числом -> 400', 'GET', `/api/v0/contracts/data/${project.id}?price=free`, undefined, [400]);
    await platform.call('ліміт більший за колекцію -> 400', 'GET',
        `/api/v0/contracts/data/${project.id}?maxPerMint=${NFT_COUNT + 1}`, undefined, [400]);

    // Whitelist: корінь порожнього - нульовий; з однією адресою корінь - це її листок, доказ порожній.
    const ZERO_ROOT = `0x${'0'.repeat(64)}`;
    let { merkleRoot } = await platform.call('корінь порожнього whitelist', 'GET', `/api/v0/whitelist/${project.id}/merkle-root`);
    if (merkleRoot !== ZERO_ROOT) throw new Error(`корінь порожнього whitelist: ${merkleRoot}`);
    const buyer = ethers.Wallet.createRandom();
    await platform.addToWhitelist(project.id, [{ name: 'buyer', walletAddress: buyer.address }]);
    ({ merkleRoot } = await platform.call('корінь whitelist з однією адресою', 'GET', `/api/v0/whitelist/${project.id}/merkle-root`));
    if (merkleRoot !== ethers.utils.keccak256(buyer.address)) throw new Error(`корінь з однією адресою: ${merkleRoot}`);
    const proof = await platform.call('доказ єдиної адреси', 'GET', `/api/v0/whitelist/${project.id}/${buyer.address}/merkle-proof`);
    if (!Array.isArray(proof) || proof.length !== 0) throw new Error(`доказ єдиної адреси: ${JSON.stringify(proof)}`);

    // Чужий гаманець не має доступу до маршрутів колекції.
    const stranger = new Platform(process.env.BASE);
    const strangerWallet = ethers.Wallet.createRandom();
    await stranger.login(strangerWallet);
    for (const [what, method, url, body] of [
        ['проєкт', 'GET', `/api/v0/collections/${project.id}`],
        ['видалення проєкту', 'DELETE', `/api/v0/collections/${project.id}`],
        ['Header мінт-сайту', 'PATCH', `/api/v0/collections/${project.id}/mint-site/header-section`, { title: 'x', description: 'x' }],
        ['шари', 'GET', `/api/v0/collections/${project.id}/layers`],
        ['whitelist', 'GET', `/api/v0/whitelist/${project.id}`],
        ['корінь whitelist', 'GET', `/api/v0/whitelist/${project.id}/merkle-root`],
        ['додавання до whitelist', 'POST', '/api/v0/whitelist', { projectId: project.id, name: 'x', walletAddress: strangerWallet.address }],
        ['дані розгортання', 'GET', `/api/v0/contracts/data/${project.id}`],
    ]) {
        await stranger.call(`чужий гаманець: ${what} -> 403`, method, url, body, [403]);
    }
    await platform.call('власник і далі бачить проєкт', 'GET', `/api/v0/collections/${project.id}`);

    // Вхід: некоректний підпис - 400 / 401, підпис не можна
    // використати двічі; вихід видаляє cookie, після нього - 401.
    const guest = new Platform(process.env.BASE);
    const LOGIN = '/api/v0/auth/metamask/login';
    const nonce = await guest.call('nonce', 'GET', `/api/v0/auth/metamask/nonce?address=${strangerWallet.address}`);
    const signature = await strangerWallet.signMessage(nonce);
    await guest.call('вхід: підпис не hex -> 400', 'POST', LOGIN, { nonce, signature: '0xzz' }, [400]);
    await guest.call('вхід: обрізаний підпис -> 401', 'POST', LOGIN, { nonce, signature: signature.slice(0, -2) }, [401]);
    await guest.call('вхід: підпис іншого гаманця -> 401', 'POST', LOGIN,
        { nonce, signature: await ethers.Wallet.createRandom().signMessage(nonce) }, [401]);
    await guest.call('повторна реєстрація -> 409', 'POST', '/api/v0/auth/metamask/register',
        await strangerWallet.signMessage('Register with metamask'), [409]);
    await guest.call('вхід', 'POST', LOGIN, { nonce, signature });
    await guest.call('повторний вхід тим самим підписом -> 401', 'POST', LOGIN, { nonce, signature }, [401]);
    await guest.call('обліковий запис після входу', 'GET', '/api/v0/accounts');
    await guest.call('вихід', 'POST', '/api/v0/auth/logout');
    await guest.call('обліковий запис після виходу -> 401', 'GET', '/api/v0/accounts', undefined, [401]);

    console.log(`\nУСПІШНО. Проєкт ${project.id}; base URI контракту: ipfs://${published.metadataCid}/ (без IPFS - ${platform.base}/metadata/${project.id}/)`);
})().catch((err) => {
    console.error(`ПОМИЛКА: ${err.message}`);
    process.exit(1);
});
