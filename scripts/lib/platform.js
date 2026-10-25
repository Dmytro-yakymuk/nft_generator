// Клієнт API платформи для скриптів перевірки й підготовки даних (smoke.js, seed-demo.js).
// Потрібен ethers з web/app/node_modules (cd web/app && npm ci --legacy-peer-deps).

const fs = require('fs');
const path = require('path');
const zlib = require('zlib');

function loadEthers() {
    try {
        return require(path.join(__dirname, '../../web/app/node_modules/ethers')).ethers;
    } catch {
        console.error('Не знайдено ethers: виконайте cd web/app && npm ci --legacy-peer-deps');
        process.exit(1);
    }
}

const ethers = loadEthers();

// Мінімальний PNG-кодувальник (RGBA, 8 біт), щоб не залежати від графічних бібліотек.
const CRC_TABLE = Array.from({ length: 256 }, (_, n) => {
    let c = n;
    for (let k = 0; k < 8; k++) c = c & 1 ? 0xedb88320 ^ (c >>> 1) : c >>> 1;
    return c >>> 0;
});

function crc32(buf) {
    let c = 0xffffffff;
    for (const byte of buf) c = CRC_TABLE[(c ^ byte) & 0xff] ^ (c >>> 8);
    return (c ^ 0xffffffff) >>> 0;
}

function chunk(type, data) {
    const out = Buffer.alloc(12 + data.length);
    out.writeUInt32BE(data.length, 0);
    out.write(type, 4, 'ascii');
    data.copy(out, 8);
    out.writeUInt32BE(crc32(out.subarray(4, 8 + data.length)), 8 + data.length);
    return out;
}

function png(size, pixel) {
    const raw = Buffer.alloc((size * 4 + 1) * size);
    for (let y = 0; y < size; y++) {
        const row = y * (size * 4 + 1);
        for (let x = 0; x < size; x++) raw.set(pixel(x, y), row + 1 + x * 4);
    }
    const header = Buffer.alloc(13);
    header.writeUInt32BE(size, 0);
    header.writeUInt32BE(size, 4);
    header.set([8, 6, 0, 0, 0], 8);
    return Buffer.concat([
        Buffer.from([0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a]),
        chunk('IHDR', header),
        chunk('IDAT', zlib.deflateSync(raw)),
        chunk('IEND', Buffer.alloc(0)),
    ]);
}

const inCircle = (x, y, cx, cy, r) => (x - cx) ** 2 + (y - cy) ** 2 <= r * r;
const TRANSPARENT = [0, 0, 0, 0];

// Тестові шари 500x500: 3 фони x 3 тіла x 2 варіанти очей = 18 комбінацій.
const DEMO_LAYERS = [
    ['Background', [[230, 80, 80], [80, 160, 230], [90, 200, 120]].map((c) => () => [...c, 255])],
    ['Body', [[250, 220, 60], [160, 90, 220], [60, 60, 60]].map((c) => (x, y) =>
        (inCircle(x, y, 250, 250, 150) ? [...c, 255] : TRANSPARENT))],
    ['Eyes', [[255, 255, 255], [255, 0, 0]].map((c) => (x, y) =>
        (inCircle(x, y, 205, 225, 25) || inCircle(x, y, 295, 225, 25) ? [...c, 255] : TRANSPARENT))],
];

class Platform {
    constructor(base) {
        this.base = (base || 'http://localhost:9012').replace(/\/$/, '');
        this.cookie = '';
    }

    async api(method, url, body) {
        const res = await fetch(this.base + url, {
            method,
            headers: { 'Content-Type': 'application/json', 'ngrok-skip-browser-warning': '1', Cookie: this.cookie },
            body: body === undefined ? undefined : JSON.stringify(body),
        });
        const setCookie = res.headers.get('set-cookie');
        if (setCookie) this.cookie = setCookie.split(';')[0];
        const text = await res.text();
        let data;
        try {
            data = JSON.parse(text);
        } catch {
            data = text;
        }
        return { status: res.status, data };
    }

    /** Виконує запит і перевіряє код відповіді. */
    async call(what, method, url, body, ok = [200]) {
        const r = await this.api(method, url, body);
        if (!ok.includes(r.status)) throw new Error(`${what}: HTTP ${r.status} ${JSON.stringify(r.data)}`);
        console.log(`ok   ${what}`);
        return r.data;
    }

    /** Вхід підписом гаманця; новий гаманець спершу реєструється. */
    async login(wallet) {
        const nonceUrl = `/api/v0/auth/metamask/nonce?address=${wallet.address}`;
        let r = await this.api('GET', nonceUrl);
        if (r.status === 404) {
            await this.call('реєстрація підписом', 'POST', '/api/v0/auth/metamask/register', await wallet.signMessage('Register with metamask'));
            r = await this.api('GET', nonceUrl);
        }
        if (r.status !== 200) throw new Error(`nonce: HTTP ${r.status} ${JSON.stringify(r.data)}`);
        await this.call(`вхід ${wallet.address}`, 'POST', '/api/v0/auth/metamask/login', { nonce: r.data, signature: await wallet.signMessage(r.data) });
    }

    async findProject(name) {
        const projects = await this.call('список проєктів', 'GET', '/api/v0/collections');
        return (projects || []).find((p) => p.name === name);
    }

    async createProject(name, nftCount) {
        await this.call(`проєкт "${name}"`, 'POST', '/api/v0/collections', { name, nftCount });
        return this.findProject(name);
    }

    async uploadLayers(projectId, layers = DEMO_LAYERS) {
        for (const [i, [name]] of layers.entries()) {
            await this.call(`шар ${name}`, 'POST', `/api/v0/collections/${projectId}/layers`, { name, priority: i + 1 });
        }
        const created = await this.call('список шарів', 'GET', `/api/v0/collections/${projectId}/layers`);
        for (const [name, variants] of layers) {
            const layer = created.find((l) => l.name === name);
            for (const [i, pixel] of variants.entries()) {
                const image = png(500, pixel).toString('base64');
                const r = await this.api('POST', `/api/v0/collections/${projectId}/layers/${layer.id}/images`, { image, name: `v${i + 1}.png` });
                if (r.status !== 200) throw new Error(`зображення ${name}/v${i + 1}.png: HTTP ${r.status}`);
            }
            console.log(`ok   зображення шару ${name}: ${variants.length}`);
        }
    }

    /**
     * Завантажує процедурний набір шарів (scripts/demo-art/render.js): шари в порядку manifest.json,
     * зображення варіантів, ваги варіантів і вага відсутності для необов'язкових шарів.
     */
    async uploadArt(projectId, dir) {
        const manifest = JSON.parse(fs.readFileSync(path.join(dir, 'manifest.json'), 'utf8'));
        for (const [i, layer] of manifest.entries()) {
            const r = await this.api('POST', `/api/v0/collections/${projectId}/layers`, { name: layer.name, priority: i });
            if (r.status !== 200) throw new Error(`шар ${layer.name}: HTTP ${r.status} ${JSON.stringify(r.data)}`);
        }
        const created = await this.call('список шарів', 'GET', `/api/v0/collections/${projectId}/layers`);
        for (const layer of manifest) {
            const { id } = created.find((l) => l.name === layer.name);
            const weights = {};
            for (const variant of layer.variants) {
                const image = fs.readFileSync(path.join(dir, layer.name, variant.file)).toString('base64');
                const r = await this.api('POST', `/api/v0/collections/${projectId}/layers/${id}/images`, { image, name: variant.file });
                if (r.status !== 200) throw new Error(`зображення ${layer.name}/${variant.file}: HTTP ${r.status}`);
                weights[variant.file] = variant.weight;
            }
            await this.call(`шар ${layer.name}: ${layer.variants.length} варіантів, ваги${layer.none ? `, відсутність ${layer.none}` : ''}`,
                'PUT', `/api/v0/collections/${projectId}/layers/${id}/weights`, { noneWeight: layer.none, weights });
        }
        return manifest;
    }

    async generate(projectId, count) {
        await this.call('запуск генерації', 'POST', `/api/v0/collections/${projectId}/generate-images`);
        let generated = 0;
        for (let i = 0; i < 240 + count && generated < count; i++) {
            await new Promise((resolve) => setTimeout(resolve, 500));
            const r = await this.api('GET', `/api/v0/collections/${projectId}/generated-images-count`);
            if (r.status === 200) generated = r.data;
        }
        if (generated < count) throw new Error(`згенеровано ${generated} з ${count} зображень`);
        console.log(`ok   згенеровано ${generated} зображень`);
    }

    async addToWhitelist(projectId, entries) {
        for (const { name, walletAddress } of entries) {
            await this.call(`whitelist ${name} ${walletAddress}`, 'POST', '/api/v0/whitelist', { projectId, name, walletAddress });
        }
    }
}

module.exports = { ethers, png, DEMO_LAYERS, Platform };
