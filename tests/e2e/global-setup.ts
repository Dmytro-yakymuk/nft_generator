// Збирає init-скрипт з імітацією MetaMask (scripts/build-metamask-mock.js) у tests/e2e/.mock/.
// За замовчуванням - тестовий акаунт №1 Anvil, мережа Anvil (0x7a69) і вузол http://localhost:8545.
import { execFileSync } from 'child_process';
import path from 'path';

export const MOCK_PATH = path.join(__dirname, '.mock', 'metamask-mock.js');

export default function globalSetup(): void {
    execFileSync(process.execPath, [path.join(__dirname, '..', '..', 'scripts', 'build-metamask-mock.js'), MOCK_PATH], {
        stdio: 'inherit',
        env: process.env,
    });
}
