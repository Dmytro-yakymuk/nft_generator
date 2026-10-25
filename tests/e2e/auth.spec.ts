// Вхід підписом гаманця (імітація MetaMask) і вихід з облікового запису.
import { expect, login, test } from './fixtures';

test.describe('Вхід через MetaMask', () => {
    test('вхід з головної сторінки веде на сторінку колекцій', async ({ page }) => {
        await login(page);
        await expect(page.getByRole('button', { name: 'Connect', exact: true })).toHaveCount(0);
        expect((await page.request.get('/api/v0/accounts')).ok()).toBeTruthy();
    });

    test('налаштування облікового запису зберігаються й підвантажуються у формі', async ({ page }) => {
        await login(page);
        await page.goto('/account-settings');
        const name = `E2E ${Date.now()}`;
        await page.locator('#name').fill(name);
        await page.locator('#email').fill('e2e@example.com');
        await page.locator('#subscription').check();
        await page.getByText('Save', { exact: true }).click();
        await expect.poll(async () => (await (await page.request.get('/api/v0/accounts')).json()).name).toBe(name);

        await page.reload();
        await expect(page.locator('#name')).toHaveValue(name);
        await expect(page.locator('#email')).toHaveValue('e2e@example.com');
        await expect(page.locator('#subscription')).toBeChecked();
        // одноразовий nonce входу й хеш пароля сервер не віддає.
        const account = await (await page.request.get('/api/v0/accounts')).json();
        expect(account).not.toHaveProperty('nonce');
        expect(account).not.toHaveProperty('passwordHash');
    });

    test('посилання навігаційної панелі ведуть до розділів головної сторінки, а останній розділ - до студії', async ({ page }) => {
        await login(page);
        // колишні сторінки входу за email (без сервера за ними) перенаправляють на головну.
        await page.goto('/sign-in');
        await expect(page).toHaveURL(/\/$/);
        await page.locator('.navbar__items').getByRole('link', { name: 'Creators' }).click();
        await expect(page).toHaveURL(/\/#how-it-works$/);
        await expect(page.locator('#how-it-works')).toBeInViewport();
        await page.locator('.navbar__items').getByRole('link', { name: 'Help' }).click();
        const start = page.locator('#contact-us');
        await expect(start).toBeInViewport();
        await start.getByRole('link', { name: 'Open Studio' }).click();
        await expect(page).toHaveURL(/\/collections$/);
    });

    test('без входу сторінка просить підключити гаманець, а вихід закриває сесію', async ({ page }) => {
        await page.goto('/collections');
        const modal = page.locator('.metamask-modal');
        await expect(modal).toBeVisible();
        await modal.getByRole('button', { name: 'Connect Wallet' }).click();
        await expect(modal).toHaveCount(0);
        expect((await page.request.get('/api/v0/accounts')).ok()).toBeTruthy();

        await page.goto('/account-settings');
        await page.getByText('Log out', { exact: true }).click();
        await expect(page).toHaveURL(/\/$/);
        expect((await page.request.get('/api/v0/accounts')).status()).toBe(401);
        await page.goto('/collections');
        await expect(page.locator('.metamask-modal')).toBeVisible();
    });
});
