import axios from 'axios';

const API_BASE_URL = '/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authAPI = {
  login: (username: string, password: string) =>
    api.post('/auth/login', { username, password }),
  register: (username: string, email: string, password: string, role: string) =>
    api.post('/auth/register', { username, email, password, role }),
};

export const packageAPI = {
  list: (page = 1, pageSize = 20) =>
    api.get(`/packages?page=${page}&page_size=${pageSize}`),
  get: (id: string) =>
    api.get(`/packages/${id}`),
  create: (formData: FormData) =>
    api.post('/packages', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),
  update: (id: string, data: any) =>
    api.put(`/packages/${id}`, data),
  delete: (id: string) =>
    api.delete(`/packages/${id}`),
  download: (id: string) =>
    api.get(`/packages/${id}/download`, { responseType: 'blob' }),
};

export const passwordAPI = {
  list: (packageId: string) =>
    api.get(`/passwords/packages/${packageId}/passwords`),
  create: (packageId: string, data: any) =>
    api.post(`/passwords/packages/${packageId}/passwords`, data),
  update: (id: string, data: any) =>
    api.put(`/passwords/${id}`, data),
  delete: (id: string) =>
    api.delete(`/passwords/${id}`),
  activate: (id: string) =>
    api.post(`/passwords/${id}/activate`),
  deactivate: (id: string) =>
    api.post(`/passwords/${id}/deactivate`),
};

export const verifyAPI = {
  verify: (packageId: string, password: string) =>
    api.post('/verify', { package_id: packageId, password }),
  batchVerify: (packageId: string, passwords: string[]) =>
    api.post('/verify/batch', { package_id: packageId, passwords }),
};

export const auditAPI = {
  list: (params: any) =>
    api.get('/audit', { params }),
  listByPackage: (packageId: string, page = 1, pageSize = 20) =>
    api.get(`/audit/package/${packageId}?page=${page}&page_size=${pageSize}`),
  export: (format = 'json') =>
    api.get(`/audit/export?format=${format}`),
};

export const adminAPI = {
  getEncryptionConfig: () =>
    api.get('/admin/encryption/config'),
  updateEncryptionConfig: (data: any) =>
    api.put('/admin/encryption/config', data),
  getDashboardStats: () =>
    api.get('/admin/stats/dashboard'),
  getRoles: () =>
    api.get('/admin/roles'),
  getSystemInfo: () =>
    api.get('/admin/system'),
};

export default api;
