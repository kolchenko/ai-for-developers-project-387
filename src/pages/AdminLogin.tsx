import {
  Alert,
  Anchor,
  Button,
  Card,
  Center,
  PasswordInput,
  Stack,
  Text,
  TextInput,
  Title,
} from '@mantine/core';
import { IconInfoCircle, IconShieldLock } from '@tabler/icons-react';
import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAdminLogin } from '../api/hooks';
import { useAdminAuth } from '../auth';

export function AdminLogin() {
  const { login } = useAdminAuth();
  const loginMutation = useAdminLogin();
  const navigate = useNavigate();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const error = loginMutation.error?.message;
  const isValid = username.trim().length > 0 && password.length > 0;

  const submit = () => {
    loginMutation.mutate(
      { username: username.trim(), password },
      {
        onSuccess: () => {
          login();
          navigate('/admin', { replace: true });
        },
      },
    );
  };

  return (
    <Center style={{ minHeight: '60vh' }}>
      <Card w={360} withBorder shadow="sm" padding="lg">
        <Stack gap="md">
          <Stack gap={4} align="center">
            <IconShieldLock size={36} stroke={1.5} />
            <Title order={3}>Вход для администратора</Title>
          </Stack>

          {error && (
            <Alert color="red" icon={<IconInfoCircle size={16} />}>
              {error}
            </Alert>
          )}

          <TextInput
            label="Логин"
            placeholder="admin"
            required
            value={username}
            onChange={(e) => setUsername(e.currentTarget.value)}
            disabled={loginMutation.isPending}
          />
          <PasswordInput
            label="Пароль"
            placeholder="••••••••"
            required
            value={password}
            onChange={(e) => setPassword(e.currentTarget.value)}
            onKeyDown={(e) => e.key === 'Enter' && isValid && submit()}
            disabled={loginMutation.isPending}
          />

          <Button fullWidth loading={loginMutation.isPending} disabled={!isValid} onClick={submit}>
            Войти
          </Button>

          <Text size="sm" c="dimmed" ta="center">
            <Anchor component={Link} to="/">
              На главную
            </Anchor>
          </Text>
        </Stack>
      </Card>
    </Center>
  );
}
