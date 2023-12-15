<template>
  <n-config-provider :locale="i18n.main" :date-locale="i18n.date" :theme-overrides="theme">
    <n-dialog-provider>
      <n-message-provider>
        <global-api />
        <n-back-top :right="10" />
      </n-message-provider>
    </n-dialog-provider>
    <div id="app-bg"></div>
    <div id="app-center">
      <n-scrollbar @scroll="listenScroll" style="max-height: 100vh">
        <router-view />
      </n-scrollbar>
    </div>
  </n-config-provider>
</template>

<script>
import { NThemeEditor, zhCN, dateZhCN } from 'naive-ui'
import GlobalApi from './components/globalApi.vue'
import theme from './theme.json'

export default {
  name: "App",
  components: { GlobalApi },
  data: () => ({
    theme,
    i18n: {
      main: zhCN,
      date: dateZhCN
    },
    scrollPro: 0
  }),
  methods: {
    listenScroll(e) {
      let steps = parseInt((e.target.clientHeight * 0.1).toFixed(0))
      let pro = parseInt((e.target.scrollTop / steps).toFixed(0))
      if (pro != this.scrollPro) {
        this.scrollPro = pro
        window.dispatchEvent(new CustomEvent('scroll-pro', { detail: { number: pro } }));
      }
    }
  }
};
</script>