import { defineConfig } from '@playwright/test';

const FRONTEND_URL = 'http://localhost:5173';
const API_URL = 'http://127.0.0.1:4010';
const DB_PATH = `/tmp/callcalendar-e2e-${process.pid}.db`;

export default defineConfig({
  testDir: './e2e',
  fullyParallel: false,
  workers: 1,
  retries: process.env.CI ? 2 : 0,
  reporter: process.env.CI ? 'github' : [['list']],
  use: {
    baseURL: FRONTEND_URL,
    trace: 'on-first-retry',
  },
  webServer: [
    {
      command: 'go run ./cmd/server',
      cwd: './backend',
      env: {
        ADDR: API_URL.replace('http://', ''),
        DB_PATH,
      },
      url: `${API_URL}/event-types`,
      reuseExistingServer: false,
      timeout: 120_000,
    },
    {
      command: 'npm run dev',
      url: FRONTEND_URL,
      reuseExistingServer: false,
      timeout: 120_000,
    },
  ],
});
