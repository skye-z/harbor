<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="sub-title">控制台</div>
            <n-select id="refresh" v-model:value="refresh" :options="refreshOptions" />
        </div>
        <div class="pa-10">
            <div class="card mb-10 pa-10">
                <div>
                    <div>本机名称#123123123</div>
                    <div></div>
                </div>
                <n-grid :x-gap="10" :y-gap="10" cols="3">
                    <n-grid-item>
                        <n-progress type="dashboard" gap-position="bottom" :percentage="80">
                            <div class="text-center">
                                <div class="text-big">10%</div>
                                <div class="text-gray text-small">CPU</div>
                            </div>
                        </n-progress>
                    </n-grid-item>
                    <n-grid-item>
                        <n-progress type="dashboard" gap-position="bottom" :percentage="80">
                            <div class="text-center">
                                <div class="text-big">10%</div>
                                <div class="text-gray text-small">内存</div>
                            </div>
                        </n-progress>
                    </n-grid-item>
                    <n-grid-item>
                        <n-progress type="dashboard" gap-position="bottom" :percentage="80">
                            <div class="text-center">
                                <div class="text-big">10%</div>
                                <div class="text-gray text-small">磁盘</div>
                            </div>
                        </n-progress>
                    </n-grid-item>
                </n-grid>
            </div>
            <n-grid :x-gap="10" :y-gap="10" cols="1 600:3 1200:6">
                <n-grid-item>
                    <div class="card">Docker信息</div>
                </n-grid-item>
                <n-grid-item>
                    <div class="card">容器数量</div>
                </n-grid-item>
                <n-grid-item>
                    <div class="card">镜像数量</div>
                </n-grid-item>
                <n-grid-item>
                    <div class="card">网络数量</div>
                </n-grid-item>
                <n-grid-item>
                    <div class="card">存储卷数量</div>
                </n-grid-item>
                <n-grid-item>
                    <div class="card">用户数量</div>
                </n-grid-item>
            </n-grid>
        </div>
    </div>
</template>
  
<script>
export default {
    name: "Console",
    data: () => ({
        refresh: 0,
        refreshOptions: [
            {
                label: "手动刷新",
                value: 0
            },
            {
                label: "每1秒",
                value: 1
            },
            {
                label: "每3秒",
                value: 3
            },
            {
                label: "每5秒",
                value: 5
            }
        ]
    }),
    methods: {
        init() {
            console.log('[Init] Console')
        },
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped>
#refresh {
    width: 100px;
}
</style>
  
