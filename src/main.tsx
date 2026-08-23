import '@mantine/core/styles.css';
import '@mantine/notifications/styles.css';
import './index.css';

import { MantineProvider } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter } from 'react-router-dom';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';
import { AdminAuthProvider } from './auth';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <MantineProvider>
      <Notifications position="top-right" />
      <QueryClientProvider client={queryClient}>
        <AdminAuthProvider>
          <BrowserRouter>
            <App />
          </BrowserRouter>
        </AdminAuthProvider>
      </QueryClientProvider>
    </MantineProvider>
  </StrictMode>,
);
