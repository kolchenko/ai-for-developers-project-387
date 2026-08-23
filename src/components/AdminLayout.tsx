import { Button, Group } from '@mantine/core';
import { IconCalendarEvent, IconClipboardList, IconLogout } from '@tabler/icons-react';
import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useAdminAuth } from '../auth';

export function AdminLayout() {
  const { pathname } = useLocation();
  const { logout } = useAdminAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/admin/login', { replace: true });
  };

  return (
    <>
      <Group mb="md" gap="xs" justify="space-between">
        <Group gap="xs">
          <Button
            variant={pathname === '/admin' ? 'filled' : 'subtle'}
            component={Link}
            to="/admin"
            leftSection={<IconCalendarEvent size={16} />}
          >
            Предстоящие встречи
          </Button>
          <Button
            variant={pathname.startsWith('/admin/event-types') ? 'filled' : 'subtle'}
            component={Link}
            to="/admin/event-types"
            leftSection={<IconClipboardList size={16} />}
          >
            Типы событий
          </Button>
        </Group>
        <Button
          variant="subtle"
          color="gray"
          onClick={handleLogout}
          leftSection={<IconLogout size={16} />}
        >
          Выйти
        </Button>
      </Group>
      <Outlet />
    </>
  );
}
