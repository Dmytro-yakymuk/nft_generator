// Whitelist: додавання адреси через форму, пошук і видалення.
import { Project, createProject, deleteProject, expect, login, test } from './fixtures';

const ADDRESS = '0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC';

test('адреса додається до whitelist, шукається за іменем і видаляється', async ({ page }) => {
    let project: Project | undefined;
    try {
        await login(page);
        project = await createProject(page.request, 3, false);
        await page.goto(`/whitelist/${project.id}`);

        await page.getByRole('button', { name: 'Add to Whitelist' }).click();
        await page.getByPlaceholder('Wallet name').fill('E2E buyer');
        await page.locator('#wallet-address').fill(ADDRESS);
        await page.getByRole('button', { name: 'Create' }).click();

        const rows = page.locator('.dashboard__white-list__website');
        await expect(rows).toHaveCount(1);
        await expect(rows.first()).toContainText('E2E buyer');
        await expect(rows.first()).toContainText(ADDRESS.toLowerCase());

        await page.getByPlaceholder('Search...').fill('nobody');
        await expect(rows).toHaveCount(0);
        await page.getByPlaceholder('Search...').fill('E2E');
        await expect(rows).toHaveCount(1);

        await rows.first().locator('.dashboard__white-list__website__close').click();
        await expect(rows).toHaveCount(0);
    } finally {
        await deleteProject(page.request, project);
    }
});
