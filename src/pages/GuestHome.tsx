import { Alert, Button, Card, Group, LoadingOverlay, SimpleGrid, Stack, Text, Title } from '@mantine/core';
import { IconClock, IconInfoCircle } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import { useEventTypes } from '../api/hooks';

const DURATION_LABELS: Record<number, string> = {
  15: '15 минут',
  30: '30 минут',
  45: '45 минут',
  60: '1 час',
};

export function GuestHome() {
  const { data: eventTypes, isLoading, error } = useEventTypes();
  const navigate = useNavigate();

  return (
    <Stack gap="md">
      <Title order={2}>Виды брони</Title>
      <Text c="dimmed">Выберите тип встречи, чтобы посмотреть свободные слоты.</Text>

      {error && (
        <Alert color="red" title="Не удалось загрузить типы событий" icon={<IconInfoCircle size={16} />}>
          {error.message}
        </Alert>
      )}

      <div style={{ position: 'relative' }}>
        <LoadingOverlay visible={isLoading} zIndex={1000} />
        <SimpleGrid cols={{ base: 1, sm: 2, lg: 3 }}>
          {(eventTypes ?? []).map((eventType) => (
            <Card
              key={eventType.id}
              data-testid={`event-type-card-${eventType.id}`}
              shadow="sm"
              padding="lg"
              radius="md"
              withBorder
            >
              <Stack gap="xs">
                <Title order={3}>{eventType.name}</Title>
                <Text>{eventType.description}</Text>
                <Group gap="xs">
                  <IconClock size={16} />
                  <Text size="sm" c="dimmed">
                    {DURATION_LABELS[eventType.durationMinutes]} ·{' '}
                    {eventType.availableFrom.slice(0, 5)}–{eventType.availableTo.slice(0, 5)}
                  </Text>
                </Group>
                <Button fullWidth mt="sm" onClick={() => navigate(`/event-types/${eventType.id}`)}>
                  Выбрать слот
                </Button>
              </Stack>
            </Card>
          ))}
        </SimpleGrid>
      </div>
    </Stack>
  );
}
