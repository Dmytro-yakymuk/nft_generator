// Експеримент E7 (п. 3.2.10 роботи): завантаження головної сторінки клієнтського застосунку в Chromium
// (Playwright) з вимкненим кешем - без обмеження мережі та за профілем "повільний 4G" Lighthouse
// (затримка 150 мс, 1.6 Мбіт/с донизу, 750 кбіт/с догори). Міряються передані (стиснені) й розпаковані
// байти ресурсів, First Contentful Paint, Largest Contentful Paint і подія load; 5 повторів.
//
//   LABEL=before node experiments/web-vitals.js     # дописує рядки в experiments/results/web.csv
//   BASE=... REPEATS=... OUT=... node experiments/web-vitals.js

const fs = require('fs');
const path = require('path');
const { chromium } = require('@playwright/test');

const BASE = (process.env.BASE || 'http://localhost:9012').replace(/\/$/, '');
const LABEL = process.env.LABEL || 'current';
const REPEATS = Number(process.env.REPEATS || 5);
const OUT = process.env.OUT || path.join(__dirname, 'results', 'web.csv');
const PROFILES = [
    { name: 'unthrottled', conditions: { offline: false, latency: 0, downloadThroughput: -1, uploadThroughput: -1 } },
    { name: 'slow-4g', conditions: { offline: false, latency: 150, downloadThroughput: 1.6e6 / 8, uploadThroughput: 750e3 / 8 } },
];

// LCP уточнюється, поки сторінка завантажується, тому записується спостерігачем з самого початку.
const OBSERVER = `
    window.__lcp = 0;
    new PerformanceObserver((list) => {
        for (const entry of list.getEntries()) window.__lcp = entry.renderTime || entry.loadTime || entry.startTime;
    }).observe({ type: 'largest-contentful-paint', buffered: true });
`;

async function measure(browser, profile) {
    const context = await browser.newContext({ viewport: { width: 1440, height: 900 } });
    const page = await context.newPage();
    await page.addInitScript(OBSERVER);
    const cdp = await context.newCDPSession(page);
    await cdp.send('Network.enable');
    await cdp.send('Network.setCacheDisabled', { cacheDisabled: true });
    await cdp.send('Network.emulateNetworkConditions', profile.conditions);
    let transferred = 0;
    let decoded = 0;
    cdp.on('Network.loadingFinished', (event) => { transferred += event.encodedDataLength; });
    cdp.on('Network.dataReceived', (event) => { decoded += event.dataLength; });

    await page.goto(`${BASE}/`, { waitUntil: 'load', timeout: 120_000 });
    await page.waitForTimeout(1500);
    const timings = await page.evaluate(() => ({
        fcp: (performance.getEntriesByName('first-contentful-paint')[0] || {}).startTime || 0,
        lcp: window.__lcp,
        load: performance.getEntriesByType('navigation')[0].loadEventEnd,
    }));
    await context.close();
    return { transferred, decoded, ...timings };
}

(async () => {
    const browser = await chromium.launch();
    if (!fs.existsSync(OUT)) fs.writeFileSync(OUT, 'label,profile,repeat,transferred_bytes,decoded_bytes,fcp_ms,lcp_ms,load_ms\n');
    for (const profile of PROFILES) {
        for (let repeat = 1; repeat <= REPEATS; repeat++) {
            const r = await measure(browser, profile);
            fs.appendFileSync(OUT, `${[LABEL, profile.name, repeat, r.transferred, r.decoded, r.fcp.toFixed(0), r.lcp.toFixed(0), r.load.toFixed(0)].join(',')}\n`);
            console.log(`${LABEL} ${profile.name} #${repeat}: ${(r.transferred / 2 ** 20).toFixed(2)} МіБ передано (${(r.decoded / 2 ** 20).toFixed(2)} розпаковано), FCP ${r.fcp.toFixed(0)} мс, LCP ${r.lcp.toFixed(0)} мс, load ${r.load.toFixed(0)} мс`);
        }
    }
    await browser.close();
})().catch((err) => {
    console.error(err);
    process.exit(1);
});
