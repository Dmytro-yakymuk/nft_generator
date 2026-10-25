// Рендер процедурного набору шарів у PNG 500x500 з прозорістю через headless Chromium (Playwright).
// Шари й ваги - scripts/demo-art/layers.js; результат - data/demo-art/<шар>/<варіант>.png і manifest.json
// з порядком шарів, вагами варіантів і вагою відсутності (його читає scripts/seed-demo.js).
//
//   node scripts/demo-art/render.js [--sheet] [--samples N] [шар...]
//
// --sheet - аркуш з усіма варіантами кожного шару (sheets/<шар>.png), --samples N - N випадкових
// зображень, складених за вагами (sheets/samples.png); список шарів обмежує рендер.
// Playwright - з PLAYWRIGHT_CORE, з node_modules репозиторію (npm ci; браузер - make e2e) або глобальний.

const fs = require('fs');
const path = require('path');

const LAYERS = require('./layers');

const OUT = process.env.OUT || path.join(__dirname, '../../data/demo-art');

function loadPlaywright() {
    const candidates = [
        process.env.PLAYWRIGHT_CORE,
        path.join(__dirname, '../../node_modules/playwright-core'),
        path.join(path.dirname(process.execPath), '../lib/node_modules/playwright-core'),
        'playwright-core',
    ].filter(Boolean);
    for (const candidate of candidates) {
        try {
            return require(candidate);
        } catch {
            // наступний варіант
        }
    }
    console.error('Не знайдено playwright-core: виконайте npm ci у корені репозиторію або задайте PLAYWRIGHT_CORE');
    process.exit(1);
}

const fileName = (name) => `${name}.png`;

async function renderLayer(page, layer) {
    const dir = path.join(OUT, layer.name);
    fs.rmSync(dir, { recursive: true, force: true });
    fs.mkdirSync(dir, { recursive: true });
    for (const variant of layer.variants) {
        await page.setContent(`<html><body style="margin:0;background:transparent">${variant.image}</body></html>`);
        const png = await page.locator('svg').screenshot({ omitBackground: true });
        fs.writeFileSync(path.join(dir, fileName(variant.name)), png);
    }
}

const fileUrl = (layer, variant) =>
    `data:image/png;base64,${fs.readFileSync(path.join(OUT, layer.name, fileName(variant.name))).toString('base64')}`;

/** Аркуш варіантів шару на шаховому тлі з підписами, вагою й часткою. */
async function sheet(page, layer) {
    const total = layer.variants.reduce((sum, variant) => sum + variant.weight, layer.none);
    const cells = layer.variants.map((variant) =>
        `<figure><div class="img"><img src="${fileUrl(layer, variant)}"></div>` +
        `<figcaption>${variant.name}<br>${variant.weight} · ${(variant.weight / total * 100).toFixed(1)}%</figcaption></figure>`).join('');
    const none = layer.none ? `<figure><div class="img none">None</div><figcaption>None<br>${layer.none} · ${(layer.none / total * 100).toFixed(1)}%</figcaption></figure>` : '';
    await page.setContent(`<html><body style="margin:0;font:14px sans-serif;background:#fff">
        <style>body{display:flex;flex-wrap:wrap;gap:10px;padding:10px;width:1080px}figure{margin:0;width:200px}
        .img{width:200px;height:200px;background:repeating-conic-gradient(#ddd 0% 25%,#f6f6f6 0% 50%) 50%/20px 20px}
        .img img{width:200px;height:200px}.none{display:flex;align-items:center;justify-content:center;border:2px dashed #aaa;box-sizing:border-box}
        figcaption{text-align:center}</style><h2 style="width:100%;margin:0">${layer.name}</h2>${none}${cells}</body></html>`);
    await page.waitForLoadState('load');
    fs.mkdirSync(path.join(OUT, 'sheets'), { recursive: true });
    await page.screenshot({ path: path.join(OUT, 'sheets', `${layer.name}.png`), fullPage: true });
}

function drawVariant(layer, random) {
    const total = layer.variants.reduce((sum, variant) => sum + variant.weight, layer.none);
    let point = random() * total;
    for (const variant of layer.variants) {
        point -= variant.weight;
        if (point < 0) return variant;
    }
    return null;
}

/** Випадкові зображення за вагами (накладання шарів у порядку пріоритету, як на сервері). */
async function samples(page, count, seed) {
    const random = require('./lib').rng(seed);
    const tiles = Array.from({ length: count }, () => {
        const parts = LAYERS.map((layer) => [layer, drawVariant(layer, random)]).filter(([, variant]) => variant);
        return `<div class="tile">${parts.map(([layer, variant]) => `<img src="${fileUrl(layer, variant)}">`).join('')}</div>`;
    }).join('');
    await page.setContent(`<html><body style="margin:0;background:#fff"><style>
        body{display:grid;grid-template-columns:repeat(4,250px);gap:8px;padding:8px;width:1032px}
        .tile{position:relative;width:250px;height:250px}.tile img{position:absolute;inset:0;width:250px;height:250px}
        </style>${tiles}</body></html>`);
    await page.waitForLoadState('load');
    fs.mkdirSync(path.join(OUT, 'sheets'), { recursive: true });
    await page.screenshot({ path: path.join(OUT, 'sheets', 'samples.png'), fullPage: true });
}

(async () => {
    const args = process.argv.slice(2);
    const withSheets = args.includes('--sheet');
    const samplesAt = args.indexOf('--samples');
    const sampleCount = samplesAt >= 0 ? Number(args[samplesAt + 1]) : 0;
    const only = args.filter((arg, i) => !arg.startsWith('--') && !(samplesAt >= 0 && i === samplesAt + 1));

    const { chromium } = loadPlaywright();
    const browser = await chromium.launch();
    const page = await browser.newPage({ viewport: { width: 500, height: 500 } });
    page.on('pageerror', (err) => console.error(err));

    const selected = only.length ? LAYERS.filter((layer) => only.includes(layer.name)) : LAYERS;
    for (const layer of selected) {
        await renderLayer(page, layer);
        console.log(`${layer.name}: ${layer.variants.length} варіантів${layer.none ? `, відсутність ${layer.none}` : ''}`);
    }

    fs.writeFileSync(path.join(OUT, 'manifest.json'), JSON.stringify(LAYERS.map((layer) => ({
        name: layer.name,
        none: layer.none,
        variants: layer.variants.map((variant) => ({ file: fileName(variant.name), weight: variant.weight })),
    })), null, 2));

    if (withSheets) {
        await page.setViewportSize({ width: 1100, height: 600 });
        for (const layer of selected) await sheet(page, layer);
    }
    if (sampleCount) {
        await page.setViewportSize({ width: 1050, height: 600 });
        await samples(page, sampleCount, Number(process.env.SEED || 1));
    }
    await browser.close();

    const variants = LAYERS.reduce((sum, layer) => sum + layer.variants.length, 0);
    const combinations = LAYERS.reduce((product, layer) => product * BigInt(layer.variants.length + (layer.none ? 1 : 0)), 1n);
    console.log(`Шарів: ${LAYERS.length}, варіантів: ${variants}, комбінацій: ${combinations} (~${Number(combinations).toExponential(2)})`);
    console.log(`Результат: ${OUT}`);
})().catch((err) => {
    console.error(err);
    process.exit(1);
});
