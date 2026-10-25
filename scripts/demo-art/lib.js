// Спільні помічники процедурного набору шарів: полотно SVG 500x500, градієнти, фільтри, фігури,
// детермінований генератор випадкових чисел і геометрія персонажа (усі шари малюються в одних координатах).

const SIZE = 500;
const OL = '#2b2140'; // колір контуру персонажа й аксесуарів
const SW = 6; // товщина основного контуру

// Геометрія персонажа: голова-еліпс на плечах, очі, рот, щоки.
const HEAD = { cx: 250, cy: 245, rx: 138, ry: 126 };
const TORSO = 'M108 500 C108 420 160 360 250 360 C340 360 392 420 392 500 Z';
const EYES = [[200, 250], [300, 250]];
const MOUTH = [250, 322];
const NOSE = [250, 292];
const CHEEKS = [[160, 302], [340, 302]];
const HEAD_ELLIPSE = `cx="${HEAD.cx}" cy="${HEAD.cy}" rx="${HEAD.rx}" ry="${HEAD.ry}"`;

/** Детермінований генератор (mulberry32), щоб набір був однаковим при кожному запуску. */
function rng(seed) {
    let a = seed >>> 0;
    return () => {
        a = (a + 0x6d2b79f5) >>> 0;
        let t = a;
        t = Math.imul(t ^ (t >>> 15), t | 1);
        t ^= t + Math.imul(t ^ (t >>> 7), t | 61);
        return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
    };
}

const svg = (body, defs = '') =>
    `<svg xmlns="http://www.w3.org/2000/svg" width="${SIZE}" height="${SIZE}" viewBox="0 0 ${SIZE} ${SIZE}">` +
    (defs ? `<defs>${defs}</defs>` : '') + body + '</svg>';

const stopsOf = (stops) => stops.map(([offset, color, opacity = 1]) =>
    `<stop offset="${offset}" stop-color="${color}" stop-opacity="${opacity}"/>`).join('');

/** Лінійний градієнт; за замовчуванням - згори донизу. */
const lin = (id, stops, x1 = 0, y1 = 0, x2 = 0, y2 = 1) =>
    `<linearGradient id="${id}" x1="${x1}" y1="${y1}" x2="${x2}" y2="${y2}">${stopsOf(stops)}</linearGradient>`;

/** Радіальний градієнт у частках фігури (fx, fy - фокус для відблиску). */
const rad = (id, stops, cx = 0.5, cy = 0.5, r = 0.5, fx = cx, fy = cy) =>
    `<radialGradient id="${id}" cx="${cx}" cy="${cy}" r="${r}" fx="${fx}" fy="${fy}">${stopsOf(stops)}</radialGradient>`;

/** Об'ємне зафарбовування: світла пляма вгорі ліворуч, тінь унизу праворуч. */
const shade = (id, light, base, dark) => rad(id, [[0, light], [0.55, base], [1, dark]], 0.4, 0.35, 0.75, 0.32, 0.25);

const blur = (id, sd) =>
    `<filter id="${id}" x="-50%" y="-50%" width="200%" height="200%"><feGaussianBlur stdDeviation="${sd}"/></filter>`;

const dropShadow = (id, dy = 6, sd = 5, opacity = 0.3) =>
    `<filter id="${id}" x="-30%" y="-30%" width="160%" height="160%">` +
    `<feDropShadow dx="0" dy="${dy}" stdDeviation="${sd}" flood-color="#1a1030" flood-opacity="${opacity}"/></filter>`;

/** Світіння: розмита копія під фігурою. */
const glow = (id, sd, color) =>
    `<filter id="${id}" x="-60%" y="-60%" width="220%" height="220%">` +
    `<feGaussianBlur in="SourceAlpha" stdDeviation="${sd}" result="b"/>` +
    `<feFlood flood-color="${color}"/><feComposite in2="b" operator="in" result="g"/>` +
    '<feMerge><feMergeNode in="g"/><feMergeNode in="g"/><feMergeNode in="SourceGraphic"/></feMerge></filter>';

/** Маска "усе, крім голови": одяг і прикраси на шиї лишаються позаду голови разом з її контуром. */
const behindHead = (id) =>
    `<mask id="${id}"><rect width="${SIZE}" height="${SIZE}" fill="#fff"/>` +
    `<ellipse cx="${HEAD.cx}" cy="${HEAD.cy}" rx="${HEAD.rx + SW / 2}" ry="${HEAD.ry + SW / 2}" fill="#000"/></mask>`;

const clipHead = (id) => `<clipPath id="${id}"><ellipse ${HEAD_ELLIPSE}/></clipPath>`;
const clipTorso = (id) => `<clipPath id="${id}"><path d="${TORSO}"/></clipPath>`;

/** Зірка з points променями. */
function star(cx, cy, outer, inner, points = 5, rotation = -90) {
    const d = [];
    for (let i = 0; i < points * 2; i++) {
        const r = i % 2 ? inner : outer;
        const a = ((rotation + (i * 180) / points) * Math.PI) / 180;
        d.push(`${(cx + r * Math.cos(a)).toFixed(1)} ${(cy + r * Math.sin(a)).toFixed(1)}`);
    }
    return `M${d.join(' L')} Z`;
}

/** Серце розміру s з центром у (cx, cy). */
const heart = (cx, cy, s) =>
    `M${cx} ${cy + s * 0.9} C${cx - s * 1.4} ${cy - s * 0.1} ${cx - s * 0.7} ${cy - s * 1.1} ${cx} ${cy - s * 0.45} ` +
    `C${cx + s * 0.7} ${cy - s * 1.1} ${cx + s * 1.4} ${cy - s * 0.1} ${cx} ${cy + s * 0.9} Z`;

/** Чотирипроменевий відблиск. */
const sparkle = (cx, cy, r) =>
    `M${cx} ${cy - r} Q${cx + r * 0.18} ${cy - r * 0.18} ${cx + r} ${cy} Q${cx + r * 0.18} ${cy + r * 0.18} ${cx} ${cy + r} ` +
    `Q${cx - r * 0.18} ${cy + r * 0.18} ${cx - r} ${cy} Q${cx - r * 0.18} ${cy - r * 0.18} ${cx} ${cy - r} Z`;

/** Віддзеркалення лівої половини праворуч (крила, сережки тощо). */
const mirror = (content) => `${content}<g transform="translate(${SIZE} 0) scale(-1 1)">${content}</g>`;

/** Варіант шару: ім'я (стає значенням атрибута), вага і SVG. */
const v = (name, weight, image) => ({ name, weight, image });

module.exports = {
    SIZE, OL, SW, HEAD, TORSO, EYES, MOUTH, NOSE, CHEEKS, HEAD_ELLIPSE,
    rng, svg, lin, rad, shade, blur, dropShadow, glow, behindHead, clipHead, clipTorso,
    star, heart, sparkle, mirror, v,
};
