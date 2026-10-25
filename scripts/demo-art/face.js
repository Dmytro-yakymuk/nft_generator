// Шари обличчя: щоки, позначки, ніс, рот, предмет у роті, очі, райдужка, зачіска, брови, окуляри.

const {
    OL, HEAD, EYES, MOUTH, NOSE, CHEEKS, rng, svg, lin, rad, shade, blur, glow, star, heart, sparkle, mirror, v,
} = require('./lib');

const [L, R] = EYES;
const [MX, MY] = MOUTH;
const [NX, NY] = NOSE;

// ---------- Щоки ----------

const bothCheeks = (draw, defs = '') => svg(CHEEKS.map(([x, y]) => draw(x, y)).join(''), defs);

const CHEEK = {
    name: 'Cheeks',
    none: 400,
    variants: [
        v('Pink Blush', 100, bothCheeks((x, y) => `<ellipse cx="${x}" cy="${y}" rx="30" ry="17" fill="#ff6fa8" opacity=".45" filter="url(#b)"/>`, blur('b', 4))),
        v('Coral Blush', 70, bothCheeks((x, y) => `<ellipse cx="${x}" cy="${y}" rx="28" ry="15" fill="#ff7a5c" opacity=".5" filter="url(#b)"/>` +
            `<path d="M${x - 12} ${y - 4} l6 8 M${x - 2} ${y - 4} l6 8 M${x + 8} ${y - 4} l6 8" stroke="#e8505b" stroke-width="2.5" stroke-linecap="round" opacity=".7"/>`, blur('b', 3))),
        v('Freckles', 70, bothCheeks((x, y) => {
            const r = rng(x);
            return Array.from({ length: 7 }, () => `<circle cx="${(x - 18 + r() * 36).toFixed(1)}" cy="${(y - 10 + r() * 18).toFixed(1)}" r="${(1.8 + r() * 1.4).toFixed(1)}" fill="#a0522d" opacity=".55"/>`).join('');
        })),
        v('Heart Cheeks', 40, bothCheeks((x, y) => `<path d="${heart(x, y, 11)}" fill="#ff4d6d" stroke="${OL}" stroke-width="2.5" opacity=".9"/>`)),
        v('Spirals', 20, bothCheeks((x, y) =>
            `<path d="M${x} ${y} m-2 0 a2 2 0 1 1 4 0 a5 5 0 1 1 -9 0 a8 8 0 1 1 15 0 a11 11 0 1 1 -20 0" fill="none" stroke="#ff6fa8" stroke-width="3" stroke-linecap="round" opacity=".8"/>`)),
        v('Whiskers', 50, svg(mirror(`<g stroke="${OL}" stroke-width="3.5" stroke-linecap="round" opacity=".8">` +
            '<path d="M178 300 L118 288"/><path d="M178 310 L116 312"/><path d="M180 320 L122 334"/></g>'))),
        v('Star Cheeks', 25, bothCheeks((x, y) => `<path d="${star(x, y, 12, 5)}" fill="#ffd23f" stroke="${OL}" stroke-width="2.5" stroke-linejoin="round"/>`)),
    ],
};

// ---------- Позначки на обличчі ----------

const FACE_MARK = {
    name: 'Face Mark',
    none: 700,
    variants: [
        v('Scar', 40, svg(`<path d="M${R[0] + 18} ${R[1] - 62} L${R[0] + 4} ${R[1] + 52}" stroke="#b5485a" stroke-width="6" stroke-linecap="round" opacity=".85"/>` +
            `<g stroke="#b5485a" stroke-width="3.5" stroke-linecap="round" opacity=".85">${[-40, -14, 24].map((d) =>
                `<path d="M${R[0] + 2 + d * -0.12} ${R[1] + d} l18 3"/>`).join('')}</g>`)),
        v('Star Tattoo', 50, svg(`<path d="${star(R[0] + 26, R[1] + 52, 11, 4.5)}" fill="#2f6fd6" opacity=".85"/>`)),
        v('Tear Drop', 40, svg(`<path d="M${L[0] - 6} ${L[1] + 46} q-10 14 0 20 q10 -6 0 -20 Z" fill="#4fb3ff" stroke="${OL}" stroke-width="2.5"/>`)),
        v('Bandage', 50, svg([-30, 30].map((angle) => `<g transform="translate(332 168) rotate(${angle})">` +
            `<rect x="-30" y="-10" width="60" height="20" rx="8" fill="#ffe0bd" stroke="${OL}" stroke-width="3"/>` +
            '<rect x="-9" y="-10" width="18" height="20" fill="#f2c9a0"/></g>').join('') +
            '<g fill="#d4a373"><circle cx="328" cy="164" r="1.6"/><circle cx="336" cy="164" r="1.6"/><circle cx="328" cy="172" r="1.6"/><circle cx="336" cy="172" r="1.6"/></g>')),
        v('War Paint', 30, svg(mirror(`<g stroke="#e8243c" stroke-width="7" stroke-linecap="round" opacity=".85"><path d="M150 288 L188 284"/><path d="M152 304 L190 300"/></g>`))),
        v('Lightning', 20, svg(`<path d="M262 128 L244 162 L256 162 L240 196 L272 152 L258 152 L270 128 Z" fill="#ffd23f" stroke="${OL}" stroke-width="3" stroke-linejoin="round"/>`)),
        v('Heart Tattoo', 40, svg(`<path d="${heart(L[0] - 34, L[1] + 50, 9)}" fill="#ff4d6d" opacity=".85"/>`)),
    ],
};

// ---------- Ніс ----------

const NOSE_LAYER = {
    name: 'Nose',
    none: 450,
    variants: [
        v('Button', 120, svg(`<ellipse cx="${NX}" cy="${NY}" rx="11" ry="8" fill="#2b1838" opacity=".22"/><ellipse cx="${NX - 3}" cy="${NY - 3}" rx="4" ry="2.5" fill="#fff" opacity=".5"/>`)),
        v('Cat Nose', 80, svg(`<path d="M${NX - 11} ${NY - 6} Q${NX} ${NY - 10} ${NX + 11} ${NY - 6} Q${NX + 4} ${NY + 6} ${NX} ${NY + 7} Q${NX - 4} ${NY + 6} ${NX - 11} ${NY - 6} Z" fill="#ff8fb3" stroke="${OL}" stroke-width="3" stroke-linejoin="round"/>`)),
        v('Dot', 100, svg(`<circle cx="${NX}" cy="${NY}" r="4" fill="${OL}" opacity=".7"/>`)),
        v('Clown', 15, svg(`<circle cx="${NX}" cy="${NY}" r="17" fill="url(#c)" stroke="${OL}" stroke-width="3.5"/>`,
            rad('c', [[0, '#ff8a8a'], [0.6, '#e8243c'], [1, '#a3121f']], 0.4, 0.35, 0.75))),
        v('Pig Snout', 30, svg(`<ellipse cx="${NX}" cy="${NY + 2}" rx="20" ry="14" fill="#ffb3c8" stroke="${OL}" stroke-width="3.5"/>` +
            `<ellipse cx="${NX - 7}" cy="${NY + 2}" rx="3.5" ry="5" fill="${OL}"/><ellipse cx="${NX + 7}" cy="${NY + 2}" rx="3.5" ry="5" fill="${OL}"/>`)),
    ],
};

// ---------- Рот ----------

const lipLine = (d, width = 6) => `<path d="${d}" fill="none" stroke="${OL}" stroke-width="${width}" stroke-linecap="round" stroke-linejoin="round"/>`;
const OPEN_MOUTH = `M${MX - 38} ${MY - 8} Q${MX} ${MY - 2} ${MX + 38} ${MY - 8} Q${MX + 30} ${MY + 38} ${MX} ${MY + 40} Q${MX - 30} ${MY + 38} ${MX - 38} ${MY - 8} Z`;

const MOUTH_LAYER = {
    name: 'Mouth',
    none: 0,
    variants: [
        v('Smile', 140, svg(lipLine(`M${MX - 32} ${MY - 6} Q${MX} ${MY + 26} ${MX + 32} ${MY - 6}`))),
        v('Big Grin', 90, svg(`<path d="${OPEN_MOUTH}" fill="#5a1e2c" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            `<clipPath id="m"><path d="${OPEN_MOUTH}"/></clipPath>` +
            `<g clip-path="url(#m)"><rect x="${MX - 40}" y="${MY - 12}" width="80" height="16" fill="#fff"/>` +
            `<path d="M${MX - 12} ${MY - 8} v12 M${MX + 12} ${MY - 8} v12" stroke="#d9d4e3" stroke-width="2"/></g>`)),
        v('Laugh', 90, svg(`<path d="${OPEN_MOUTH}" fill="#5a1e2c" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            `<clipPath id="m"><path d="${OPEN_MOUTH}"/></clipPath>` +
            `<g clip-path="url(#m)"><ellipse cx="${MX}" cy="${MY + 36}" rx="24" ry="16" fill="#ff6f8e"/></g>`)),
        v('Cat Mouth', 80, svg(lipLine(`M${MX - 26} ${MY} Q${MX - 13} ${MY + 16} ${MX} ${MY} Q${MX + 13} ${MY + 16} ${MX + 26} ${MY}`, 5))),
        v('Surprised', 60, svg(`<ellipse cx="${MX}" cy="${MY + 6}" rx="15" ry="19" fill="#5a1e2c" stroke="${OL}" stroke-width="5"/>` +
            `<ellipse cx="${MX}" cy="${MY + 16}" rx="8" ry="5" fill="#ff6f8e"/>`)),
        v('Smirk', 70, svg(lipLine(`M${MX - 26} ${MY + 4} Q${MX + 4} ${MY + 14} ${MX + 30} ${MY - 10}`) +
            lipLine(`M${MX + 26} ${MY - 14} L${MX + 34} ${MY - 6}`, 4))),
        v('Fangs', 40, svg(`<path d="M${MX - 32} ${MY - 6} Q${MX} ${MY + 26} ${MX + 32} ${MY - 6}" fill="#5a1e2c" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            `<path d="M${MX - 22} ${MY - 2} l7 16 l6 -14 Z M${MX + 22} ${MY - 2} l-7 16 l-6 -14 Z" fill="#fff" stroke="${OL}" stroke-width="2.5" stroke-linejoin="round"/>`)),
        v('Tongue Out', 60, svg(`<path d="M${MX - 10} ${MY + 6} Q${MX - 12} ${MY + 36} ${MX + 4} ${MY + 36} Q${MX + 18} ${MY + 34} ${MX + 14} ${MY + 6} Z" fill="#ff6f8e" stroke="${OL}" stroke-width="4"/>` +
            `<path d="M${MX + 2} ${MY + 10} v14" stroke="#d9466b" stroke-width="2.5" stroke-linecap="round"/>` +
            lipLine(`M${MX - 30} ${MY - 2} Q${MX} ${MY + 18} ${MX + 30} ${MY - 2}`))),
        v('Flat', 50, svg(lipLine(`M${MX - 22} ${MY + 6} L${MX + 22} ${MY + 6}`))),
        v('Frown', 30, svg(lipLine(`M${MX - 26} ${MY + 16} Q${MX} ${MY - 6} ${MX + 26} ${MY + 16}`))),
        v('Buck Teeth', 30, svg(lipLine(`M${MX - 30} ${MY - 4} Q${MX} ${MY + 14} ${MX + 30} ${MY - 4}`) +
            `<path d="M${MX - 12} ${MY + 6} h24 v14 a3 3 0 0 1 -3 3 h-18 a3 3 0 0 1 -3 -3 Z" fill="#fff" stroke="${OL}" stroke-width="3"/><path d="M${MX} ${MY + 6} v17" stroke="${OL}" stroke-width="2"/>`)),
        v('Zigzag', 20, svg(lipLine(`M${MX - 30} ${MY + 6} l10 -8 l10 8 l10 -8 l10 8 l10 -8 l10 8`, 5))),
        v('Gold Grill', 6, svg(`<path d="${OPEN_MOUTH}" fill="#5a1e2c" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            `<clipPath id="m"><path d="${OPEN_MOUTH}"/></clipPath>` +
            `<g clip-path="url(#m)"><rect x="${MX - 40}" y="${MY - 12}" width="80" height="18" fill="url(#g)"/>` +
            `<path d="M${MX - 18} ${MY - 8} v12 M${MX} ${MY - 8} v14 M${MX + 18} ${MY - 8} v12" stroke="#a8701a" stroke-width="2"/></g>` +
            `<path d="${sparkle(MX + 20, MY - 2, 7)}" fill="#fff"/>`,
        lin('g', [[0, '#fff1a8'], [0.5, '#f2c14e'], [1, '#c78a22']]))),
    ],
};

// ---------- Предмет у роті ----------

const MOUTH_ITEM = {
    name: 'Mouth Item',
    none: 850,
    variants: [
        v('Bubble Gum', 40, svg(`<circle cx="${MX + 26}" cy="${MY + 14}" r="30" fill="url(#g)" stroke="${OL}" stroke-width="4" opacity=".95"/>` +
            `<ellipse cx="${MX + 16}" cy="${MY + 2}" rx="9" ry="6" fill="#fff" opacity=".7" transform="rotate(-30 ${MX + 16} ${MY + 2})"/>`,
        rad('g', [[0, '#ffd6ea'], [1, '#ff6fb5']], 0.4, 0.35, 0.7))),
        v('Lollipop', 35, svg(`<path d="M${MX + 14} ${MY + 6} L${MX + 70} ${MY + 54}" stroke="#fff" stroke-width="7" stroke-linecap="round"/>` +
            `<path d="M${MX + 14} ${MY + 6} L${MX + 70} ${MY + 54}" stroke="${OL}" stroke-width="2" stroke-linecap="round" opacity=".4"/>` +
            `<circle cx="${MX + 78}" cy="${MY + 62}" r="26" fill="#ff5fa2" stroke="${OL}" stroke-width="4"/>` +
            `<path d="M${MX + 78} ${MY + 62} m-2 0 a3 3 0 1 1 5 0 a8 8 0 1 1 -14 0 a13 13 0 1 1 23 0 a18 18 0 1 1 -32 0" fill="none" stroke="#fff" stroke-width="4" stroke-linecap="round"/>`)),
        v('Rose', 25, svg(`<path d="M${MX + 4} ${MY + 6} C${MX + 40} ${MY - 4} ${MX + 70} ${MY - 20} ${MX + 100} ${MY - 40}" stroke="#3e8e41" stroke-width="6" fill="none" stroke-linecap="round"/>` +
            `<path d="M${MX + 50} ${MY - 6} q14 -16 26 -8 q-14 12 -26 8 Z" fill="#5cb85c" stroke="${OL}" stroke-width="2"/>` +
            `<g transform="translate(${MX + 112} ${MY - 50})"><circle r="20" fill="url(#r)" stroke="${OL}" stroke-width="3.5"/>` +
            '<path d="M-8 -2 C-6 -12 8 -12 8 -2 C8 6 -2 8 -4 2" fill="none" stroke="#8b0f2a" stroke-width="3" stroke-linecap="round"/></g>',
        rad('r', [[0, '#ff6b81'], [1, '#c4122f']], 0.4, 0.35, 0.75))),
        v('Pacifier', 25, svg(`<ellipse cx="${MX}" cy="${MY + 10}" rx="34" ry="16" fill="#7fd3ff" stroke="${OL}" stroke-width="4"/>` +
            `<circle cx="${MX}" cy="${MY + 10}" r="8" fill="#fff" stroke="${OL}" stroke-width="3"/>` +
            `<path d="M${MX - 16} ${MY + 24} A16 16 0 0 0 ${MX + 16} ${MY + 24}" fill="none" stroke="#ffd23f" stroke-width="6"/>`)),
        v('Leaf', 30, svg(`<path d="M${MX + 10} ${MY + 6} L${MX + 40} ${MY - 4}" stroke="#3e8e41" stroke-width="4" stroke-linecap="round"/>` +
            `<path d="M${MX + 38} ${MY - 4} C${MX + 50} ${MY - 30} ${MX + 86} ${MY - 30} ${MX + 96} ${MY - 14} C${MX + 82} ${MY + 6} ${MX + 52} ${MY + 8} ${MX + 38} ${MY - 4} Z" fill="url(#l)" stroke="${OL}" stroke-width="3"/>` +
            `<path d="M${MX + 40} ${MY - 4} C${MX + 60} ${MY - 12} ${MX + 78} ${MY - 14} ${MX + 94} ${MY - 14}" fill="none" stroke="#2f6b2f" stroke-width="2"/>`,
        lin('l', [[0, '#8be07a'], [1, '#3e8e41']]))),
    ],
};

// ---------- Очі (білки) і райдужка ----------

const SCLERA = `${shade('s', '#ffffff', '#ffffff', '#dfe3f2')}`;
const eyes = (shape) => svg(EYES.map(([x, y]) => shape(x, y)).join(''), SCLERA);
const sclera = (d) => `<path d="${d}" fill="url(#s)" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>`;
const circle = (x, y, r) => `M${x - r} ${y} a${r} ${r} 0 1 0 ${2 * r} 0 a${r} ${r} 0 1 0 ${-2 * r} 0 Z`;

const EYE_LAYER = {
    name: 'Eyes',
    none: 0,
    variants: [
        v('Round', 150, eyes((x, y) => sclera(circle(x, y, 33)))),
        v('Big', 110, eyes((x, y) => sclera(circle(x, y, 41)) +
            `<path d="M${x - 30} ${y - 34} l-8 -10 M${x - 18} ${y - 41} l-4 -12" stroke="${OL}" stroke-width="4" stroke-linecap="round"/>`)),
        v('Wide', 100, eyes((x, y) => `<ellipse cx="${x}" cy="${y}" rx="40" ry="29" fill="url(#s)" stroke="${OL}" stroke-width="5"/>`)),
        v('Almond', 90, eyes((x, y) => sclera(`M${x - 42} ${y + 4} Q${x - 6} ${y - 44} ${x + 40} ${y - 4} Q${x + 8} ${y + 40} ${x - 42} ${y + 4} Z`))),
        v('Droopy', 60, eyes((x, y) => sclera(`M${x - 33} ${y - 20} L${x + 33} ${y - 20} A33 33 0 0 1 ${x - 33} ${y - 20} Z`) +
            `<path d="M${x - 38} ${y - 20} L${x + 38} ${y - 20}" stroke="${OL}" stroke-width="8" stroke-linecap="round"/>`)),
        v('Sharp', 50, svg(
            sclera(`M${L[0] - 38} ${L[1] - 26} L${L[0] + 34} ${L[1] - 22} Q${L[0] + 40} ${L[1] + 30} ${L[0]} ${L[1] + 32} Q${L[0] - 40} ${L[1] + 30} ${L[0] - 38} ${L[1] - 26} Z`) +
            sclera(`M${R[0] + 38} ${R[1] - 26} L${R[0] - 34} ${R[1] - 22} Q${R[0] - 40} ${R[1] + 30} ${R[0]} ${R[1] + 32} Q${R[0] + 40} ${R[1] + 30} ${R[0] + 38} ${R[1] - 26} Z`),
            SCLERA)),
    ],
};

const IR = 19;
const SHINE = (x, y) => `<circle cx="${x - 7}" cy="${y - 7}" r="5.5" fill="#fff"/><circle cx="${x + 6}" cy="${y + 6}" r="2.5" fill="#fff" opacity=".85"/>`;
const iris = (inner, outer, { pupil = (x, y) => `<circle cx="${x}" cy="${y}" r="8.5" fill="#120a1c"/>`, extra = () => '', defs = '', right } = {}) => svg(
    EYES.map(([x, y], i) => {
        const id = i && right ? 'ir2' : 'ir';
        return `<circle cx="${x}" cy="${y}" r="${IR}" fill="url(#${id})" stroke="${OL}" stroke-width="2.5"/>` + pupil(x, y) + SHINE(x, y) + extra(x, y);
    }).join(''),
    rad('ir', [[0, inner], [0.75, outer], [1, '#1d1430']]) + (right ? rad('ir2', [[0, right[0]], [0.75, right[1]], [1, '#1d1430']]) : '') + defs);

const IRIS = {
    name: 'Iris',
    none: 0,
    variants: [
        v('Brown', 140, iris('#c48a5a', '#6b3f1f')),
        v('Hazel', 100, iris('#c9b45a', '#6b6b2a')),
        v('Blue', 120, iris('#7fd0ff', '#1f5fbf')),
        v('Green', 100, iris('#a6f08a', '#2f8a3a')),
        v('Violet', 70, iris('#d6a6ff', '#6a2fbf')),
        v('Amber', 60, iris('#ffd36b', '#c26a12')),
        v('Red', 40, iris('#ff8a8a', '#a3121f')),
        v('Ice', 40, iris('#ffffff', '#9fd8ff')),
        v('Heterochromia', 20, iris('#7fd0ff', '#1f5fbf', { right: ['#ffd36b', '#c26a12'] })),
        v('Cat Slit', 25, iris('#fff27a', '#d4a10a', { pupil: (x, y) => `<ellipse cx="${x}" cy="${y}" rx="3.5" ry="15" fill="#120a1c"/>` })),
        v('Galaxy', 12, iris('#8f6bff', '#2a145f', {
            extra: (x, y) => `<path d="${sparkle(x + 6, y - 2, 4)}" fill="#fff"/><circle cx="${x - 3}" cy="${y + 9}" r="1.3" fill="#fff"/><circle cx="${x + 10}" cy="${y + 8}" r="1" fill="#fff"/>`,
        })),
        v('Rainbow', 6, svg(EYES.map(([x, y]) => `<circle cx="${x}" cy="${y}" r="${IR}" fill="#fff" stroke="${OL}" stroke-width="2.5"/>` +
            ['#ff5f6d', '#ffa64d', '#ffe66d', '#6dff9e', '#5fc8ff', '#a77bff'].map((c, i) =>
                `<path d="M${x} ${y} L${(x + IR * Math.cos((i * 60 * Math.PI) / 180)).toFixed(1)} ${(y + IR * Math.sin((i * 60 * Math.PI) / 180)).toFixed(1)} A${IR} ${IR} 0 0 1 ${(x + IR * Math.cos(((i + 1) * 60 * Math.PI) / 180)).toFixed(1)} ${(y + IR * Math.sin(((i + 1) * 60 * Math.PI) / 180)).toFixed(1)} Z" fill="${c}"/>`).join('') +
            `<circle cx="${x}" cy="${y}" r="${IR}" fill="none" stroke="${OL}" stroke-width="2.5"/><circle cx="${x}" cy="${y}" r="8.5" fill="#120a1c"/>` + SHINE(x, y)).join(''))),
        v('Laser', 4, svg(`<g filter="url(#g)">` + EYES.map(([x, y]) => `<circle cx="${x}" cy="${y}" r="${IR}" fill="#ff2a2a" stroke="${OL}" stroke-width="2.5"/><circle cx="${x}" cy="${y}" r="9" fill="#fff"/>` +
            `<path d="M${x} ${y} L${x + (x < 250 ? -260 : 260)} ${y + 120}" stroke="#ff3b3b" stroke-width="10" stroke-linecap="round" opacity=".85"/>` +
            `<path d="M${x} ${y} L${x + (x < 250 ? -260 : 260)} ${y + 120}" stroke="#fff" stroke-width="3" stroke-linecap="round"/>`).join('') + '</g>',
        glow('g', 6, '#ff2a2a'))),
    ],
};

// ---------- Зачіска ----------

/** Точки еліпса голови, збільшеного на grow, для кутів (у градусах, 0 - праворуч, 90 - вгору). */
function headArc(from, to, steps, grow = 0) {
    const points = [];
    for (let i = 0; i <= steps; i++) {
        const a = ((from + ((to - from) * i) / steps) * Math.PI) / 180;
        points.push([HEAD.cx + (HEAD.rx + grow) * Math.cos(a), HEAD.cy - (HEAD.ry + grow) * Math.sin(a)]);
    }
    return points;
}
const pt = ([x, y]) => `${x.toFixed(1)} ${y.toFixed(1)}`;

function bumps(from, to, count, grow, bump) {
    // хвиляста лінія з напівкругів уздовж верху голови.
    const points = headArc(from, to, count, grow);
    let d = `M${pt(points[0])}`;
    for (let i = 1; i < points.length; i++) d += ` A${bump} ${bump} 0 0 1 ${pt(points[i])}`;
    return d;
}

const hair = (body, defs = '') => svg(body, defs);
const hairFill = (id, light, dark) => lin(id, [[0, light], [1, dark]]);

function spikes(from, to, count, grow, height) {
    const base = headArc(from, to, count * 2, grow);
    let d = `M${pt(base[0])}`;
    for (let i = 1; i < base.length; i += 2) {
        const a = ((from + ((to - from) * i) / (count * 2)) * Math.PI) / 180;
        const tip = [HEAD.cx + (HEAD.rx + grow + height) * Math.cos(a), HEAD.cy - (HEAD.ry + grow + height) * Math.sin(a)];
        d += ` L${pt(tip)} L${pt(base[i + 1] || base[i])}`;
    }
    return d;
}

const HAIR = {
    name: 'Hair',
    none: 300,
    variants: [
        v('Tuft', 80, hair(`<path d="M244 122 C230 96 238 72 262 70 C252 82 252 96 268 104 C262 110 254 116 250 124 Z" fill="url(#h)" stroke="${OL}" stroke-width="4.5" stroke-linejoin="round"/>`,
            hairFill('h', '#5a4636', '#2e2219'))),
        v('Mohawk', 50, hair(`<path d="${spikes(122, 58, 5, -4, 46)} Z" fill="url(#h)" stroke="${OL}" stroke-width="4.5" stroke-linejoin="round"/>`,
            hairFill('h', '#ff8fd0', '#d6248f'))),
        v('Curly', 70, hair(`<path d="${bumps(176, 4, 11, 4, 17)} C340 190 160 190 ${pt(headArc(176, 176, 1, 4)[0])} Z" fill="url(#h)" stroke="${OL}" stroke-width="4.5" stroke-linejoin="round"/>` +
            '<g fill="none" stroke="#2e1a10" stroke-width="3" opacity=".35"><path d="M190 140 a10 10 0 1 1 14 8"/><path d="M290 132 a10 10 0 1 1 14 8"/><path d="M240 124 a10 10 0 1 1 14 8"/></g>',
        hairFill('h', '#9a5a32', '#5a3018'))),
        v('Spiky', 60, hair(`<path d="${spikes(168, 12, 7, 2, 34)} C330 200 170 200 ${pt(headArc(168, 168, 1, 2)[0])} Z" fill="url(#h)" stroke="${OL}" stroke-width="4.5" stroke-linejoin="round"/>`,
            hairFill('h', '#fff2a6', '#f2b92e'))),
        v('Bowl Cut', 60, hair(`<path d="M${pt(headArc(190, 190, 1, 6)[0])} ${headArc(190, -10, 16, 8).map((p) => `L${pt(p)}`).join(' ')} ` +
            'L384 205 L360 190 L340 200 L320 186 L300 198 L280 184 L260 196 L240 184 L220 198 L200 186 L180 200 L160 188 L140 202 L116 205 Z" ' +
            `fill="url(#h)" stroke="${OL}" stroke-width="4.5" stroke-linejoin="round"/>` +
            '<path d="M170 150 C200 132 240 128 270 132" stroke="#fff" stroke-width="5" fill="none" opacity=".25" stroke-linecap="round"/>',
        hairFill('h', '#3a3550', '#14121c'))),
        v('Side Swoop', 60, hair(`<path d="M${pt(headArc(200, 200, 1, 4)[0])} ${headArc(200, 30, 14, 6).map((p) => `L${pt(p)}`).join(' ')} ` +
            'C330 150 280 160 240 196 C220 214 190 222 160 214 C140 210 122 220 114 232 Z" ' +
            `fill="url(#h)" stroke="${OL}" stroke-width="4.5" stroke-linejoin="round"/>`,
        hairFill('h', '#ffb066', '#e0611f'))),
        v('Space Buns', 50, hair(mirror(`<circle cx="160" cy="132" r="38" fill="url(#h)" stroke="${OL}" stroke-width="4.5"/>` +
            '<path d="M140 118 a20 20 0 0 1 30 -6" stroke="#fff" stroke-width="4" fill="none" opacity=".3"/>' +
            `<path d="M184 160 l12 -12" stroke="#ff6fa8" stroke-width="7" stroke-linecap="round"/>`),
        hairFill('h', '#c49bff', '#7a3fd6'))),
        v('Leaf Sprout', 50, hair(`<path d="M250 122 C250 100 252 86 256 74" stroke="#3e8e41" stroke-width="6" fill="none" stroke-linecap="round"/>` +
            `<path d="M256 78 C240 50 206 52 196 66 C214 84 240 86 256 78 Z" fill="url(#l)" stroke="${OL}" stroke-width="3.5"/>` +
            `<path d="M256 76 C270 46 304 44 314 58 C298 78 272 82 256 76 Z" fill="url(#l)" stroke="${OL}" stroke-width="3.5"/>`,
        lin('l', [[0, '#a6f08a'], [1, '#3e8e41']]))),
        v('Afro', 30, hair(`<path d="${bumps(196, -16, 15, 22, 22)} C360 230 140 230 ${pt(headArc(196, 196, 1, 22)[0])} Z" fill="url(#h)" stroke="${OL}" stroke-width="4.5" stroke-linejoin="round"/>`,
            hairFill('h', '#4a3426', '#20140c'))),
        v('Flame Hair', 8, hair(`<g filter="url(#g)"><path d="M140 170 C120 120 150 90 170 60 C175 90 185 100 200 104 C196 70 214 40 240 20 C236 60 250 80 266 86 C270 56 290 40 316 30 C304 66 318 90 340 96 C350 80 352 66 350 50 C380 90 384 140 360 170 C320 140 180 140 140 170 Z" fill="url(#f)"/></g>`,
            lin('f', [[0, '#fff3a6'], [0.4, '#ffb21f'], [1, '#ff3d1f']], 0, 0, 0, 1) + glow('g', 7, '#ff7a1a'))),
    ],
};

// ---------- Брови ----------

const brow = (d, width = 7) => `<path d="${d}" fill="none" stroke="${OL}" stroke-width="${width}" stroke-linecap="round"/>`;

const BROWS = {
    name: 'Brows',
    none: 250,
    variants: [
        v('Neutral', 120, svg(brow(`M${L[0] - 24} ${L[1] - 50} Q${L[0]} ${L[1] - 60} ${L[0] + 24} ${L[1] - 52}`) +
            brow(`M${R[0] - 24} ${R[1] - 52} Q${R[0]} ${R[1] - 60} ${R[0] + 24} ${R[1] - 50}`))),
        v('Raised', 80, svg(brow(`M${L[0] - 26} ${L[1] - 58} Q${L[0]} ${L[1] - 76} ${L[0] + 24} ${L[1] - 62}`) +
            brow(`M${R[0] - 24} ${R[1] - 62} Q${R[0]} ${R[1] - 76} ${R[0] + 26} ${R[1] - 58}`))),
        v('Angry', 70, svg(brow(`M${L[0] - 26} ${L[1] - 62} L${L[0] + 26} ${L[1] - 44}`, 8) + brow(`M${R[0] + 26} ${R[1] - 62} L${R[0] - 26} ${R[1] - 44}`, 8))),
        v('Worried', 60, svg(brow(`M${L[0] - 26} ${L[1] - 46} Q${L[0]} ${L[1] - 52} ${L[0] + 24} ${L[1] - 64}`) +
            brow(`M${R[0] - 24} ${R[1] - 64} Q${R[0]} ${R[1] - 52} ${R[0] + 26} ${R[1] - 46}`))),
        v('Thick', 50, svg(brow(`M${L[0] - 22} ${L[1] - 52} Q${L[0]} ${L[1] - 62} ${L[0] + 22} ${L[1] - 54}`, 14) +
            brow(`M${R[0] - 22} ${R[1] - 54} Q${R[0]} ${R[1] - 62} ${R[0] + 22} ${R[1] - 52}`, 14))),
        v('Dots', 40, svg(`<ellipse cx="${L[0]}" cy="${L[1] - 56}" rx="11" ry="7" fill="${OL}"/><ellipse cx="${R[0]}" cy="${R[1] - 56}" rx="11" ry="7" fill="${OL}"/>`)),
        v('Unibrow', 10, svg(brow(`M${L[0] - 26} ${L[1] - 50} Q${L[0]} ${L[1] - 64} 250 ${L[1] - 54} Q${R[0]} ${R[1] - 64} ${R[0] + 26} ${R[1] - 50}`, 11))),
    ],
};

// ---------- Окуляри ----------

const BRIDGE = `<path d="M${L[0] + 40} ${L[1] - 4} Q250 ${L[1] - 18} ${R[0] - 40} ${R[1] - 4}" fill="none" stroke="${OL}" stroke-width="6"/>`;
const TEMPLES = `<path d="M${L[0] - 44} ${L[1] - 8} L118 ${L[1] - 18} M${R[0] + 44} ${R[1] - 8} L382 ${R[1] - 18}" stroke="${OL}" stroke-width="6" stroke-linecap="round"/>`;
const lenses = (shape, defs = '') => svg(TEMPLES + BRIDGE + EYES.map(([x, y]) => shape(x, y)).join(''), defs);

const EYEWEAR = {
    name: 'Eyewear',
    none: 700,
    variants: [
        v('Round Glasses', 70, lenses((x, y) => `<circle cx="${x}" cy="${y}" r="44" fill="#d9f2ff" fill-opacity=".18" stroke="${OL}" stroke-width="6"/>` +
            `<path d="M${x - 24} ${y - 22} l14 -10" stroke="#fff" stroke-width="5" stroke-linecap="round" opacity=".7"/>`)),
        v('Aviators', 60, lenses((x, y) => `<path d="M${x - 46} ${y - 24} Q${x} ${y - 34} ${x + 46} ${y - 24} Q${x + 46} ${y + 34} ${x} ${y + 36} Q${x - 46} ${y + 34} ${x - 46} ${y - 24} Z" fill="url(#a)" stroke="#c78a22" stroke-width="5"/>` +
            `<path d="M${x - 30} ${y - 14} l18 -6" stroke="#fff" stroke-width="5" stroke-linecap="round" opacity=".55"/>`,
        lin('a', [[0, '#2b2140', 0.95], [1, '#6a4fb3', 0.85]]))),
        v('Heart Shades', 40, lenses((x, y) => `<path d="${heart(x, y, 38)}" fill="url(#h)" stroke="${OL}" stroke-width="5"/>` +
            `<path d="M${x - 24} ${y - 22} l12 -6" stroke="#fff" stroke-width="5" stroke-linecap="round" opacity=".6"/>`,
        lin('h', [[0, '#ff6fa8', 0.92], [1, '#c2185b', 0.92]]))),
        v('Cyber Visor', 30, svg(`<g filter="url(#g)"><path d="M118 222 Q250 196 382 222 L378 282 Q250 300 122 282 Z" fill="url(#v)" stroke="${OL}" stroke-width="5"/></g>` +
            '<path d="M140 240 Q250 222 360 240" stroke="#fff" stroke-width="4" fill="none" opacity=".6"/>',
        lin('v', [[0, '#4ff3ff', 0.85], [1, '#8f5bff', 0.85]], 0, 0, 1, 0) + glow('g', 6, '#4ff3ff'))),
        v('Monocle', 25, svg(`<circle cx="${R[0]}" cy="${R[1]}" r="44" fill="#fffbe0" fill-opacity=".2" stroke="#d4972a" stroke-width="7"/>` +
            `<path d="M${R[0] + 30} ${R[1] + 32} C${R[0] + 40} ${R[1] + 90} ${R[0] + 10} ${R[1] + 120} ${R[0] + 30} ${R[1] + 170}" stroke="#d4972a" stroke-width="3" fill="none" stroke-dasharray="6 4"/>`)),
        v('3D Glasses', 30, svg(`<path d="M114 212 L386 212 L382 290 L264 290 L250 268 L236 290 L118 290 Z" fill="#fff" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            `<rect x="136" y="224" width="96" height="54" rx="8" fill="#ff3b55" fill-opacity=".85"/><rect x="268" y="224" width="96" height="54" rx="8" fill="#33d6ff" fill-opacity=".85"/>`)),
        v('Star Shades', 20, lenses((x, y) => `<path d="${star(x, y + 2, 48, 22)}" fill="url(#s)" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>`,
            lin('s', [[0, '#ffe066', 0.92], [1, '#ff9f1c', 0.92]]))),
        v('Nerd Glasses', 50, lenses((x, y) => `<rect x="${x - 44}" y="${y - 32}" width="88" height="64" rx="12" fill="#d9f2ff" fill-opacity=".18" stroke="${OL}" stroke-width="10"/>`)),
    ],
};

module.exports = { CHEEK, FACE_MARK, NOSE_LAYER, MOUTH_LAYER, MOUTH_ITEM, EYE_LAYER, IRIS, HAIR, BROWS, EYEWEAR };
