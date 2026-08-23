import axios from 'axios';
import { ApiError, getErrorMessage } from './types';

export const api = axios.create({
  baseURL: '/api',
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (axios.isAxiosError(error) && error.response) {
      const status = error.response.status;
      const message =
        typeof error.response.data?.message === 'string'
          ? error.response.data.message
          : getErrorMessage(status);
      return Promise.reject(new ApiError(status, message));
    }
    return Promise.reject(error);
  },
);
