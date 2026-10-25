// Генератор: ваги варіантів і необов'язковий шар, попередній перегляд, генерація, метадані, публікація в IPFS.
import { Page } from '@playwright/test';

import { Project, createProject, deleteProject, expect, login, test } from './fixtures';

const NFT_COUNT = 6;

test.describe.serial('Генерація колекції з вагами', () => {
    let page: Page;
    let project: Project;

    test.beforeAll(async ({ browser }) => {
        const context = await browser.newContext();
        page = await context.newPage();
        await page.addInitScript({ path: require('path').join(__dirname, '.mock', 'metamask-mock.js') });
        await login(page);
        project = await createProject(page.request, NFT_COUNT);
    });

    test.afterAll(async () => {
        await deleteProject(page.request, project);
        await page.context().close();
    });

    test('ваги, частки й необов\'язковий шар зберігаються', async () => {
        await page.goto(`/collection-generator/${project.id}`);
        await expect(page.locator('.generator__aside__layer, .generator__aside__layer-active')).toHaveCount(3);
        await page.locator('.generator__aside__layer, .generator__aside__layer-active').filter({ hasText: 'Hat' }).click();
        await expect(page.locator('.collection-content__heading')).toHaveText('Hat');

        const weights = page.locator('.collection-image__weight__input');
        await expect(weights).toHaveCount(2);
        await expect(page.locator('.collection-image__weight__percent').first()).toHaveText('50.0%');

        // 300 і 100: частки 75% і 25%; потім необов'язковий шар з вагою відсутності 100 - 60/20/20%.
        await weights.nth(0).fill('300');
        await expect(page.locator('.collection-image__weight__percent').nth(0)).toHaveText('75.0%');
        await page.getByLabel('Optional layer').check();
        await expect(page.locator('.collection-image__none')).toBeVisible();
        await expect(page.locator('.collection-image__weight__percent').nth(1)).toHaveText('60.0%');
        await expect(page.locator('.rarity').nth(1)).toHaveText('Common');

        await page.getByRole('button', { name: 'Save Weights' }).click();
        await expect(page.getByRole('button', { name: 'Save Weights' })).toBeDisabled();
        const layers = await (await page.request.get(`/api/v0/collections/${project.id}/layers`)).json();
        const hat = layers.find((layer: { name: string }) => layer.name === 'Hat');
        expect(hat.noneWeight).toBe(100);
        expect(hat.weights).toEqual({ 'v1.png': 300, 'v2.png': 100 });
        // активний шар лишився тим самим.
        await expect(page.locator('.collection-content__heading')).toHaveText('Hat');
    });

    test('попередній перегляд складає вісім зображень за вагами', async () => {
        await page.getByRole('button', { name: 'Preview' }).click();
        const tiles = page.locator('.layers-preview__image');
        await expect(tiles).toHaveCount(8);
        await expect.poll(() => tiles.locator('img').evaluateAll((images) =>
            images.length > 0 && images.every((image) => (image as HTMLImageElement).complete && (image as HTMLImageElement).naturalWidth > 0))).toBe(true);
        await page.getByRole('button', { name: 'Close', exact: true }).click();
        await expect(page.locator('.layers-preview')).toHaveCount(0);
    });

    test('генерація створює колекцію, а метадані йдуть у порядку шарів', async () => {
        await page.getByRole('button', { name: 'Generate And Continue' }).click();
        await expect(page).toHaveURL(new RegExp(`/generate-metadata/${project.id}$`), { timeout: 60_000 });

        await expect(page.locator('.token-gallery__token')).toHaveCount(NFT_COUNT);
        for (let token = 0; token < NFT_COUNT; token++) {
            const metadata = await (await page.request.get(`/metadata/${project.id}/${token}`)).json();
            expect(metadata.attributes.map((attribute: { trait_type: string }) => attribute.trait_type)).toEqual(['Background', 'Body', 'Hat']);
            expect(metadata.dna).toMatch(/^[0-9]+$/);
            expect((await page.request.get(metadata.image)).ok()).toBeTruthy();
        }
        expect((await page.request.get(`/metadata/${project.id}/${NFT_COUNT}`)).status()).toBe(404);

        // картка колекції показує, що зображення згенеровано, а Generate Again веде назад у генератор.
        const collection = await (await page.request.get(`/api/v0/collections/${project.id}`)).json();
        expect(collection.areImagesGenerated).toBe(true);
        await page.getByRole('button', { name: 'Generate Again' }).click();
        await expect(page).toHaveURL(new RegExp(`/collection-generator/${project.id}$`));
        await page.goBack();
        await expect(page).toHaveURL(new RegExp(`/generate-metadata/${project.id}$`));
    });

    test('колекція публікується в IPFS', async () => {
        await page.getByRole('button', { name: 'Publish to IPFS' }).click();
        await expect(page.locator('.generate-metadata__ipfs')).toContainText('Contract base URI: ipfs://');
        const updated = await (await page.request.get(`/api/v0/collections/${project.id}`)).json();
        expect(updated.metadataCid).toMatch(/^bafy/);
        expect(updated.imagesCid).toMatch(/^bafy/);
    });
});
