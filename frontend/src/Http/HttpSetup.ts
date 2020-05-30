import {ReqMethod} from "../Interfaces/Enums/ReqMethod";

const BASE_URL = 'http://localhost';
const PORT = '8080';
// This port is only used for posting images.
const IMAGE_PORT = '8081';
const API_BASE = 'v1';

// This function created the URL used to communicate with the API.
// The image boolean exists because it the images are posted to a different port then other content.
// If you set images to true, you'll send data to a different port then if it is false.
const makeUrl = (url: string, image?: boolean): string => {
    return `${BASE_URL}:${image ? IMAGE_PORT : PORT}/${API_BASE}/${url}`;
};

const getReqConfig = (reqMethod: ReqMethod, body?: any): any => {
    return {
        method: reqMethod, // *GET, POST, PUT, DELETE, etc.
        mode: 'cors', // no-cors, *cors, same-origin
        cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
        credentials: 'include', // include, *same-origin, omit
        headers: {
            'Content-Type': 'application/json',
        },
        redirect: 'follow', // manual, *follow, error
        referrerPolicy: 'no-referrer',
        body
    }
};

export const httpGet = async (url: string): Promise<any> => {

    let response: Response = await fetch(makeUrl(url), getReqConfig(ReqMethod.GET));

    return await response.json();
};

// If you want to post an image, just set image to true when calling this function. The right PORT will be set automatically.
// Also, if you are sending an image, the body will not be stringified.
export const httpPost = async (url: string, body: any, image?: boolean): Promise<any> => {

    let response: Response = await fetch(makeUrl(url, image), getReqConfig(ReqMethod.POST, image ? body : JSON.stringify(body)));

    return await response.json();
};

export const httpPut = async (url: string, body: any): Promise<any> => {

    let response: Response = await fetch(makeUrl(url), getReqConfig(ReqMethod.PUT, JSON.stringify(body)));

    return await response.json();
};

export const httpDelete = async (url: string): Promise<any> => {

    let response: Response = await fetch(makeUrl(url), getReqConfig(ReqMethod.DELETE));

    return await response.json();
};
