<template>
  <div class="login-container">
    <div class="login-background"></div>
    <div class="login-content">
      <div class="login-header">
        <div class="logo">
          <img src="../assets/icon/icon-dark@2x.png" />
        </div>
        <div>
          <div style="font-size: 32px;font-weight: bold;line-height: 32px;">Harbor</div>
          <div style="font-size: 12px;color: #999;">轻量级容器管理平台</div>
        </div>
      </div>
      <n-card class="login-card">
        <n-form ref="formRef" :model="formValue" :show-label="false">
          <n-form-item path="username">
            <n-input v-model:value="formValue.username" placeholder="用户名" size="large">
              <template #prefix>
                <n-icon>
                  <User />
                </n-icon>
              </template>
            </n-input>
          </n-form-item>
          <n-form-item path="password">
            <n-input v-model:value="formValue.password" type="password" placeholder="密码" @keydown.enter="handleLogin" size="large">
              <template #prefix>
                <n-icon>
                  <Lock />
                </n-icon>
              </template>
            </n-input>
          </n-form-item>
        </n-form>
        <n-button type="primary" @click="handleLogin" :loading="loading" block size="large">
          登录
        </n-button>
      </n-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { authApi } from '../plugins/api'
import { useUserStore } from '../plugins/stores/user'
import { User,Lock } from '@vicons/fa'
import MD5 from 'crypto-js/md5'

const router = useRouter()
const message = useMessage()
const userStore = useUserStore()
const loading = ref(false)

const formValue = ref({
  username: '',
  password: ''
})

const handleLogin = async () => {
  if (!formValue.value.username || !formValue.value.password) {
    message.error('请输入用户名和密码')
    return
  }

  loading.value = true
  try {
    const encryptedPassword = MD5(formValue.value.password).toString()
    const response = await authApi.login({
      username: formValue.value.username,
      password: encryptedPassword
    })

    if (response && response.token) {
      userStore.setUserFromLogin(response)

      await nextTick()
      message.success('登录成功')
      router.replace({ name: 'Dashboard' })
    } else {
      throw new Error('响应数据格式错误')
    }
  } catch (error: any) {
    message.error('登录失败: ' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  overflow: hidden;
}

.login-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
}

.login-background::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(99, 102, 241, 0.1) 0%, transparent 50%);
  animation: rotate 20s linear infinite;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
}

.login-content {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 40px;
}

.login-header {
  color: white;
  display: flex;
  align-items: center;
}

.logo {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  justify-content: center;
  border-radius: 20px;
  align-items: center;
  margin-right: 10px;
  font-weight: bold;
  font-size: 40px;
  color: white;
  padding: 10px;
  display: flex;
  height: 50px;
  width: 50px;
}

.logo img {
  height: 50px;
  width: 50px;
}

.login-card {
  width: 360px;
  padding: 10px;
  backdrop-filter: blur(10px);
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  transition: background-color 0.3s ease, box-shadow 0.3s ease;
}
</style>
