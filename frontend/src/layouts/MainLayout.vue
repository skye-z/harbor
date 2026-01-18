<template>
   <div>
     <n-scrollbar style="height: 100vh">
       <n-el tag="div" class="header">
         <div class="header-content">
           <div class="logo">Harbor</div>
           <div class="menu-section">
             <n-tooltip v-for="option in menuOptions" :key="option.key" :show-arrow="false" placement="bottom">
               <template #trigger>
                 <div class="menu-item" :class="{ 'active': activeKey === option.key }"
                   @click="handleMenuSelect(option.key)">
                   <n-icon :component="option.icon" :size="20" />
                 </div>
               </template>
               {{ option.label }}
             </n-tooltip>
           </div>
           <div class="user-section">
             <n-dropdown :options="userOptions" @select="handleUserSelect">
               <n-button text>
                 <template #icon>
                   <n-icon :component="PersonCircleOutline" />
                 </template>
                 <span>{{ userStore.user?.username || 'Admin' }}</span>
               </n-button>
             </n-dropdown>
           </div>
         </div>
       </n-el>
       <main class="content">
         <router-view v-slot="{ Component, route }">
           <transition name="page" mode="out-in">
             <component :is="Component" :key="route.path" />
           </transition>
         </router-view>
       </main>
      </n-scrollbar>
    </div>
</template>

<script setup lang="ts">
import { computed, h, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NIcon, useLoadingBar } from 'naive-ui'
import { useUserStore } from '../plugins/stores/user'
import {
  PersonCircleOutline,
  SpeedometerOutline,
  CubeOutline,
  ImagesOutline,
  ServerOutline,
  MoonOutline,
  LogOutOutline,
  SettingsOutline
} from '@vicons/ionicons5'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const loadingBar = useLoadingBar()

// 根据当前路由路径的第一段动态计算 activeKey
const activeKey = computed(() => {
  const path = route.path
  const firstSegment = path.split('/')[1] || 'dashboard'
  return firstSegment
})

const menuOptions = [
  {
    label: '概览',
    key: 'dashboard',
    icon: SpeedometerOutline
  },
  {
    label: '容器',
    key: 'containers',
    icon: CubeOutline
  },
  {
    label: '镜像',
    key: 'images',
    icon: ImagesOutline
  },
  {
    label: '连接',
    key: 'storage',
    icon: ServerOutline
  },
  {
    label: '系统',
    key: 'system',
    icon: SettingsOutline
  }
]

const userOptions = computed(() => [
  {
    label: '切换主题',
    key: 'toggle-theme',
    icon() {
      return h(NIcon, null, {
        default: () => h(MoonOutline)
      })
    },
  },
  {
    type: 'divider'
  },
  {
    label: '退出登录',
    key: 'logout',
    icon() {
      return h(NIcon, null, {
        default: () => h(LogOutOutline)
      })
    },
  }
])

const toggleTheme = () => {
  ;(window as any).toggleTheme()
}

const handleMenuSelect = (key: string) => {
  router.push({ name: key.charAt(0).toUpperCase() + key.slice(1) })
}

const handleUserSelect = (key: string) => {
  if (key === 'toggle-theme') {
    toggleTheme()
  } else if (key === 'logout') {
    userStore.logout()
    router.push({ name: 'Login' })
  }
}

const setupRouterGuards = () => {
  router.beforeEach((to, from, next) => {
    loadingBar.start()
    next()
  })

  router.afterEach(() => {
    loadingBar.finish()
  })
}

onMounted(() => {
  setupRouterGuards()
})
</script>

<style scoped>
.header {
  border-bottom: 1px solid var(--n-border-color);
  background-color: transparent;
  backdrop-filter: blur(10px);
  align-items: center;
  position: fixed;
  display: flex;
  height: 50px;
  z-index: 99;
  width: 100%;
  left: 0;
  top: 0;
}

.header-content {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
}

.logo {
  font-size: 20px;
  font-weight: bold;
  color: #18a058;
  white-space: nowrap;
}

.menu-section {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
}

.menu-item {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  color: #666;
}

.menu-item:hover {
  background-color: var(--n-color-target);
}

.menu-item.active {
  color: #18a058;
  background-color: var(--n-color-target);
}

.content {
  max-width: 1200px;
  padding-top: 50px;
  margin: 0 auto;
}

.user-section {
  flex-shrink: 0;
}

/* 页面切换过渡动画 */
.page-enter-active,
.page-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.page-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

.page-enter-to,
.page-leave-from {
  opacity: 1;
  transform: translateY(0);
}
</style>
