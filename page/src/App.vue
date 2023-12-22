<template>
  <n-config-provider :theme="lightTheme" :locale="i18n.main" :date-locale="i18n.date" :theme-overrides="theme">
    <n-dialog-provider>
      <n-message-provider>
        <n-notification-provider>
          <global-api />
          <n-back-top :right="10" />
        </n-notification-provider>
      </n-message-provider>
    </n-dialog-provider>
    <div id="app-bg"></div>
    <div id="app-center">
      <router-view />
    </div>
  </n-config-provider>
</template>

<script>
import { lightTheme, zhCN, dateZhCN } from 'naive-ui'
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
    }
  }),
  methods: {
  }
};
</script>