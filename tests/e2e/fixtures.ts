// Спільне для наскрізних тестів: сторінка з імітацією MetaMask, вхід і звернення до API від імені гаманця.
import { APIRequestContext, Page, expect, test as base } from '@playwright/test';
import path from 'path';
import zlib from 'zlib';

export const MOCK_PATH = path.join(__dirname, '.mock', 'metamask-mock.js');

export const test = base.extend({
    page: async ({ page }, use) => {
        await page.addInitScript({ path: MOCK_PATH });
        await use(page);
    },
});
export { expect };

/** Вхід кнопкою Connect на головній: імітація підписує nonce, застосунок переходить на сторінку колекцій. */
export async function login(page: Page): Promise<void> {
    await page.goto('/');
    await page.getByRole('button', { name: 'Connect', exact: true }).click();
    await expect(page).toHaveURL(/\/collections$/);
}

/** Мінімальний PNG-кодувальник (RGBA): тестові шари без графічних бібліотек. */
const CRC_TABLE = Array.from({ length: 256 }, (_, n) => {
    let c = n;
    for (let k = 0; k < 8; k++) c = c & 1 ? 0xedb88320 ^ (c >>> 1) : c >>> 1;
    return c >>> 0;
});

function chunk(type: string, data: Buffer): Buffer {
    const out = Buffer.alloc(12 + data.length);
    out.writeUInt32BE(data.length, 0);
    out.write(type, 4, 'ascii');
    data.copy(out, 8);
    let crc = 0xffffffff;
    for (const byte of out.subarray(4, 8 + data.length)) crc = CRC_TABLE[(crc ^ byte) & 0xff] ^ (crc >>> 8);
    out.writeUInt32BE((crc ^ 0xffffffff) >>> 0, 8 + data.length);
    return out;
}

export function png(size: number, pixel: (x: number, y: number) => number[]): Buffer {
    const raw = Buffer.alloc((size * 4 + 1) * size);
    for (let y = 0; y < size; y++) {
        for (let x = 0; x < size; x++) raw.set(pixel(x, y), y * (size * 4 + 1) + 1 + x * 4);
    }
    const header = Buffer.alloc(13);
    header.writeUInt32BE(size, 0);
    header.writeUInt32BE(size, 4);
    header.set([8, 6, 0, 0, 0], 8);
    return Buffer.concat([
        Buffer.from([0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a]),
        chunk('IHDR', header), chunk('IDAT', zlib.deflateSync(raw)), chunk('IEND', Buffer.alloc(0)),
    ]);
}

const disc = (cx: number, cy: number, r: number, color: number[]) => (x: number, y: number) =>
    ((x - cx) ** 2 + (y - cy) ** 2 <= r * r ? [...color, 255] : [0, 0, 0, 0]);

/** Тестові шари 500x500: фон (3 кольори), тіло (3), капелюх (2, необов'язковий у сценарії генератора). */
export const LAYERS: [string, Array<(x: number, y: number) => number[]>][] = [
    ['Background', [[230, 80, 80], [80, 160, 230], [90, 200, 120]].map((c) => () => [...c, 255])],
    ['Body', [[250, 220, 60], [160, 90, 220], [60, 60, 60]].map((c) => disc(250, 290, 150, c))],
    ['Hat', [[20, 20, 20], [240, 240, 240]].map((c) => disc(250, 120, 60, c))],
];

export interface Project {
    id: string;
    name: string;
    nftCount: number;
}

/** Створює проєкт з тестовими шарами через API від імені гаманця, що увійшов у браузері. */
export async function createProject(request: APIRequestContext, nftCount: number, withLayers = true): Promise<Project> {
    const name = `E2E ${Date.now()} ${Math.random().toString(36).slice(2, 6)}`;
    expect((await request.post('/api/v0/collections', { data: { name, nftCount } })).ok()).toBeTruthy();
    const projects: Project[] = await (await request.get('/api/v0/collections')).json();
    const project = projects.find((candidate) => candidate.name === name)!;
    if (!withLayers) return project;

    for (const [i, [layer]] of LAYERS.entries()) {
        expect((await request.post(`/api/v0/collections/${project.id}/layers`, { data: { name: layer, priority: i } })).ok()).toBeTruthy();
    }
    const created: { id: string; name: string }[] = await (await request.get(`/api/v0/collections/${project.id}/layers`)).json();
    for (const [layer, variants] of LAYERS) {
        const { id } = created.find((candidate) => candidate.name === layer)!;
        for (const [i, pixel] of variants.entries()) {
            const response = await request.post(`/api/v0/collections/${project.id}/layers/${id}/images`, {
                data: { image: png(500, pixel).toString('base64'), name: `v${i + 1}.png` },
            });
            expect(response.ok()).toBeTruthy();
        }
    }
    return project;
}

export async function deleteProject(request: APIRequestContext, project?: Project): Promise<void> {
    if (project) await request.delete(`/api/v0/collections/${project.id}`);
}
