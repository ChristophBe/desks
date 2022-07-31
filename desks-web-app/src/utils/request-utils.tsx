type RequestMethod = 'GET' | 'POST' | 'PATCH' | 'PUT' | 'DELETE';

async function request<B>(method: RequestMethod, url: string, body?: B): Promise<Response> {
    const requestInit: RequestInit = {
        method: method,
        mode: 'cors',
        cache: 'no-cache',
        credentials: 'same-origin',
        headers: {
            'Content-Type': 'application/json'
        },
        redirect: 'follow',
        referrerPolicy: 'no-referrer',
    }
    if (body) {
        requestInit.body = JSON.stringify(body)
    }
    return await fetch(url, requestInit);
}

export async function postData<B>(url: string, body: B): Promise<Response> {
    return request<B>("POST", url, body)
}

export async function patchData<B>(url: string, body: B): Promise<Response> {
    return request<B>("PATCH", url, body)
}

export async function deleteData(url: string): Promise<Response> {
    return request("DELETE", url)
}

export async function getData(url: string): Promise<Response> {
    return request("GET", url)
}


