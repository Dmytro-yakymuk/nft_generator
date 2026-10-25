// Контракт: розгортання з інтерфейсу і переведення базового URI вже розгорнутого контракту на колекцію,
// опубліковану в IPFS. Транзакції виконуються лише тоді, коли платформа налаштована на мережу імітації
// MetaMask (за замовчуванням Anvil, make chain); інакше сценарії пропускаються з поясненням.
import { Page } from '@playwright/test';

import { Project, createProject, deleteProject, expect, login, test } from './fixtures';

const GATEWAY = process.env.IPFS_GATEWAY || 'http://localhost:8080';

async function generatedProject(page: Page): Promise<Project> {
    const project = await createProject(page.request, 3);
    expect((await page.request.post(`/api/v0/collections/${project.id}/generate-images`)).ok()).toBeTruthy();
    await expect.poll(async () => (await page.request.get(`/api/v0/collections/${project.id}/generated-images-count`)).json(),
        { timeout: 60_000 }).toBe(3);
    return project;
}

async function skipUnlessMockNetwork(page: Page): Promise<void> {
    const appChain: string = await page.evaluate(() => (window as any).APP_CONFIG?.chain?.chainId || '');
    const mockChain: string = await page.evaluate(() => (window as any).ethereum.chainId);
    test.skip(appChain.toLowerCase() !== mockChain.toLowerCase(), `платформа налаштована на мережу ${appChain}, імітація - на ${mockChain}`);
}

/** Умови продажу, які власник вводить на сторінці контракту. */
interface SaleTerms {
    price: string;
    maxPerMint: string;
}

/** Розгортає контракт колекції кнопкою Publish і повертає його адресу, збережену в проєкті. */
async function deployFromUI(page: Page, project: Project, terms?: SaleTerms): Promise<string> {
    await page.goto(`/create-contract/${project.id}`);
    await expect(page.getByRole('button', { name: 'Publish' })).toBeVisible();
    expect((await page.request.get(`/api/v0/contracts/data/${project.id}`)).ok()).toBeTruthy();
    await skipUnlessMockNetwork(page);

    // кількість NFT береться з колекції й не редагується.
    await expect(page.locator('#nft-count')).toHaveValue(String(project.nftCount));
    await page.locator('input').first().fill('E2E Collection');
    if (terms) {
        await page.locator('#mint').fill(terms.price);
        await page.locator('#max-per-mint').fill(terms.maxPerMint);
    }
    await page.getByRole('button', { name: 'Publish' }).click();
    let address = '';
    await expect.poll(async () => {
        address = (await (await page.request.get(`/api/v0/collections/${project.id}`)).json()).contractAddress;
        return address;
    }, { timeout: 60_000 }).toMatch(/^0x(?!0{40})[0-9a-fA-F]{40}$/);
    return address;
}

test('контракт колекції розгортається з інтерфейсу з умовами продажу', async ({ page }) => {
    let project: Project | undefined;
    try {
        await login(page);
        project = await generatedProject(page);
        const contractAddress = await deployFromUI(page, project, { price: '0.002', maxPerMint: '2' });

        // умови продажу записано в контракт: максимум - розмір колекції.
        const terms = await page.evaluate(async (address) => {
            const { ethers, APP_CONFIG } = window as any;
            const contract = new ethers.Contract(address, [
                'function name() view returns (string)',
                'function MAX_SUPPLY() view returns (uint256)',
                'function PRICE() view returns (uint256)',
                'function MAX_PER_MINT() view returns (uint256)',
            ], new ethers.providers.JsonRpcProvider(APP_CONFIG.chain.rpcUrl));
            return {
                name: await contract.name(),
                maxSupply: (await contract.MAX_SUPPLY()).toNumber(),
                price: ethers.utils.formatEther(await contract.PRICE()),
                maxPerMint: (await contract.MAX_PER_MINT()).toNumber(),
            };
        }, contractAddress);
        expect(terms).toEqual({ name: 'E2E Collection', maxSupply: project.nftCount, price: '0.002', maxPerMint: 2 });
    } finally {
        await deleteProject(page.request, project);
    }
});

test('розгорнутий контракт переводиться на колекцію в IPFS', async ({ page }) => {
    let project: Project | undefined;
    try {
        await login(page);
        project = await generatedProject(page);
        // контракт розгорнуто до публікації: його базовий URI - адреса платформи.
        const contractAddress = await deployFromUI(page, project);

        // власник відкриває продаж і мінтить токен 0 - tokenURI працює лише для наявного токена.
        await page.evaluate(async (address) => {
            const { ethers, ethereum } = window as any;
            const contract = new ethers.Contract(address,
                ['function setSaleState(bool newState)', 'function mint(uint numberOfTokens) payable'],
                new ethers.providers.Web3Provider(ethereum).getSigner());
            await (await contract.setSaleState(true)).wait();
            await (await contract.mint(1, { value: ethers.utils.parseEther('0.001') })).wait();
        }, contractAddress);

        const published = await (await page.request.post(`/api/v0/collections/${project.id}/ipfs`)).json();
        await page.goto(`/generate-metadata/${project.id}`);
        await expect(page.locator('.contract-base-uri--outdated')).toContainText('tokenURI(0) = http');
        await page.getByRole('button', { name: 'Point the contract to IPFS' }).click();
        await expect(page.locator('.contract-base-uri--ipfs')).toContainText(`ipfs://${published.metadataCid}/0`);

        // tokenURI у блокчейні вказує на IPFS, а шлюз віддає метадані й зображення токена.
        const tokenURI: string = await page.evaluate(async (address) => {
            const { ethers, APP_CONFIG } = window as any;
            const contract = new ethers.Contract(address, ['function tokenURI(uint256) view returns (string)'],
                new ethers.providers.JsonRpcProvider(APP_CONFIG.chain.rpcUrl));
            return contract.tokenURI(0);
        }, contractAddress);
        expect(tokenURI).toBe(`ipfs://${published.metadataCid}/0`);

        const metadata = await (await page.request.get(`${GATEWAY}/ipfs/${published.metadataCid}/0`)).json();
        expect(metadata.name).toMatch(/#0$/);
        expect(metadata.image).toBe(`ipfs://${published.imagesCid}/0.png`);
        const image = await page.request.get(`${GATEWAY}/ipfs/${published.imagesCid}/0.png`);
        expect(image.ok()).toBeTruthy();
        expect((await image.body()).subarray(1, 4).toString()).toBe('PNG');
    } finally {
        await deleteProject(page.request, project);
    }
});
