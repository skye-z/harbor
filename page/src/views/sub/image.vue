<template>
    <div class="pa-10">
        <div>Image</div>
        <div>搜索远程镜像、列表、拉取、构建、上传、删除、导入、导出</div>
    </div>
</template>
  
<script>
export default {
    name: "Image",
    data: () => ({
    }),
    methods: {
        init() {
            console.log('[Init] Image')
        },
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
