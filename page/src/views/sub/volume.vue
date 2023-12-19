<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="sub-title">存储卷</div>
            <n-button strong secondary type="primary">创建卷</n-button>
        </div>
        <div>存储卷列表、创建、删除、挂载、卸载</div>
    </div>
</template>
  
<script>
import { volume } from "../../plugins/api";
export default {
    name: "Volume",
    data: () => ({
    }),
    methods: {
        init() {
            console.log('[Init] Volume')
            this.getList();
        },
        getList(){
            volume.getList().then(res => {
                console.log(res)
            }).catch(err => {
                console.log(err)
            })
        }
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
