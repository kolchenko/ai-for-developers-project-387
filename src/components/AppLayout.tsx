import { AppShell, Button, Group, Title } from '@mantine/core';
import { IconCalendarEvent, IconShield } from '@tabler/icons-react';
import { Link, Outlet, useLocation } from 'react-router-dom';

export function AppLayout() {
  const { pathname } = useLocation();

  return (
    <AppShell header={{ height: 60 }} padding="md">
      <AppShell.Header>
        <Group h="100%" px="md" justify="space-between">
          <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
            <Group gap="xs" style={{ cursor: 'pointer' }}>
              <IconCalendarEvent size={24} />
              <Title order={4}>Календарь звонков</Title>
            </Group>
          </Link>
          <Group gap="sm">
            <Button
              variant={pathname === '/' ? 'filled' : 'light'}
              component={Link}
              to="/"
              leftSection={<IconCalendarEvent size={16} />}
            >
              Выбрать встречу
            </Button>
            <Button
              variant={pathname.startsWith('/admin') ? 'filled' : 'light'}
              component={Link}
              to="/admin"
              leftSection={<IconShield size={16} />}
            >
              Админ
            </Button>
          </Group>
        </Group>
      </AppShell.Header>
      <AppShell.Main>
        <Outlet />
      </AppShell.Main>
    </AppShell>
  );
}
