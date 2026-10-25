// Розкладання часу одного зображення для E12 (п. 3.2.11 роботи): скільки HashLips Art Engine витрачає на завантаження
// й декодування шарів, накладання і кодування PNG - тими самими викликами node-canvas, що й src/main.js HashLips
// (loadImage для кожного шару кожного зображення, drawImage, toBuffer("image/png")). Варіанти обираються за вагами
// з назв файлів ("Назва#вага.png", None - прозорий варіант, який HashLips теж малює) з робочих каталогів run.js.
//
//   node experiments/hashlips/breakdown.js 300 data/hashlips-runs/E1/layers/L0 ... L4     # Node 18

const fs = require('fs');
const path = require('path');
const { createCanvas, loadImage } = require(path.join(__dirname, '..', '..', 'data', 'hashlips', 'node_modules', 'canvas'));

const SIZE = 500;

function pick(variants) {
    let r = Math.floor(Math.random() * variants.reduce((sum, v) => sum + v.weight, 0));
    for (const v of variants) {
        r -= v.weight;
        if (r < 0) return v;
    }
    return variants[variants.length - 1];
}

async function main() {
    const [count, ...dirs] = process.argv.slice(2);
    const layers = dirs.map((dir) => fs.readdirSync(dir).filter((name) => name.endsWith('.png')).sort()
        .map((name) => ({ file: path.join(dir, name), weight: Number(name.slice(name.lastIndexOf('#') + 1, -4)) })));
    const canvas = createCanvas(SIZE, SIZE);
    const ctx = canvas.getContext('2d');
    ctx.imageSmoothingEnabled = false;

    const time = { load: 0n, draw: 0n, encode: 0n };
    let bytes = 0;
    for (let i = 0; i < Number(count); i++) {
        const chosen = layers.map(pick);
        let start = process.hrtime.bigint();
        const images = await Promise.all(chosen.map((v) => loadImage(v.file)));
        time.load += process.hrtime.bigint() - start;

        start = process.hrtime.bigint();
        ctx.clearRect(0, 0, SIZE, SIZE);
        for (const image of images) ctx.drawImage(image, 0, 0, SIZE, SIZE);
        time.draw += process.hrtime.bigint() - start;

        start = process.hrtime.bigint();
        bytes += canvas.toBuffer('image/png').length;
        time.encode += process.hrtime.bigint() - start;
    }

    const ms = (t) => (Number(t) / 1e6 / Number(count)).toFixed(2);
    console.log(`Node ${process.versions.node}, Cairo ${require(path.join(__dirname, '..', '..', 'data', 'hashlips', 'node_modules', 'canvas')).cairoVersion}: ` +
        `на зображення: завантаження й декодування шарів ${ms(time.load)} мс, накладання ${ms(time.draw)} мс, ` +
        `кодування PNG ${ms(time.encode)} мс; ${(bytes / Number(count) / 1024).toFixed(1)} КіБ`);
}

main().catch((err) => {
    console.error(`ПОМИЛКА: ${err.message}`);
    process.exit(1);
});
