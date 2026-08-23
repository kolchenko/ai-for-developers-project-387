import { expect, test } from '@playwright/test';
import { API_BASE, createEventType, getFreeSlots, resetData } from './helpers';

test.describe('Контракт API через фронтенд-прокси', () => {
  test.beforeEach(async ({ request }) => {
    await resetData(request);
  });

  test('E1: слоты несуществующего типа события — 404', async ({ request }) => {
    const res = await request.get(`${API_BASE}/event-types/et-missing/slots`);
    expect(res.status()).toBe(404);
  });

  test('E2: бронь вне сетки слотов — 422', async ({ request }) => {
    const et = await createEventType(request, { durationMinutes: 30 });
    const slots = await getFreeSlots(request, et.id);
    expect(slots.length).toBeGreaterThan(0);

    const startsAt = new Date(new Date(slots[0].startsAt).getTime() + 60_000).toISOString();
    const res = await request.post(`${API_BASE}/bookings`, {
      data: { eventTypeId: et.id, startsAt, guestName: 'Иван', guestEmail: 'ivan@example.com' },
    });
    expect(res.status()).toBe(422);
  });

  test('E3: повторная бронь того же слота — 409', async ({ request }) => {
    const et = await createEventType(request, { durationMinutes: 30 });
    const slots = await getFreeSlots(request, et.id);
    expect(slots.length).toBeGreaterThan(0);

    const payload = {
      eventTypeId: et.id,
      startsAt: slots[0].startsAt,
      guestName: 'Иван',
      guestEmail: 'ivan@example.com',
    };
    const first = await request.post(`${API_BASE}/bookings`, { data: payload });
    expect(first.status()).toBe(201);

    const second = await request.post(`${API_BASE}/bookings`, { data: payload });
    expect(second.status()).toBe(409);
  });

  test('E4: неверный логин администратора — 401', async ({ request }) => {
    const res = await request.post(`${API_BASE}/admin/login`, {
      data: { username: 'admin', password: 'wrong' },
    });
    expect(res.status()).toBe(401);
  });
});
