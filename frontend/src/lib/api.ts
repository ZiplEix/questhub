import axios from "axios";
import { PUBLIC_BASE_API_URL } from '$env/static/public';
import { goto } from "$app/navigation";

export const api = axios.create({
    baseURL: PUBLIC_BASE_API_URL,
    withCredentials: false,
    headers: {
        'Content-Type': 'application/json'
    }
})

api.interceptors.response.use(
    response => response,
    error => {
        if (error.response && error.response.status === 401) {
            goto('/login');
        }
        return Promise.reject(error);
    }
)
