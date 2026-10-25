// Експеримент E12 (п. 3.2.11 роботи): час генерації тих самих наборів шарів відкритим генератором
// HashLips Art Engine (Node.js, node-canvas / Cairo, ліцензія MIT) - для порівняння з платформою (E1, E10, E11).
// HashLips закріплено на коміті COMMIT і встановлено в data/hashlips (не в git); кожен набір збирається в окремому
// робочому каталозі data/hashlips-runs/<набір>: код HashLips - символьними посиланнями, свої src/config.js
// і layers/, результат - build/images і build/json, тож клон лишається незмінним.
// Набори:
//   E1  - 5 шарів по 7 варіантів з рівними вагами (C = 16807), ті самі пікселі, що в E1 платформи (layerImage);
//   ART - набір з 26 шарів (data/demo-art, manifest.json): вага варіанта - у назві файлу "Назва#вага.png";
//         необов'язкових шарів HashLips не має, тож відсутність шару - прозорий варіант "None#вага.png".
// Налаштування HashLips, крім кількості, шарів і розміру 500x500, - стандартні; випадковий фон вимкнено
// (у наборах є власний фон). Час - від запуску процесу node index.js до його завершення (разом із записом
// усіх PNG і JSON); попередній build/ видаляється до запуску. Після прогону перевіряється, що файлів N
// і всі DNA різні; кількість відкинутих повторів DNA ("DNA exists!") - у стовпці dna_retries.
//
//   node experiments/hashlips/run.js [E1|ART]      # Node 18; дописує рядки в experiments/results/hashlips.csv
//   AMOUNTS=10000 REPEATS=3 node experiments/hashlips/run.js E1

const fs = require('fs');
const path = require('path');
const { execFileSync, spawn } = require('child_process');
const { png } = require('../../scripts/lib/platform');
const { layerImage } = require('../generation-time');

const REPO = 'https://github.com/HashLips/hashlips_art_engine.git';
const COMMIT = 'd8ee279043d2d4a8de3bdfac0d89d0e966fb04a2';
const ROOT = path.join(__dirname, '..', '..');
const HASHLIPS = path.join(ROOT, 'data', 'hashlips');
const RUNS = path.join(ROOT, 'data', 'hashlips-runs');
const ART = process.env.ART || path.join(ROOT, 'data', 'demo-art');
const OUT = process.env.OUT || path.join(__dirname, '..', 'results', 'hashlips.csv');
const REPEATS = Number(process.env.REPEATS || 5);
const AMOUNTS = (process.env.AMOUNTS || '100,1000').split(',').map(Number);
const SIZE = 500;

/** Клонує HashLips на закріпленому коміті й ставить залежності (canvas ^2.8.0 -> 2.11.2 з готовими бінарниками). */
function setup() {
    if (!fs.existsSync(HASHLIPS)) {
        execFileSync('git', ['clone', '--quiet', REPO, HASHLIPS], { stdio: 'inherit' });
        execFileSync('git', ['checkout', '--quiet', COMMIT], { cwd: HASHLIPS, stdio: 'inherit' });
    }
    if (!fs.existsSync(path.join(HASHLIPS, 'node_modules', 'canvas'))) {
        // package-lock.json закріплює canvas 2.8.0 без бінарників під Node 18 - ставимо за діапазоном package.json.
        execFileSync('npm', ['install', '--no-package-lock', '--no-audit', '--no-fund'], { cwd: HASHLIPS, stdio: 'inherit' });
    }
    const head = execFileSync('git', ['rev-parse', 'HEAD'], { cwd: HASHLIPS, encoding: 'utf8' }).trim();
    if (head !== COMMIT) throw new Error(`data/hashlips на коміті ${head}, потрібен ${COMMIT}`);
}

/** Шари набору у форматі HashLips: каталог на шар, файл "Назва#вага.png" на варіант. */
function prepareLayers(set, layersDir) {
    fs.rmSync(layersDir, { recursive: true, force: true });
    if (set === 'E1') {
        const layers = 5;
        const images = 7;
        for (let k = 0; k < layers; k++) {
            fs.mkdirSync(path.join(layersDir, `L${k}`), { recursive: true });
            for (let v = 0; v < images; v++) {
                fs.writeFileSync(path.join(layersDir, `L${k}`, `v${v + 1}#100.png`), layerImage(k, v, images));
            }
        }
        return { order: Array.from({ length: layers }, (_, k) => `L${k}`), combinations: String(images ** layers) };
    }

    const manifest = JSON.parse(fs.readFileSync(path.join(ART, 'manifest.json'), 'utf8'));
    const transparent = png(SIZE, () => [0, 0, 0, 0]);
    let combinations = 1n;
    for (const layer of manifest) {
        const dir = path.join(layersDir, layer.name);
        fs.mkdirSync(dir, { recursive: true });
        for (const variant of layer.variants) {
            fs.copyFileSync(path.join(ART, layer.name, variant.file), path.join(dir, `${variant.file.replace(/\.png$/, '')}#${variant.weight}.png`));
        }
        if (layer.none) fs.writeFileSync(path.join(dir, `None#${layer.none}.png`), transparent);
        combinations *= BigInt(layer.variants.length + (layer.none ? 1 : 0));
    }
    return { order: manifest.map((layer) => layer.name), combinations: combinations.toString() };
}

/** Робочий каталог набору: код HashLips - посиланнями, src/config.js - стандартний з нашими N, шарами й розміром. */
function prepareWorkdir(dir) {
    fs.mkdirSync(path.join(dir, 'src'), { recursive: true });
    for (const name of ['index.js', 'constants', 'modules', 'node_modules', path.join('src', 'main.js')]) {
        const link = path.join(dir, name);
        if (!fs.existsSync(link)) fs.symlinkSync(path.join(HASHLIPS, name), link);
    }
}

function writeConfig(dir, order, amount) {
    let config = fs.readFileSync(path.join(HASHLIPS, 'src', 'config.js'), 'utf8');
    const replacements = [
        [/growEditionSizeTo: \d+,/, `growEditionSizeTo: ${amount},`],
        [/layersOrder: \[[\s\S]*?\n {4}\],/, `layersOrder: ${JSON.stringify(order.map((name) => ({ name })))},`],
        [/width: 512,/, `width: ${SIZE},`],
        [/height: 512,/, `height: ${SIZE},`],
        [/generate: true,/, 'generate: false,'],
    ];
    for (const [pattern, value] of replacements) {
        if (!pattern.test(config)) throw new Error(`у src/config.js HashLips не знайдено ${pattern}`);
        config = config.replace(pattern, value);
    }
    fs.writeFileSync(path.join(dir, 'src', 'config.js'), config);
}

/** Запускає node index.js у робочому каталозі; повертає час і кількість створених та відкинутих DNA. */
function runOnce(dir) {
    return new Promise((resolve, reject) => {
        const counts = { created: 0, retries: 0 };
        let rest = '';
        const start = process.hrtime.bigint();
        const child = spawn(process.execPath, ['index.js'], { cwd: dir, stdio: ['ignore', 'pipe', 'inherit'] });
        child.stdout.setEncoding('utf8');
        child.stdout.on('data', (text) => {
            const lines = (rest + text).split('\n');
            rest = lines.pop();
            for (const line of lines) {
                if (line.startsWith('Created edition')) counts.created++;
                else if (line === 'DNA exists!') counts.retries++;
            }
        });
        child.on('error', reject);
        child.on('close', (code) => {
            const seconds = Number(process.hrtime.bigint() - start) / 1e9;
            if (code !== 0) reject(new Error(`node index.js завершився з кодом ${code}`));
            else resolve({ seconds, ...counts });
        });
    });
}

function sizeOf(dir, filter) {
    return fs.readdirSync(dir).filter(filter).reduce((total, name) => total + fs.statSync(path.join(dir, name)).size, 0);
}

async function main() {
    if (process.versions.modules !== '108') throw new Error(`потрібен Node 18 (бінарники canvas), а не ${process.version}`);
    setup();
    const sets = process.argv[2] ? [process.argv[2]] : ['E1', 'ART'];
    if (!fs.existsSync(OUT)) fs.writeFileSync(OUT, 'set,layers,combinations,amount,repeat,seconds,png_bytes,json_bytes,dna_retries\n');

    for (const set of sets) {
        const dir = path.join(RUNS, set);
        prepareWorkdir(dir);
        const { order, combinations } = prepareLayers(set, path.join(dir, 'layers'));
        for (const amount of AMOUNTS) {
            writeConfig(dir, order, amount);
            for (let repeat = 1; repeat <= REPEATS; repeat++) {
                const build = path.join(dir, 'build');
                fs.rmSync(build, { recursive: true, force: true });
                const { seconds, created, retries } = await runOnce(dir);

                const images = path.join(build, 'images');
                const json = path.join(build, 'json');
                const pngCount = fs.readdirSync(images).length;
                const metadata = JSON.parse(fs.readFileSync(path.join(json, '_metadata.json'), 'utf8'));
                const distinct = new Set(metadata.map((token) => token.dna)).size;
                if (created !== amount || pngCount !== amount || distinct !== amount) {
                    throw new Error(`${set} N=${amount}: створено ${created}, PNG ${pngCount}, різних DNA ${distinct}`);
                }
                const pngBytes = sizeOf(images, (name) => name.endsWith('.png'));
                const jsonBytes = sizeOf(json, (name) => name !== '_metadata.json');
                fs.appendFileSync(OUT, `${[set, order.length, combinations, amount, repeat, seconds.toFixed(2), pngBytes, jsonBytes, retries].join(',')}\n`);
                console.error(`HashLips ${set} N=${amount} #${repeat}: ${seconds.toFixed(2)} с (${(seconds / amount * 1000).toFixed(1)} мс/зобр.), ` +
                    `PNG ${(pngBytes / 2 ** 20).toFixed(1)} МіБ, повторів DNA ${retries}`);
            }
        }
        fs.rmSync(path.join(dir, 'build'), { recursive: true, force: true });
    }
}

main().catch((err) => {
    console.error(`ПОМИЛКА: ${err.message}`);
    process.exit(1);
});
