export const API_CONFIG = {
  BASE_URL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  TIMEOUT: 30000,
  HEADERS: {
    'Content-Type': 'application/json',
  },
};

export const ENDPOINTS = {
  CONTENTS: {
    SEARCH: '/contents/search',
    GET_BY_ID: (id: string) => `/contents/${id}`,
    STATS: '/contents/stats',
  },
};


