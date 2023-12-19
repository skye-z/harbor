import request from './request'

function post(url, data) {
    return request({
        url: url,
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        data: JSON.stringify(data)
    })
}

function get(url){
    return request({
        url: url,
        method: 'GET'
    })
}

export const docker = {
    getInfo: () => get('/docker/info'),
}

export const container = {
    getList: () => get('/container/list'),
}

export const image = {
    getList: () => get('/image/list'),
}

export const network = {
    getList: () => get('/network/list'),
}

export const volume = {
    getList: () => get('/volume/list'),
}