import { Alert, Button, Group, Modal, Stack, Text, TextInput, Title } from '@mantine/core';
import { IconCircleCheck, IconInfoCircle } from '@tabler/icons-react';
import { useState } from 'react';
import { useCreateBooking } from '../api/hooks';
import type { Slot } from '../api/types';
import { formatTime } from '../utils/time';

interface BookingModalProps {
  opened: boolean;
  onClose: () => void;
  slot: Slot;
  eventTypeId: string;
  eventTypeName: string;
}

export function BookingModal({ opened, onClose, slot, eventTypeId, eventTypeName }: BookingModalProps) {
  const [guestName, setGuestName] = useState('');
  const [guestEmail, setGuestEmail] = useState('');
  const createBooking = useCreateBooking();

  const error = createBooking.error?.message;

  const submit = () => {
    createBooking.mutate({ eventTypeId, startsAt: slot.startsAt, guestName, guestEmail });
  };

  const handleClose = () => {
    createBooking.reset();
    setGuestName('');
    setGuestEmail('');
    onClose();
  };

  if (createBooking.isSuccess) {
    return (
      <Modal opened={opened} onClose={handleClose} title="Бронирование" centered>
        <Stack gap="md" align="center">
          <IconCircleCheck size={48} color="var(--mantine-color-green-6)" />
          <Title order={3}>Запись подтверждена</Title>
          <Text ta="center" c="dimmed">
            {eventTypeName}
            <br />
            {formatTime(slot.startsAt)}–{formatTime(slot.endsAt)}
          </Text>
          <Button fullWidth onClick={handleClose}>
            Готово
          </Button>
        </Stack>
      </Modal>
    );
  }

  return (
    <Modal opened={opened} onClose={handleClose} title="Бронирование" centered>
      <Stack gap="md">
        <Group justify="space-between">
          <Text size="sm" c="dimmed">
            {eventTypeName}
          </Text>
          <Text size="sm" fw={600}>
            {formatTime(slot.startsAt)}–{formatTime(slot.endsAt)}
          </Text>
        </Group>

        {error && (
          <Alert color="red" icon={<IconInfoCircle size={16} />}>
            {error}
          </Alert>
        )}

        <TextInput
          label="Ваше имя"
          placeholder="Иван Петров"
          required
          value={guestName}
          onChange={(e) => setGuestName(e.currentTarget.value)}
          disabled={createBooking.isPending}
        />
        <TextInput
          label="Email"
          placeholder="ivan@example.com"
          required
          type="email"
          value={guestEmail}
          onChange={(e) => setGuestEmail(e.currentTarget.value)}
          disabled={createBooking.isPending}
        />

        <Button
          fullWidth
          loading={createBooking.isPending}
          disabled={!guestName.trim() || !guestEmail.trim()}
          onClick={submit}
        >
          Подтвердить бронирование
        </Button>
      </Stack>
    </Modal>
  );
}
