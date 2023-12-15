<template>
    <div class="full-width flex">
        <div id="app-menu" class="border-right">
            <div class="flex align-center justify-center pa-10 border-bottom">
                <img id="menu-logo" src="../assets/icon/icon-light@1x.png" />
                <div id="menu-title">Harbor</div>
            </div>
            <n-scrollbar style="max-height: calc(100vh - 160px)">
                <n-menu :value="select" :options="menu" @update:value="updateMenu" />
            </n-scrollbar>
            <div class="border-top pa-10 text-center text-gray">v1.0.0 build 123</div>
            <div class="border-top pa-10 text-center">
                <n-button tertiary type="error" class="full-width" @click="exit">
                    <template #icon>
                        <n-icon>
                            <Power />
                        </n-icon>
                    </template>
                    退出登陆
                </n-button>
            </div>
        </div>
        <div id="app-content" class="full-width">
            <n-scrollbar style="max-height: 100vh">
                <router-view />
            </n-scrollbar>
        </div>
    </div>
</template>
  
<script>
import {
    Home, StackedScrolling1, ContainerRegistry, BareMetalServer, Network4,
    VolumeBlockStorage, Blog, CloudMonitoring, User, Settings, Power
} from '@vicons/carbon';
import { NIcon } from "naive-ui";
export default {
    name: "Home",
    components: { Power },
    data: () => ({
        menu: [],
        select: ''
    }),
    methods: {
        init() {
            console.log('[Init] Home')
            this.menu = [
                {
                    label: "控制台",
                    key: "console",
                    icon: this.renderIcon(Home)
                },
                {
                    key: "divider-1",
                    type: "divider"
                },
                {
                    label: "服务编排",
                    key: "stack",
                    icon: this.renderIcon(StackedScrolling1)
                },
                {
                    label: "应用容器",
                    key: "container",
                    icon: this.renderIcon(ContainerRegistry)
                },
                {
                    label: "镜像仓库",
                    key: "image",
                    icon: this.renderIcon(BareMetalServer)
                },
                {
                    label: "内部网络",
                    key: "network",
                    icon: this.renderIcon(Network4)
                },
                {
                    label: "存储卷",
                    key: "volume",
                    icon: this.renderIcon(VolumeBlockStorage)
                },
                {
                    key: "divider-1",
                    type: "divider"
                },
                {
                    label: "容器监控",
                    key: "monitor",
                    icon: this.renderIcon(CloudMonitoring)
                },
                {
                    label: "平台日志",
                    key: "logs",
                    icon: this.renderIcon(Blog)
                },
                {
                    key: "divider-1",
                    type: "divider"
                },
                {
                    label: "用户",
                    key: "user",
                    icon: this.renderIcon(User)
                },
                {
                    label: "设置",
                    key: "setting",
                    icon: this.renderIcon(Settings)
                },
            ]
            this.updateSelect(this.$route.fullPath);
        },
        updateMenu(key, _item) {
            this.$router.push('/' + key);
        },
        renderIcon(icon) {
            return () => h(NIcon, null, { default: () => h(icon) });
        },
        updateSelect(path) {
            if (path === '/') {
                this.select = 'console';
                this.$router.push('/console')
            }
            else {
                path = path.substring(1);
                if (path.indexOf('/') !== -1) path = path.substring(0, path.indexOf('/'));
                this.select = path;
            }
        },
        exit() {
            window.$message.success("登陆已退出, Bye!")
            localStorage.removeItem("user:access:token")
            setTimeout(() => {
                this.$router.push('/auth')
            }, 1000)
        }
    },
    mounted() {
        this.init()
    },
    watch: {
        $route: {
            handler(to) {
                this.updateSelect(to.fullPath);
            },
            deep: true,
        },
    },
};
</script>
  
<style scoped>
#app-menu {
    height: 100vh;
}

#menu-logo {
    margin-left: 10px;
    width: 32px;
}

#menu-title {
    margin: 5px 10px 5.5px 0;
    line-height: 24px;
    font-size: 24px;
}

#app-menu:deep(.n-menu-item-content) {
    padding-left: 15px !important;
}

#app-menu:deep(.n-menu-divider) {
    margin-left: 0 !important;
    margin-right: 0 !important;
}
</style>
  
