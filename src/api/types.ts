export const SLOT_DURATIONS = [15, 30, 45, 60] as const;

export type SlotDuration = (typeof SLOT_DURATIONS)[number];

export interface EventType {
  id: string;
  name: string;
  description: string;
  durationMinutes: SlotDuration;
  availableFrom: string;
  availableTo: string;
}

export type EventTypeCreate = Omit<EventType, 'id'>;

export type EventTypeUpdate = Partial<EventTypeCreate>;

export interface Booking {
  id: string;
  eventTypeId: string;
  startsAt: string;
  endsAt: string;
  guestName: string;
  guestEmail: string;
}

export type BookingCreate = Omit<Booking, 'id' | 'endsAt'>;

export interface Slot {
  startsAt: string;
  endsAt: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  username: string;
}

export type ApiErrorKind = 401 | 404 | 409 | 422;

export class ApiError extends Error {
  readonly status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

export function getErrorMessage(status: number): string {
  switch (status) {
    case 401:
      return 'Неверный логин или пароль';
    case 404:
      return 'Ресурс не найден';
    case 409:
      return 'Выбранное время уже занято';
    case 422:
      return 'Слот недопустим';
    default:
      return 'Произошла ошибка при обращении к серверу';
  }
}
