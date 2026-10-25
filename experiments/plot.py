"""Графіки й зведені таблиці експериментів (розділ 3 роботи) з CSV у experiments/results.

    python3 experiments/plot.py              # потрібен matplotlib; PNG - у ../документи/thesis-images/ (поза репозиторієм), зведення - у stdout
    python3 experiments/plot.py hashlips     # лише вибрані: назви функцій без plot_ / summarize_

Кольори - перші три слоти перевіреної категорійної палітри (синій, помаранчевий, бірюзовий);
для чорно-білого друку серії додатково розрізняються маркером і типом лінії та підписані на графіку.
"""

import csv
import math
import statistics
import sys
from collections import defaultdict
from pathlib import Path

import matplotlib

matplotlib.use("Agg")
import matplotlib.pyplot as plt  # noqa: E402
from matplotlib.ticker import FuncFormatter, LogLocator  # noqa: E402

ROOT = Path(__file__).resolve().parent.parent
RESULTS = ROOT / "experiments" / "results"
# Рисунки роботи лежать поза репозиторієм - у каталозі документів поруч із ним
IMAGES = ROOT.parent / "документи" / "thesis-images"

INK, MUTED, GRID = "#0b0b0b", "#52514e", "#dddcd7"
SERIES = [  # колір, маркер, лінія
    ("#2a78d6", "o", "-"),
    ("#eb6834", "s", "--"),
    ("#1baf7a", "^", ":"),
]

plt.rcParams.update({
    "font.family": "Liberation Serif",  # метрично сумісний з Times New Roman
    "font.size": 12,
    "axes.edgecolor": MUTED,
    "axes.labelcolor": INK,
    "axes.spines.top": False,
    "axes.spines.right": False,
    "xtick.color": MUTED,
    "ytick.color": MUTED,
    "xtick.labelcolor": INK,
    "ytick.labelcolor": INK,
    "axes.grid": True,
    "grid.color": GRID,
    "grid.linewidth": 0.6,
    "lines.linewidth": 2,
    "lines.markersize": 7,
    "savefig.dpi": 200,
    "mathtext.fontset": "custom",  # степені десяти тим самим шрифтом, що й решта підписів
    "mathtext.rm": "Liberation Serif",
    "mathtext.it": "Liberation Serif:italic",
})


def read(name):
    with open(RESULTS / name, newline="", encoding="utf-8") as file:
        return list(csv.DictReader(file))


def summarize(rows, keys, value):
    """Середнє, стандартне відхилення й кількість повторів value для кожної комбінації keys."""
    groups = defaultdict(list)
    for row in rows:
        if row.get("status", "ok") == "ok" and row[value] != "":
            groups[tuple(row[k] for k in keys)].append(float(row[value]))
    return {
        key: (statistics.mean(v), statistics.stdev(v) if len(v) > 1 else 0.0, len(v))
        for key, v in groups.items()
    }


def uk_number(value, digits=0):
    """Число з пробілом між тисячами й комою в дробовій частині."""
    text = f"{value:,.{digits}f}".replace(",", " ").replace(".", ",")
    return text


def power_of_ten(value, _):
    exponent = round(math.log10(value))
    return f"$10^{{{exponent}}}$"


def label_end(ax, x, y, text, color, dy=0):
    ax.annotate(text, (x, y), xytext=(6, dy), textcoords="offset points", color=INK,
                va="center", ha="left", fontsize=11,
                bbox=dict(boxstyle="round,pad=0.15", fc="white", ec=color, lw=1))


def save(fig, name):
    IMAGES.mkdir(exist_ok=True)
    path = IMAGES / name
    fig.savefig(path, bbox_inches="tight", facecolor="white")
    plt.close(fig)
    print(f"  -> {path.relative_to(ROOT.parent)}")


def plot_conversion():
    stats = summarize(read("conversion.csv"), ["combinations", "method"], "ns_per_conversion")
    totals = sorted({int(c) for c, _ in stats})

    print("\nE3. Перетворення номера комбінації, нс на одне перетворення (середнє ± σ, 5 повторів)")
    print("C;покрокове;пряме;прискорення")
    for c in totals:
        inc, dir_ = stats[(str(c), "increment")], stats[(str(c), "direct")]
        print(f"10^{len(str(c)) - 1};{inc[0]:.4g} ± {inc[1]:.2g};{dir_[0]:.1f} ± {dir_[1]:.1f};{inc[0] / dir_[0]:.3g}")

    fig, ax = plt.subplots(figsize=(6.3, 3.9))
    for (name, title), (color, marker, style) in zip(
            [("increment", "покрокове збільшення, O(C)"), ("direct", "пряме обчислення, O(n)")], SERIES):
        means = [stats[(str(c), name)][0] / 1e6 for c in totals]  # мс
        ax.plot(totals, means, color=color, marker=marker, linestyle=style, label=title)
        label_end(ax, totals[-1], means[-1], title.split(",")[0], color)
    ax.set_xscale("log")
    ax.set_yscale("log")
    ax.xaxis.set_major_formatter(FuncFormatter(power_of_ten))
    ax.yaxis.set_major_locator(LogLocator(base=10, numticks=12))
    ax.yaxis.set_major_formatter(FuncFormatter(power_of_ten))
    ax.set_xlabel("Кількість комбінацій C (n шарів по 10 варіантів)")
    ax.set_ylabel("Час на одне зображення, мс")
    ax.legend(loc="upper left", frameon=False)
    ax.set_xlim(totals[0] / 2, totals[-1] * 40)
    save(fig, "exp-conversion.png")


def plot_sampling():
    rows = read("sampling.csv")
    stats_ms = summarize(rows, ["amount", "method"], "ms")
    stats_draws = summarize(rows, ["amount", "method"], "draws")
    amounts = sorted({int(a) for a, _ in stats_ms})
    methods = [("reseed", "повторне ініціювання генератора (було)"),
               ("rejection", "відкидання повторів, одне ініціювання"),
               ("floyd", "алгоритм Флойда (стало)")]

    print("\nE3b. Вибір N різних комбінацій з C = 100 000 (середнє з 5 повторів)")
    print("N;" + ";".join(f"{m} мс;{m} спроб" for m, _ in methods))
    for a in amounts:
        cells = []
        for m, _ in methods:
            if (str(a), m) in stats_ms:
                cells += [f"{stats_ms[(str(a), m)][0]:.4g}", f"{stats_draws[(str(a), m)][0]:.4g}"]
            else:
                cells += ["—", "—"]
        print(f"{a};" + ";".join(cells))

    fig, ax = plt.subplots(figsize=(6.3, 3.9))
    for (name, title), (color, marker, style) in zip(methods, SERIES):
        points = [(a, stats_ms[(str(a), name)][0]) for a in amounts if (str(a), name) in stats_ms]
        xs, ys = zip(*points)
        ax.plot(xs, ys, color=color, marker=marker, linestyle=style, label=title)
    ax.set_yscale("log")
    ax.yaxis.set_major_locator(LogLocator(base=10, numticks=12))
    ax.yaxis.set_major_formatter(FuncFormatter(power_of_ten))
    ax.xaxis.set_major_formatter(FuncFormatter(lambda v, _: uk_number(v)))
    ax.set_xlabel("Кількість зображень N (з C = 100 000 комбінацій)")
    ax.set_ylabel("Час вибору комбінацій, мс")
    ax.legend(loc="lower right", frameon=True, facecolor="white", edgecolor="white", framealpha=1, fontsize=11)
    save(fig, "exp-sampling.png")


def plot_gas():
    rows = read("gas.csv")
    data = {(r["whitelist_size"], r["variant"]): r for r in rows}
    sizes = sorted({int(r["whitelist_size"]) for r in rows})

    print("\nE4/E5. Газ whitelist (Anvil): запис у контракт, транзакцій; мінт 1 токена; довжина доказу")
    print("N;Меркл запис;Меркл мінт;доказ;mapping запис;транзакцій;mapping мінт;різниця мінту")
    for n in sizes:
        m, l = data[(str(n), "merkle")], data[(str(n), "mapping")]
        print(f"{n};{m['setup_gas']};{m['mint_gas_mean']};{m['proof_length']};{l['setup_gas']};{l['setup_txs']};"
              f"{l['mint_gas_mean']};{int(m['mint_gas_mean']) - int(l['mint_gas_mean'])}")

    fig, ax = plt.subplots(figsize=(6.3, 3.9))
    for (variant, title), (color, marker, style) in zip(
            [("mapping", "перелік адрес у mapping"), ("merkle", "корінь дерева Меркла")], SERIES[1::-1]):
        ys = [int(data[(str(n), variant)]["setup_gas"]) for n in sizes]
        ax.plot(sizes, ys, color=color, marker=marker, linestyle=style, label=title)
    ax.set_xscale("log")
    ax.set_yscale("log")
    ax.xaxis.set_major_formatter(FuncFormatter(lambda v, _: uk_number(v)))
    ax.yaxis.set_major_formatter(FuncFormatter(power_of_ten))
    ax.set_xlabel("Кількість адрес у whitelist N")
    ax.set_ylabel("Газ запису whitelist у контракт")
    ax.legend(loc="upper left", frameon=False)
    save(fig, "exp-gas.png")


def plot_generation():
    path = RESULTS / "generation.csv"
    if not path.exists():
        print("\n(generation.csv ще немає - E1/E2 пропущено)")
        return
    rows = read("generation.csv")
    seconds = summarize(rows, ["experiment", "layers", "amount"], "seconds")
    sizes = summarize(rows, ["experiment", "layers", "amount"], "bytes")
    variants = [("", "послідовно, без кешу шарів"), ("-parallel", "12 воркерів, кеш шарів")]

    for experiment, key_index, xlabel, name, log in [
            ("E1", 2, "Кількість зображень N (5 шарів по 7 варіантів)", "exp-generation-amount.png", True),
            ("E2", 1, "Кількість шарів n (по 10 варіантів, N = 1000)", "exp-generation-layers.png", False)]:
        fig, ax = plt.subplots(figsize=(6.3, 3.9))
        for (suffix, title), (color, marker, style) in zip(variants, SERIES):
            keys = sorted((k for k in seconds if k[0] == experiment + suffix), key=lambda k: int(k[key_index]))
            if not keys:
                continue
            print(f"\n{experiment}{suffix}. Час генерації, с (середнє ± σ), мс на зображення, обсяг, МіБ")
            for k in keys:
                mean, sd, count = seconds[k]
                amount = int(k[2])
                print(f"n={k[1]} N={amount};{mean:.2f} ± {sd:.2f} ({count});{mean / amount * 1000:.1f};"
                      f"{sizes[k][0] / 2 ** 20:.1f}")
            xs = [int(k[key_index]) for k in keys]
            ys = [seconds[k][0] for k in keys]
            ax.plot(xs, ys, color=color, marker=marker, linestyle=style, label=title)
            for x, y in zip(xs, ys):
                ax.annotate(uk_number(y, 1), (x, y), xytext=(0, 8), textcoords="offset points",
                            ha="center", fontsize=10, color=INK)
        if log:
            ax.set_xscale("log")
            ax.set_yscale("log")
            ax.yaxis.set_major_formatter(FuncFormatter(lambda v, _: uk_number(v, 0 if v >= 1 else 1)))
        else:
            ax.set_ylim(0, None)
        ax.set_xticks(xs)
        ax.xaxis.set_major_formatter(FuncFormatter(lambda v, _: uk_number(v)))
        ax.set_xlabel(xlabel)
        ax.set_ylabel("Час генерації, с")
        ax.legend(loc="upper left", frameon=False, fontsize=11)
        save(fig, name)


def plot_frequencies():
    rows = read("frequencies.csv")
    print("\nE8. Відповідність частот вагам: хі-квадрат для кожного з шарів набору, 5 прогонів")
    print("N;тестів;p < 0,05;частка;мін. p;вибір, мс (середнє);спроб на зображення")
    for amount in sorted({int(r["amount"]) for r in rows}):
        group = [r for r in rows if int(r["amount"]) == amount]
        rejected = sum(float(r["p_value"]) < 0.05 for r in group)
        runs = {r["repeat"]: (float(r["choice_ms"]), int(r["draws"])) for r in group}
        choice = statistics.mean(ms for ms, _ in runs.values())
        draws = statistics.mean(d for _, d in runs.values()) / amount
        print(f"{amount};{len(group)};{rejected};{rejected / len(group):.3f};"
              f"{min(float(r['p_value']) for r in group):.4f};{choice:.2f};{draws:.3f}")

    detail = read("frequencies-detail.csv")
    amount = sum(int(r["count"]) for r in detail if r["layer"] == detail[0]["layer"])
    expected = [float(r["expected_share"]) for r in detail]
    observed = [float(r["observed_share"]) for r in detail]
    worst = max(detail, key=lambda r: abs(float(r["observed_share"]) - float(r["expected_share"]))
                / math.sqrt(float(r["expected_share"]) * (1 - float(r["expected_share"])) / amount))
    print(f"Варіантів: {len(detail)}, N = {amount}; найбільше відхилення - {worst['layer']}/{worst['variant']}: "
          f"{float(worst['observed_share']):.4f} проти {float(worst['expected_share']):.4f}")

    fig, ax = plt.subplots(figsize=(6.3, 4.4))
    low, high = 0.002, 1
    grid = [low * (high / low) ** (i / 100) for i in range(101)]
    # смуга ±2σ біноміального розподілу частки для N зображень
    ax.fill_between(grid, [max(p - 2 * math.sqrt(p * (1 - p) / amount), low / 2) for p in grid],
                    [p + 2 * math.sqrt(p * (1 - p) / amount) for p in grid], color=GRID, lw=0, label="±2σ для N = 10 000")
    ax.plot([low, high], [low, high], color=MUTED, lw=1, linestyle="--", label="частота = ймовірність")
    color, marker, _ = SERIES[0]
    ax.scatter(expected, observed, s=16, color=color, edgecolors="white", linewidths=0.6, zorder=3,
               label=f"{len(detail)} варіантів {len({r['layer'] for r in detail})} шарів набору")
    ax.set_xscale("log")
    ax.set_yscale("log")
    ax.set_xlim(low, high)
    ax.set_ylim(low, high)
    percent = FuncFormatter(lambda v, _: uk_number(v * 100, 1 if v < 0.01 else 0) + " %")
    ax.xaxis.set_major_formatter(percent)
    ax.yaxis.set_major_formatter(percent)
    ax.set_xlabel("Задана ймовірність варіанта (вага / сума ваг шару)")
    ax.set_ylabel("Частка зображень з варіантом")
    ax.legend(loc="upper left", frameon=False, fontsize=11)
    save(fig, "exp-frequencies.png")


# Нерівномірність ваг s - порядкова величина, тому серії - кроки однієї синьої гами (світла -> темна).
SKEW = [("0.0", "#86b6ef", "o", "-"), ("0.5", "#3987e5", "s", "--"), ("1.0", "#1c5cab", "^", ":"),
        ("2.0", "#0d366b", "D", "-.")]


def plot_uniqueness():
    rows = read("uniqueness.csv")
    draws = defaultdict(list)
    top = defaultdict(list)
    ok = defaultdict(int)
    for r in rows:
        key = (r["zipf_s"], float(r["ratio"]))
        draws[key].append(int(r["draws"]) / int(r["amount"]))
        top[key].append(float(r["top_share"]))
        ok[key] += r["ok"] == "true"
    expected_top = {r["zipf_s"]: float(r["expected_top_share"]) for r in rows}
    ratios = sorted({float(r["ratio"]) for r in rows})

    print("\nE9. Унікальність за вагами, C = 1296: спроб на зображення (середнє з 5); частка найімовірнішого варіанта")
    print("s;" + ";".join(f"N/C={uk_number(x, 2)}" for x in ratios))
    for s, *_ in SKEW:
        print(f"{s} спроб;" + ";".join(f"{statistics.mean(draws[(s, x)]):.2f}{'' if ok[(s, x)] == 5 else ' (ліміт)'}" for x in ratios))
        print(f"{s} частка (задана {expected_top[s]:.3f});" + ";".join(f"{statistics.mean(top[(s, x)]):.3f}" for x in ratios))

    fig, (left, right) = plt.subplots(1, 2, figsize=(9.6, 3.9))
    for s, color, marker, style in SKEW:
        title = "рівні ваги" if s == "0.0" else f"s = {uk_number(float(s), 1)}"
        done = [x for x in ratios if ok[(s, x)] == 5]
        left.plot(done, [statistics.mean(draws[(s, x)]) for x in done], color=color, marker=marker, linestyle=style, label=title)
        right.plot(ratios, [statistics.mean(top[(s, x)]) / expected_top[s] for x in ratios],
                   color=color, marker=marker, linestyle=style, label=title)
    left.set_yscale("log")
    left.yaxis.set_major_formatter(FuncFormatter(lambda v, _: uk_number(v, 0 if v >= 1 else 1)))
    left.set_xlabel("Частка вибраних комбінацій N / C")
    left.set_ylabel("Спроб на одне зображення")
    left.legend(loc="upper left", frameon=False, fontsize=11)
    right.axhline(1, color=MUTED, lw=1, linestyle="--")
    right.set_ylim(0, 1.15)
    right.yaxis.set_major_formatter(FuncFormatter(lambda v, _: uk_number(v, 1)))
    right.set_xlabel("Частка вибраних комбінацій N / C")
    right.set_ylabel("Частка найімовірнішого варіанта /\nзадана ймовірність")
    for ax in (left, right):
        ax.xaxis.set_major_formatter(FuncFormatter(lambda v, _: uk_number(v, 1)))
    fig.tight_layout()
    save(fig, "exp-uniqueness.png")


def plot_parallel():
    stats = summarize(read("parallel.csv"), ["workers"], "seconds")
    workers = sorted(int(w) for (w,) in stats)
    base = stats[("1",)][0]
    amount = int(read("parallel.csv")[0]["amount"])
    print(f"\nE11. Паралельна генерація {amount} зображень набору з 26 шарів: час, с (середнє ± σ); мс на зображення; прискорення")
    for w in workers:
        mean, sd, count = stats[(str(w),)]
        print(f"{w};{mean:.2f} ± {sd:.2f} ({count});{mean / amount * 1000:.1f};{base / mean:.2f}")

    fig, ax = plt.subplots(figsize=(6.3, 3.9))
    ax.plot([1, max(workers)], [1, max(workers)], color=MUTED, lw=1, linestyle="--", label="лінійне прискорення")
    color, marker, style = SERIES[0]
    speedups = [base / stats[(str(w),)][0] for w in workers]
    ax.plot(workers, speedups, color=color, marker=marker, linestyle=style, label="виміряне прискорення")
    for w, speedup in zip(workers, speedups):
        ax.annotate(uk_number(speedup, 1), (w, speedup), xytext=(0, -16), textcoords="offset points",
                    ha="center", fontsize=10, color=INK)
    ax.set_xticks(workers)
    ax.set_ylim(0, max(workers) + 0.5)
    ax.set_xlabel("Кількість воркерів (процесор: 8 ядер, 12 потоків)")
    ax.set_ylabel("Прискорення відносно одного воркера")
    ax.legend(loc="upper left", frameon=False)
    save(fig, "exp-parallel.png")


def summarize_web():
    rows = read("web.csv")
    print("\nE7. Головна сторінка (середнє з 5): передано / розпаковано, МіБ; FCP, LCP, load, с")
    for label in dict.fromkeys(r["label"] for r in rows):
        for profile in ("unthrottled", "slow-4g"):
            group = [r for r in rows if r["label"] == label and r["profile"] == profile]
            if not group:
                continue
            mean = lambda key: statistics.mean(float(r[key]) for r in group)
            print(f"{label};{profile};{mean('transferred_bytes') / 2 ** 20:.2f} / {mean('decoded_bytes') / 2 ** 20:.2f};"
                  f"{mean('fcp_ms') / 1000:.2f};{mean('lcp_ms') / 1000:.2f};{mean('load_ms') / 1000:.2f}")


def gamma_q(a, x):
    """Регуляризована верхня неповна гамма-функція Q(a, x) - хвіст розподілу хі-квадрат (як у тесті Go)."""
    if x <= 0:
        return 1.0
    if x < a + 1:
        total = term = 1 / a
        n = 1
        while abs(term) > abs(total) * 1e-15 and n < 1000:
            term *= x / (a + n)
            total += term
            n += 1
        return 1 - total * math.exp(-x + a * math.log(x) - math.lgamma(a))
    tiny = 1e-300
    b = x + 1 - a
    c, d = 1 / tiny, 1 / b
    h = d
    for n in range(1, 1000):
        an = -n * (n - a)
        b += 2
        d = an * d + b
        d = d if abs(d) > tiny else tiny
        c = b + an / c
        c = c if abs(c) > tiny else tiny
        d = 1 / d
        h *= d * c
        if abs(d * c - 1) < 1e-15:
            break
    return math.exp(-x + a * math.log(x) - math.lgamma(a)) * h


def summarize_art_generation():
    rows = [r for r in read("generation.csv") if r["experiment"] == "ART"]
    print("\nE10. Генерація демонстраційної колекції (26 шарів): час, с (середнє ± σ); мс на зображення; КіБ на зображення")
    for amount in sorted({int(r["amount"]) for r in rows}):
        group = [r for r in rows if int(r["amount"]) == amount]
        seconds = [float(r["seconds"]) for r in group]
        size = statistics.mean(int(r["bytes"]) for r in group)
        sd = statistics.stdev(seconds) if len(seconds) > 1 else 0
        print(f"N={amount} ({len(group)});{statistics.mean(seconds):.2f} ± {sd:.2f};"
              f"{statistics.mean(seconds) / amount * 1000:.1f};{size / amount / 1024:.0f};{size / 2 ** 20:.1f} МіБ")

    path = RESULTS / "frequencies-platform.csv"
    if not path.exists():
        return
    layers = defaultdict(list)
    for r in read("frequencies-platform.csv"):
        layers[r["layer"]].append((int(r["count"]), float(r["expected_share"]), int(r["tokens"])))
    rejected, worst = 0, (1, "")
    for layer, variants in layers.items():
        tokens = variants[0][2]
        statistic = sum((count - tokens * share) ** 2 / (tokens * share) for count, share, _ in variants)
        p_value = gamma_q((len(variants) - 1) / 2, statistic / 2)
        rejected += p_value < 0.05
        worst = min(worst, (p_value, layer))
    print(f"Частоти на платформі ({tokens} токенів): {rejected} з {len(layers)} шарів з p < 0,05; "
          f"найменше p = {worst[0]:.4f} ({worst[1]})")


def plot_hashlips():
    """E12: HashLips Art Engine проти платформи на тих самих наборах шарів."""
    hashlips = read("hashlips.csv")
    seconds = summarize(hashlips, ["set", "amount"], "seconds")
    png_bytes = summarize(hashlips, ["set", "amount"], "png_bytes")
    json_bytes = summarize(hashlips, ["set", "amount"], "json_bytes")
    retries = summarize(hashlips, ["set", "amount"], "dna_retries")
    platform = read("generation.csv")
    platform_seconds = summarize(platform, ["experiment", "amount"], "seconds")
    platform_bytes = summarize(platform, ["experiment", "amount"], "bytes")
    variants = [("", "послідовно, без кешу"), ("-1worker", "1 воркер"), ("-parallel", "12 воркерів")]

    print("\nE12. HashLips Art Engine (1 процес Node.js) і платформа: час генерації, с (середнє ± σ, повторів); "
          "мс на зображення; HashLips / платформа")
    for name, experiment in (("E1", "E1"), ("ART", "ART")):
        for amount in sorted({int(a) for s, a in seconds if s == name}):
            mean, sd, count = seconds[(name, str(amount))]
            cells = [f"HashLips {mean:.2f} ± {sd:.2f} ({count}), {mean / amount * 1000:.1f} мс"]
            for suffix, title in variants:
                key = (experiment + suffix, str(amount))
                if key in platform_seconds:
                    p_mean, p_sd, p_count = platform_seconds[key]
                    cells.append(f"{title} {p_mean:.2f} ± {p_sd:.2f} ({p_count}), {p_mean / amount * 1000:.1f} мс, "
                                 f"×{mean / p_mean:.2f}")
            print(f"{name} N={amount}; " + "; ".join(cells))
            key = (experiment + "-parallel", str(amount))
            if key in platform_bytes:
                print(f"    PNG, КіБ на зображення: HashLips {png_bytes[(name, str(amount))][0] / amount / 1024:.1f}, "
                      f"платформа {platform_bytes[key][0] / amount / 1024:.1f}; JSON HashLips "
                      f"{json_bytes[(name, str(amount))][0] / amount:.0f} Б; повторів DNA {retries[(name, str(amount))][0]:.0f}")

    amount = "1000"
    rows = [("HashLips Art Engine\n(Node.js, один процес)", None), ("Платформа\n(послідовно, без кешу)", ""),
            ("Платформа,\n1 воркер", "-1worker"), ("Платформа,\n12 воркерів", "-parallel")]
    fig, axes = plt.subplots(2, 1, figsize=(6.3, 4.6))
    for ax, (name, title) in zip(axes, [("E1", "5 синтетичних шарів (як в E1)"), ("ART", "26 художніх шарів (як в E10)")]):
        labels, values, colors = [], [], []
        for label, suffix in rows:
            value = seconds[(name, amount)][0] if suffix is None else platform_seconds.get((name + suffix, amount), (None,))[0]
            if value is None:
                continue
            labels.append(label)
            values.append(value)
            colors.append(SERIES[1][0] if suffix is None else SERIES[0][0])
        ys = list(range(len(values)))[::-1]
        ax.barh(ys, values, height=0.55, color=colors)
        for y, value in zip(ys, values):
            ax.annotate(f"{uk_number(value, 1)} с", (value, y), xytext=(4, 0), textcoords="offset points",
                        va="center", ha="left", fontsize=11, color=INK)
        ax.set_yticks(ys)
        ax.set_yticklabels(labels, fontsize=10)
        ax.set_xlim(0, max(values) * 1.18)
        ax.grid(axis="y", visible=False)
        ax.set_title(title, loc="left", fontsize=12, color=INK)
        ax.xaxis.set_major_formatter(FuncFormatter(lambda v, _: uk_number(v)))
    axes[-1].set_xlabel("Час генерації 1000 зображень 500×500 пікселів, с")
    fig.tight_layout()
    save(fig, "exp-hashlips.png")


if __name__ == "__main__":
    IMAGES.mkdir(exist_ok=True)
    plots = (plot_conversion, plot_sampling, plot_gas, plot_generation, plot_frequencies, plot_uniqueness,
             summarize_art_generation, plot_parallel, summarize_web, plot_hashlips)
    only = set(sys.argv[1:])
    for plot in plots:
        if only and plot.__name__.split("_", 1)[1] not in only:
            continue
        try:
            plot()
        except FileNotFoundError as err:
            print(f"\n(немає {err.filename} - пропущено)", file=sys.stderr)
