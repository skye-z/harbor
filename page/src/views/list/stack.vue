<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="sub-title">服务编排</div>
            <n-button strong secondary type="primary">创建服务</n-button>
        </div>
        <div>编排列表、部署、删除、更新</div>
    </div>
</template>
  
<script>
export default {
    name: "Stack",
    data: () => ({
    }),
    methods: {
        init() {
            console.log('[Init] Stack')
        },
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
