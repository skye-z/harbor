<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="sub-title">容器监控</div>
            <n-button strong secondary type="primary">创建监控</n-button>
        </div>
        <div>选定容器记录性能指标、指标图表、告警配置</div>
    </div>
</template>
  
<script>
export default {
    name: "Monitor",
    data: () => ({
    }),
    methods: {
        init() {
            console.log('[Init] Monitor')
        },
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
