<template>
    <div class="pa-10">
        <div>Monitor</div>
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
  
