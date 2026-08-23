import {
  ActionIcon,
  Alert,
  Button,
  Group,
  Modal,
  Stack,
  Table,
  Text,
  Title,
  Tooltip,
} from '@mantine/core';
import { IconInfoCircle, IconPencil, IconPlus, IconTrash } from '@tabler/icons-react';
import { useState } from 'react';
import { useDeleteEventType, useEventTypes } from '../api/hooks';
import type { EventType } from '../api/types';
import { EventTypeFormModal } from '../components/EventTypeFormModal';

const DURATION_LABELS: Record<number, string> = {
  15: '15 минут',
  30: '30 минут',
  45: '45 минут',
  60: '1 час',
};

export function AdminEventTypes() {
  const { data: eventTypes, isLoading, error } = useEventTypes();
  const deleteEventType = useDeleteEventType();

  const [formOpened, setFormOpened] = useState(false);
  const [editing, setEditing] = useState<EventType | undefined>(undefined);
  const [deleting, setDeleting] = useState<EventType | undefined>(undefined);

  const openCreate = () => {
    setEditing(undefined);
    setFormOpened(true);
  };

  const openEdit = (eventType: EventType) => {
    setEditing(eventType);
    setFormOpened(true);
  };

  const confirmDelete = () => {
    if (!deleting) return;
    deleteEventType.mutate(deleting.id, {
      onSuccess: () => setDeleting(undefined),
    });
  };

  return (
    <Stack gap="md">
      <Group justify="space-between">
        <Title order={2}>Типы событий</Title>
        <Button leftSection={<IconPlus size={16} />} onClick={openCreate}>
          Новый тип события
        </Button>
      </Group>

      {error && (
        <Alert color="red" title="Не удалось загрузить типы событий" icon={<IconInfoCircle size={16} />}>
          {error.message}
        </Alert>
      )}

      <Table striped highlightOnHover withTableBorder>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Название</Table.Th>
            <Table.Th>Описание</Table.Th>
            <Table.Th>Длительность</Table.Th>
            <Table.Th>Окно записи</Table.Th>
            <Table.Th style={{ width: 90 }} />
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {(eventTypes ?? []).map((eventType) => (
            <Table.Tr key={eventType.id}>
              <Table.Td fw={600}>{eventType.name}</Table.Td>
              <Table.Td>{eventType.description}</Table.Td>
              <Table.Td>{DURATION_LABELS[eventType.durationMinutes]}</Table.Td>
              <Table.Td>
                {eventType.availableFrom.slice(0, 5)}–{eventType.availableTo.slice(0, 5)}
              </Table.Td>
              <Table.Td>
                <Group gap={4} wrap="nowrap">
                  <Tooltip label="Редактировать">
                    <ActionIcon variant="subtle" onClick={() => openEdit(eventType)}>
                      <IconPencil size={16} />
                    </ActionIcon>
                  </Tooltip>
                  <Tooltip label="Удалить">
                    <ActionIcon variant="subtle" color="red" onClick={() => setDeleting(eventType)}>
                      <IconTrash size={16} />
                    </ActionIcon>
                  </Tooltip>
                </Group>
              </Table.Td>
            </Table.Tr>
          ))}
        </Table.Tbody>
      </Table>

      {!isLoading && (eventTypes ?? []).length === 0 && !error && (
        <Text c="dimmed">Типов событий пока нет. Создайте первый.</Text>
      )}

      {formOpened && (
        <EventTypeFormModal
          key={editing?.id ?? 'new'}
          opened
          onClose={() => setFormOpened(false)}
          eventType={editing}
        />
      )}

      <Modal
        opened={Boolean(deleting)}
        onClose={() => setDeleting(undefined)}
        title="Удалить тип события?"
        centered
      >
        <Stack gap="md">
          <Text>
            Бронирования типа «{deleting?.name}» будут удалены вместе с ним.
          </Text>
          <Group justify="flex-end">
            <Button variant="default" onClick={() => setDeleting(undefined)}>
              Отмена
            </Button>
            <Button color="red" loading={deleteEventType.isPending} onClick={confirmDelete}>
              Удалить
            </Button>
          </Group>
        </Stack>
      </Modal>
    </Stack>
  );
}
