"""Колаж з перших згенерованих зображень колекції для рисунка роботи (без добору вручну).

    python3 scripts/demo-art/collage.py <каталог з PNG> <результат.png> [стовпців] [рядків]

Зображення беруться в порядку номерів токенів (імена файлів - номери токенів, як у TokenMetadata);
потрібен Pillow (є в середовищі з matplotlib).
"""

import sys
from pathlib import Path

from PIL import Image

source, target = Path(sys.argv[1]), Path(sys.argv[2])
columns = int(sys.argv[3]) if len(sys.argv) > 3 else 4
rows = int(sys.argv[4]) if len(sys.argv) > 4 else 4
tile, gap = 240, 8

files = sorted(source.glob("*.png"), key=lambda f: (len(f.stem), f.stem))[: columns * rows]
sheet = Image.new("RGB", (columns * tile + (columns - 1) * gap, rows * tile + (rows - 1) * gap), "white")
for i, file in enumerate(files):
    image = Image.open(file).convert("RGB").resize((tile, tile), Image.LANCZOS)
    sheet.paste(image, ((i % columns) * (tile + gap), (i // columns) * (tile + gap)))
sheet.save(target, optimize=True)
print(f"{len(files)} зображень -> {target}")
