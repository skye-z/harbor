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

// Docker接口
export const docker = {
    getInfo: () => get('/docker/info'),
}