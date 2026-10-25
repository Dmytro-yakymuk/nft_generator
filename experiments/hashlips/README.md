# E12. Порівняння з HashLips Art Engine за часом генерації

[HashLips Art Engine](https://github.com/HashLips/hashlips_art_engine) (ліцензія MIT) - найпоширеніший відкритий генератор NFT-колекцій: скрипт на Node.js, що накладає PNG-шари бібліотекою node-canvas (Cairo) і пише метадані JSON. Скрипт `run.js` генерує ним ті самі набори шарів, що й платформа в E1 і E10, на тому самому комп'ютері.

```sh
node experiments/hashlips/run.js E1                       # N = 100 і 1000, по 5 повторів
AMOUNTS=10000 REPEATS=3 node experiments/hashlips/run.js E1
node experiments/hashlips/run.js ART                      # набір з 26 шарів: спершу render.js (див. ../README.md)
```

- Потрібен Node 18 (готові бінарники `canvas` 2.11.2 - Cairo 1.16.0; збирати Cairo не треба).
- Перший запуск клонує HashLips на закріпленому коміті `d8ee2790` у `data/hashlips` і ставить залежності (`npm install --no-package-lock`: `package-lock.json` закріплює `canvas` 2.8.0 без бінарників під Node 18). Каталог `data/` не в git.
- Кожен набір - окремий робочий каталог `data/hashlips-runs/<набір>`: код HashLips підключено символьними посиланнями, `src/config.js` - стандартний з іншими N, порядком шарів, розміром 500x500 і вимкненим випадковим фоном, `layers/` - шари набору.
- **E1:** 5 шарів по 7 варіантів з рівними вагами (`v1#100.png`), пікселі - функцією `layerImage` з `experiments/generation-time.js`.
- **ART:** набір `data/demo-art` за `manifest.json`, вага - у назві файлу (`Sunset#140.png`). HashLips не має необов'язкових шарів, тому відсутність шару - прозорий варіант `None#<вага відсутності>.png` того ж розміру, як роблять його користувачі.
- Час - від запуску `node index.js` до завершення процесу (усі PNG і JSON записано); перевіряється, що файлів N і всі DNA різні. Результати - `experiments/results/hashlips.csv`: час, обсяг PNG, обсяг JSON окремих токенів (без зведеного `_metadata.json`), кількість відкинутих повторів DNA.

Для порівняння в однакових умовах платформа вимірюється з одним воркером на поточній збірці (мітки `E1-1worker`, `ART-1worker` у `results/generation.csv`):

```sh
GENERATION_WORKERS=1 docker compose up -d --wait app
LABEL=E1-1worker AMOUNTS=100,1000 node experiments/generation-time.js E1
LABEL=ART-1worker AMOUNTS=100,1000 node experiments/generation-time.js ART   # перезаписує results/frequencies-platform.csv
docker compose up -d --wait app                                               # знову всі ядра
```
