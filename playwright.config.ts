// Наскрізні тести Playwright проти запущеної платформи (make up): make e2e або
//   BASE=http://localhost:9012 node node_modules/@playwright/test/cli.js test
// Імітацію MetaMask збирає tests/e2e/global-setup.ts (MOCK_CHAIN_ID / MOCK_RPC_URL / MOCK_PRIVATE_KEY).
import { defineConfig } from '@playwright/test';

export default defineConfig({
    testDir: './tests/e2e',
    globalSetup: './tests/e2e/global-setup.ts',
    timeout: 180_000,
    expect: { timeout: 20_000 },
    // сценарії створюють проєкти одного гаманця, тому виконуються по черзі.
    workers: 1,
    retries: 0,
    reporter: [['list'], ['json', { outputFile: 'tests/e2e/results/report.json' }]],
    outputDir: 'tests/e2e/results/artifacts',
    use: {
        baseURL: process.env.BASE || 'http://localhost:9012',
        viewport: { width: 1440, height: 900 },
        headless: true,
        trace: 'retain-on-failure',
        screenshot: 'only-on-failure',
    },
});
