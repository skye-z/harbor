import {
    createRouter,
    createWebHashHistory
} from 'vue-router'

const Home = () => import('../views/home.vue')
const Auth = () => import('../views/auth.vue')

const Console = () => import('../views/sub/console.vue')
const Stack = () => import('../views/sub/stack.vue')
const Container = () => import('../views/sub/container.vue')
const Image = () => import('../views/sub/image.vue')
const Network = () => import('../views/sub/network.vue')
const Volume = () => import('../views/sub/volume.vue')
const Monitor = () => import('../views/sub/monitor.vue')
const Logs = () => import('../views/sub/logs.vue')
const User = () => import('../views/sub/user.vue')
const Setting = () => import('../views/sub/setting.vue')

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            name: 'Auth',
            path: '/auth',
            component: Auth,
            meta: {
                title: 'Auth',
                auth: false
            }
        },
        {
            name: 'Home',
            path: '/',
            component: Home,
            meta: {
                title: 'Home',
                auth: true
            },
            children: [
                {
                    name: 'Console',
                    path: '/console',
                    component: Console,
                    meta: {
                        title: '控制台',
                        auth: true
                    }
                },
                {
                    name: 'Stack',
                    path: '/stack',
                    component: Stack,
                    meta: {
                        title: '服务编排',
                        auth: true
                    }
                },
                {
                    name: 'Container',
                    path: '/container',
                    component: Container,
                    meta: {
                        title: '应用容器',
                        auth: true
                    }
                },
                {
                    name: 'Image',
                    path: '/image',
                    component: Image,
                    meta: {
                        title: '镜像仓库',
                        auth: true
                    }
                },
                {
                    name: 'Network',
                    path: '/network',
                    component: Network,
                    meta: {
                        title: '内部网络',
                        auth: true
                    }
                },
                {
                    name: 'Volume',
                    path: '/volume',
                    component: Volume,
                    meta: {
                        title: '存储卷',
                        auth: true
                    }
                },
                {
                    name: 'Monitor',
                    path: '/monitor',
                    component: Monitor,
                    meta: {
                        title: '容器监控',
                        auth: true
                    }
                },
                {
                    name: 'Logs',
                    path: '/logs',
                    component: Logs,
                    meta: {
                        title: '平台日志',
                        auth: true
                    }
                },
                {
                    name: 'User',
                    path: '/user',
                    component: User,
                    meta: {
                        title: '用户管理',
                        auth: true
                    }
                },
                {
                    name: 'Setting',
                    path: '/setting',
                    component: Setting,
                    meta: {
                        title: '平台设置',
                        auth: true
                    }
                },
            ]
        },
    ],
})

router.beforeEach((to, _from, next) => {
    if (!to.meta.auth || Boolean(localStorage.getItem("user:access:token")) || to.name == "Auth") {
        document.title = (to.meta.title === undefined ? 'Unknown Page - ' : to.meta.title + ' - ') + 'Harbor'
        next()
    } else {
        next('/auth')
    }
})

export default router