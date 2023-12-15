<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="sub-title">镜像仓库</div>
            <div>
                <n-button strong secondary type="primary" class="mr-10">创建镜像</n-button>
                <n-button strong secondary type="primary">拉取</n-button>
            </div>
        </div>
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
  
