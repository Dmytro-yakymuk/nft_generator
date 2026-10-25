#!/bin/sh
# Закріплення опублікованої колекції у Filebase, щоб її файли були доступні в IPFS і без вузла платформи.
# Каталоги зображень і метаданих вивантажуються з вузла IPFS платформи як CAR-файли (корінь - той самий CID,
# що в базовому URI контракту) і завантажуються в IPFS-бакет Filebase через S3 API з метаданими import=car;
# Filebase повертає CID кореня в заголовку x-amz-meta-cid.
#
#   make pin-filebase PROJECT=<id>   # FILEBASE_BUCKET, FILEBASE_S3_KEY, FILEBASE_S3_SECRET - у .env
set -eu
cd "$(dirname "$0")/.."

PROJECT=${1:-}
case "$PROJECT" in
  '' | *[!0-9a-f-]*) echo "Вкажіть ідентифікатор колекції: make pin-filebase PROJECT=<id>" >&2; exit 1 ;;
esac

# Лише потрібні змінні з .env, без виконання файлу; ключі передаються curl через stdin, а не аргументами.
env_value() { grep -E "^$1=" .env 2>/dev/null | tail -n 1 | cut -d= -f2- | sed -e 's/^"//' -e 's/"$//'; }
BUCKET=${FILEBASE_BUCKET:-$(env_value FILEBASE_BUCKET)}
KEY=${FILEBASE_S3_KEY:-$(env_value FILEBASE_S3_KEY)}
SECRET=${FILEBASE_S3_SECRET:-$(env_value FILEBASE_S3_SECRET)}
if [ -z "$BUCKET" ] || [ -z "$KEY" ] || [ -z "$SECRET" ]; then
  echo "Потрібні FILEBASE_BUCKET, FILEBASE_S3_KEY і FILEBASE_S3_SECRET (сторінка Access Keys консолі Filebase)" >&2
  exit 1
fi

CIDS=$(docker compose exec -T db psql -U postgres generator -tAc \
  "select images_cid || ' ' || metadata_cid from collections where convert_from(id, 'utf8') = '$PROJECT'")
if [ -z "$(echo "$CIDS" | tr -d ' ')" ]; then
  echo "Колекцію не знайдено або її ще не опубліковано в IPFS (кнопка Publish to IPFS)" >&2
  exit 1
fi

WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
for cid in $CIDS; do
  docker compose exec -T ipfs ipfs dag export "$cid" > "$WORK/$cid.car"
  printf 'user = "%s:%s"\n' "$KEY" "$SECRET" | curl -sS --fail-with-body -K - --aws-sigv4 "aws:amz:us-east-1:s3" \
    -H "x-amz-meta-import: car" -T "$WORK/$cid.car" -D "$WORK/$cid.headers" -o "$WORK/$cid.body" \
    "https://s3.filebase.com/$BUCKET/$cid.car" || { cat "$WORK/$cid.body" >&2; exit 1; }
  pinned=$(grep -i '^x-amz-meta-cid:' "$WORK/$cid.headers" | cut -d: -f2- | tr -d ' \r')
  size=$(wc -c < "$WORK/$cid.car")
  if [ "$pinned" = "$cid" ]; then
    echo "закріплено: $cid ($size Б)"
  else
    echo "Filebase повернув CID ${pinned:-(немає)} замість $cid" >&2
    exit 1
  fi
done
echo "Перевірка: https://ipfs.filebase.io/ipfs/<CID>/0 або будь-який публічний шлюз IPFS"
