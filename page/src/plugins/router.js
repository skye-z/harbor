import {
    createRouter,
    createWebHashHistory
} from 'vue-router'

const Home = () => import('../views/home.vue')
const Auth = () => import('../views/auth.vue')

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            name: 'Home',
            path: '/',
            component: Home,
            meta: {
                title: 'Home',
                auth: true
            }
        },
        {
            name: 'Auth',
            path: '/auth',
            component: Auth,
            meta: {
                title: 'Auth',
                auth: false
            }
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