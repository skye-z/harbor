<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div>
                <div class="sub-title">{{ name }}</div>
                <div class="text-small">{{ image.name }}</div>
            </div>
            <div class="flex align-center">
                <n-button quaternary circle type="primary" class="mr-10" :loading="loading" @click="getInfo">
                    <template #icon>
                        <n-icon>
                            <ArrowSync20Filled />
                        </n-icon>
                    </template>
                </n-button>
            </div>
        </div>
    </div>
</template>
  
<script>
import { container } from "../../plugins/api";
import { ArrowSync20Filled } from '@vicons/fluent';

export default {
    name: "ContainerLog",
    components: { ArrowSync20Filled },
    data: () => ({
        loading: true,
        id: '',
        name: '正在加载中',
        image: {
            id: '',
            name: 'loading...'
        }
    }),
    methods: {
        init() {
            this.id = this.$route.params.id;
            this.getInfo();
        },
        getInfo() {
            this.loading = true;
            container.getInfo(this.id).then(res => {
                if (res.state) {
                    let item = res.data
                    if (item.Name != undefined && item.Name != null) this.name = item.Name.substring(1);
                    this.image = {
                        id: item.Image,
                        name: item.Config.Image,
                    }
                }
                this.getLogs();
            }).catch(err => {
                console.log(err)
                setTimeout(() => {
                    this.loading = false;
                }, 500);
            })
        },
        getLogs() {
            setTimeout(() => {
                this.loading = false;
            }, 500);
        }
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped>
.sub-title {
    margin-bottom: 3px;
}
</style>
  
