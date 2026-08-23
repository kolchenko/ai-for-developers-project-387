import { api } from './client';
import type {
  Booking,
  BookingCreate,
  EventType,
  EventTypeCreate,
  EventTypeUpdate,
  LoginRequest,
  LoginResponse,
  Slot,
} from './types';

export async function listEventTypes(): Promise<EventType[]> {
  const { data } = await api.get<EventType[]>('/event-types');
  return data;
}

export async function getSlots(eventTypeId: string): Promise<Slot[]> {
  const { data } = await api.get<Slot[]>(`/event-types/${eventTypeId}/slots`);
  return data;
}

export async function createBooking(booking: BookingCreate): Promise<Booking> {
  const { data } = await api.post<Booking>('/bookings', booking);
  return data;
}

export async function adminLogin(credentials: LoginRequest): Promise<LoginResponse> {
  const { data } = await api.post<LoginResponse>('/admin/login', credentials);
  return data;
}

export async function adminCreateEventType(eventType: EventTypeCreate): Promise<EventType> {
  const { data } = await api.post<EventType>('/admin/event-types', eventType);
  return data;
}

export async function adminUpdateEventType(
  eventTypeId: string,
  patch: EventTypeUpdate,
): Promise<EventType> {
  const { data } = await api.patch<EventType>(`/admin/event-types/${eventTypeId}`, patch);
  return data;
}

export async function adminDeleteEventType(eventTypeId: string): Promise<void> {
  await api.delete(`/admin/event-types/${eventTypeId}`);
}

export async function adminUpcomingBookings(): Promise<Booking[]> {
  const { data } = await api.get<Booking[]>('/admin/bookings');
  return data;
}

export async function adminCancelBooking(bookingId: string): Promise<void> {
  await api.delete(`/admin/bookings/${bookingId}`);
}
