import axios, { AxiosRequestConfig } from 'axios';

export const requestGet = async <T>(config: AxiosRequestConfig): Promise<T | undefined> => {
    if (!config?.url) {
        return undefined;
    }

    try {
        const response = await axios.get(config.url);
        return response ? response.data : undefined;
    } catch (err: any) {
        console.error('=====> ERR', err.message);
        return undefined;
    }
};
