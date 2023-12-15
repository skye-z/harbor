<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="sub-title">内部网络</div>
            <n-button strong secondary type="primary">创建网络</n-button>
        </div>
        <div>网络列表、创建、删除、连接、断开、编辑</div>
    </div>
</template>
  
<script>
export default {
    name: "Network",
    data: () => ({
    }),
    methods: {
        init() {
            console.log('[Init] Network')
        },
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
