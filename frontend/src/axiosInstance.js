import axios from "axios";

const baseURL = import.meta.env.VITE_API_BASE_URL || "http://localhost:3000";

const axiosInstance = axios.create({
    baseURL,
    // headers: { 'Authorization': 'Bearer token' }
    withCredentials: true, // If needed for cookies/session
});

export default axiosInstance;
