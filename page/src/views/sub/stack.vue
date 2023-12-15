<template>
    <div class="pa-10">
        <div>Stack</div>
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
  
