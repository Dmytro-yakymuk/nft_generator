// Шари тіла персонажа: шкіра (голова й плечі), візерунок шкіри, одяг, значок, намисто, сережки.

const {
    OL, SW, HEAD, TORSO, HEAD_ELLIPSE, rng, svg, lin, rad, shade, blur, glow, behindHead, star, heart, sparkle, mirror, v,
} = require('./lib');

// ---------- Шкіра ----------

const HIGHLIGHT = '<ellipse cx="188" cy="170" rx="44" ry="20" fill="#fff" opacity=".38" transform="rotate(-24 188 170)"/>' +
    '<ellipse cx="160" cy="205" rx="10" ry="6" fill="#fff" opacity=".3" transform="rotate(-40 160 205)"/>';

function skin(light, base, dark, { headFill = 'url(#head)', torsoFill = 'url(#torso)', extra = '', defs = '', opacity = 1 } = {}) {
    return svg(
        `<g opacity="${opacity}">` +
        `<path d="${TORSO}" fill="${torsoFill}" stroke="${OL}" stroke-width="${SW}" stroke-linejoin="round"/>` +
        `<ellipse cx="250" cy="372" rx="108" ry="22" fill="${dark}" opacity=".55" clip-path="url(#t)"/>` +
        `<ellipse ${HEAD_ELLIPSE} fill="${headFill}" stroke="${OL}" stroke-width="${SW}"/>` +
        extra + HIGHLIGHT + '</g>',
        shade('head', light, base, dark) + lin('torso', [[0, base], [1, dark]]) +
        `<clipPath id="t"><path d="${TORSO}"/></clipPath>` + defs);
}

function metallic(id, bands) {
    return lin(id, bands.map((color, i) => [i / (bands.length - 1), color]), 0, 0, 0.3, 1);
}

function galaxyInside(seed) {
    const r = rng(seed);
    let stars = '';
    for (let i = 0; i < 70; i++) {
        stars += `<circle cx="${(110 + r() * 280).toFixed(0)}" cy="${(115 + r() * 385).toFixed(0)}" r="${(0.6 + r() * 1.8).toFixed(1)}" fill="#fff" opacity="${(0.4 + r() * 0.6).toFixed(2)}"/>`;
    }
    return '<g clip-path="url(#body)">' +
        '<g filter="url(#nb)"><ellipse cx="200" cy="200" rx="90" ry="60" fill="#ff5fd2" opacity=".45"/>' +
        '<ellipse cx="320" cy="300" rx="90" ry="70" fill="#4f8bff" opacity=".45"/>' +
        '<ellipse cx="250" cy="450" rx="120" ry="50" fill="#9b5fff" opacity=".4"/></g>' +
        stars + `<path d="${sparkle(300, 180, 12)}" fill="#fff"/><path d="${sparkle(170, 420, 9)}" fill="#fff"/></g>`;
}

const BODY_CLIP = `<clipPath id="body"><ellipse ${HEAD_ELLIPSE}/><path d="${TORSO}"/></clipPath>`;

const SKIN = {
    name: 'Skin',
    none: 0,
    variants: [
        v('Peach', 130, skin('#ffeadb', '#ffcfae', '#ee9f7b')),
        v('Mint', 110, skin('#e2fff2', '#a8edd0', '#5fc39c')),
        v('Lavender', 110, skin('#f2eaff', '#d2bdff', '#9c7fe0')),
        v('Sky', 110, skin('#e6f5ff', '#a6d8ff', '#5fa6e3')),
        v('Lemon', 100, skin('#fffbe0', '#ffe58a', '#e9b93c')),
        v('Coral', 100, skin('#ffe0db', '#ff9e8f', '#dd6150')),
        v('Snow', 90, skin('#ffffff', '#eef1fb', '#bcc5e3')),
        v('Charcoal', 70, skin('#8c87a6', '#5d5873', '#36324a')),
        v('Alien', 60, skin('#e2ffc8', '#9ee36a', '#56a82f')),
        v('Ghost', 30, skin('#ffffff', '#e8f0ff', '#a9bce8', { opacity: 0.78 })),
        v('Chrome', 20, skin('#ffffff', '#c3ccd6', '#6b7682', {
            headFill: 'url(#metal)', torsoFill: 'url(#metal)',
            defs: metallic('metal', ['#f7f9fb', '#a9b4bf', '#eef2f5', '#7d8894', '#dfe5ea', '#66717d']),
        })),
        v('Galaxy', 10, skin('#6a4fb3', '#2c1d5a', '#120b2b', {
            extra: galaxyInside(77), defs: BODY_CLIP + blur('nb', 22),
        })),
        v('Gold', 8, skin('#fff6c4', '#f2c14e', '#a8701a', {
            headFill: 'url(#metal)', torsoFill: 'url(#metal)',
            defs: metallic('metal', ['#fff7cf', '#e2a93a', '#ffe58f', '#b97d1f', '#ffd76a', '#9c6414']),
        })),
    ],
};

// ---------- Візерунок шкіри (у межах голови й плечей) ----------

const onBody = (body, defs = '') => svg(`<g clip-path="url(#body)">${body}</g>`, BODY_CLIP + defs);
const DARK = '#2b1838';

function spots(seed) {
    const r = rng(seed);
    let out = '';
    for (let i = 0; i < 26; i++) {
        const x = 120 + r() * 260;
        const y = 125 + r() * 370;
        out += `<ellipse cx="${x.toFixed(0)}" cy="${y.toFixed(0)}" rx="${(7 + r() * 13).toFixed(0)}" ry="${(6 + r() * 10).toFixed(0)}" transform="rotate(${(r() * 180).toFixed(0)} ${x.toFixed(0)} ${y.toFixed(0)})"/>`;
    }
    return `<g fill="${DARK}" opacity=".17">${out}</g>`;
}

function stripes() {
    let out = '';
    for (const [y, len] of [[150, 70], [190, 85], [230, 70], [270, 60], [410, 80], [450, 95]]) {
        out += `<path d="M105 ${y} q${len * 0.6} 10 ${len} 28 q-${len * 0.5} -4 -${len} 6 Z"/>`;
    }
    return `<g fill="${DARK}" opacity=".22">${mirror(out)}</g>` +
        `<path d="M235 128 q15 22 30 0 q-6 30 -15 40 q-9 -10 -15 -40 Z" fill="${DARK}" opacity=".22"/>`;
}

function sparkleSkin(seed) {
    const r = rng(seed);
    let out = '';
    for (let i = 0; i < 22; i++) out += `<path d="${sparkle(120 + r() * 260, 130 + r() * 360, 3 + r() * 7)}"/>`;
    return `<g fill="#fff" opacity=".75">${out}</g>`;
}

const patches = () => `<g fill="${DARK}" opacity=".16">` +
    '<path d="M120 200 C140 160 200 170 195 210 C190 250 140 260 120 240 Z"/>' +
    '<path d="M300 140 C340 130 380 170 360 200 C340 225 300 200 300 170 Z"/>' +
    '<path d="M280 420 C320 400 380 430 370 470 C360 500 300 500 285 470 Z"/></g>';

const scales = () => onBody(
    '<pattern id="sc" width="28" height="18" patternUnits="userSpaceOnUse">' +
    `<path d="M0 18 A14 14 0 0 1 28 18" fill="none" stroke="${DARK}" stroke-width="2.5"/>` +
    `<path d="M-14 9 A14 14 0 0 1 14 9 M14 9 A14 14 0 0 1 42 9" fill="none" stroke="${DARK}" stroke-width="2.5"/></pattern>` +
    '<rect width="500" height="500" fill="url(#sc)" opacity=".14"/>');

const tribal = () => '<g fill="none" stroke="#fff" stroke-width="5" stroke-linecap="round" stroke-linejoin="round" opacity=".55">' +
    '<path d="M205 150 L220 170 L235 150 L250 170 L265 150 L280 170 L295 150"/>' +
    mirror('<path d="M140 380 L170 410 L140 440 L170 470"/>') + '</g>';

function constellation() {
    const points = [[170, 170], [205, 140], [245, 160], [285, 135], [330, 165], [355, 210], [150, 420], [190, 450], [230, 430]];
    const lines = [[0, 1], [1, 2], [2, 3], [3, 4], [4, 5], [6, 7], [7, 8]];
    return '<g stroke="#fff" stroke-width="1.6" opacity=".7">' +
        lines.map(([a, b]) => `<line x1="${points[a][0]}" y1="${points[a][1]}" x2="${points[b][0]}" y2="${points[b][1]}"/>`).join('') + '</g>' +
        '<g fill="#fff">' + points.map(([x, y]) => `<path d="${sparkle(x, y, 7)}"/>`).join('') + '</g>';
}

const cracks = () => `<g fill="none" stroke="${DARK}" stroke-width="2.5" stroke-linecap="round" opacity=".45">` +
    '<path d="M130 230 L160 240 L170 270 L200 280 M160 240 L175 225"/>' +
    '<path d="M370 180 L345 200 L350 225 L325 245 M345 200 L330 190"/>' +
    '<path d="M300 470 L320 440 L350 445 L360 420"/></g>';

const SKIN_PATTERN = {
    name: 'Skin Pattern',
    none: 450,
    variants: [
        v('Spots', 60, onBody(spots(5))),
        v('Tiger Stripes', 50, onBody(stripes())),
        v('Sparkle Skin', 40, onBody(sparkleSkin(9))),
        v('Patches', 50, onBody(patches())),
        v('Scales', 30, scales()),
        v('Tribal', 30, onBody(tribal())),
        v('Constellation', 20, onBody(constellation())),
        v('Cracks', 15, onBody(cracks())),
    ],
};

// ---------- Одяг (позаду голови) ----------

const wear = (body, defs = '') => svg(
    `<g mask="url(#nohead)"><g clip-path="url(#torso)">${body}</g>` +
    `<path d="${TORSO}" fill="none" stroke="${OL}" stroke-width="${SW}" stroke-linejoin="round"/></g>`,
    behindHead('nohead') + `<clipPath id="torso"><path d="${TORSO}"/></clipPath>` + defs);

const fill = (paint) => `<rect width="500" height="500" fill="${paint}"/>`;
const CREASES = `<path d="M150 470 C165 440 175 430 180 420 M350 470 C335 440 325 430 320 420" stroke="${OL}" stroke-width="3" fill="none" opacity=".25"/>`;

const hoodie = (light, base, dark) => wear(
    fill('url(#c)') + CREASES +
    `<path d="M150 410 C170 370 330 370 350 410 C330 395 170 395 150 410 Z" fill="${dark}" stroke="${OL}" stroke-width="4"/>` +
    `<path d="M225 395 C222 420 220 440 218 455 M275 395 C278 420 280 440 282 455" stroke="#fff" stroke-width="5" fill="none" stroke-linecap="round"/>` +
    `<circle cx="218" cy="459" r="5" fill="#fff" stroke="${OL}" stroke-width="2"/><circle cx="282" cy="459" r="5" fill="#fff" stroke="${OL}" stroke-width="2"/>` +
    `<path d="M185 500 L195 462 C230 455 270 455 305 462 L315 500" fill="${dark}" stroke="${OL}" stroke-width="4" opacity=".9"/>`,
    lin('c', [[0, light], [0.5, base], [1, dark]]));

function stripedTee() {
    let bands = '';
    for (let y = 380; y < 500; y += 26) bands += `<rect y="${y}" width="500" height="12" fill="#2f4a8a"/>`;
    return wear(fill('#fbfbff') + bands +
        `<path d="M195 372 C215 400 285 400 305 372" fill="none" stroke="${OL}" stroke-width="5"/>` + CREASES);
}

const whiteTee = () => wear(fill('url(#c)') +
    `<path d="M200 370 C220 398 280 398 300 370" fill="none" stroke="#c9cfe0" stroke-width="10"/>` +
    `<path d="M200 370 C220 398 280 398 300 370" fill="none" stroke="${OL}" stroke-width="4"/>` + CREASES,
lin('c', [[0, '#ffffff'], [1, '#d9deee']]));

const tuxedo = () => wear(fill('url(#c)') +
    '<path d="M215 370 L250 460 L285 370 Z" fill="#fff"/>' +
    `<path d="M150 400 L215 370 L250 460 L230 500 L160 500 Z" fill="#1d1b26" stroke="${OL}" stroke-width="3"/>` +
    `<path d="M350 400 L285 370 L250 460 L270 500 L340 500 Z" fill="#1d1b26" stroke="${OL}" stroke-width="3"/>` +
    '<path d="M215 370 L240 430 L222 425 Z M285 370 L260 430 L278 425 Z" fill="#3a3746"/>' +
    `<path d="M232 395 L250 405 L268 395 L268 415 L250 405 L232 415 Z" fill="#111" stroke="${OL}" stroke-width="2"/>` +
    '<circle cx="250" cy="440" r="3.5" fill="#222"/><circle cx="250" cy="470" r="3.5" fill="#222"/>',
lin('c', [[0, '#3b3948'], [1, '#14121c']]));

const kimono = () => wear(fill('url(#c)') +
    '<g fill="#fff" opacity=".85">' + [[150, 440], [330, 420], [370, 480], [190, 490]].map(([x, y]) =>
        `<g transform="translate(${x} ${y})">${[0, 72, 144, 216, 288].map((a) => `<ellipse cx="0" cy="-8" rx="5" ry="9" transform="rotate(${a})"/>`).join('')}</g>`).join('') + '</g>' +
    `<path d="M205 372 L265 500 L240 500 L190 385 Z" fill="#f2c14e" stroke="${OL}" stroke-width="3"/>` +
    `<path d="M295 372 L235 500 L260 500 L310 385 Z" fill="#f2c14e" stroke="${OL}" stroke-width="3"/>` +
    `<rect x="120" y="470" width="260" height="30" fill="#2b2140" stroke="${OL}" stroke-width="3"/>`,
lin('c', [[0, '#e8475a'], [1, '#a81f3a']]));

const spaceSuit = () => wear(fill('url(#c)') + CREASES +
    `<path d="M165 395 C190 370 310 370 335 395 L330 405 C300 385 200 385 170 405 Z" fill="#ff8a3d" stroke="${OL}" stroke-width="4"/>` +
    `<rect x="200" y="425" width="100" height="60" rx="10" fill="#dfe6f0" stroke="${OL}" stroke-width="4"/>` +
    '<circle cx="225" cy="445" r="7" fill="#ff5f6d"/><circle cx="250" cy="445" r="7" fill="#ffd23f"/><circle cx="275" cy="445" r="7" fill="#4fd1ff"/>' +
    '<rect x="215" y="463" width="70" height="8" rx="4" fill="#7d8a99"/>' +
    `<path d="${star(150, 450, 12, 5)}" fill="#4f8bff"/>`,
lin('c', [[0, '#ffffff'], [1, '#c9d2e0']]));

const armor = () => wear(fill('url(#c)') +
    `<g fill="none" stroke="${OL}" stroke-width="4"><path d="M120 440 C180 420 320 420 380 440"/><path d="M118 475 C180 455 320 455 382 475"/><path d="M250 390 L250 500"/></g>` +
    `<g fill="#ffe08a" stroke="${OL}" stroke-width="2">${[[200, 430], [300, 430], [180, 465], [320, 465]].map(([x, y]) => `<circle cx="${x}" cy="${y}" r="5"/>`).join('')}</g>` +
    `<path d="M140 400 C150 380 190 372 215 372 L205 400 Z M360 400 C350 380 310 372 285 372 L295 400 Z" fill="#c3ccd6" stroke="${OL}" stroke-width="4"/>`,
lin('c', [[0, '#eef2f6'], [0.4, '#9aa6b2'], [0.7, '#d5dce3'], [1, '#6b7682']]));

const varsity = () => wear(fill('#2e7d5b') +
    `<path d="M108 500 C108 440 130 400 165 385 L185 500 Z M392 500 C392 440 370 400 335 385 L315 500 Z" fill="#f3ead2" stroke="${OL}" stroke-width="4"/>` +
    `<path d="M200 372 C220 395 280 395 300 372" fill="none" stroke="#f3ead2" stroke-width="12"/>` +
    '<g fill="#f3ead2"><circle cx="250" cy="420" r="4"/><circle cx="250" cy="445" r="4"/><circle cx="250" cy="470" r="4"/></g>' +
    `<text x="210" y="470" font-family="Georgia, serif" font-weight="bold" font-size="44" fill="#f3ead2" stroke="${OL}" stroke-width="2">G</text>`);

const leather = () => wear(fill('url(#c)') +
    `<path d="M215 372 L245 500 L255 500 L285 372 Z" fill="#e8e3da"/>` +
    `<path d="M160 395 L215 372 L238 440 L200 430 Z M340 395 L285 372 L262 440 L300 430 Z" fill="#2b2a30" stroke="${OL}" stroke-width="3"/>` +
    '<path d="M244 440 L244 500" stroke="#c8ccd2" stroke-width="3" stroke-dasharray="4 3"/>' +
    '<rect x="300" y="450" width="40" height="5" rx="2" fill="#c8ccd2"/>',
lin('c', [[0, '#4a4852'], [1, '#18171c']]));

const overalls = () => wear(fill('#ffd54f') +
    `<path d="M180 420 L320 420 L330 500 L170 500 Z" fill="url(#d)" stroke="${OL}" stroke-width="4"/>` +
    `<path d="M175 380 L195 420 M325 380 L305 420" stroke="#3b6fb6" stroke-width="16" stroke-linecap="round"/>` +
    `<circle cx="196" cy="425" r="6" fill="#f2c14e" stroke="${OL}" stroke-width="2"/><circle cx="304" cy="425" r="6" fill="#f2c14e" stroke="${OL}" stroke-width="2"/>` +
    `<rect x="225" y="440" width="50" height="34" rx="4" fill="none" stroke="#9cc3ff" stroke-width="3" stroke-dasharray="5 3"/>`,
lin('d', [[0, '#5b8fd6'], [1, '#2f5b9c']]));

const turtleneck = () => wear(fill('url(#c)') +
    `<path d="M175 380 C190 355 310 355 325 380 L325 410 C300 395 200 395 175 410 Z" fill="#e6a531" stroke="${OL}" stroke-width="4"/>` +
    '<g stroke="#b97d1f" stroke-width="3" opacity=".6">' +
    Array.from({ length: 13 }, (_, i) => `<line x1="${185 + i * 10}" y1="378" x2="${185 + i * 10}" y2="402"/>`).join('') + '</g>' +
    '<g stroke="#b97d1f" stroke-width="2" opacity=".35" fill="none">' +
    [140, 180, 220, 260, 300, 340].map((x) => `<path d="M${x} 420 l8 10 l-8 10 l8 10 l-8 10 l8 10"/>`).join('') + '</g>',
lin('c', [[0, '#ffc95c'], [1, '#d1901f']]));

const goldRobe = () => wear(fill('url(#c)') +
    `<path d="M200 372 L250 480 L300 372" fill="none" stroke="#7a1f3d" stroke-width="18"/>` +
    `<path d="M200 372 L250 480 L300 372" fill="none" stroke="#fff3b0" stroke-width="4"/>` +
    `<g fill="#7a1f3d">${[[160, 440], [340, 440], [250, 495]].map(([x, y]) => `<path d="${star(x, y, 14, 6)}"/>`).join('')}</g>` +
    '<rect width="500" height="500" fill="url(#sheen)"/>',
lin('c', [[0, '#fff1a8'], [0.35, '#e6b23a'], [0.6, '#ffe58f'], [1, '#a8701a']], 0, 0, 1, 1) +
    lin('sheen', [[0, '#fff', 0], [0.5, '#fff', 0.35], [0.55, '#fff', 0], [1, '#fff', 0]], 0, 0, 1, 1));

const CLOTHING = {
    name: 'Clothing',
    none: 120,
    variants: [
        v('Red Hoodie', 90, hoodie('#ff7a85', '#e94b5a', '#b22a3c')),
        v('Navy Hoodie', 80, hoodie('#5a6db0', '#2f3e75', '#1c2550')),
        v('Striped Tee', 80, stripedTee()),
        v('White Tee', 90, whiteTee()),
        v('Tuxedo', 40, tuxedo()),
        v('Kimono', 40, kimono()),
        v('Space Suit', 30, spaceSuit()),
        v('Knight Armor', 20, armor()),
        v('Varsity Jacket', 60, varsity()),
        v('Leather Jacket', 50, leather()),
        v('Overalls', 50, overalls()),
        v('Turtleneck', 60, turtleneck()),
        v('Gold Robe', 6, goldRobe()),
    ],
};

// ---------- Значок на грудях ----------

const PIN = [322, 452];
const pin = (body, defs = '') => svg(`<g transform="translate(${PIN[0]} ${PIN[1]})" filter="url(#ds)">${body}</g>`,
    '<filter id="ds"><feDropShadow dx="0" dy="3" stdDeviation="2" flood-color="#1a1030" flood-opacity=".35"/></filter>' + defs);
const enamel = (d, color) => `<path d="${d}" fill="${color}" stroke="${OL}" stroke-width="3.5" stroke-linejoin="round"/>` +
    `<path d="${d}" fill="none" stroke="#ffe08a" stroke-width="1.5" transform="scale(.82)" opacity=".9"/>`;

const BADGE = {
    name: 'Badge',
    none: 520,
    variants: [
        v('Star Pin', 60, pin(enamel(star(0, 0, 24, 10), '#ffd23f'))),
        v('Heart Pin', 60, pin(enamel(heart(0, 2, 20), '#ff4d6d'))),
        v('Lightning Pin', 50, pin(enamel('M6 -26 L-14 4 L-2 4 L-8 26 L14 -6 L2 -6 Z', '#ffd23f'))),
        v('Smiley Pin', 50, pin(`<circle r="22" fill="#ffe14d" stroke="${OL}" stroke-width="3.5"/>` +
            `<circle cx="-7" cy="-5" r="3" fill="${OL}"/><circle cx="7" cy="-5" r="3" fill="${OL}"/>` +
            `<path d="M-11 5 Q0 16 11 5" fill="none" stroke="${OL}" stroke-width="3" stroke-linecap="round"/>`)),
        v('Planet Pin', 30, pin(`<circle r="15" fill="#8f7bff" stroke="${OL}" stroke-width="3"/>` +
            `<ellipse rx="26" ry="8" fill="none" stroke="${OL}" stroke-width="6" transform="rotate(-20)"/>` +
            '<ellipse rx="26" ry="8" fill="none" stroke="#ffc2e2" stroke-width="3" transform="rotate(-20)"/>')),
        v('Flame Pin', 30, pin(enamel('M0 -26 C10 -12 20 -4 16 10 C13 22 4 26 0 26 C-4 26 -13 22 -16 10 C-19 -2 -8 -8 -6 -18 C-2 -10 2 -12 0 -26 Z', '#ff7a1a'))),
        v('Diamond Pin', 15, pin(enamel('M-20 -8 L-10 -20 L10 -20 L20 -8 L0 22 Z', '#7fe8ff') +
            `<path d="M-20 -8 L20 -8 M-10 -20 L-4 -8 L0 22 M10 -20 L4 -8" stroke="${OL}" stroke-width="2" fill="none"/>`, '')),
        v('GM Pin', 30, pin(`<rect x="-26" y="-16" width="52" height="32" rx="16" fill="#1d1b26" stroke="${OL}" stroke-width="3"/>` +
            '<text x="0" y="8" text-anchor="middle" font-family="Arial Black, Arial, sans-serif" font-weight="900" font-size="20" fill="#7cffb2">gm</text>')),
        v('Peace Pin', 25, pin(`<circle r="21" fill="#4fd1c5" stroke="${OL}" stroke-width="3.5"/>` +
            '<g stroke="#fff" stroke-width="3.5" fill="none" stroke-linecap="round"><circle r="13"/><path d="M0 -13 L0 13 M0 2 L-9 10 M0 2 L9 10"/></g>')),
    ],
};

// ---------- Намисто (позаду голови) ----------

const neck = (body, defs = '') => svg(`<g mask="url(#nohead)">${body}</g>`, behindHead('nohead') + defs);
const NECK_CURVE = 'M162 372 C190 420 310 420 338 372';

function beads(count, r, color, stroke = OL) {
    let out = '';
    for (let i = 0; i <= count; i++) {
        const t = i / count;
        const x = (1 - t) ** 3 * 162 + 3 * (1 - t) ** 2 * t * 190 + 3 * (1 - t) * t ** 2 * 310 + t ** 3 * 338;
        const y = (1 - t) ** 3 * 372 + 3 * (1 - t) ** 2 * t * 420 + 3 * (1 - t) * t ** 2 * 420 + t ** 3 * 372;
        out += `<circle cx="${x.toFixed(1)}" cy="${y.toFixed(1)}" r="${r}" fill="${color}" stroke="${stroke}" stroke-width="2"/>`;
    }
    return out;
}

const NECKLACE = {
    name: 'Necklace',
    none: 450,
    variants: [
        v('Gold Chain', 70, neck(`<path d="${NECK_CURVE}" fill="none" stroke="${OL}" stroke-width="11"/>` +
            `<path d="${NECK_CURVE}" fill="none" stroke="url(#g)" stroke-width="7" stroke-dasharray="9 3"/>`,
        lin('g', [[0, '#fff1a8'], [0.5, '#e2a93a'], [1, '#ffe58f']], 0, 0, 1, 0))),
        v('Pearls', 60, neck(beads(16, 7, 'url(#p)'), rad('p', [[0, '#ffffff'], [0.7, '#f3ecf7'], [1, '#c9bcd6']], 0.35, 0.3, 0.7))),
        v('Bow Tie', 70, neck(`<g transform="translate(250 392)">` +
            `<path d="M0 0 L-38 -20 C-46 -6 -46 6 -38 20 Z M0 0 L38 -20 C46 -6 46 6 38 20 Z" fill="url(#b)" stroke="${OL}" stroke-width="4" stroke-linejoin="round"/>` +
            `<rect x="-9" y="-11" width="18" height="22" rx="5" fill="#c2185b" stroke="${OL}" stroke-width="4"/></g>`,
        lin('b', [[0, '#ff6fa8'], [1, '#c2185b']]))),
        v('Red Scarf', 50, neck(`<path d="M150 370 C180 425 320 425 350 370 L355 395 C320 450 180 450 145 395 Z" fill="url(#s)" stroke="${OL}" stroke-width="4"/>` +
            `<path d="M300 415 L330 495 L300 500 L280 425 Z" fill="url(#s)" stroke="${OL}" stroke-width="4"/>` +
            '<path d="M160 395 C200 430 300 430 340 395" fill="none" stroke="#fff" stroke-width="3" stroke-dasharray="10 8" opacity=".6"/>',
        lin('s', [[0, '#ff5a5a'], [1, '#b3202c']]))),
        v('Medal', 30, neck(`<path d="M195 372 L245 440 L255 440 L305 372" fill="none" stroke="#2f6fd6" stroke-width="16"/>` +
            '<path d="M195 372 L245 440 L255 440 L305 372" fill="none" stroke="#fff" stroke-width="4" stroke-dasharray="1 0"/>' +
            `<circle cx="250" cy="455" r="22" fill="url(#m)" stroke="${OL}" stroke-width="4"/><path d="${star(250, 455, 11, 5)}" fill="#fff6c4"/>`,
        rad('m', [[0, '#fff1a8'], [0.6, '#f2c14e'], [1, '#b97d1f']], 0.4, 0.35, 0.75))),
        v('Crystal Pendant', 25, neck(`<path d="${NECK_CURVE}" fill="none" stroke="#c8ccd2" stroke-width="3"/>` +
            `<g filter="url(#gl)"><path d="M250 410 L264 428 L250 462 L236 428 Z" fill="url(#c)" stroke="${OL}" stroke-width="3"/></g>`,
        lin('c', [[0, '#e6fbff'], [1, '#6fd3ff']]) + glow('gl', 5, '#7fe8ff'))),
        v('Spiked Collar', 30, neck(`<path d="M160 372 C190 412 310 412 340 372" fill="none" stroke="${OL}" stroke-width="20"/>` +
            `<path d="M160 372 C190 412 310 412 340 372" fill="none" stroke="#24222b" stroke-width="14"/>` +
            [[190, 397], [220, 406], [250, 409], [280, 406], [310, 397]].map(([x, y]) =>
                `<path d="M${x - 6} ${y + 4} L${x} ${y + 22} L${x + 6} ${y + 4} Z" fill="#dfe5ea" stroke="${OL}" stroke-width="2"/>`).join(''))),
        v('Bandana', 40, neck(`<path d="M165 372 C200 400 300 400 335 372 L250 470 Z" fill="url(#b)" stroke="${OL}" stroke-width="4" stroke-linejoin="round"/>` +
            '<g fill="#fff" opacity=".85"><circle cx="250" cy="420" r="4"/><circle cx="225" cy="400" r="3"/><circle cx="275" cy="400" r="3"/><circle cx="250" cy="447" r="3"/></g>',
        lin('b', [[0, '#ff5a5a'], [1, '#c41f2f']]))),
    ],
};

// ---------- Сережки (біля країв голови) ----------

const EAR = [121, 292];
const ears = (left, right = left, defs = '') => svg(
    `<g transform="translate(${EAR[0]} ${EAR[1]})">${left}</g><g transform="translate(${500 - EAR[0]} ${EAR[1]}) scale(-1 1)">${right}</g>`, defs);

const EARRINGS = {
    name: 'Earrings',
    none: 650,
    variants: [
        v('Gold Hoops', 80, ears(`<circle cx="-4" cy="22" r="18" fill="none" stroke="${OL}" stroke-width="8"/><circle cx="-4" cy="22" r="18" fill="none" stroke="url(#g)" stroke-width="4.5"/>`,
            undefined, lin('g', [[0, '#fff1a8'], [1, '#d4972a']]))),
        v('Diamond Studs', 60, ears(`<g filter="url(#gl)"><path d="M-4 -8 L6 2 L-4 12 L-14 2 Z" fill="#e6fbff" stroke="${OL}" stroke-width="3"/></g>`,
            undefined, glow('gl', 3, '#bff3ff'))),
        v('Pearl Drops', 50, ears(`<line x1="-4" y1="2" x2="-4" y2="20" stroke="#d4972a" stroke-width="3"/><circle cx="-4" cy="30" r="10" fill="url(#p)" stroke="${OL}" stroke-width="3"/>`,
            undefined, rad('p', [[0, '#ffffff'], [1, '#d8cde4']], 0.35, 0.3, 0.7))),
        v('Silver Cross', 30, ears(`<line x1="-4" y1="2" x2="-4" y2="14" stroke="#c8ccd2" stroke-width="3"/><path d="M-7 14 h6 v8 h8 v6 h-8 v14 h-6 v-14 h-8 v-6 h8 Z" fill="#e6ebf0" stroke="${OL}" stroke-width="2.5"/>`)),
        v('Feathers', 40, ears(`<line x1="-4" y1="2" x2="-4" y2="14" stroke="#d4972a" stroke-width="3"/>` +
            `<path d="M-4 14 C-18 30 -14 56 -4 70 C6 56 10 30 -4 14 Z" fill="url(#f)" stroke="${OL}" stroke-width="2.5"/>` +
            '<path d="M-4 18 L-4 68" stroke="#fff" stroke-width="1.5" opacity=".7"/>', undefined,
        lin('f', [[0, '#4fd1c5'], [0.5, '#8f7bff'], [1, '#ff6fa8']]))),
        v('Moon & Star', 25, ears(
            `<line x1="-4" y1="2" x2="-4" y2="12" stroke="#d4972a" stroke-width="3"/><path d="M0 12 A13 13 0 1 0 6 34 A10 10 0 1 1 0 12 Z" fill="#ffe58f" stroke="${OL}" stroke-width="2.5"/>`,
            `<line x1="-4" y1="2" x2="-4" y2="12" stroke="#d4972a" stroke-width="3"/><path d="${star(-4, 26, 13, 5.5)}" fill="#ffe58f" stroke="${OL}" stroke-width="2.5"/>`)),
        v('Cherries', 15, ears(`<path d="M-4 2 C-4 14 -10 20 -14 28 M-4 2 C-2 14 4 20 6 28" fill="none" stroke="#3e8e41" stroke-width="3"/>` +
            `<circle cx="-14" cy="34" r="8" fill="#e8243c" stroke="${OL}" stroke-width="2.5"/><circle cx="6" cy="34" r="8" fill="#e8243c" stroke="${OL}" stroke-width="2.5"/>` +
            '<circle cx="-16" cy="31" r="2.5" fill="#fff" opacity=".8"/><circle cx="4" cy="31" r="2.5" fill="#fff" opacity=".8"/>')),
    ],
};

module.exports = { SKIN, SKIN_PATTERN, CLOTHING, BADGE, NECKLACE, EARRINGS, HEAD };
