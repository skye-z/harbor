<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="sub-title">平台设置</div>
            <n-button strong secondary type="primary">Github</n-button>
        </div>
        <div>平台配置</div>
    </div>
</template>
  
<script>
export default {
    name: "Setting",
    data: () => ({
    }),
    methods: {
        init() {
            console.log('[Init] Setting')
        },
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
