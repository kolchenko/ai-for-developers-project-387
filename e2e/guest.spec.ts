import { expect, test } from '@playwright/test';
import {
  bookSlotViaModal,
  createEventType,
  formatSlotTime,
  getFreeSlots,
  pickFirstSlot,
  resetData,
  slotButton,
  slotButtons,
} from './helpers';

test.describe('Гость', () => {
  test.beforeEach(async ({ request }) => {
    await resetData(request);
  });

  test('G1: гость видит список видов брони и открывает слоты выбранного типа', async ({
    page,
    request,
  }) => {
    const et = await createEventType(request, {
      name: 'Консультация',
      description: 'Разбор вашего проекта',
      durationMinutes: 30,
    });
    await createEventType(request, {
      name: 'Разбор кода',
      description: 'Код-ревью',
      durationMinutes: 45,
      availableFrom: '10:00:00',
      availableTo: '17:00:00',
    });

    await page.goto('/');
    await expect(page.getByRole('heading', { name: 'Виды брони' })).toBeVisible();

    const consultCard = page.getByTestId(`event-type-card-${et.id}`);
    await expect(consultCard).toBeVisible();
    await expect(consultCard).toContainText('Разбор вашего проекта');
    await expect(consultCard).toContainText('30 минут');
    await expect(consultCard).toContainText('09:00–18:00');

    await expect(page.getByText('45 минут')).toBeVisible();
    await expect(page.getByText('10:00–17:00')).toBeVisible();

    await consultCard.getByRole('button', { name: 'Выбрать слот' }).click();
    await expect(page).toHaveURL(new RegExp(`/event-types/${et.id}$`));
    await expect(page.getByRole('heading', { name: 'Консультация' })).toBeVisible();
  });

  test('G2: полный путь бронирования — подтверждение и исчезновение слота', async ({
    page,
    request,
  }) => {
    const et = await createEventType(request, { name: 'Консультация', durationMinutes: 15 });

    await page.goto(`/event-types/${et.id}`);

    const target = (await getFreeSlots(request, et.id))[0];
    const time = formatSlotTime(target.startsAt);
    const slotBtn = slotButton(page, time).first();
    await expect(slotBtn).toBeVisible();
    const before = await slotButton(page, time).count();
    expect(before).toBeGreaterThan(0);

    await bookSlotViaModal(page, slotBtn, 'Иван Петров', 'ivan@example.com');

    await expect(page.getByRole('heading', { name: 'Запись подтверждена' })).toBeVisible();
    await page.getByRole('button', { name: 'Готово' }).click();
    await expect(page.getByRole('dialog')).toHaveCount(0);

    await expect(slotButton(page, time)).toHaveCount(before - 1);

    const slots = await getFreeSlots(request, et.id);
    expect(slots.some((s) => s.startsAt === target.startsAt)).toBe(false);
  });

  test('G3: повторная бронь того же слота — конфликт 409', async ({ page, request }) => {
    const et = await createEventType(request, { name: 'Консультация', durationMinutes: 30 });

    const pageB = await page.context().newPage();

    await page.goto(`/event-types/${et.id}`);
    await pageB.goto(`/event-types/${et.id}`);

    const { time, button } = await pickFirstSlot(page);
    await expect(pageB.getByRole('button', { name: time, exact: true }).first()).toBeVisible();

    await bookSlotViaModal(page, button, 'Иван Петров', 'ivan@example.com');
    await expect(page.getByRole('heading', { name: 'Запись подтверждена' })).toBeVisible();

    await bookSlotViaModal(pageB, slotButtons(pageB).first(), 'Пётр Иванов', 'petr@example.com');
    await expect(pageB.getByText('Выбранное время уже занято')).toBeVisible();

    await pageB.close();
  });
});
