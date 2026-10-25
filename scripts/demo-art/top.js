// Верхні шари: головний убір, улюбленець, частинки, освітлення, рамка.

const { OL, HEAD, rng, svg, lin, rad, blur, glow, star, heart, sparkle, mirror, v } = require('./lib');

// ---------- Головний убір ----------

const metal = (id, light, base, dark) => lin(id, [[0, light], [0.5, base], [1, dark]]);
const knit = (x1, x2, y1, y2, step = 10) => Array.from({ length: Math.floor((x2 - x1) / step) + 1 },
    (_, i) => `<line x1="${x1 + i * step}" y1="${y1}" x2="${x1 + i * step}" y2="${y2}"/>`).join('');

const HEADGEAR = {
    name: 'Headgear',
    none: 500,
    variants: [
        v('Crown', 50, svg(`<path d="M180 146 L174 84 L208 112 L230 66 L250 100 L270 66 L292 112 L326 84 L320 146 Q250 130 180 146 Z" fill="url(#g)" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            '<path d="M182 132 Q250 118 318 132" stroke="#b97d1f" stroke-width="4" fill="none"/>' +
            `<circle cx="250" cy="128" r="9" fill="#e8243c" stroke="${OL}" stroke-width="3"/><circle cx="212" cy="132" r="6" fill="#2f6fd6" stroke="${OL}" stroke-width="2.5"/><circle cx="288" cy="132" r="6" fill="#2fbf71" stroke="${OL}" stroke-width="2.5"/>` +
            `<circle cx="174" cy="84" r="6" fill="#fff6c4" stroke="${OL}" stroke-width="2.5"/><circle cx="230" cy="66" r="6" fill="#fff6c4" stroke="${OL}" stroke-width="2.5"/><circle cx="270" cy="66" r="6" fill="#fff6c4" stroke="${OL}" stroke-width="2.5"/><circle cx="326" cy="84" r="6" fill="#fff6c4" stroke="${OL}" stroke-width="2.5"/>`,
        metal('g', '#fff3b0', '#f2c14e', '#c78a22'))),
        v('Beanie', 70, svg(`<path d="M116 196 C116 70 384 70 384 196 Z" fill="url(#k)" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            `<g stroke="#1f6f6a" stroke-width="2" opacity=".35">${knit(140, 360, 110, 190, 14)}</g>` +
            `<path d="M110 188 Q250 158 390 188 L390 218 Q250 190 110 218 Z" fill="#2a9d8f" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            `<g stroke="#1f6f6a" stroke-width="3" opacity=".5">${knit(124, 376, 178, 210, 12)}</g>` +
            `<circle cx="250" cy="92" r="24" fill="url(#p)" stroke="${OL}" stroke-width="5"/>`,
        lin('k', [[0, '#5fd4c4'], [1, '#2a9d8f']]) + rad('p', [[0, '#fff'], [1, '#ffd1dc']], 0.4, 0.35, 0.7))),
        v('Baseball Cap', 70, svg(`<path d="M300 186 C340 176 400 180 444 200 C420 214 360 214 300 206 Z" fill="#c41f2f" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            `<path d="M120 196 C120 76 380 76 380 196 Q250 176 120 196 Z" fill="url(#c)" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            `<path d="M250 100 L250 186 M190 108 C175 140 170 170 172 190 M310 108 C325 140 330 170 328 190" stroke="${OL}" stroke-width="3" fill="none" opacity=".35"/>` +
            `<circle cx="250" cy="98" r="8" fill="#c41f2f" stroke="${OL}" stroke-width="3"/>` +
            `<path d="${star(214, 150, 18, 8)}" fill="#fff" stroke="${OL}" stroke-width="3" stroke-linejoin="round"/>`,
        lin('c', [[0, '#ff5a5a'], [1, '#c41f2f']]))),
        v('Wizard Hat', 30, svg(`<path d="M160 168 Q220 100 262 16 Q290 30 296 46 Q300 110 340 168 Z" fill="url(#w)" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            `<ellipse cx="250" cy="170" rx="132" ry="24" fill="#4b2a8f" stroke="${OL}" stroke-width="5"/>` +
            `<g fill="#ffe066">${[[240, 120, 10], [282, 80, 8], [205, 150, 7], [300, 140, 6], [268, 44, 5]].map(([x, y, r]) => `<path d="${star(x, y, r, r / 2.3)}"/>`).join('')}</g>` +
            `<path d="M196 146 Q250 132 304 146" stroke="#ffe066" stroke-width="8" fill="none"/>`,
        lin('w', [[0, '#8f6bff'], [1, '#4b2a8f']]))),
        v('Top Hat', 40, svg(`<ellipse cx="250" cy="148" rx="118" ry="18" fill="#1d1b26" stroke="${OL}" stroke-width="5"/>` +
            `<path d="M178 148 L184 40 Q250 28 316 40 L322 148 Q250 162 178 148 Z" fill="url(#t)" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            `<path d="M181 118 Q250 128 319 118 L320 140 Q250 152 180 140 Z" fill="#c41f2f" stroke="${OL}" stroke-width="3"/>` +
            '<path d="M200 50 L198 110" stroke="#fff" stroke-width="5" opacity=".18" stroke-linecap="round"/>',
        lin('t', [[0, '#3b3948'], [1, '#14121c']], 0, 0, 1, 0))),
        v('Halo', 30, svg(`<g filter="url(#g)"><ellipse cx="250" cy="82" rx="82" ry="18" fill="none" stroke="#ffe066" stroke-width="11"/></g>` +
            '<ellipse cx="250" cy="82" rx="82" ry="18" fill="none" stroke="#fffbe0" stroke-width="4"/>',
        glow('g', 8, '#ffd23f'))),
        v('Party Hat', 50, svg(`<g transform="rotate(12 250 140)"><path d="M200 146 L262 26 L306 146 Q253 160 200 146 Z" fill="url(#p)" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            `<clipPath id="cone"><path d="M200 146 L262 26 L306 146 Q253 160 200 146 Z"/></clipPath>` +
            `<g clip-path="url(#cone)" stroke="#fff" stroke-width="12" opacity=".8"><path d="M190 120 L320 70"/><path d="M190 160 L330 105"/><path d="M190 80 L320 30"/></g>` +
            `<path d="M200 146 L262 26 L306 146 Q253 160 200 146 Z" fill="none" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            `<circle cx="262" cy="24" r="16" fill="#ffe066" stroke="${OL}" stroke-width="4"/></g>`,
        lin('p', [[0, '#ff6fa8'], [1, '#8f5bff']], 0, 0, 1, 1))),
        v('Viking Helmet', 25, svg(mirror(`<path d="M148 160 C110 150 86 110 96 64 C108 92 124 112 156 126 Z" fill="url(#h)" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>`) +
            `<path d="M124 186 C124 70 376 70 376 186 Z" fill="url(#m)" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            `<path d="M120 178 Q250 156 380 178 L380 200 Q250 178 120 200 Z" fill="#c78a22" stroke="${OL}" stroke-width="4"/>` +
            `<path d="M250 92 L250 170" stroke="#c78a22" stroke-width="12"/>` +
            `<g fill="#fff6c4" stroke="${OL}" stroke-width="2">${[150, 200, 250, 300, 350].map((x) => `<circle cx="${x}" cy="${186 - Math.abs(x - 250) * 0.06}" r="4"/>`).join('')}</g>`,
        metal('m', '#eef2f6', '#a9b4bf', '#6b7682') + lin('h', [[0, '#fffbe8'], [1, '#d9cfa6']]))),
        v('Cat Ears', 60, svg(`<path d="M150 150 Q250 100 350 150" stroke="${OL}" stroke-width="12" fill="none" stroke-linecap="round"/>` +
            mirror(`<path d="M150 160 L162 70 L226 128 Z" fill="#2b2140" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
                '<path d="M166 138 L172 96 L204 126 Z" fill="#ff8fb3"/>'))),
        v('Bunny Ears', 50, svg(mirror(`<g transform="rotate(-12 205 130)"><ellipse cx="205" cy="62" rx="26" ry="72" fill="#fff" stroke="${OL}" stroke-width="5"/>` +
            '<ellipse cx="205" cy="66" rx="12" ry="52" fill="#ffb3c8"/></g>') +
            `<path d="M140 156 Q250 106 360 156" stroke="#ff8fb3" stroke-width="12" fill="none" stroke-linecap="round"/>`)),
        v('Devil Horns', 40, svg(mirror(`<path d="M178 140 C160 110 160 80 176 58 C182 88 198 108 222 122 Z" fill="url(#d)" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>`),
            lin('d', [[0, '#ff5a5a'], [1, '#9c1020']]))),
        v('Antennae', 40, svg(mirror(`<path d="M222 124 C214 96 196 76 180 64" stroke="${OL}" stroke-width="5" fill="none" stroke-linecap="round"/>` +
            `<g filter="url(#g)"><circle cx="178" cy="60" r="13" fill="#7cffb2" stroke="${OL}" stroke-width="4"/></g>`),
        glow('g', 5, '#7cffb2'))),
        v('Flower Crown', 40, svg(headTopArc(156, 24, 8, 4).map(([x, y], i) => {
            const colors = ['#ff6fa8', '#ffd23f', '#ffffff', '#b98bff'];
            return `<ellipse cx="${x - 14}" cy="${y + 4}" rx="12" ry="6" fill="#5cb85c" stroke="${OL}" stroke-width="2" transform="rotate(-30 ${x - 14} ${y + 4})"/>` +
                `<g transform="translate(${x} ${y})">${[0, 72, 144, 216, 288].map((a) => `<ellipse cy="-9" rx="7" ry="10" fill="${colors[i % 4]}" stroke="${OL}" stroke-width="2" transform="rotate(${a})"/>`).join('')}` +
                `<circle r="5" fill="#ffb21f" stroke="${OL}" stroke-width="2"/></g>`;
        }).join(''))),
        v('Propeller Cap', 25, svg(`<path d="M120 196 C120 76 380 76 380 196 Q250 176 120 196 Z" fill="#ffd23f" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            `<path d="M250 102 C220 120 200 160 196 190 L250 182 Z" fill="#e8243c"/><path d="M250 102 C280 120 300 160 304 190 L250 182 Z" fill="#2f6fd6"/>` +
            `<path d="M120 196 C120 76 380 76 380 196 Q250 176 120 196 Z" fill="none" stroke="${OL}" stroke-width="5" stroke-linejoin="round"/>` +
            `<path d="M250 100 L250 76" stroke="${OL}" stroke-width="5"/>` +
            `<ellipse cx="216" cy="72" rx="36" ry="9" fill="#ff6fa8" stroke="${OL}" stroke-width="4"/><ellipse cx="284" cy="72" rx="36" ry="9" fill="#4fd1ff" stroke="${OL}" stroke-width="4"/>` +
            `<circle cx="250" cy="74" r="7" fill="#fff" stroke="${OL}" stroke-width="3"/>`)),
        v('Headphones', 60, svg(`<path d="M120 250 C110 110 390 110 380 250" stroke="${OL}" stroke-width="18" fill="none" stroke-linecap="round"/>` +
            '<path d="M120 250 C110 110 390 110 380 250" stroke="#4a4852" stroke-width="10" fill="none" stroke-linecap="round"/>' +
            mirror(`<rect x="94" y="214" width="42" height="78" rx="18" fill="url(#h)" stroke="${OL}" stroke-width="5"/><rect x="88" y="230" width="14" height="46" rx="6" fill="#ff6fa8" stroke="${OL}" stroke-width="3"/>`),
        lin('h', [[0, '#5a5870'], [1, '#24222b']]))),
    ],
};

function headTopArc(from, to, steps, grow) {
    const points = [];
    for (let i = 0; i <= steps; i++) {
        const a = ((from + ((to - from) * i) / steps) * Math.PI) / 180;
        points.push([HEAD.cx + (HEAD.rx + grow) * Math.cos(a), HEAD.cy - (HEAD.ry + grow) * Math.sin(a)]);
    }
    return points;
}

// ---------- Улюбленець (праворуч угорі) ----------

const PET = [412, 142];
const pet = (body, defs = '') => svg(`<g transform="translate(${PET[0]} ${PET[1]})" filter="url(#ds)">${body}</g>`,
    '<filter id="ds" x="-50%" y="-50%" width="200%" height="200%"><feDropShadow dx="0" dy="8" stdDeviation="5" flood-color="#1a1030" flood-opacity=".25"/></filter>' + defs);
const petEyes = (dx = 9, y = -4, r = 4.5) => `<circle cx="${-dx}" cy="${y}" r="${r}" fill="${OL}"/><circle cx="${dx}" cy="${y}" r="${r}" fill="${OL}"/>` +
    `<circle cx="${-dx - 1.5}" cy="${y - 1.5}" r="1.5" fill="#fff"/><circle cx="${dx - 1.5}" cy="${y - 1.5}" r="1.5" fill="#fff"/>`;

const PET_LAYER = {
    name: 'Pet',
    none: 750,
    variants: [
        v('Ghost', 60, pet(`<path d="M-30 30 L-30 -6 C-30 -50 30 -50 30 -6 L30 30 L20 22 L10 30 L0 22 L-10 30 L-20 22 Z" fill="#fff" fill-opacity=".92" stroke="${OL}" stroke-width="4" stroke-linejoin="round"/>` +
            petEyes(10, -8) + '<ellipse cx="0" cy="8" rx="5" ry="6" fill="#2b2140"/><ellipse cx="-18" cy="4" rx="5" ry="3" fill="#ff8fb3" opacity=".6"/><ellipse cx="18" cy="4" rx="5" ry="3" fill="#ff8fb3" opacity=".6"/>')),
        v('Bird', 50, pet(`<path d="M24 -2 L44 4 L24 10 Z" fill="#ffb21f" stroke="${OL}" stroke-width="3" stroke-linejoin="round"/>` +
            `<circle r="30" fill="url(#b)" stroke="${OL}" stroke-width="4"/>` +
            `<path d="M-30 6 C-48 -4 -50 22 -26 22 Z" fill="#2f6fd6" stroke="${OL}" stroke-width="3"/>` +
            '<ellipse cx="2" cy="14" rx="16" ry="11" fill="#e6f5ff"/>' + petEyes(10, -6, 4) +
            `<path d="M-6 -30 C-2 -44 8 -44 6 -30" fill="#5fa6e3" stroke="${OL}" stroke-width="3"/>`,
        rad('b', [[0, '#a6d8ff'], [1, '#2f6fd6']], 0.4, 0.35, 0.75))),
        v('Slime', 50, pet(`<path d="M-36 26 C-40 -10 -20 -36 0 -36 C20 -36 40 -10 36 26 Q0 34 -36 26 Z" fill="url(#s)" stroke="${OL}" stroke-width="4"/>` +
            '<ellipse cx="-14" cy="-20" rx="9" ry="5" fill="#fff" opacity=".6" transform="rotate(-30 -14 -20)"/>' + petEyes(11, 0, 5) +
            `<path d="M-7 12 Q0 18 7 12" stroke="${OL}" stroke-width="3" fill="none" stroke-linecap="round"/>`,
        rad('s', [[0, '#c8ffb0'], [1, '#4fbf3a']], 0.4, 0.3, 0.8))),
        v('Robot Drone', 30, pet(`<path d="M-34 -36 L34 -36 M0 -36 L0 -26" stroke="${OL}" stroke-width="4"/>` +
            `<ellipse cx="-22" cy="-38" rx="16" ry="4" fill="#c8ccd2" stroke="${OL}" stroke-width="2.5"/><ellipse cx="22" cy="-38" rx="16" ry="4" fill="#c8ccd2" stroke="${OL}" stroke-width="2.5"/>` +
            `<circle r="28" fill="url(#m)" stroke="${OL}" stroke-width="4"/>` +
            `<g filter="url(#g)"><circle cx="0" cy="2" r="11" fill="#4ff3ff" stroke="${OL}" stroke-width="3"/></g><circle cx="-3" cy="-1" r="3" fill="#fff"/>`,
        metal('m', '#f2f5f8', '#a9b4bf', '#6b7682') + glow('g', 4, '#4ff3ff'))),
        v('Kitten', 40, pet(`<path d="M-30 -10 L-26 -42 L-6 -26 Z M30 -10 L26 -42 L6 -26 Z" fill="#ff9f43" stroke="${OL}" stroke-width="4" stroke-linejoin="round"/>` +
            `<ellipse rx="34" ry="28" fill="url(#k)" stroke="${OL}" stroke-width="4"/>` +
            '<path d="M-24 -18 L-22 -34 L-12 -24 Z M24 -18 L22 -34 L12 -24 Z" fill="#ffb3c8"/>' + petEyes(12, -2, 4.5) +
            `<path d="M-4 8 L4 8 L0 12 Z" fill="#ff6f8e"/><path d="M0 12 Q-5 18 -10 14 M0 12 Q5 18 10 14" stroke="${OL}" stroke-width="2.5" fill="none" stroke-linecap="round"/>` +
            `<path d="M-14 -20 l4 8 M0 -24 l0 8 M14 -20 l-4 8" stroke="#d9772b" stroke-width="3" stroke-linecap="round"/>`,
        rad('k', [[0, '#ffd1a0'], [1, '#ff9f43']], 0.4, 0.35, 0.75))),
        v('Jellyfish', 30, pet(`<g stroke="#ff8fd0" stroke-width="4" fill="none" stroke-linecap="round" opacity=".9"><path d="M-18 6 q-6 14 0 26 q6 12 0 22"/><path d="M-6 8 q6 14 0 26 q-6 12 0 22"/><path d="M6 8 q-6 14 0 26 q6 12 0 22"/><path d="M18 6 q6 14 0 26 q-6 12 0 22"/></g>` +
            `<path d="M-34 8 C-34 -40 34 -40 34 8 Q17 0 0 8 Q-17 0 -34 8 Z" fill="url(#j)" stroke="${OL}" stroke-width="4" stroke-linejoin="round"/>` + petEyes(10, -10, 4),
        rad('j', [[0, '#ffe0f4'], [1, '#ff6fc8']], 0.4, 0.3, 0.8))),
        v('Mini Dragon', 10, pet(`<path d="M-20 -8 C-50 -40 -60 -10 -44 6 C-36 -6 -28 -6 -20 0 Z M20 -8 C50 -40 60 -10 44 6 C36 -6 28 -6 20 0 Z" fill="#7fd36b" stroke="${OL}" stroke-width="3.5" stroke-linejoin="round"/>` +
            `<circle r="28" fill="url(#d)" stroke="${OL}" stroke-width="4"/>` +
            `<path d="M-14 -24 L-18 -40 L-6 -28 Z M14 -24 L18 -40 L6 -28 Z" fill="#fff6c4" stroke="${OL}" stroke-width="3" stroke-linejoin="round"/>` +
            '<ellipse cx="0" cy="12" rx="14" ry="9" fill="#d4f7b8"/>' + petEyes(11, -4, 5) +
            '<circle cx="-4" cy="12" r="1.8" fill="#2f7a4a"/><circle cx="4" cy="12" r="1.8" fill="#2f7a4a"/>' +
            `<path d="M30 10 q12 -4 14 -16" stroke="#ff7a1a" stroke-width="5" fill="none" stroke-linecap="round" filter="url(#f)"/>`,
        rad('d', [[0, '#c8ffb0'], [1, '#3e9b57']], 0.4, 0.35, 0.75) + glow('f', 3, '#ffb21f'))),
        v('UFO', 20, pet('<path d="M-20 18 L-44 70 L44 70 L20 18 Z" fill="url(#beam)"/>' +
            `<ellipse cx="0" cy="-6" rx="20" ry="18" fill="#bff3ff" fill-opacity=".85" stroke="${OL}" stroke-width="4"/>` +
            `<ellipse cx="0" cy="8" rx="46" ry="15" fill="url(#u)" stroke="${OL}" stroke-width="4"/>` +
            '<g fill="#ffe066"><circle cx="-28" cy="9" r="4"/><circle cx="-10" cy="12" r="4"/><circle cx="10" cy="12" r="4"/><circle cx="28" cy="9" r="4"/></g>',
        metal('u', '#eef2f6', '#a9b4bf', '#6b7682') + lin('beam', [[0, '#fff6a6', 0.7], [1, '#fff6a6', 0]]))),
    ],
};

// ---------- Частинки (поза обличчям) ----------

function scatter(seed, count, draw) {
    const r = rng(seed);
    let out = '';
    let placed = 0;
    for (let tries = 0; placed < count && tries < 2000; tries++) {
        const x = 15 + r() * 470;
        const y = 15 + r() * 470;
        const inHead = ((x - HEAD.cx) / (HEAD.rx + 18)) ** 2 + ((y - HEAD.cy) / (HEAD.ry + 18)) ** 2 < 1;
        const inBody = y > 350 && Math.abs(x - 250) < 160;
        if (inHead || inBody) continue;
        out += draw(x, y, r, placed++);
    }
    return out;
}

const note = (x, y, s) => `<g transform="translate(${x.toFixed(0)} ${y.toFixed(0)}) scale(${s.toFixed(2)})"><ellipse cx="-6" cy="10" rx="8" ry="6" transform="rotate(-20 -6 10)"/><path d="M1 8 L1 -18 L14 -12 L14 -4 L4 -8 L4 8 Z"/></g>`;
const flake = (x, y, s) => `<g transform="translate(${x.toFixed(0)} ${y.toFixed(0)}) scale(${s.toFixed(2)})">${[0, 60, 120].map((a) => `<path d="M0 -12 L0 12 M-4 -9 L0 -5 L4 -9 M-4 9 L0 5 L4 9" transform="rotate(${a})"/>`).join('')}</g>`;

const PARTICLES = {
    name: 'Particles',
    none: 450,
    variants: [
        v('Sparkles', 80, svg(`<g filter="url(#g)">${scatter(3, 16, (x, y, r, i) => `<path d="${sparkle(x, y, 6 + r() * 10)}" fill="${i % 3 ? '#fff' : '#ffe066'}"/>`)}</g>`, glow('g', 3, '#fff6c4'))),
        v('Bubbles', 70, svg(scatter(5, 14, (x, y, r) => {
            const s = 8 + r() * 14;
            return `<circle cx="${x.toFixed(0)}" cy="${y.toFixed(0)}" r="${s.toFixed(0)}" fill="#e8fbff" fill-opacity=".25" stroke="#fff" stroke-width="2.5" stroke-opacity=".85"/>` +
                `<path d="M${(x - s * 0.5).toFixed(0)} ${(y - s * 0.2).toFixed(0)} a${(s * 0.55).toFixed(0)} ${(s * 0.55).toFixed(0)} 0 0 1 ${(s * 0.45).toFixed(0)} ${(-s * 0.35).toFixed(0)}" stroke="#fff" stroke-width="2.5" fill="none" stroke-linecap="round"/>`;
        }))),
        v('Hearts', 60, svg(scatter(7, 12, (x, y, r, i) => `<path d="${heart(x, y, 8 + r() * 8)}" fill="${['#ff4d6d', '#ff8fb3', '#ffb3c8'][i % 3]}" stroke="${OL}" stroke-width="2" opacity=".9" transform="rotate(${((r() - 0.5) * 40).toFixed(0)} ${x.toFixed(0)} ${y.toFixed(0)})"/>`))),
        v('Snowflakes', 50, svg(`<g stroke="#fff" stroke-width="2.5" stroke-linecap="round" fill="none" opacity=".9">${scatter(11, 16, (x, y, r) => flake(x, y, 0.6 + r() * 0.8))}</g>`)),
        v('Fireflies', 50, svg(`<g filter="url(#g)">${scatter(13, 22, (x, y, r) => `<circle cx="${x.toFixed(0)}" cy="${y.toFixed(0)}" r="${(2 + r() * 3).toFixed(1)}" fill="#fff6a6"/>`)}</g>`, glow('g', 5, '#ffe066'))),
        v('Music Notes', 40, svg(`<g fill="${OL}" opacity=".85">${scatter(17, 9, (x, y, r) => note(x, y, 0.9 + r() * 0.6))}</g>`)),
        v('Coins', 20, svg(scatter(19, 11, (x, y, r) => `<g transform="translate(${x.toFixed(0)} ${y.toFixed(0)}) rotate(${((r() - 0.5) * 50).toFixed(0)})">` +
            `<ellipse rx="13" ry="13" fill="url(#c)" stroke="${OL}" stroke-width="2.5"/><path d="${star(0, 0, 6, 2.6)}" fill="#fff6c4"/></g>`),
        rad('c', [[0, '#fff1a8'], [0.6, '#f2c14e'], [1, '#c78a22']], 0.4, 0.35, 0.75))),
        v('Petals', 50, svg(scatter(23, 16, (x, y, r) => `<ellipse cx="${x.toFixed(0)}" cy="${y.toFixed(0)}" rx="${(5 + r() * 4).toFixed(1)}" ry="${(9 + r() * 5).toFixed(1)}" fill="url(#p)" opacity=".95" transform="rotate(${(r() * 180).toFixed(0)} ${x.toFixed(0)} ${y.toFixed(0)})"/>`),
        lin('p', [[0, '#ffe0ef'], [1, '#ff8fb8']]))),
    ],
};

// ---------- Освітлення ----------

function glitch(seed) {
    const r = rng(seed);
    let bands = '';
    for (let i = 0; i < 9; i++) {
        const y = r() * 480;
        bands += `<rect x="${(r() * 80 - 40).toFixed(0)}" y="${y.toFixed(0)}" width="${(200 + r() * 340).toFixed(0)}" height="${(3 + r() * 12).toFixed(0)}" fill="${i % 2 ? '#ff2fd0' : '#2fe8ff'}" opacity=".35"/>`;
    }
    return bands + '<rect width="500" height="500" fill="url(#scan)" opacity=".35"/>';
}

const LIGHTING = {
    name: 'Lighting',
    none: 650,
    variants: [
        v('Vignette', 80, svg('<rect width="500" height="500" fill="url(#v)"/>', rad('v', [[0, '#120a1c', 0], [0.62, '#120a1c', 0], [1, '#120a1c', 0.6]], 0.5, 0.5, 0.72))),
        v('Light Leak', 50, svg('<g filter="url(#b)"><ellipse cx="40" cy="60" rx="170" ry="120" fill="#ff9a3c" opacity=".38"/><ellipse cx="490" cy="470" rx="150" ry="110" fill="#ff3d8b" opacity=".3"/></g>', blur('b', 40))),
        v('Lens Flare', 40, svg('<g filter="url(#b)"><circle cx="410" cy="70" r="46" fill="#fffbe0" opacity=".8"/></g>' +
            '<circle cx="410" cy="70" r="14" fill="#fff"/>' +
            '<path d="M330 70 L490 70 M410 -10 L410 150" stroke="#fff" stroke-width="3" opacity=".7" filter="url(#b2)"/>' +
            [[330, 150, 18, '#9fefff'], [270, 210, 9, '#ffd1f7'], [180, 300, 28, '#c2f0ff'], [120, 360, 12, '#fff6a6']].map(([x, y, rr, c]) =>
                `<circle cx="${x}" cy="${y}" r="${rr}" fill="${c}" opacity=".35"/>`).join(''),
        blur('b', 14) + blur('b2', 1.5))),
        v('Neon Rim', 30, svg(`<g filter="url(#g)" fill="none" stroke="#4ff3ff" stroke-width="3"><ellipse cx="${HEAD.cx}" cy="${HEAD.cy}" rx="${HEAD.rx + 4}" ry="${HEAD.ry + 4}"/>` +
            '<path d="M104 500 C104 418 158 356 250 356 C342 356 396 418 396 500"/></g>', glow('g', 6, '#4ff3ff'))),
        v('Spotlight', 40, svg('<path d="M200 0 L300 0 L470 500 L30 500 Z" fill="url(#s)" filter="url(#b)"/>',
            lin('s', [[0, '#fffbe0', 0.45], [1, '#fffbe0', 0]]) + blur('b', 16))),
        v('Sunbeams', 40, svg(`<g filter="url(#b)" fill="#fff" opacity=".2">${[0, 1, 2, 3].map((i) =>
            `<path d="M${-60 + i * 70} -20 L${10 + i * 70} -20 L${260 + i * 110} 520 L${180 + i * 110} 520 Z"/>`).join('')}</g>`, blur('b', 10))),
        v('Glitch', 10, svg(glitch(29),
            '<pattern id="scan" width="4" height="4" patternUnits="userSpaceOnUse"><rect width="4" height="1.5" fill="#000"/></pattern>')),
    ],
};

// ---------- Рамка ----------

const border = (stroke, width, extra = '', defs = '', inset = width / 2) => svg(
    `<rect x="${inset}" y="${inset}" width="${500 - 2 * inset}" height="${500 - 2 * inset}" fill="none" stroke="${stroke}" stroke-width="${width}"/>` + extra, defs);
const corners = (draw) => [[30, 30], [470, 30], [30, 470], [470, 470]].map(([x, y]) => draw(x, y)).join('');

function pixelBorder() {
    const colors = ['#ff5f6d', '#ffd23f', '#4fd1c5', '#8f7bff'];
    let out = '';
    for (let i = 0; i < 25; i++) {
        const c = colors[i % 4];
        out += `<rect x="${i * 20}" y="0" width="20" height="20" fill="${c}"/><rect x="${480 - i * 20}" y="480" width="20" height="20" fill="${c}"/>`;
        out += `<rect x="0" y="${480 - i * 20}" width="20" height="20" fill="${colors[(i + 2) % 4]}"/><rect x="480" y="${i * 20}" width="20" height="20" fill="${colors[(i + 2) % 4]}"/>`;
    }
    return out;
}

const FRAME = {
    name: 'Frame',
    none: 750,
    variants: [
        v('Gold Frame', 40, border('url(#g)', 26,
            '<rect x="26" y="26" width="448" height="448" fill="none" stroke="#8a5a12" stroke-width="3"/><rect x="3" y="3" width="494" height="494" fill="none" stroke="#8a5a12" stroke-width="3"/>' +
            corners((x, y) => `<circle cx="${x - Math.sign(x - 250) * 17}" cy="${y - Math.sign(y - 250) * 17}" r="9" fill="#fff3b0" stroke="#8a5a12" stroke-width="2.5"/>`),
            lin('g', [[0, '#fff3b0'], [0.3, '#e2a93a'], [0.6, '#ffe58f'], [1, '#b97d1f']], 0, 0, 1, 1))),
        v('Silver Frame', 40, border('url(#s)', 22,
            '<rect x="22" y="22" width="456" height="456" fill="none" stroke="#5d6672" stroke-width="3"/>',
            lin('s', [[0, '#ffffff'], [0.4, '#a9b4bf'], [0.7, '#eef2f5'], [1, '#6b7682']], 0, 0, 1, 1))),
        v('Neon Frame', 30, svg('<g filter="url(#g)"><rect x="16" y="16" width="468" height="468" rx="26" fill="none" stroke="#ff4fd8" stroke-width="7"/></g>' +
            '<rect x="16" y="16" width="468" height="468" rx="26" fill="none" stroke="#ffe0fa" stroke-width="2.5"/>', glow('g', 8, '#ff4fd8'))),
        v('Polaroid', 30, svg('<path d="M0 0 H500 V500 H0 Z M20 20 V430 H480 V20 Z" fill="url(#p)" fill-rule="evenodd"/>' +
            '<path d="M20 20 V430 H480 V20 Z" fill="none" stroke="#d8d2c4" stroke-width="2"/>' +
            '<path d="M150 466 C190 452 220 476 260 462 S330 456 350 466" stroke="#2f4a8a" stroke-width="3" fill="none" stroke-linecap="round"/>',
        lin('p', [[0, '#ffffff'], [1, '#f1ece0']]))),
        v('Pixel Border', 25, svg(pixelBorder())),
        v('Rounded Card', 40, svg('<path d="M0 0 H500 V500 H0 Z M60 14 H440 A46 46 0 0 1 486 60 V440 A46 46 0 0 1 440 486 H60 A46 46 0 0 1 14 440 V60 A46 46 0 0 1 60 14 Z" fill="#1d1b26" fill-rule="evenodd"/>' +
            '<rect x="14" y="14" width="472" height="472" rx="46" fill="none" stroke="#fff" stroke-width="3" opacity=".8"/>')),
        v('Wood Frame', 30, border('url(#w)', 28,
            `<g stroke="#5a3418" stroke-width="2" opacity=".45" fill="none">${[6, 12, 20].map((d) => `<rect x="${d}" y="${d}" width="${500 - 2 * d}" height="${500 - 2 * d}" rx="2"/>`).join('')}</g>` +
            '<rect x="28" y="28" width="444" height="444" fill="none" stroke="#3b2210" stroke-width="3"/>',
            lin('w', [[0, '#c98b52'], [0.5, '#9a5a2e'], [1, '#c98b52']], 0, 0, 1, 1))),
        v('Holo Frame', 5, svg('<rect x="11" y="11" width="478" height="478" rx="14" fill="none" stroke="url(#h)" stroke-width="22"/>' +
            corners((x, y) => `<path d="${sparkle(x, y, 18)}" fill="#fff"/>`),
        lin('h', [[0, '#ff9ad5'], [0.25, '#9ad7ff'], [0.5, '#a1ffce'], [0.75, '#fff59a'], [1, '#d6a6ff']], 0, 0, 1, 1))),
    ],
};

module.exports = { HEADGEAR, PET_LAYER, PARTICLES, LIGHTING, FRAME };
