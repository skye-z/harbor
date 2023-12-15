<template>
    <div class="pa-10">
        <div>Container</div>
        <div>容器列表、创建、删除、编辑、重命名、停止、关闭、暂停、恢复、接入网络、资源限制</div>
    </div>
</template>
  
<script>
export default {
    name: "Container",
    data: () => ({
    }),
    methods: {
        init() {
            console.log('[Init] Container')
        },
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
