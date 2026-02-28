<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { darkTheme as naiveDarkTheme, type GlobalThemeOverrides } from 'naive-ui'
import GlobalAPI from './components/GlobalAPI.vue'

// 导入 highlight.js 核心和语言模块
import hljs from 'highlight.js/lib/core'
import json from 'highlight.js/lib/languages/json'
import bash from 'highlight.js/lib/languages/bash'
import plaintext from 'highlight.js/lib/languages/plaintext'

// 注册语言
hljs.registerLanguage('json', json)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('plaintext', plaintext)

const getInitialTheme = () => {
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark' || savedTheme === 'light') {
    return savedTheme === 'dark'
  }

  const prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
  return prefersDark
}

const isDark = ref(getInitialTheme())

const theme = computed(() => isDark.value ? naiveDarkTheme : undefined)

const lightThemeOverrides: GlobalThemeOverrides = {
  common: {
    bodyColor: '#f5f5f5',
    borderRadius: '8px',
    borderRadiusSmall: '6px'
  }
}

const darkThemeOverrides: GlobalThemeOverrides = {
  common: {
    bodyColor: '#101014',
    borderRadius: '8px',
    borderRadiusSmall: '6px'
  }
}

const toggleTheme = () => {
  isDark.value = !isDark.value
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
  document.documentElement.setAttribute('data-theme', isDark.value ? 'dark' : 'light')
}

if (typeof window !== 'undefined') {
  (window as any).toggleTheme = toggleTheme
}

onMounted(() => {
  document.documentElement.setAttribute('data-theme', isDark.value ? 'dark' : 'light')
})

watch(isDark, (val) => {
  localStorage.setItem('theme', val ? 'dark' : 'light')
  document.documentElement.setAttribute('data-theme', val ? 'dark' : 'light')
})
</script>

<template>
  <n-config-provider :theme="theme" :theme-overrides="isDark ? darkThemeOverrides : lightThemeOverrides" :hljs="hljs">
    <n-loading-bar-provider>
      <n-message-provider>
        <n-dialog-provider>
          <n-notification-provider>
            <GlobalAPI />
            <router-view />
          </n-notification-provider>
        </n-dialog-provider>
      </n-message-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<style>
@import 'vfonts/FiraCode.css';

html {
  transition: background-color 0.3s, color 0.3s;
}

body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-touch-callout: none;
  -webkit-user-select: none;
  -khtml-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
}

[data-theme='light'] body {
  background-color: #f5f5f5;
  color: #333;
}

[data-theme='dark'] body {
  background-color: #101014;
  color: #fff;
}

.n-card-header{
  padding: 10px !important;
}

.n-card__action{
  padding: 5px !important;
}
</style>
