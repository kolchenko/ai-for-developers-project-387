import { expect, test } from '@playwright/test';
import {
  API_BASE,
  adminLogin,
  createEventType,
  getFreeSlots,
  resetData,
  uniqueName,
} from './helpers';

test.describe('Администратор', () => {
  test.beforeEach(async ({ request }) => {
    await resetData(request);
  });

  test('A1: доступ к админке требует входа; неверный и верный логин, выход', async ({ page }) => {
    await page.goto('/admin');
    await expect(page).toHaveURL(/\/admin\/login$/);

    await page.getByLabel('Логин').fill('admin');
    await page.getByLabel('Пароль').fill('wrong');
    await page.getByRole('button', { name: 'Войти' }).click();
    await expect(page.getByText('Неверный логин или пароль')).toBeVisible();

    await page.getByLabel('Пароль').fill('admin');
    await page.getByRole('button', { name: 'Войти' }).click();
    await expect(page).toHaveURL(/\/admin$/);
    await expect(page.getByRole('heading', { name: 'Предстоящие встречи' })).toBeVisible();

    await page.getByRole('button', { name: 'Выйти' }).click();
    await expect(page).toHaveURL(/\/admin\/login$/);
  });

  test('A2+A3: бронь гостя видна админу и отменяется через подтверждение', async ({
    page,
    request,
  }) => {
    const et = await createEventType(request, { name: 'Консультация' });

    const slots = await getFreeSlots(request, et.id);
    expect(slots.length).toBeGreaterThan(0);

    const res = await request.post(`${API_BASE}/bookings`, {
      data: {
        eventTypeId: et.id,
        startsAt: slots[0].startsAt,
        guestName: 'Иван Петров',
        guestEmail: 'ivan@example.com',
      },
    });
    expect(res.status()).toBe(201);

    await adminLogin(page);
    await expect(page.getByRole('heading', { name: 'Предстоящие встречи' })).toBeVisible();

    await expect(page.getByText('Иван Петров')).toBeVisible();
    await expect(page.getByText('ivan@example.com')).toBeVisible();
    await expect(page.getByText('Консультация')).toBeVisible();

    await page.getByRole('button', { name: 'Отменить' }).first().click();
    const dialog = page.getByRole('dialog');
    await expect(dialog).toBeVisible();
    await expect(dialog).toContainText('Иван Петров');
    await dialog.getByRole('button', { name: 'Отменить' }).click();

    await expect(page.getByText('Предстоящих встреч нет.')).toBeVisible();
  });

  test('A4: админ создаёт тип события через форму — он появляется у гостя', async ({ page }) => {
    await adminLogin(page);
    await page.getByRole('link', { name: 'Типы событий' }).click();
    await expect(page.getByRole('heading', { name: 'Типы событий' })).toBeVisible();

    const name = uniqueName('Вебинар');
    await page.getByRole('button', { name: 'Новый тип события' }).click();
    const dialog = page.getByRole('dialog');
    await expect(dialog).toBeVisible();
    await dialog.getByLabel('Название').fill(name);
    await dialog.getByLabel('Описание').fill('Онлайн-встреча');
    await dialog.getByRole('button', { name: 'Создать' }).click();

    await expect(
      page.locator('table tbody tr').filter({ hasText: name }),
    ).toContainText('Онлайн-встреча');

    await page.goto('/');
    const card = page.locator('.mantine-Card-root').filter({ hasText: name });
    await expect(card).toBeVisible();
    await expect(card).toContainText('Онлайн-встреча');
  });

  test('A5: админ редактирует тип события', async ({ page, request }) => {
    const et = await createEventType(request, { name: 'Старое название' });
    const newName = uniqueName('Новое название');

    await adminLogin(page);
    await page.getByRole('link', { name: 'Типы событий' }).click();

    const row = page.locator('table tbody tr').filter({ hasText: 'Старое название' });
    await expect(row).toBeVisible();
    await row.getByRole('button').first().click();

    const dialog = page.getByRole('dialog');
    await expect(dialog).toBeVisible();
    await dialog.getByLabel('Название').fill(newName);
    await dialog.getByRole('button', { name: 'Сохранить' }).click();

    await expect(
      page.locator('table tbody tr').filter({ hasText: newName }),
    ).toBeVisible();
    await expect(
      page.locator('table tbody tr').filter({ hasText: 'Старое название' }),
    ).toHaveCount(0);
  });

  test('A6: админ удаляет тип события — брони удаляются каскадно', async ({ page, request }) => {
    const et = await createEventType(request, { name: 'Временный тип' });

    const slots = await getFreeSlots(request, et.id);
    const res = await request.post(`${API_BASE}/bookings`, {
      data: {
        eventTypeId: et.id,
        startsAt: slots[0].startsAt,
        guestName: 'Иван Петров',
        guestEmail: 'ivan@example.com',
      },
    });
    expect(res.status()).toBe(201);

    await adminLogin(page);
    await page.getByRole('link', { name: 'Типы событий' }).click();

    const row = page.locator('table tbody tr').filter({ hasText: 'Временный тип' });
    await expect(row).toBeVisible();
    await row.getByRole('button').nth(1).click();

    const dialog = page.getByRole('dialog');
    await expect(dialog).toBeVisible();
    await expect(dialog).toContainText('Бронирования типа «Временный тип» будут удалены');
    await dialog.getByRole('button', { name: 'Удалить' }).click();

    await expect(
      page.locator('table tbody tr').filter({ hasText: 'Временный тип' }),
    ).toHaveCount(0);

    const bookings = await (
      await request.get(`${API_BASE}/admin/bookings`)
    ).json();
    expect(bookings).toEqual([]);
  });
});
