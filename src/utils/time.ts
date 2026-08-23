import dayjs from 'dayjs';
import 'dayjs/locale/ru';

dayjs.locale('ru');

export function formatTime(iso: string): string {
  return dayjs(iso).format('HH:mm');
}

export function formatDate(iso: string): string {
  return dayjs(iso).format('DD.MM.YYYY');
}

export function formatDateLong(iso: string): string {
  return dayjs(iso).format('dddd, D MMMM');
}

export function groupSlotsByDay<T extends { startsAt: string }>(slots: T[]): Array<{
  date: string;
  items: T[];
}> {
  const groups = new Map<string, T[]>();
  for (const slot of slots) {
    const date = dayjs(slot.startsAt).format('YYYY-MM-DD');
    const items = groups.get(date) ?? [];
    items.push(slot);
    groups.set(date, items);
  }
  return [...groups.entries()]
    .map(([date, items]) => ({
      date,
      items: [...items].sort((a, b) => a.startsAt.localeCompare(b.startsAt)),
    }))
    .sort((a, b) => a.date.localeCompare(b.date));
}

export function toTimeValue(value: string): string {
  return value.length === 5 ? `${value}:00` : value;
}
