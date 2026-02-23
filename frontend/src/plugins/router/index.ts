import { createRouter, createWebHashHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '../stores/user'

// 路由组件懒加载（支持 webpackChunkName）
const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import(/* webpackChunkName: "login" */ '@/views/Login.vue')
  },
  {
    path: '/',
    component: () => import(/* webpackChunkName: "layout" */ '@/layouts/MainLayout.vue'),
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import(/* webpackChunkName: "dashboard" */ '@/views/Dashboard.vue')
      },
      {
        path: 'containers',
        name: 'Containers',
        component: () => import(/* webpackChunkName: "containers" */ '@/views/Containers.vue')
      },
      {
        path: 'containers/:id',
        name: 'ContainerDetail',
        component: () => import(/* webpackChunkName: "container-detail" */ '@/views/ContainerDetail.vue')
      },
      {
        path: 'containers/:id/terminal',
        name: 'ContainerTerminal',
        component: () => import(/* webpackChunkName: "container-terminal" */ '@/views/ContainerTerminal.vue')
      },
      {
        path: 'containers/:id/logs',
        name: 'ContainerLogs',
        component: () => import(/* webpackChunkName: "container-logs" */ '@/views/ContainerLogs.vue')
      },
      {
        path: 'containers/:id/file',
        name: 'ContainerFile',
        component: () => import(/* webpackChunkName: "container-file" */ '@/views/ContainerFile.vue')
      },
      {
        path: 'images',
        name: 'Images',
        component: () => import(/* webpackChunkName: "images" */ '@/views/Images.vue')
      },
      {
        path: 'images/pull',
        name: 'ImagePull',
        component: () => import(/* webpackChunkName: "image-pull" */ '@/views/ImagePull.vue')
      },
      {
        path: 'images/:id',
        name: 'ImageDetail',
        component: () => import(/* webpackChunkName: "image-detail" */ '@/views/ImageDetail.vue')
      },
      {
        path: 'containers/create',
        name: 'ContainerCreate',
        component: () => import(/* webpackChunkName: "container-create" */ '@/views/ContainerCreate.vue')
      },
      {
        path: 'storage',
        name: 'Storage',
        component: () => import(/* webpackChunkName: "storage" */ '@/views/Storage.vue')
      },
      {
        path: 'system',
        name: 'System',
        component: () => import(/* webpackChunkName: "system" */ '@/views/System.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()

  // 允许访问 login 页面
  if (to.name === 'Login') {
    next()
    return
  }

  // 检查用户是否已认证
  if (!userStore.isAuthenticated) {
    next({ name: 'Login' })
    return
  }

  next()
})

export default router


