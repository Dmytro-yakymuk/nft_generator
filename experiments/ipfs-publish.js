// Експеримент E6 (п. 3.2.5 роботи): публікація згенерованої колекції в IPFS через REST API платформи.
// Для N = 100 / 1000 / 10000 зображень (5 шарів по 7 варіантів) вимірюється час запиту
// POST .../ipfs, перевіряється, що повторна публікація дає ті самі CID (адресація за вмістом),
// і що метадані та зображення токенів віддаються HTTP-шлюзом вузла за адресами з контракту.
//
//   node experiments/ipfs-publish.js     # дописує рядки в experiments/results/ipfs.csv
//   GATEWAY=http://localhost:8080 REPEATS=5 node experiments/ipfs-publish.js

const fs = require('fs');
const path = require('path');
const { execFileSync } = require('child_process');
const { ethers, Platform } = require('../scripts/lib/platform');
const { uploadLayers, generate } = require('./generation-time');

const REPEATS = Number(process.env.REPEATS || 5);
const GATEWAY = (process.env.GATEWAY || 'http://localhost:8080').replace(/\/$/, '');
const OUT = process.env.OUT || path.join(__dirname, 'results', 'ipfs.csv');
const AMOUNTS = (process.env.AMOUNTS || '100,1000,10000').split(',').map(Number);

function inContainer(service, command) {
    return execFileSync('docker', ['compose', 'exec', '-T', service, 'sh', '-c', command], { encoding: 'utf8' }).trim();
}

async function timed(fn) {
    const start = process.hrtime.bigint();
    const result = await fn();
    return [Number(process.hrtime.bigint() - start) / 1e9, result];
}

async function fetchOk(url) {
    const [seconds, response] = await timed(() => fetch(url));
    if (response.status !== 200) throw new Error(`${url}: HTTP ${response.status}`);
    const body = Buffer.from(await response.arrayBuffer());
    return { seconds, body };
}

(async () => {
    const platform = new Platform(process.env.BASE);
    await platform.login(ethers.Wallet.createRandom());
    if (!fs.existsSync(OUT)) fs.writeFileSync(OUT, 'amount,repeat,seconds,images_cid,metadata_cid,same_cid,metadata_fetch_s,image_fetch_s\n');

    for (const amount of AMOUNTS) {
        const project = await platform.createProject(`E6 N=${amount} ${Date.now()}`, amount);
        await uploadLayers(platform, project.id, 5, 7);
        await generate(platform, project.id, amount);

        let first;
        for (let repeat = 1; repeat <= REPEATS; repeat++) {
            const [seconds, r] = await timed(() => platform.api('POST', `/api/v0/collections/${project.id}/ipfs`));
            if (r.status !== 200) throw new Error(`публікація: HTTP ${r.status} ${JSON.stringify(r.data)}`);
            const { imagesCid, metadataCid } = r.data;
            first = first || { imagesCid, metadataCid };
            const same = imagesCid === first.imagesCid && metadataCid === first.metadataCid;

            // токен з останнім номером: tokenURI = ipfs://<CID метаданих>/<номер>
            const tokenId = amount - 1;
            const metadata = await fetchOk(`${GATEWAY}/ipfs/${metadataCid}/${tokenId}`);
            const { image } = JSON.parse(metadata.body.toString());
            if (image !== `ipfs://${imagesCid}/${tokenId}.png`) throw new Error(`image у метаданих: ${image}`);
            const picture = await fetchOk(`${GATEWAY}/ipfs/${imagesCid}/${tokenId}.png`);
            if (picture.body.subarray(1, 4).toString() !== 'PNG') throw new Error('шлюз віддав не PNG');

            fs.appendFileSync(OUT, `${[amount, repeat, seconds.toFixed(2), imagesCid, metadataCid, same, metadata.seconds.toFixed(3), picture.seconds.toFixed(3)].join(',')}\n`);
            console.error(`N=${amount} #${repeat}: ${seconds.toFixed(2)} с, CID ${same ? 'ті самі' : 'ІНШІ'}, метадані ${metadata.seconds.toFixed(3)} с, зображення ${picture.seconds.toFixed(3)} с`);
        }

        await platform.api('DELETE', `/api/v0/collections/${project.id}`);
        inContainer('app', `rm -rf data/files/generated/${project.id} data/files/${project.id} data/files/metadata/${project.id}`);
    }
    console.error(`Обсяг сховища вузла IPFS: ${inContainer('ipfs', 'ipfs repo stat --human | grep RepoSize')}`);
})().catch((err) => {
    console.error(`ПОМИЛКА: ${err.message}`);
    process.exit(1);
});
