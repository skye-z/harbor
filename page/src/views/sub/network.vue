<template>
    <div class="pa-10">
        <div>Network</div>
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
  
