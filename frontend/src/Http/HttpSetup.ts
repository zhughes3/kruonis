const BASE_URL = 'http://timelines.dev';
const PORT = '3000';
const API_BASE = 'v1';

const makeUrl = (url: string): string => {
    return `${BASE_URL}:${PORT}/${API_BASE}/${url}`;
}

export const httpGet = async (url: string): Promise<any> => {
    let response: Response = await fetch(makeUrl(url));
    let data = await response.json()
    return data;
}

export const httpPost = async (url: string, body: any): Promise<any> => {
    
    let response: Response = await fetch(makeUrl(url), {
        method: 'POST',
        mode: 'no-cors',
        credentials: 'omit',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(body)
    });
    
    let data = await response.json()
    return data;
}