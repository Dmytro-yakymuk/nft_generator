// Шари сцени позаду персонажа: фон, візерунок, небесне тіло, аура, крила.

const { OL, HEAD, rng, svg, lin, rad, blur, glow, star, sparkle, mirror, v } = require('./lib');

const LIGHT = rad('light', [[0, '#fff', 0.3], [1, '#fff', 0]], 0.5, 0.42, 0.62);

const sky = (stops, extra = '', defs = '') => svg(
    `<rect width="500" height="500" fill="url(#sky)"/>${extra}<rect width="500" height="500" fill="url(#light)"/>`,
    lin('sky', stops) + LIGHT + defs);

function starField(seed, count, maxR, opacity = 0.9) {
    const r = rng(seed);
    let out = '';
    for (let i = 0; i < count; i++) {
        out += `<circle cx="${(r() * 500).toFixed(1)}" cy="${(r() * 500).toFixed(1)}" r="${(0.5 + r() * maxR).toFixed(2)}" fill="#fff" opacity="${(0.3 + r() * (opacity - 0.3)).toFixed(2)}"/>`;
    }
    return out;
}

const hills = (color, opacity) =>
    `<path d="M0 430 C80 390 150 400 230 425 C320 452 410 380 500 400 L500 500 L0 500 Z" fill="${color}" opacity="${opacity}"/>` +
    `<path d="M0 465 C110 430 190 470 290 460 C380 450 440 430 500 445 L500 500 L0 500 Z" fill="${color}" opacity="${opacity * 1.4}"/>`;

const BACKGROUND = {
    name: 'Background',
    none: 0,
    variants: [
        v('Sunset', 140, sky([[0, '#ff5f8f'], [0.5, '#ff9a8b'], [1, '#ffd29d']], hills('#7a2e6b', 0.25))),
        v('Aurora', 120, sky([[0, '#071a33'], [0.55, '#0f3b4f'], [1, '#2d6a5a']],
            starField(7, 50, 1.2, 0.8) +
            '<g filter="url(#b)" fill="none" stroke-linecap="round">' +
            '<path d="M-40 170 C90 90 190 230 300 130 S470 90 560 150" stroke="#4dffb5" stroke-width="46" opacity=".5"/>' +
            '<path d="M-40 230 C120 150 230 280 340 190 S480 170 560 210" stroke="#b56bff" stroke-width="34" opacity=".4"/>' +
            '<path d="M-40 110 C60 70 160 160 260 90 S430 40 560 90" stroke="#5ad1ff" stroke-width="22" opacity=".45"/></g>',
            blur('b', 14))),
        v('Deep Space', 110, sky([[0, '#0b0820'], [0.6, '#22184a'], [1, '#3a1f5c']],
            '<g filter="url(#b)"><ellipse cx="120" cy="140" rx="150" ry="90" fill="#7b3fe4" opacity=".35"/>' +
            '<ellipse cx="410" cy="330" rx="140" ry="100" fill="#2f7be4" opacity=".3"/></g>' + starField(11, 90, 1.6),
            blur('b', 40))),
        v('Mint Dream', 130, sky([[0, '#e3ffd6'], [1, '#8fe3b5']], hills('#2f9c78', 0.18))),
        v('Peach Sky', 130, sky([[0, '#fff1df'], [1, '#fbb39a']])),
        v('Ocean', 120, sky([[0, '#1d8fbf'], [0.6, '#5cc8e8'], [1, '#bdf1ff']],
            '<g fill="none" stroke="#fff" stroke-width="3" opacity=".25">' +
            '<path d="M0 420 Q60 405 120 420 T240 420 T360 420 T480 420 T600 420"/>' +
            '<path d="M0 455 Q60 440 120 455 T240 455 T360 455 T480 455 T600 455"/></g>')),
        v('Lavender Haze', 120, sky([[0, '#e6ccff'], [1, '#93c5fd']])),
        v('Golden Hour', 90, sky([[0, '#ffd86f'], [0.6, '#fda769'], [1, '#f6735f']], hills('#8f3b2a', 0.2))),
        v('Cherry Blossom', 90, sky([[0, '#ffd1ec'], [1, '#b6c8f5']])),
        v('Midnight', 60, sky([[0, '#0f1a2e'], [1, '#2b4a74']], starField(23, 60, 1.3, 0.7))),
        v('Bubblegum', 60, sky([[0, '#ff8fb3'], [1, '#ffd6ea']])),
        v('Holographic', 8, svg(
            '<rect width="500" height="500" fill="url(#holo)"/>' +
            '<rect x="-200" y="160" width="900" height="70" fill="#fff" opacity=".35" transform="rotate(-30 250 250)" filter="url(#b)"/>' +
            '<rect x="-200" y="280" width="900" height="26" fill="#fff" opacity=".3" transform="rotate(-30 250 250)" filter="url(#b)"/>' +
            '<rect width="500" height="500" fill="url(#light)"/>',
            lin('holo', [[0, '#ffc8f0'], [0.25, '#c2e4ff'], [0.5, '#b4ffd9'], [0.75, '#fff6c2'], [1, '#e4c6ff']], 0, 0, 1, 1) +
            LIGHT + blur('b', 10))),
    ],
};

const pattern = (body, defs = '') => svg(body, defs);

function sunburst() {
    let rays = '';
    for (let i = 0; i < 24; i += 2) {
        const a1 = (i * 15 * Math.PI) / 180;
        const a2 = ((i + 1) * 15 * Math.PI) / 180;
        rays += `<path d="M250 260 L${250 + 600 * Math.cos(a1)} ${260 + 600 * Math.sin(a1)} L${250 + 600 * Math.cos(a2)} ${260 + 600 * Math.sin(a2)} Z"/>`;
    }
    return `<g fill="#fff" opacity=".16">${rays}</g>`;
}

function hexagons() {
    let out = '';
    const r = 34;
    const w = Math.sqrt(3) * r;
    for (let row = -1; row < 11; row++) {
        for (let col = -1; col < 10; col++) {
            const cx = col * w + (row % 2 ? w / 2 : 0);
            const cy = row * r * 1.5;
            const points = [];
            for (let k = 0; k < 6; k++) {
                const a = ((60 * k - 30) * Math.PI) / 180;
                points.push(`${(cx + r * Math.cos(a)).toFixed(1)},${(cy + r * Math.sin(a)).toFixed(1)}`);
            }
            out += `<polygon points="${points.join(' ')}"/>`;
        }
    }
    return `<g fill="none" stroke="#fff" stroke-width="2.5" opacity=".22">${out}</g>`;
}

function cloud(x, y, s) {
    return `<g transform="translate(${x} ${y}) scale(${s})" filter="url(#cs)">` +
        '<path d="M-60 20 C-80 20 -85 -5 -65 -10 C-70 -35 -35 -45 -20 -28 C-12 -55 30 -55 38 -25 C60 -35 82 -15 70 5 C88 10 82 22 66 20 Z" fill="#fff" opacity=".92"/></g>';
}

function confettiStars(seed) {
    const r = rng(seed);
    let out = '';
    for (let i = 0; i < 46; i++) {
        out += `<path d="${sparkle(r() * 500, r() * 500, 3 + r() * 9)}" fill="#fff" opacity="${(0.45 + r() * 0.5).toFixed(2)}"/>`;
    }
    return out;
}

function bokeh(seed) {
    const r = rng(seed);
    const colors = ['#fff', '#ffe6a8', '#ffc2e2', '#c2f0ff'];
    let out = '';
    for (let i = 0; i < 26; i++) {
        out += `<circle cx="${(r() * 500).toFixed(0)}" cy="${(r() * 500).toFixed(0)}" r="${(10 + r() * 34).toFixed(0)}" fill="${colors[i % 4]}" opacity="${(0.12 + r() * 0.25).toFixed(2)}"/>`;
    }
    return `<g filter="url(#b)">${out}</g>`;
}

function synthwave() {
    let lines = '';
    for (let i = 0; i <= 16; i++) {
        const x = -400 + i * 81;
        lines += `<line x1="250" y1="300" x2="${x}" y2="520"/>`;
    }
    for (const y of [312, 328, 350, 380, 420, 470]) lines += `<line x1="0" y1="${y}" x2="500" y2="${y}"/>`;
    return '<rect y="300" width="500" height="200" fill="#2a0845" opacity=".35"/>' +
        `<g stroke="#ff4fd8" stroke-width="2" opacity=".75" filter="url(#g)">${lines}</g>` +
        '<rect y="294" width="500" height="8" fill="#ffd1f7" opacity=".7" filter="url(#b)"/>';
}

const BACKDROP = {
    name: 'Backdrop',
    none: 300,
    variants: [
        v('Stars', 90, pattern(confettiStars(31))),
        v('Polka Dots', 90, pattern(
            '<pattern id="p" width="50" height="50" patternUnits="userSpaceOnUse"><circle cx="12" cy="12" r="7" fill="#fff"/><circle cx="37" cy="37" r="7" fill="#fff"/></pattern>' +
            '<rect width="500" height="500" fill="url(#p)" opacity=".22"/>')),
        v('Sunburst', 80, pattern(sunburst())),
        v('Synthwave Grid', 40, pattern(synthwave(), glow('g', 3, '#ff4fd8') + blur('b', 4))),
        v('Waves', 80, pattern('<g fill="none" stroke="#fff" stroke-width="10" opacity=".16">' +
            [60, 130, 200, 270, 340, 410, 480].map((y) => `<path d="M-20 ${y} Q40 ${y - 28} 100 ${y} T220 ${y} T340 ${y} T460 ${y} T580 ${y}"/>`).join('') + '</g>')),
        v('Hexagons', 60, pattern(hexagons())),
        v('Clouds', 90, pattern(cloud(90, 110, 1.1) + cloud(420, 70, 0.8) + cloud(60, 420, 1.3) + cloud(450, 430, 1.2),
            '<filter id="cs"><feDropShadow dx="0" dy="5" stdDeviation="4" flood-color="#5b4b8a" flood-opacity=".2"/></filter>')),
        v('Checkerboard', 50, pattern(
            '<pattern id="p" width="62.5" height="62.5" patternUnits="userSpaceOnUse"><rect width="31.25" height="31.25" fill="#fff"/><rect x="31.25" y="31.25" width="31.25" height="31.25" fill="#fff"/></pattern>' +
            '<rect width="500" height="500" fill="url(#p)" opacity=".13"/>')),
        v('Diagonal Stripes', 70, pattern(
            '<pattern id="p" width="40" height="40" patternUnits="userSpaceOnUse" patternTransform="rotate(35)"><rect width="18" height="40" fill="#fff"/></pattern>' +
            '<rect width="500" height="500" fill="url(#p)" opacity=".14"/>')),
        v('Bokeh', 60, pattern(bokeh(41), blur('b', 3))),
    ],
};

const at = (x, y, body) => `<g transform="translate(${x} ${y})">${body}</g>`;

function craters() {
    return '<circle cx="-14" cy="-10" r="9" fill="#d9cfa6" opacity=".7"/><circle cx="16" cy="12" r="12" fill="#d9cfa6" opacity=".6"/>' +
        '<circle cx="12" cy="-20" r="5" fill="#d9cfa6" opacity=".7"/><circle cx="-18" cy="20" r="6" fill="#d9cfa6" opacity=".6"/>';
}

const CELESTIAL = {
    name: 'Celestial',
    none: 500,
    variants: [
        v('Full Moon', 90, svg(at(95, 92,
            '<circle r="70" fill="#fff8d6" opacity=".35" filter="url(#b)"/>' +
            `<circle r="46" fill="url(#m)" stroke="${OL}" stroke-width="4"/>` + craters()),
        rad('m', [[0, '#fffbe8'], [1, '#efe2a8']], 0.4, 0.35, 0.7) + blur('b', 14))),
        v('Crescent Moon', 110, svg(at(95, 92,
            '<circle r="64" fill="#fff6c9" opacity=".3" filter="url(#b)"/>' +
            `<path d="M10 -46 A46 46 0 1 0 40 30 A38 38 0 1 1 10 -46 Z" fill="url(#m)" stroke="${OL}" stroke-width="4" stroke-linejoin="round"/>` +
            `<path d="${sparkle(48, -40, 10)}" fill="#fff"/><path d="${sparkle(-50, 45, 7)}" fill="#fff"/>`),
        lin('m', [[0, '#fff6b0'], [1, '#ffc86b']]) + blur('b', 12))),
        v('Ringed Planet', 70, svg(at(98, 95,
            `<ellipse rx="78" ry="20" fill="none" stroke="url(#r)" stroke-width="9" transform="rotate(-18)"/>` +
            `<circle r="42" fill="url(#p)" stroke="${OL}" stroke-width="4"/>` +
            '<path d="M-38 -12 Q0 -24 40 -6" stroke="#fff" stroke-width="5" fill="none" opacity=".35"/>' +
            '<path d="M-40 10 Q0 0 41 16" stroke="#7a3a12" stroke-width="5" fill="none" opacity=".3"/>' +
            '<clipPath id="front"><rect x="-90" y="0" width="180" height="60" transform="rotate(-18)"/></clipPath>' +
            `<ellipse rx="78" ry="20" fill="none" stroke="${OL}" stroke-width="15" transform="rotate(-18)" clip-path="url(#front)"/>` +
            `<ellipse rx="78" ry="20" fill="none" stroke="url(#r)" stroke-width="9" transform="rotate(-18)" clip-path="url(#front)"/>`),
        rad('p', [[0, '#ffd59e'], [0.6, '#f59e5b'], [1, '#c2603a']], 0.4, 0.35, 0.75) +
            lin('r', [[0, '#fff1c9'], [0.5, '#e8b4ff'], [1, '#9ad7ff']], 0, 0, 1, 0))),
        v('Sun', 80, svg(at(95, 92,
            '<circle r="80" fill="#ffd54f" opacity=".45" filter="url(#b)"/>' +
            `<g stroke="#ffb300" stroke-width="7" stroke-linecap="round">${Array.from({ length: 12 }, (_, i) => {
                const a = (i * 30 * Math.PI) / 180;
                return `<line x1="${(56 * Math.cos(a)).toFixed(1)}" y1="${(56 * Math.sin(a)).toFixed(1)}" x2="${(72 * Math.cos(a)).toFixed(1)}" y2="${(72 * Math.sin(a)).toFixed(1)}"/>`;
            }).join('')}</g>` +
            `<circle r="44" fill="url(#s)" stroke="${OL}" stroke-width="4"/>`),
        rad('s', [[0, '#fff7c2'], [0.6, '#ffd23f'], [1, '#ff9f1c']], 0.4, 0.35, 0.75) + blur('b', 16))),
        v('Comet', 60, svg(
            '<path d="M150 70 L10 0 L0 30 Z" fill="url(#t)" opacity=".9" filter="url(#b)"/>' +
            '<path d="M150 70 L30 -10 L-10 60 Z" fill="url(#t)" opacity=".5" filter="url(#b)"/>' +
            `<circle cx="150" cy="72" r="18" fill="#fff" stroke="${OL}" stroke-width="4"/>` +
            '<circle cx="150" cy="72" r="30" fill="#bff3ff" opacity=".45" filter="url(#b)"/>',
        lin('t', [[0, '#7fe3ff', 0], [1, '#e9fbff', 1]], 0, 0, 1, 1) + blur('b', 6))),
        v('Twin Moons', 40, svg(
            at(70, 80, `<circle r="50" fill="#e9e4ff" opacity=".3" filter="url(#b)"/><circle r="30" fill="url(#a)" stroke="${OL}" stroke-width="4"/><circle cx="-8" cy="-6" r="6" fill="#b9b0e8"/><circle cx="10" cy="9" r="8" fill="#b9b0e8"/>`) +
            at(140, 130, `<circle r="20" fill="url(#c)" stroke="${OL}" stroke-width="4"/><circle cx="5" cy="-4" r="5" fill="#f2a7b8"/>`),
        rad('a', [[0, '#ffffff'], [1, '#cfc6ff']], 0.4, 0.35, 0.7) + rad('c', [[0, '#fff0f4'], [1, '#ffb8c9']], 0.4, 0.35, 0.7) + blur('b', 10))),
        v('Earth', 40, svg(at(95, 92,
            '<circle r="62" fill="#9fd8ff" opacity=".35" filter="url(#b)"/>' +
            '<clipPath id="e"><circle r="44"/></clipPath>' +
            '<circle r="44" fill="url(#o)"/>' +
            '<g clip-path="url(#e)" fill="#62c46b"><path d="M-40 -20 C-25 -35 -5 -30 0 -15 C5 0 -15 5 -20 18 C-28 30 -45 10 -40 -20 Z"/>' +
            '<path d="M10 -40 C25 -38 40 -20 30 -5 C22 6 12 -8 8 -20 Z"/><path d="M8 20 C20 12 38 18 34 34 C25 46 6 40 8 20 Z"/></g>' +
            `<circle r="44" fill="url(#hl)" stroke="${OL}" stroke-width="4"/>`),
        rad('o', [[0, '#5fb8ff'], [1, '#1f5fbf']], 0.4, 0.35, 0.75) + rad('hl', [[0, '#fff', 0.45], [0.5, '#fff', 0], [1, '#000', 0.25]], 0.35, 0.3, 0.8) + blur('b', 12))),
        v('Black Hole', 10, svg(at(100, 100,
            '<ellipse rx="92" ry="30" fill="none" stroke="#ff9d3c" stroke-width="16" opacity=".8" transform="rotate(-15)" filter="url(#b)"/>' +
            '<circle r="60" fill="#ff7a1a" opacity=".35" filter="url(#bb)"/>' +
            '<circle r="38" fill="#05030a"/>' +
            '<circle r="40" fill="none" stroke="#ffd38a" stroke-width="3" opacity=".9"/>' +
            '<clipPath id="front"><rect x="-120" y="0" width="240" height="80" transform="rotate(-15)"/></clipPath>' +
            '<ellipse rx="92" ry="30" fill="none" stroke="#ffcf7a" stroke-width="9" transform="rotate(-15)" clip-path="url(#front)"/>'),
        blur('b', 5) + blur('bb', 18))),
    ],
};

const aroundHead = (fill, opacity = 1) =>
    `<ellipse cx="${HEAD.cx}" cy="${HEAD.cy + 40}" rx="210" ry="220" fill="${fill}" opacity="${opacity}"/>`;

function flames(seed, colors) {
    const r = rng(seed);
    let out = '';
    for (let i = 0; i < 30; i++) {
        const a = Math.PI + (i / 29) * Math.PI;
        const x = 250 + Math.cos(a) * (150 + r() * 10);
        const y = 270 + Math.sin(a) * (150 + r() * 10);
        const h = 50 + r() * 60;
        out += `<path d="M${x - 18} ${y + 10} Q${x - 10} ${y - h / 2} ${x + (r() - 0.5) * 20} ${y - h} Q${x + 12} ${y - h / 2} ${x + 18} ${y + 10} Z" fill="${colors[i % colors.length]}" transform="rotate(${(a * 180) / Math.PI + 90} ${x} ${y})"/>`;
    }
    return out;
}

function lightning(seed) {
    const r = rng(seed);
    let out = '';
    for (let i = 0; i < 9; i++) {
        const a = Math.PI * 0.95 + (i / 8) * Math.PI * 1.1;
        let x = 250 + Math.cos(a) * 165;
        let y = 270 + Math.sin(a) * 160;
        let d = `M${x.toFixed(0)} ${y.toFixed(0)}`;
        for (let k = 0; k < 4; k++) {
            x += Math.cos(a) * 18 + (r() - 0.5) * 26;
            y += Math.sin(a) * 18 + (r() - 0.5) * 26;
            d += ` L${x.toFixed(0)} ${y.toFixed(0)}`;
        }
        out += `<path d="${d}"/>`;
    }
    return out;
}

const AURA = {
    name: 'Aura',
    none: 600,
    variants: [
        v('Pink Glow', 80, svg(aroundHead('url(#a)'), rad('a', [[0, '#ff7ad9', 0.8], [0.6, '#ff7ad9', 0.35], [1, '#ff7ad9', 0]]))),
        v('Cyan Glow', 80, svg(aroundHead('url(#a)'), rad('a', [[0, '#5ff4ff', 0.8], [0.6, '#5ff4ff', 0.35], [1, '#5ff4ff', 0]]))),
        v('Gold Glow', 60, svg(aroundHead('url(#a)'), rad('a', [[0, '#ffe066', 0.9], [0.6, '#ffc233', 0.4], [1, '#ffc233', 0]]))),
        v('Rainbow Ring', 40, svg(
            `<g fill="none" stroke-width="10" filter="url(#b)" opacity=".85">${['#ff5f6d', '#ffa64d', '#ffe66d', '#6dff9e', '#5fc8ff', '#a77bff']
                .map((c, i) => `<circle cx="250" cy="245" r="${200 - i * 9}" stroke="${c}"/>`).join('')}</g>`,
            blur('b', 2))),
        v('Fire Aura', 50, svg(`<g filter="url(#b)" opacity=".9">${flames(5, ['#ff5e1a', '#ff9a1f', '#ffd23f'])}</g>`, blur('b', 6))),
        v('Electric Aura', 40, svg(
            `<g fill="none" stroke="#9ff6ff" stroke-width="5" stroke-linecap="round" stroke-linejoin="round" filter="url(#g)">${lightning(9)}</g>` +
            aroundHead('url(#a)'),
            glow('g', 6, '#3fd7ff') + rad('a', [[0, '#3fd7ff', 0.35], [0.7, '#3fd7ff', 0.1], [1, '#3fd7ff', 0]]))),
        v('Shadow Aura', 50, svg(`<g filter="url(#b)">${aroundHead('#2a0f45', 0.55)}${flames(13, ['#3b1466', '#5a1f8f'])}</g>`, blur('b', 12))),
    ],
};

// Крила: ліве крило від плеча вгору-ліворуч, праве - віддзеркалення.
const wing = (body, defs) => svg(mirror(`<g transform="translate(165 405) scale(1.3) translate(-165 -405)">${body}</g>`), defs);

/** Крило з пір'я: три ряди пір'їн віялом від плеча (довгі махові позаду, короткі покривні попереду). */
function feathers(stroke = OL) {
    const [sx, sy] = [165, 400];
    const rows = [
        { count: 7, from: -170, to: -100, length: 150, width: 26 },
        { count: 6, from: -163, to: -108, length: 110, width: 27 },
        { count: 5, from: -155, to: -115, length: 70, width: 25 },
    ];
    let out = '';
    for (const [k, row] of rows.entries()) {
        for (let i = 0; i < row.count; i++) {
            const a = ((row.from + ((row.to - row.from) * i) / (row.count - 1)) * Math.PI) / 180;
            const [dx, dy] = [Math.cos(a), Math.sin(a)];
            const [bx, by] = [sx + dx * 12, sy + dy * 12];
            const [tx, ty] = [sx + dx * row.length, sy + dy * row.length];
            const [mx, my] = [(bx + tx) / 2, (by + ty) / 2];
            const [px, py] = [-dy * row.width, dx * row.width];
            out += `<path d="M${bx.toFixed(1)} ${by.toFixed(1)} Q${(mx + px).toFixed(1)} ${(my + py).toFixed(1)} ${tx.toFixed(1)} ${ty.toFixed(1)} ` +
                `Q${(mx - px).toFixed(1)} ${(my - py).toFixed(1)} ${bx.toFixed(1)} ${by.toFixed(1)} Z" fill="url(#w${k})" stroke="${stroke}" stroke-width="3.5" stroke-linejoin="round"/>`;
        }
    }
    return out;
}

const feathersDefs = (colors) => colors.map(([light, dark], k) => lin(`w${k}`, [[0, light], [1, dark]], 1, 1, 0, 0)).join('');

const ANGEL = feathers();

/** Перетинчасте крило: плече, зап'ястя й три кінчики пальців, між ними - увігнуті краї перетинки. */
const membrane = (bones, extra = '') =>
    '<path d="M165 400 L112 300 L34 196 Q72 252 22 296 Q74 318 56 378 Q112 368 165 400 Z" fill="url(#w)" stroke="' + OL + '" stroke-width="5" stroke-linejoin="round"/>' +
    `<path d="M165 400 L112 300 L34 196 M112 300 L22 296 M112 300 L56 378" stroke="${bones}" stroke-width="5" fill="none" stroke-linecap="round" stroke-linejoin="round"/>` +
    `<path d="M112 300 l-8 -16 l16 6 Z" fill="#f4e6c8" stroke="${OL}" stroke-width="2.5"/>` + extra;

const BAT = membrane(OL);

const BUTTERFLY = `<path d="M170 380 C120 300 40 240 40 300 C40 350 110 370 170 380 Z" fill="url(#w)" stroke="${OL}" stroke-width="5"/>` +
    `<path d="M170 395 C120 400 70 450 100 470 C130 488 160 440 170 395 Z" fill="url(#w2)" stroke="${OL}" stroke-width="5"/>` +
    `<circle cx="88" cy="300" r="14" fill="#fff" opacity=".85" stroke="${OL}" stroke-width="3"/><circle cx="125" cy="335" r="8" fill="#2b2140" opacity=".7"/>` +
    '<circle cx="118" cy="440" r="9" fill="#fff" opacity=".8"/>';

const FAIRY = '<g filter="url(#g)">' +
    `<ellipse cx="100" cy="290" rx="40" ry="100" transform="rotate(-35 100 290)" fill="url(#w)" stroke="#e8fbff" stroke-width="3"/>` +
    `<ellipse cx="115" cy="420" rx="26" ry="60" transform="rotate(-60 115 420)" fill="url(#w)" stroke="#e8fbff" stroke-width="3"/></g>` +
    '<path d="M160 380 C130 330 100 270 70 215 M160 395 C120 410 95 425 75 440" stroke="#fff" stroke-width="2" fill="none" opacity=".8"/>';

const DRAGON = membrane('#1f5a36', '<g fill="#fff" opacity=".18"><circle cx="90" cy="300" r="10"/><circle cx="70" cy="330" r="7"/><circle cx="100" cy="345" r="6"/></g>');

const MECHA = `<g stroke="${OL}" stroke-width="4" stroke-linejoin="round">` +
    '<path d="M162 392 L60 210 L80 205 L170 360 Z" fill="url(#w)"/>' +
    '<path d="M160 398 L38 280 L60 270 L166 375 Z" fill="url(#w)"/>' +
    '<path d="M158 405 L40 360 L52 345 L162 388 Z" fill="url(#w)"/></g>' +
    '<g stroke="#4ff3ff" stroke-width="3" filter="url(#g)"><path d="M150 360 L78 230"/><path d="M150 380 L60 292"/><path d="M150 396 L58 360"/></g>';

const PHOENIX = `<g filter="url(#g)">${feathers('#b3200f')}</g>`;

const WINGS = {
    name: 'Wings',
    none: 750,
    variants: [
        v('Angel', 60, wing(ANGEL, feathersDefs([['#ffffff', '#c9d3f5'], ['#ffffff', '#dde4fb'], ['#ffffff', '#eef2ff']]))),
        v('Bat', 50, wing(BAT, lin('w', [[0, '#6a3d8f'], [1, '#33184d']]))),
        v('Butterfly', 50, wing(BUTTERFLY, lin('w', [[0, '#ff9ad5'], [1, '#8f7bff']], 0, 0, 1, 1) + lin('w2', [[0, '#ffd27a'], [1, '#ff7aa8']]))),
        v('Fairy', 50, wing(FAIRY, lin('w', [[0, '#c9f6ff', 0.75], [1, '#e7d4ff', 0.55]]) + glow('g', 5, '#9fefff'))),
        v('Dragon', 25, wing(DRAGON, lin('w', [[0, '#7fd36b'], [1, '#2f7a4a']]))),
        v('Mecha', 15, wing(MECHA, lin('w', [[0, '#d9e1ea'], [1, '#7d8a99']]) + glow('g', 3, '#4ff3ff'))),
        v('Phoenix', 5, wing(PHOENIX, feathersDefs([['#fff3a6', '#ff3d1f'], ['#fff7c2', '#ff7a1a'], ['#fffbe0', '#ffb21f']]) + glow('g', 6, '#ff7a1a'))),
    ],
};

module.exports = { BACKGROUND, BACKDROP, CELESTIAL, AURA, WINGS };
