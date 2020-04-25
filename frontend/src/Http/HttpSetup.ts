const BASE_URL = 'http://localhost';
const PORT = '8080';
const API_BASE = 'v1';

const makeUrl = (url: string): string => {
    return `${BASE_URL}:${PORT}/${API_BASE}/${url}`;
};

export const httpGet = async (url: string): Promise<any> => {
    let response: Response = await fetch(makeUrl(url));
    return await response.json();
};

export const httpPost = async (url: string, body: any): Promise<any> => {

    let response: Response = await fetch(makeUrl(url), {
        method: 'POST', // *GET, POST, PUT, DELETE, etc.
        mode: 'cors', // no-cors, *cors, same-origin
        cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
        credentials: 'same-origin', // include, *same-origin, omit
        headers: {
            'Content-Type': 'application/json'
        },
        redirect: 'follow', // manual, *follow, error
        referrerPolicy: 'no-referrer',
        body: JSON.stringify(body)
    });

    return await response.json();
};
