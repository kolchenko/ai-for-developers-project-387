import {
  Alert,
  Badge,
  Button,
  Group,
  LoadingOverlay,
  Modal,
  Stack,
  Table,
  Text,
  Title,
} from '@mantine/core';
import { IconInfoCircle, IconX } from '@tabler/icons-react';
import { useState } from 'react';
import { useCancelBooking, useEventTypes, useUpcomingBookings } from '../api/hooks';
import type { Booking } from '../api/types';
import { formatDate, formatTime } from '../utils/time';

export function AdminBookings() {
  const { data: bookings, isLoading, error } = useUpcomingBookings();
  const { data: eventTypes } = useEventTypes();
  const [cancelling, setCancelling] = useState<string | null>(null);

  const nameOf = (eventTypeId: string) =>
    eventTypes?.find((et) => et.id === eventTypeId)?.name ?? eventTypeId;

  return (
    <Stack gap="md" style={{ position: 'relative' }}>
      <Title order={2}>Предстоящие встречи</Title>

      {error && (
        <Alert color="red" title="Не удалось загрузить встречи" icon={<IconInfoCircle size={16} />}>
          {error.message}
        </Alert>
      )}

      <LoadingOverlay visible={isLoading} zIndex={1000} />

      {!isLoading && !error && (bookings ?? []).length === 0 && (
        <Text c="dimmed">Предстоящих встреч нет.</Text>
      )}

      <Table striped highlightOnHover withTableBorder>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Дата</Table.Th>
            <Table.Th>Время</Table.Th>
            <Table.Th>Тип события</Table.Th>
            <Table.Th>Гость</Table.Th>
            <Table.Th style={{ width: 100 }} />
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {(bookings ?? []).map((booking: Booking) => (
            <Table.Tr key={booking.id}>
              <Table.Td>{formatDate(booking.startsAt)}</Table.Td>
              <Table.Td>
                {formatTime(booking.startsAt)}–{formatTime(booking.endsAt)}
              </Table.Td>
              <Table.Td>
                <Badge variant="light">{nameOf(booking.eventTypeId)}</Badge>
              </Table.Td>
              <Table.Td>
                <Stack gap={0}>
                  <Text size="sm">{booking.guestName}</Text>
                  <Text size="xs" c="dimmed">
                    {booking.guestEmail}
                  </Text>
                </Stack>
              </Table.Td>
              <Table.Td>
                <Button
                  variant="light"
                  color="red"
                  size="xs"
                  leftSection={<IconX size={14} />}
                  loading={cancelling === booking.id}
                  onClick={() => setCancelling(booking.id)}
                >
                  Отменить
                </Button>
              </Table.Td>
            </Table.Tr>
          ))}
        </Table.Tbody>
      </Table>

      {cancelling && (
        <ConfirmCancelModal
          booking={bookings?.find((b) => b.id === cancelling) ?? null}
          onClose={() => setCancelling(null)}
        />
      )}
    </Stack>
  );
}

function ConfirmCancelModal({
  booking,
  onClose,
}: {
  booking: Booking | null;
  onClose: () => void;
}) {
  const cancelBooking = useCancelBooking();

  const confirm = () => {
    if (!booking) return;
    cancelBooking.mutate(booking.id, { onSuccess: onClose });
  };

  return (
    <Modal opened={Boolean(booking)} onClose={onClose} title="Отменить бронирование?" centered>
      <Stack gap="md">
        <Text>
          {booking ? `${booking.guestName}, ${formatDate(booking.startsAt)} ${formatTime(booking.startsAt)}` : ''}
        </Text>
        <Group justify="flex-end">
          <Button variant="default" onClick={onClose}>
            Назад
          </Button>
          <Button color="red" loading={cancelBooking.isPending} onClick={confirm}>
            Отменить
          </Button>
        </Group>
      </Stack>
    </Modal>
  );
}
