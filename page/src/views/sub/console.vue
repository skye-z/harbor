<template>
    <div class="pa-10">
        <div>Console</div>
        <div>本机性能监控(docker部署需安装代理)、docker信息、容器数量、镜像数量、网络数量、存储卷数量、用户数量</div>
    </div>
</template>
  
<script>
export default {
    name: "Console",
    data: () => ({
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
  
<style scoped></style>
  
