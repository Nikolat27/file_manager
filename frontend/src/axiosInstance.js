import axios from "axios";
import { useUserStore } from "./stores/user";

const baseURL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8000";

const axiosInstance = axios.create({
    baseURL,
    // headers: { 'Authorization': 'Bearer token' }
    withCredentials: true, // If needed for cookies/session
});

axiosInstance.interceptors.request.use(
    (config) => {
        const userStore = useUserStore();
        if (userStore.token) {
            config.headers["Authorization"] = userStore.token;
        }

        return config;
    },
    (error) => Promise.reject(error)
);

export default axiosInstance;
