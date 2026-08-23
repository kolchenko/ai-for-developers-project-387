import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  adminCancelBooking,
  adminCreateEventType,
  adminDeleteEventType,
  adminLogin,
  adminUpcomingBookings,
  adminUpdateEventType,
  createBooking,
  getSlots,
  listEventTypes,
} from './endpoints';
import type { EventTypeUpdate } from './types';

export const queryKeys = {
  eventTypes: ['event-types'] as const,
  slots: (eventTypeId: string) => ['slots', eventTypeId] as const,
  upcomingBookings: ['bookings', 'upcoming'] as const,
};

export function useEventTypes() {
  return useQuery({ queryKey: queryKeys.eventTypes, queryFn: listEventTypes });
}

export function useAdminLogin() {
  return useMutation({ mutationFn: adminLogin });
}

export function useSlots(eventTypeId: string) {
  return useQuery({
    queryKey: queryKeys.slots(eventTypeId),
    queryFn: () => getSlots(eventTypeId),
    enabled: Boolean(eventTypeId),
  });
}

export function useCreateBooking() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createBooking,
    onSuccess: (booking) => {
      void queryClient.invalidateQueries({ queryKey: queryKeys.slots(booking.eventTypeId) });
      void queryClient.invalidateQueries({ queryKey: queryKeys.upcomingBookings });
    },
  });
}

export function useUpcomingBookings() {
  return useQuery({ queryKey: queryKeys.upcomingBookings, queryFn: adminUpcomingBookings });
}

export function useCancelBooking() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: adminCancelBooking,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: queryKeys.upcomingBookings });
    },
  });
}

export function useCreateEventType() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: adminCreateEventType,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: queryKeys.eventTypes });
    },
  });
}

export function useUpdateEventType() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, patch }: { id: string; patch: EventTypeUpdate }) =>
      adminUpdateEventType(id, patch),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: queryKeys.eventTypes });
    },
  });
}

export function useDeleteEventType() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: adminDeleteEventType,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: queryKeys.eventTypes });
      void queryClient.invalidateQueries({ queryKey: queryKeys.upcomingBookings });
    },
  });
}
