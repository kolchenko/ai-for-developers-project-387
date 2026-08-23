import { Navigate, Outlet, Route, Routes } from 'react-router-dom';
import { AdminLayout } from './components/AdminLayout';
import { AppLayout } from './components/AppLayout';
import { useAdminAuth } from './auth';
import { AdminBookings } from './pages/AdminBookings';
import { AdminEventTypes } from './pages/AdminEventTypes';
import { AdminLogin } from './pages/AdminLogin';
import { GuestHome } from './pages/GuestHome';
import { GuestSlotSelection } from './pages/GuestSlotSelection';

function RequireAuth() {
  const { isAuthenticated } = useAdminAuth();
  if (!isAuthenticated) {
    return <Navigate to="/admin/login" replace />;
  }
  return <Outlet />;
}

export default function App() {
  return (
    <Routes>
      <Route element={<AppLayout />}>
        <Route index element={<GuestHome />} />
        <Route path="event-types/:eventTypeId" element={<GuestSlotSelection />} />
        <Route path="admin/login" element={<AdminLogin />} />
        <Route element={<RequireAuth />}>
          <Route path="admin" element={<AdminLayout />}>
            <Route index element={<AdminBookings />} />
            <Route path="event-types" element={<AdminEventTypes />} />
          </Route>
        </Route>
      </Route>
    </Routes>
  );
}
