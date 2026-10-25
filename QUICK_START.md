# Швидкий старт

Від `git clone` до публічної адреси. Потрібні Docker з Compose v2 і браузер з розширенням MetaMask; Node, Go, solc і ngrok окремо встановлювати не потрібно.

## 1. Локальний запуск

```bash
git clone <репозиторій> generator && cd generator
cp .env.example .env      # за потреби змінити TOKEN_SECRET і порти
make up                   # застосунок, PostgreSQL і вузол IPFS (Kubo); перша збірка образу - кілька хвилин
```

Платформа: http://localhost:9012, HTTP-шлюз вузла IPFS: http://localhost:8080. Кнопка Connect підключає MetaMask: перший вхід реєструє гаманець (підпис повідомлення "Register with metamask"), далі вхід виконується підписом одноразового nonce.

Перевірка без браузера (створює тестовий проєкт, генерує 5 зображень, перевіряє метадані токенів і публікує колекцію у вузол IPFS):

```bash
cd web/app && npm ci --legacy-peer-deps && cd ../..   # один раз, для ethers
make smoke
```

## 2. Локальний блокчейн (Anvil)

```bash
make chain        # http://localhost:8545, chainId 31337
make chain-keys   # 10 тестових акаунтів по 10 000 ETH
```

У MetaMask імпортувати один із приватних ключів з `make chain-keys`; мережу Anvil застосунок запропонує додати сам під час розгортання контракту. Стан блокчейну зберігається між перезапусками (`docker compose ... down -v` його видаляє).

Розгортання контракту:

```bash
node scripts/seed-demo.js   # проєкт "Demo Collection": шари, 10 зображень, whitelist з трьох адрес Anvil
```

Перед розгортанням колекцію варто опублікувати в IPFS: на сторінці згенерованої колекції - Publish to IPFS (каталоги зображень і метаданих додаються у вузол Kubo платформи, файли відкриваються через шлюз `http://localhost:8080/ipfs/<CID>/`). Тоді в конструктор контракту потрапляє base URI `ipfs://<CID метаданих>/`, і токени не залежать від сервера платформи; контракт, розгорнутий раніше, переводиться на IPFS кнопкою Point the contract to IPFS (`setBaseURI` з MetaMask). Якщо в `.env` задано `FILEBASE_BUCKET`, `FILEBASE_S3_KEY` і `FILEBASE_S3_SECRET`, кнопка Pin in Filebase закріплює обидва каталоги у Filebase з тими самими CID.

Далі в браузері: вхід гаманцем тестового акаунта №1 (власник демо-проєкту) -> Studio -> проєкт -> Contract -> назва, ціна токена й ліміт за транзакцію (кількість токенів - з колекції) -> Publish. Адреса контракту зберігається в проєкті; розгорнуті адреси - у `contracts/deployments/` і `.env` (`CONTRACT_ADDRESS`).

Резервний спосіб без браузера - з командного рядка ключем власника колекції (адреса й блок контракту так само зберігаються в проєкті):

```bash
DEPLOYER_PRIVATE_KEY=0x... make deploy-contract PROJECT=<id> DRY_RUN=1   # лише оцінка газу й вартості
DEPLOYER_PRIVATE_KEY=0x... make deploy-contract PROJECT=<id>             # розгортання
make deploy-contract PROJECT=<id> KEY_ENV=SEPOLIA_DEPLOYER_PRIVATE_KEY  # ключ з іншої змінної (наприклад, з .env)
make deploy-contract PROJECT=<id> PRICE=0.002 MAX_PER_MINT=3            # власні умови продажу
```

Мережа - та, що в конфігурації сервера (`CHAIN_SERVER_RPC_URL`, інакше `CHAIN_RPC_URL`); ключ читається лише зі змінної середовища й має належати гаманцю власника колекції, інакше команда відмовить.

Мережа задається в `.env` (`CHAIN_ID`, `CHAIN_NAME`, `CHAIN_RPC_URL`, `CHAIN_EXPLORER_URL`; за замовчуванням - Anvil, приклад для Sepolia - у `.env.example`) і застосовується після `make up` без перезбирання фронтенду. Після змін у `contracts/CollectionNFT.sol` - `make contracts` (компіляція й Go-біндинг у Docker).

## 3. Мінт-сайт

Website -> заповнити хоча б розділ Header -> Publish -> адреса (за замовчуванням - з назви проєкту). Для публікації контракт має бути розгорнуто. Опублікований сайт відкривається без входу в платформу: `http://localhost:9012/sites/<адреса>` або на власному піддомені `http://<адреса>.localhost:9012/` (`MINT_SITE_DOMAIN`); покупець підключає MetaMask і мінтить у presale (якщо гаманець у whitelist) або в публічному продажі. Після зміни whitelist розгорнутого контракту власник оновлює корінь Меркла в контракті кнопкою Update contract на сторінці whitelist (`setMerkleRoot` з MetaMask).

## 4. Публічна адреса

Без реєстрації - тимчасова адреса (нова після кожного запуску; колекція, опублікована в IPFS, від неї не залежить): Cloudflare Quick Tunnel, а якщо його сервіс недоступний з мережі - SSH-тунель srv.us (адресу задає SSH-ключ у томі `tunnelkey`, тож після перезапуску вона та сама) або, якщо й він недоступний, localhost.run (повільний для великих файлів). Галерея бере адреси зображень з `PUBLIC_URL`, тож після першого запуску задайте `PUBLIC_URL=<адреса тунелю>` у `.env` і виконайте `make up`:

```bash
make tunnel-quick   # https://<слова>.trycloudflare.com, https://<id>.srv.us або https://<id>.lhr.life
make tunnel-stop
```

Сталий домен - через ngrok:

1. Зареєструватися на https://dashboard.ngrok.com, скопіювати authtoken і створити безкоштовний сталий домен (Domains).
2. У `.env`:
   ```
   NGROK_AUTHTOKEN=<authtoken>
   NGROK_DOMAIN=<домен>.ngrok-free.app
   PUBLIC_URL=https://<домен>.ngrok-free.app
   ```
3. Застосувати й запустити тунель:
   ```bash
   make up       # застосунок отримує новий PUBLIC_URL
   make tunnel   # https://<домен>.ngrok-free.app, журнал запитів - http://localhost:4040
   BASE=https://<домен>.ngrok-free.app make smoke
   ```

`PUBLIC_URL` потрапляє в base URI контракту під час розгортання, тому домен має бути сталим, а контракт слід розгортати вже з публічною адресою. Безкоштовний ngrok показує сторінку-попередження при першому відкритті сайту в браузері; адреса працює, доки запущено комп'ютер.

## 5. Зупинка

```bash
make down                                   # зупинити, дані зберігаються
docker compose --profile tunnel --profile chain down -v   # зупинити й видалити дані
```

## 6. Без Docker

```bash
make config                                   # configs/config.json з шаблону deploy/config.json
cd web/app && npm ci --legacy-peer-deps && npm run build && cd ../..   # потрібен Node 16
make run                                      # база в Docker, сервер - go run
```

## Можливі проблеми

| Симптом | Рішення |
|---|---|
| Порт 9012 або 5432 зайнятий | змінити `APP_PORT` / `DB_PORT` у `.env` |
| `error:0308010C:digital envelope routines::unsupported` під час збірки фронтенду | збирати на Node 16 (у Docker це вже так) |
| Метадані або зображення вказують на `localhost` після запуску тунелю | задати `PUBLIC_URL` у `.env` і виконати `make up` |
