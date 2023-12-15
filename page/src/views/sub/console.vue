<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="sub-title">控制台</div>
            <n-select id="refresh" v-model:value="refresh" :options="refreshOptions" />
        </div>
        <!-- <div>本机性能监控(docker部署需安装代理)、docker信息、容器数量、镜像数量、网络数量、存储卷数量、用户数量</div> -->
    </div>
</template>
  
<script>
export default {
    name: "Console",
    data: () => ({
        refresh: 0,
        refreshOptions: [
            {
                label: "手动刷新",
                value: 0
            },
            {
                label: "每1秒",
                value: 1
            },
            {
                label: "每3秒",
                value: 3
            },
            {
                label: "每5秒",
                value: 5
            }
        ]
    }),
    methods: {
        init() {
            console.log('[Init] Console')
        },
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped>
#refresh{
    width: 100px;
}
</style>
  
