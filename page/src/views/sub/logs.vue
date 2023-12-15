<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="sub-title">平台日志</div>
        </div>
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
  
