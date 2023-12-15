<template>
    <div class="pa-10">
        <div>User</div>
        <div>用户列表、创建、编辑、删除、权限</div>
    </div>
</template>
  
<script>
export default {
    name: "User",
    data: () => ({
    }),
    methods: {
        init() {
            console.log('[Init] User')
        },
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
