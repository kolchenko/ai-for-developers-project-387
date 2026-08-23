import {
  Alert,
  Button,
  Card,
  Group,
  LoadingOverlay,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import { IconInfoCircle } from '@tabler/icons-react';
import { useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useEventTypes, useSlots } from '../api/hooks';
import type { Slot } from '../api/types';
import { BookingModal } from '../components/BookingModal';
import { formatDateLong, formatTime, groupSlotsByDay } from '../utils/time';

export function GuestSlotSelection() {
  const { eventTypeId } = useParams<{ eventTypeId: string }>();
  const navigate = useNavigate();

  const { data: eventTypes } = useEventTypes();
  const eventType = eventTypes?.find((et) => et.id === eventTypeId);

  const { data: slots, isLoading, error } = useSlots(eventTypeId ?? '');

  const groups = useMemo(() => (slots ? groupSlotsByDay(slots) : []), [slots]);

  const [selectedSlot, setSelectedSlot] = useState<Slot | null>(null);

  if (eventTypeId === undefined) {
    return null;
  }

  return (
    <Stack gap="md" style={{ position: 'relative' }}>
      <Group justify="space-between">
        <Stack gap={0}>
          <Title order={2}>{eventType?.name ?? 'Слоты'}</Title>
          <Text c="dimmed">
            Свободные слоты формируются на 14 дней. Выберите удобное время для бронирования.
          </Text>
        </Stack>
        <Button variant="subtle" onClick={() => navigate('/')}>
          ← К списку встреч
        </Button>
      </Group>

      {error && (
        <Alert color="red" title="Не удалось загрузить слоты" icon={<IconInfoCircle size={16} />}>
          {error.message}
        </Alert>
      )}

      <LoadingOverlay visible={isLoading} zIndex={1000} />

      {!isLoading && !error && groups.length === 0 && (
        <Alert color="blue" title="Свободных слотов нет" icon={<IconInfoCircle size={16} />}>
          На ближайшие 14 дней доступных слотов не найдено.
        </Alert>
      )}

      {groups.map((group) => (
        <Card key={group.date} shadow="sm" padding="lg" radius="md" withBorder>
          <Stack gap="sm">
            <Text fw={600}>{formatDateLong(group.date)}</Text>
            <Group gap="sm">
              {group.items.map((slot) => (
                <Button
                  key={slot.startsAt}
                  variant="outline"
                  onClick={() => setSelectedSlot(slot)}
                >
                  {formatTime(slot.startsAt)}
                </Button>
              ))}
            </Group>
          </Stack>
        </Card>
      ))}

      {selectedSlot && eventType && (
        <BookingModal
          opened
          onClose={() => setSelectedSlot(null)}
          slot={selectedSlot}
          eventTypeId={eventType.id}
          eventTypeName={eventType.name}
        />
      )}
    </Stack>
  );
}
