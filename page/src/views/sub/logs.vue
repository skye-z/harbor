<template>
    <div class="pa-10">
        <div>Logs</div>
        <div>平台与docker操作日志</div>
    </div>
</template>
  
<script>
export default {
    name: "Logs",
    data: () => ({
    }),
    methods: {
        init() {
            console.log('[Init] Logs')
        },
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
