import axios, { type AxiosInstance } from "axios";

const http: AxiosInstance = axios.create({
    baseURL: '/api',
    timeout: 300000,
    headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
    },
});

export default http;