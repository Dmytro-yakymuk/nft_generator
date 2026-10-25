# Запуск платформи: make up -> http://localhost:9012, make tunnel -> публічна адреса ngrok.
# Налаштування - у .env (див. .env.example); без нього діють значення за замовчуванням.

-include .env
export

APP_PORT ?= 9012
DB_PORT ?= 5432
POSTGRES_PASSWORD ?= postgres
COMPOSE := docker compose
TEST_DB := postgres://postgres:$(POSTGRES_PASSWORD)@localhost:$(DB_PORT)/generatordb_test?sslmode=disable
SOLC_VERSION := 0.8.24
GO_ETHEREUM_VERSION := v1.10.16
DOCKER_USER := $(shell id -u):$(shell id -g)
GOMODCACHE := $(shell go env GOMODCACHE 2>/dev/null)

.DEFAULT_GOAL := help

.PHONY: help
help: ## Перелік команд
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: up
up: ## Зібрати й запустити платформу (PostgreSQL, вузол IPFS і застосунок)
	$(COMPOSE) build app
	$(COMPOSE) up -d --wait db app
	@echo "Платформа: http://localhost:$(APP_PORT)"

.PHONY: down
down: ## Зупинити всі сервіси (дані в томах Docker зберігаються)
	$(COMPOSE) --profile tunnel --profile chain down

.PHONY: rebuild
rebuild: ## Перезібрати образ застосунку без кешу й перезапустити
	$(COMPOSE) build --no-cache app
	$(COMPOSE) up -d --wait app

.PHONY: logs
logs: ## Журнал застосунку
	$(COMPOSE) logs -f app

.PHONY: ps
ps: ## Стан сервісів
	$(COMPOSE) --profile tunnel --profile chain ps

.PHONY: psql
psql: ## Консоль PostgreSQL
	$(COMPOSE) exec db psql -U postgres generator

.PHONY: tunnel
tunnel: ## Публічна HTTPS-адреса через ngrok (NGROK_AUTHTOKEN і NGROK_DOMAIN у .env)
	@if [ -z "$(NGROK_AUTHTOKEN)" ] || [ -z "$(NGROK_DOMAIN)" ]; then \
		echo "Задайте NGROK_AUTHTOKEN і NGROK_DOMAIN у .env"; exit 1; fi
	@if [ "$(PUBLIC_URL)" != "https://$(NGROK_DOMAIN)" ]; then \
		echo "Увага: PUBLIC_URL у .env має дорівнювати https://$(NGROK_DOMAIN), інакше метадані й контракт вказуватимуть на іншу адресу (після зміни - make up)"; fi
	$(COMPOSE) --profile tunnel up -d ngrok
	@echo "Платформа: https://$(NGROK_DOMAIN)   журнал запитів: http://localhost:4040"

.PHONY: tunnel-quick
tunnel-quick: ## Тимчасова публічна HTTPS-адреса без реєстрації (Cloudflare Quick Tunnel - щоразу нова; інакше srv.us - стала; інакше localhost.run)
	$(COMPOSE) --profile quick-tunnel up -d --force-recreate cloudflared
	@for i in $$(seq 30); do \
		url=$$($(COMPOSE) logs cloudflared 2>&1 | grep -oE 'https://[a-z0-9]+(-[a-z0-9]+)+\.trycloudflare\.com' | tail -n 1); \
		[ -n "$$url" ] && break; \
		$(COMPOSE) logs cloudflared 2>&1 | grep -q 'failed to request quick Tunnel' && break; \
		sleep 1; \
	done; \
	if [ -z "$$url" ]; then \
		$(COMPOSE) --profile quick-tunnel stop cloudflared >/dev/null 2>&1; \
		echo "Cloudflare недоступний з цієї мережі - тунель srv.us" >&2; \
		$(COMPOSE) --profile quick-tunnel up -d --force-recreate srvus >/dev/null 2>&1; \
		for i in $$(seq 30); do \
			url=$$($(COMPOSE) logs srvus 2>&1 | grep -oE 'https://[a-z0-9]+\.srv\.us' | tail -n 1); \
			[ -n "$$url" ] && break; \
			sleep 1; \
		done; \
	fi; \
	if [ -z "$$url" ]; then \
		$(COMPOSE) --profile quick-tunnel stop srvus >/dev/null 2>&1; \
		echo "srv.us недоступний - тунель localhost.run (повільний для великих файлів)" >&2; \
		$(COMPOSE) --profile quick-tunnel up -d --force-recreate localhostrun >/dev/null 2>&1; \
		for i in $$(seq 30); do \
			url=$$($(COMPOSE) logs localhostrun 2>&1 | grep -oE 'https://[a-z0-9]+\.lhr\.life' | tail -n 1); \
			[ -n "$$url" ] && break; \
			sleep 1; \
		done; \
	fi; \
	if [ -z "$$url" ]; then echo "Публічну адресу не отримано: $(COMPOSE) logs cloudflared srvus localhostrun" >&2; exit 1; fi; \
	echo "Платформа: $$url   мінт-сайт: $$url/sites/<адреса>   зупинити: make tunnel-stop"; \
	if [ "$(PUBLIC_URL)" != "$$url" ]; then echo "Галерея й метадані беруть адреси з PUBLIC_URL: задайте PUBLIC_URL=$$url у .env і виконайте make up" >&2; fi

.PHONY: tunnel-stop
tunnel-stop: ## Зупинити публічний тунель (ngrok, Cloudflare, srv.us або localhost.run)
	$(COMPOSE) --profile tunnel --profile quick-tunnel stop ngrok cloudflared srvus localhostrun

.PHONY: chain
chain: ## Локальний блокчейн Anvil: http://localhost:8545, chainId 31337
	$(COMPOSE) --profile chain up -d anvil
	@echo "MetaMask: додати мережу http://localhost:8545 (chainId 31337); тестові ключі - make chain-keys"

.PHONY: chain-keys
chain-keys: ## Адреси й приватні ключі тестових акаунтів Anvil
	@$(COMPOSE) --profile chain logs anvil | sed -n '/Available Accounts/,/Wallet/p'

.PHONY: contracts
contracts: ## Скомпілювати CollectionNFT.sol (solc у Docker) і згенерувати Go-біндинг (abigen)
	@test -d contracts/node_modules/@openzeppelin || (cd contracts && npm ci)
	@mkdir -p contracts/bin
	docker run --rm -u $(DOCKER_USER) -v "$(CURDIR)/contracts":/sources:ro ethereum/solc:$(SOLC_VERSION) \
		--optimize --optimize-runs 200 --base-path /sources --include-path /sources/node_modules \
		--combined-json abi,bin /sources/CollectionNFT.sol > contracts/bin/CollectionNFT.json
	python3 scripts/split-solc-output.py contracts/bin/CollectionNFT.json CollectionNFT contracts/bin
	docker run --rm -u $(DOCKER_USER) -e CGO_ENABLED=0 -e HOME=/tmp -e GOPATH=/tmp/go -e GOCACHE=/tmp/gocache \
		$(if $(GOMODCACHE),-v "$(GOMODCACHE)":/gomod -e GOMODCACHE=/gomod) -v "$(CURDIR)":/src -w /tmp golang:1.17-alpine \
		go run github.com/ethereum/go-ethereum/cmd/abigen@$(GO_ETHEREUM_VERSION) \
		--abi /src/contracts/bin/CollectionNFT.abi --bin /src/contracts/bin/CollectionNFT.bin \
		--pkg nftcontracts --type CollectionNFT --out /src/internal/nftcontracts/collectionnft.go
	@chmod 644 internal/nftcontracts/collectionnft.go

.PHONY: config
config: ## Створити configs/config.json для запуску без Docker (наявний не перезаписується)
	@mkdir -p configs
	@if [ -f configs/config.json ]; then echo "configs/config.json уже є"; \
		else cp deploy/config.json configs/config.json && echo "Створено configs/config.json"; fi

.PHONY: run
run: config ## Запустити сервер без Docker (потрібні Go і зібраний web/app/dist)
	@mkdir -p web/public-pages/vendor
	cp web/app/node_modules/ethers/dist/ethers.umd.min.js web/public-pages/vendor/
	$(COMPOSE) up -d --wait db
	go run ./cmd run --config=./configs

.PHONY: test
test: ## Go-тести з детектором гонок, зокрема інтеграційні з PostgreSQL з compose
	$(COMPOSE) up -d --wait db
	@$(COMPOSE) exec -T db createdb -U postgres generatordb_test 2>/dev/null || true
	GENERATOR_TEST_DB="$(TEST_DB)" CGO_ENABLED=1 go test -race ./...

.PHONY: smoke
smoke: ## Перевірити розгорнуту платформу через API (BASE=https://... для ngrok)
	node scripts/smoke.js

.PHONY: e2e
e2e: ## Наскрізні тести Playwright проти запущеної платформи (make up; для розгортання контракту - make chain і мережа Anvil)
	@test -d node_modules/@playwright/test || npm ci
	node node_modules/@playwright/test/cli.js install chromium-headless-shell
	BASE=$${BASE:-http://localhost:$(APP_PORT)} node node_modules/@playwright/test/cli.js test

KEY_ENV ?= DEPLOYER_PRIVATE_KEY

.PHONY: pin-filebase
pin-filebase: ## Закріпити опубліковану в IPFS колекцію у Filebase (CAR через S3 API): PROJECT=<id>; FILEBASE_BUCKET, FILEBASE_S3_KEY, FILEBASE_S3_SECRET - у .env
	@test -n "$(PROJECT)" || { echo "Вкажіть PROJECT=<ідентифікатор колекції>" >&2; exit 1; }
	sh scripts/pin-filebase.sh $(PROJECT)

.PHONY: deploy-contract
deploy-contract: ## Розгорнути контракт колекції з командного рядка (резерв до сторінки контракту): PROJECT=<id>; ключ власника - у змінній KEY_ENV (DEPLOYER_PRIVATE_KEY); PRICE=<ETH>, MAX_PER_MINT=<n> - умови продажу; TIP=<wei> - надбавка за газ; DRY_RUN=1 - лише оцінка газу
	@test -n "$(PROJECT)" || { echo "Вкажіть PROJECT=<ідентифікатор колекції>" >&2; exit 1; }
	$(COMPOSE) exec -T -e $(KEY_ENV) app ./generator deploy-contract --config=./configs --project=$(PROJECT) --key-env=$(KEY_ENV) $(if $(PRICE),--price=$(PRICE)) $(if $(MAX_PER_MINT),--max-per-mint=$(MAX_PER_MINT)) $(if $(TIP),--tip=$(TIP)) $(if $(DRY_RUN),--dry-run)
