import axios from "axios";

export type UserData = {
  id: number;
  name: string;
  description: string;
};

type UpdateUserDataRequest = {
  name: string;
  description: string;
};

type LoginRequest = {
  login: string;
  password: string;
};

type LoginResponse = {
  access_token: string;
};

type RegisterRequest = {
  login: string;
  password: string;
};

type RegisterResponse = {
  access_token: string;
};

type RefreshResponse = {
  access_token: string;
};

const internalApi = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL + "/api/v1",
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

const refreshApi = {
  refresh: () =>
    internalApi
      .post<RefreshResponse>("/auth/refresh")
      .then((response) =>
        localStorage.setItem("accessToken", response.data.access_token),
      )
      .catch(() => localStorage.removeItem("accessToken"))
      .finally(() => console.log("Sent refresh request")),
};

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL + "/api/v1",
  withCredentials: true,
  timeout: 3000,
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("accessToken");

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    console.log("Failed response");
    const originalRequest = error.config;

    if (error.response?.status === 401) {
      console.log("Failed: Unauthorised");

      if (!originalRequest._retry) {
        console.log("Retrying");

        originalRequest._retry = true;

        await refreshApi.refresh();

        return api(originalRequest);
      }
    }

    console.log("Failed even after retrieval");
    return Promise.reject(error);
  },
);

export const authApi = {
  login: (data: LoginRequest) =>
    api
      .post<LoginResponse>("/auth/login", data)
      .then((response) =>
        localStorage.setItem("accessToken", response.data.access_token),
      ),

  register: (data: RegisterRequest) =>
    api
      .post<RegisterResponse>("/auth/register", data)
      .then((response) =>
        localStorage.setItem("accessToken", response.data.access_token),
      ),

  logout: () =>
    api.post("/auth/logout").then(() => localStorage.removeItem("accessToken")),
};

export const usersApi = {
  getMe: () => api.get<UserData>("/users/my"),

  getUser: (id: number) => api.get<UserData>(`/users/${id}`),

  updateUser: (id: number, data: UpdateUserDataRequest) =>
    api.put<UserData>(`/users/${id}`, data),

  deleteUser: (id: number) => api.delete<UserData>(`/users/${id}`),
};
