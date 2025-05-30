import http from '@/api';

export interface MLBotConfig {
    id: number;
    enabled: boolean;
    host: string;
    port: number;
    uid: string;
    protocol: string;
    messagePrefix: string;
}

export interface MLBotConfigRequest {
    enabled: boolean;
    host: string;
    port: number;
    uid: string;
    protocol: string;
    messagePrefix: string;
}

export interface MLBotTestRequest {
    host: string;
    port: number;
    uid: string;
    protocol: string;
    message: string;
}

// Get mlBot configuration
export const getMLBotConfig = () => {
    return http.get<MLBotConfig>('/mlbot/config');
};

// Update mlBot configuration
export const updateMLBotConfig = (config: MLBotConfigRequest) => {
    return http.post('/mlbot/config', config);
};

// Test mlBot notification
export const testMLBotNotification = (params: MLBotTestRequest) => {
    return http.post('/mlbot/test', params);
};
