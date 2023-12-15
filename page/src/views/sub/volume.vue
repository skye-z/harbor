<template>
    <div class="pa-10">
        <div>Volume</div>
        <div>存储卷列表、创建、删除、挂载、卸载</div>
    </div>
</template>
  
<script>
export default {
    name: "Volume",
    data: () => ({
    }),
    methods: {
        init() {
            console.log('[Init] Volume')
        },
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
