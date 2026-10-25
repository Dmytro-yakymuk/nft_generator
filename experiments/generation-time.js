// Експерименти E1 / E2 (розділ 3 роботи): час генерації колекції на запущеній платформі через REST API.
// E1 - залежність від кількості зображень N (5 шарів по 7 варіантів, C = 16807);
// E2 - залежність від кількості шарів n (по 10 варіантів, C = 10^n) для N = 1000;
// E3 - перевірка E3 на запущеній платформі: 20 зображень з 8 / 9 / 10 шарів по 10 варіантів
//      (запускається на образах з послідовним інкрементуванням і з прямим обчисленням, LABEL=E3-old / E3-new розрізняє рядки).
// ART - час генерації процедурного набору з 26 шарів з вагами (scripts/demo-art/render.js, data/demo-art)
//      для N = 100 і 1000; після першого прогону з N = 1000 частоти варіантів у метаданих згенерованих
//      токенів записуються в results/frequencies-platform.csv (перевірка E8 на запущеній платформі).
// Час вимірюється від запиту POST .../generate-images до появи всіх N файлів (опитування кожні 250 мс);
// обсяг результатів - du у контейнері app. Після кожного прогону файли проєкту видаляються.
//
//   node experiments/generation-time.js [E1|E2]     # дописує рядки в experiments/results/generation.csv
//   REPEATS=3 OUT=... node experiments/generation-time.js E1
//   AMOUNTS=1000 node experiments/generation-time.js ART     # AMOUNTS - і для E1
//   LABEL=E1-1worker ...                                       # мітка рядків (E12: GENERATION_WORKERS=1)

const fs = require('fs');
const path = require('path');
const { execFileSync } = require('child_process');
const { ethers, png, Platform } = require('../scripts/lib/platform');

const REPEATS = Number(process.env.REPEATS || 5);
const POLL_MS = 250;
const SIZE = 500;
const OUT = process.env.OUT || path.join(__dirname, 'results', 'generation.csv');

const PLANS = {
    E1: (process.env.AMOUNTS || '100,1000,5000,10000').split(',').map((amount) => ({ layers: 5, images: 7, amount: Number(amount) })),
    E2: [3, 5, 7, 10].map((layers) => ({ layers, images: 10, amount: 1000 })),
    E3: [8, 9, 10].map((layers) => ({ layers, images: 10, amount: 20 })),
    ART: (process.env.AMOUNTS || '100,1000').split(',').map((amount) => ({ art: true, amount: Number(amount) })),
};

const ART = process.env.ART || path.join(__dirname, '..', 'data', 'demo-art');

/** Кількість токенів з кожним значенням кожного атрибута в метаданих згенерованої колекції. */
function platformFrequencies(projectId) {
    const json = inContainer(`for f in data/files/metadata/${projectId}/*; do cat "$f"; echo; done`);
    const counts = {};
    let tokens = 0;
    for (const line of json.split('\n').filter(Boolean)) {
        tokens++;
        for (const { trait_type: layer, value } of JSON.parse(line).attributes) {
            counts[`${layer}\t${value}`] = (counts[`${layer}\t${value}`] || 0) + 1;
        }
    }
    const manifest = JSON.parse(fs.readFileSync(path.join(ART, 'manifest.json'), 'utf8'));
    const rows = ['layer,variant,weight,expected_share,count,tokens'];
    for (const layer of manifest) {
        const variants = layer.variants.map((v) => [v.file.replace(/\.png$/, ''), v.weight]);
        if (layer.none) variants.push(['None', layer.none]);
        const sum = variants.reduce((total, [, weight]) => total + weight, 0);
        for (const [name, weight] of variants) {
            rows.push([JSON.stringify(layer.name), JSON.stringify(name), weight, (weight / sum).toFixed(5), counts[`${layer.name}\t${name}`] || 0, tokens].join(','));
        }
    }
    fs.writeFileSync(path.join(__dirname, 'results', 'frequencies-platform.csv'), `${rows.join('\n')}\n`);
    console.error(`частоти варіантів ${tokens} токенів -> results/frequencies-platform.csv`);
}

/** Варіант v шару k: фон - градієнт на все полотно, інші шари - напівпрозоре коло з градієнтом. */
function layerImage(k, v, images) {
    const hue = (v / images) * 2 * Math.PI;
    const color = [0, 2, 4].map((shift) => Math.round(127 + 120 * Math.sin(hue + shift)));
    if (k === 0) {
        return png(SIZE, (x, y) => [...color.map((c, i) => (c + (i === 0 ? x : y) / 4) % 256), 255]);
    }
    const angle = (k * 2.4) % (2 * Math.PI);
    const cx = 250 + 140 * Math.cos(angle);
    const cy = 250 + 140 * Math.sin(angle);
    const r = 60 + 10 * (k % 4);
    return png(SIZE, (x, y) => {
        const d = Math.hypot(x - cx, y - cy);
        return d > r ? [0, 0, 0, 0] : [...color.map((c) => Math.round(c * (1 - d / (2 * r)))), 230];
    });
}

function inContainer(command) {
    return execFileSync('docker', ['compose', 'exec', '-T', 'app', 'sh', '-c', command], { encoding: 'utf8', maxBuffer: 256 * 2 ** 20 }).trim();
}

async function uploadLayers(platform, projectId, layers, images) {
    for (let k = 0; k < layers; k++) {
        await platform.api('POST', `/api/v0/collections/${projectId}/layers`, { name: `L${k}`, priority: k + 1 });
    }
    const created = await platform.api('GET', `/api/v0/collections/${projectId}/layers`);
    for (let k = 0; k < layers; k++) {
        const layer = created.data.find((l) => l.name === `L${k}`);
        for (let v = 0; v < images; v++) {
            const image = layerImage(k, v, images).toString('base64');
            const r = await platform.api('POST', `/api/v0/collections/${projectId}/layers/${layer.id}/images`, { image, name: `v${v + 1}.png` });
            if (r.status !== 200) throw new Error(`шар L${k}, варіант ${v + 1}: HTTP ${r.status}`);
        }
    }
}

async function generate(platform, projectId, amount) {
    const start = process.hrtime.bigint();
    const r = await platform.api('POST', `/api/v0/collections/${projectId}/generate-images`);
    if (r.status !== 200) throw new Error(`генерація: HTTP ${r.status} ${JSON.stringify(r.data)}`);
    for (;;) {
        await new Promise((resolve) => setTimeout(resolve, POLL_MS));
        const count = await platform.api('GET', `/api/v0/collections/${projectId}/generated-images-count`);
        if (count.status === 200 && count.data >= amount) break;
    }
    return Number(process.hrtime.bigint() - start) / 1e9;
}

async function main() {
    const which = process.argv[2] ? [process.argv[2]] : ['E1', 'E2'];
    const platform = new Platform(process.env.BASE);
    await platform.login(ethers.Wallet.createRandom());
    if (!fs.existsSync(OUT)) fs.writeFileSync(OUT, 'experiment,layers,images_per_layer,combinations,amount,repeat,seconds,bytes\n');

    for (const experiment of which) {
        for (const plan of PLANS[experiment]) {
            const { amount } = plan;
            let { layers, images } = plan;
            let combinations = images ** layers;
            for (let repeat = 1; repeat <= REPEATS; repeat++) {
                const name = `${experiment} n=${layers} m=${images} N=${amount} #${repeat} ${Date.now()}`;
                const project = await platform.createProject(name, amount);
                if (plan.art) {
                    // шари набору мають різну кількість варіантів: images - середня, combinations - добуток.
                    const manifest = await platform.uploadArt(project.id, ART);
                    layers = manifest.length;
                    const counts = manifest.map((layer) => layer.variants.length + (layer.none ? 1 : 0));
                    images = (counts.reduce((a, b) => a + b, 0) / layers).toFixed(1);
                    combinations = counts.reduce((product, count) => product * BigInt(count), 1n).toString();
                } else {
                    await uploadLayers(platform, project.id, layers, images);
                }
                const seconds = await generate(platform, project.id, amount);
                const bytes = Number(inContainer(`du -sb data/files/generated/${project.id} | cut -f1`));
                if (plan.art && amount === 1000 && repeat === 1) platformFrequencies(project.id);
                fs.appendFileSync(OUT, `${[process.env.LABEL || experiment, layers, images, combinations, amount, repeat, seconds.toFixed(2), bytes].join(',')}\n`);
                console.error(`${name}: ${seconds.toFixed(2)} с, ${(bytes / 2 ** 20).toFixed(1)} МіБ`);

                await platform.api('DELETE', `/api/v0/collections/${project.id}`);
                inContainer(`rm -rf data/files/generated/${project.id} data/files/${project.id} data/files/metadata/${project.id}`);
            }
        }
    }
}

module.exports = { uploadLayers, generate, layerImage };

if (require.main === module) {
    main().catch((err) => {
        console.error(`ПОМИЛКА: ${err.message}`);
        process.exit(1);
    });
}
