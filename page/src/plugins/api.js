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

function get(url) {
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
    getInfo: id => get('/container/info?id=' + id),
    getLogs: (id, tail) => get('/container/logs?id=' + id + '&tail=' + tail),
    start: id => get('/container/start?id=' + id),
    stop: id => get('/container/stop?id=' + id),
    restart: id => get('/container/restart?id=' + id),
    kill: id => get('/container/kill?id=' + id),
    pause: id => get('/container/pause?id=' + id),
    unpause: id => get('/container/unpause?id=' + id)
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