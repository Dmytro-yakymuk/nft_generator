# Generator - NFT Collection Generator

Веб-платформа для генерації та розгортання NFT-колекцій на блокчейні Ethereum: завантаження шарів, генерація унікальних зображень і метаданих ERC-721, публікація колекції в IPFS (вузол Kubo у складі платформи) і закріплення у Filebase, whitelist з деревом Меркла, розгортання смарт-контракту через MetaMask, конструктор мінт-сайту, дашборд продажів за подіями контракту в блокчейні.

Покрокова інструкція: [QUICK_START.md](QUICK_START.md).

## Запуск

Потрібен лише Docker з Compose v2.

```bash
cp .env.example .env   # необов'язково: без .env діють значення за замовчуванням
make up                # застосунок, PostgreSQL і вузол IPFS (Kubo) -> http://localhost:9012
make smoke             # перевірка через API (потрібен ethers: cd web/app && npm ci --legacy-peer-deps)
```

| Команда | Що робить |
|---|---|
| `make up` / `make down` | зібрати й запустити / зупинити (дані зберігаються в томах Docker) |
| `make logs`, `make ps`, `make psql` | журнал застосунку, стан сервісів, консоль PostgreSQL |
| `make tunnel` | публічна HTTPS-адреса через ngrok (`NGROK_AUTHTOKEN`, `NGROK_DOMAIN`, `PUBLIC_URL` у `.env`) |
| `make tunnel-quick`, `make tunnel-stop` | тимчасова публічна HTTPS-адреса без реєстрації: Cloudflare Quick Tunnel, srv.us або localhost.run (що буде доступне) / зупинити тунель |
| `make chain`, `make chain-keys` | локальний блокчейн Anvil (`http://localhost:8545`, chainId 31337) і його тестові ключі |
| `make contracts` | компіляція `CollectionNFT.sol` (solc 0.8.24) і Go-біндинг (abigen) у Docker |
| `node scripts/seed-demo.js` | демо-проєкт для розгортання контракту (шари, колекція, whitelist) |
| `make test` | Go-тести, зокрема інтеграційні з PostgreSQL |
| `make smoke`, `make e2e` | сценарій через API; наскрізні тести Playwright (для розгортання контракту - мережа Anvil) |
| `make deploy-contract PROJECT=<id>` | розгортання контракту колекції з командного рядка ключем власника зі змінної `KEY_ENV` (`DEPLOYER_PRIVATE_KEY`); `PRICE=<ETH>`, `MAX_PER_MINT=<n>` - умови продажу (за замовчуванням 0.001 і 5), `DRY_RUN=1` - лише оцінка газу |
| `make pin-filebase PROJECT=<id>` | закріплення опублікованої в IPFS колекції у Filebase з командного рядка (те саме робить кнопка Pin in Filebase на сторінці колекції): каталоги зображень і метаданих як CAR-архіви через S3 API, ті самі CID (`FILEBASE_BUCKET`, `FILEBASE_S3_KEY`, `FILEBASE_S3_SECRET` у `.env`) |
| `make run` | сервер без Docker (Go + зібраний `web/app/dist`; база - з compose) |
| `make help` | усі команди |

## Сервіси й порти

| Порт | Сервіс |
|---|---|
| 9012 | застосунок: веб-інтерфейс, REST API `/api/v0`, опубліковані мінт-сайти `/sites/{адреса}` і `http://{адреса}.localhost:9012/`, метадані токенів `/metadata/{projectId}/{tokenId}`, файли `/images/` |
| 5432 | PostgreSQL |
| 8080 | HTTP-шлюз вузла IPFS (Kubo): опубліковані колекції за адресою `/ipfs/{CID}/` (`IPFS_GATEWAY_PORT`); RPC API вузла (5001) доступний лише застосунку в мережі compose |
| 4001 | обмін вузла IPFS з мережею IPFS (TCP і UDP) |
| 4040 | інспектор запитів ngrok (`make tunnel`) |
| 8545 | Anvil (`make chain`) |

## Конфігурація

Шаблон - [deploy/config.json](deploy/config.json); той самий файл використовується в Docker і для запуску без Docker (`make config` копіює його в `configs/config.json`). Під час читання конфігу сервер підставляє змінні середовища `${VAR}` і `${VAR:-значення за замовчуванням}`; перелік змінних - у [.env.example](.env.example).

Мережа Ethereum (`CHAIN_ID`, `CHAIN_NAME`, `CHAIN_RPC_URL`, `CHAIN_EXPLORER_URL`, за замовчуванням Anvil) передається у веб-застосунок як `window.APP_CONFIG`; розгорнуті контракти - у [contracts/deployments/](contracts/deployments/).

Умови продажу (ціна токена й ліміт за транзакцію) власник задає на сторінці контракту, і вони потрапляють у конструктор; максимальна кількість токенів дорівнює кількості NFT колекції. Мінт-сайт читає ці умови із самого контракту.

`GENERATION_PNG_COMPRESSION` - стиснення PNG згенерованих зображень: `default`, `best-speed` (кодування в кілька разів швидше, файли приблизно на третину більші), `best-compression`, `none`. `MINT_SITE_DOMAIN` - домен, на піддоменах якого відкриваються опубліковані мінт-сайти (`localhost` за замовчуванням; для власного домену потрібні DNS-запис `*.домен` і сертифікат; порожнє значення в `.env` вимикає піддомени - так зручно для публічної адреси через тунель).

`PUBLIC_URL` - публічна адреса платформи. З неї формуються base URI контракту (`${PUBLIC_URL}/metadata/{projectId}/`) і посилання на зображення в метаданих, тому її потрібно задати до розгортання контракту.

## Структура

```
cmd/                 точка входу (generator run --config=<каталог>)
accounts/            облікові записи, автентифікація MetaMask
collections/         проєкти, шари, генерація зображень і метаданих
whitelist/, nfts/    whitelist (дерево Меркла), дані розгортання контракту
pkg/ipfs/, pkg/filebase/  публікація колекції у вузол IPFS, закріплення у Filebase (CAR-архіви через S3 API)
dashboard/           аналітика продажів за подіями Transfer контракту (eth_getLogs)
server/httpserver/   HTTP-сервер і обробники
storage/             PostgreSQL
contracts/           смарт-контракт CollectionNFT.sol
web/app/             клієнтський застосунок (React 17, TypeScript, webpack 5, Node 16)
web/public-pages/    шаблон мінт-сайту
deploy/              Dockerfile і шаблон конфігу
scripts/smoke.js     перевірка розгорнутої платформи
```

## Ліцензія

Приватний проект. Всі права захищені.
