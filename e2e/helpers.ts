import { expect } from '@playwright/test';
import type { APIRequestContext, Locator, Page } from '@playwright/test';
import dayjs from 'dayjs';
import type { EventType, Slot } from '../src/api/types';

export function formatSlotTime(iso: string): string {
  return dayjs(iso).format('HH:mm');
}

export const API_BASE = '/api';

export function uniqueName(prefix: string): string {
  return `${prefix} ${Date.now()}`;
}

export async function resetData(request: APIRequestContext): Promise<void> {
  const res = await request.get(`${API_BASE}/event-types`);
  const eventTypes = (await res.json()) as EventType[];
  for (const et of eventTypes) {
    await request.delete(`${API_BASE}/admin/event-types/${et.id}`);
  }
}

export interface EventTypeSeed {
  name?: string;
  description?: string;
  durationMinutes?: number;
  availableFrom?: string;
  availableTo?: string;
}

export async function createEventType(
  request: APIRequestContext,
  overrides: EventTypeSeed = {},
): Promise<EventType> {
  const payload = {
    name: overrides.name ?? uniqueName('Консультация'),
    description: overrides.description ?? 'Тестовый тип события',
    durationMinutes: overrides.durationMinutes ?? 30,
    availableFrom: overrides.availableFrom ?? '09:00:00',
    availableTo: overrides.availableTo ?? '18:00:00',
  };
  const res = await request.post(`${API_BASE}/admin/event-types`, { data: payload });
  expect(res.status(), `create event type failed: ${await res.text()}`).toBe(201);
  return (await res.json()) as EventType;
}

export async function getFreeSlots(request: APIRequestContext, eventTypeId: string): Promise<Slot[]> {
  const res = await request.get(`${API_BASE}/event-types/${eventTypeId}/slots`);
  expect(res.status()).toBe(200);
  return (await res.json()) as Slot[];
}

export function slotButtons(page: Page) {
  return page.getByRole('button', { name: /^\d{2}:\d{2}$/ });
}

export function slotButton(page: Page, time: string) {
  return page.getByRole('button', { name: time, exact: true });
}

export async function pickFirstSlot(page: Page): Promise<{ time: string; button: Locator }> {
  const button = slotButtons(page).first();
  await expect(button).toBeVisible();
  const time = (await button.textContent())?.trim() ?? '';
  expect(time).toMatch(/^\d{2}:\d{2}$/);
  return { time, button };
}

export async function bookSlotViaModal(
  page: Page,
  button: Locator,
  name: string,
  email: string,
): Promise<void> {
  await button.click();
  const dialog = page.getByRole('dialog');
  await expect(dialog).toBeVisible();
  await dialog.getByLabel('Ваше имя').fill(name);
  await dialog.getByLabel('Email').fill(email);
  await dialog.getByRole('button', { name: 'Подтвердить бронирование' }).click();
}

export async function adminLogin(page: Page): Promise<void> {
  await page.goto('/admin/login');
  await page.getByLabel('Логин').fill('admin');
  await page.getByLabel('Пароль').fill('admin');
  await page.getByRole('button', { name: 'Войти' }).click();
  await expect(page).toHaveURL(/\/admin$/);
}
