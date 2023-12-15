<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="sub-title">用户管理</div>
            <n-button strong secondary type="primary">创建用户</n-button>
        </div>
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
  
