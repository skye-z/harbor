<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div>
                <div class="sub-title">{{ name }}</div>
                <div class="text-small">{{ image.name }}</div>
            </div>
            <div class="flex align-center">
                <n-button-group class="mr-10">
                    <n-popover placement="bottom" trigger="hover">
                        <template #trigger>
                            <n-button tertiary circle type="success" @click="control('start')">
                                <template #icon>
                                    <n-icon>
                                        <Play12Regular />
                                    </n-icon>
                                </template>
                            </n-button>
                        </template>
                        <span>启动</span>
                    </n-popover>
                    <n-popover placement="bottom" trigger="hover">
                        <template #trigger>
                            <n-button tertiary circle type="error" @click="control('stop')">
                                <template #icon>
                                    <n-icon>
                                        <RecordStop12Regular />
                                    </n-icon>
                                </template>
                            </n-button>
                        </template>
                        <span>停止</span>
                    </n-popover>
                    <n-popover placement="bottom" trigger="hover">
                        <template #trigger>
                            <n-button tertiary circle type="warning" @click="control('pause')">
                                <template #icon>
                                    <n-icon>
                                        <Pause12Regular />
                                    </n-icon>
                                </template>
                            </n-button>
                        </template>
                        <span>暂停</span>
                    </n-popover>
                    <n-popover placement="bottom" trigger="hover">
                        <template #trigger>
                            <n-button tertiary circle type="success" @click="control('unpause')">
                                <template #icon>
                                    <n-icon>
                                        <Replay20Filled />
                                    </n-icon>
                                </template>
                            </n-button>
                        </template>
                        <span>恢复</span>
                    </n-popover>
                    <n-popover placement="bottom" trigger="hover">
                        <template #trigger>
                            <n-button tertiary circle type="warning" @click="control('restart')">
                                <template #icon>
                                    <n-icon>
                                        <ArrowReset24Filled />
                                    </n-icon>
                                </template>
                            </n-button>
                        </template>
                        <span>重启</span>
                    </n-popover>
                    <n-popover placement="bottom" trigger="hover">
                        <template #trigger>
                            <n-button tertiary circle type="error" @click="control('kill')">
                                <template #icon>
                                    <n-icon>
                                        <Power24Filled />
                                    </n-icon>
                                </template>
                            </n-button>
                        </template>
                        <span>强停</span>
                    </n-popover>
                    <n-popover placement="bottom" trigger="hover">
                        <template #trigger>
                            <n-button strong secondary circle type="error" @click="control('remove')">
                                <template #icon>
                                    <n-icon>
                                        <Delete16Regular />
                                    </n-icon>
                                </template>
                            </n-button>
                        </template>
                        <span>删除</span>
                    </n-popover>
                </n-button-group>
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
import { Play12Regular, RecordStop12Regular, Pause12Regular, Replay20Filled, Power24Filled, Delete16Regular, ArrowReset24Filled, ArrowSync20Filled } from '@vicons/fluent';

export default {
    name: "ContainerInfo",
    components: { Play12Regular, RecordStop12Regular, Pause12Regular, Replay20Filled, Power24Filled, Delete16Regular, ArrowReset24Filled, ArrowSync20Filled },
    data: () => ({
        id: '',
        name: '正在加载中',
        image: {
            id: '',
            name: 'loading...'
        },
        state: {},
        ports: {},
        network: {},
        project: '',
        created: '',
    }),
    methods: {
        init() {
            this.id = this.$route.params.id;
            this.getInfo();
        },
        getInfo() {
            container.getInfo(this.id).then(res => {
                if (res.state) {
                    let item = res.data
                    if (item.Name != undefined && item.Name != null) this.name = item.Name.substring(1);
                    this.image = {
                        id: item.Image,
                        name: item.Config.Image,
                    }
                    this.state = item.State
                    this.network = item.NetworkSettings
                    this.ports = item.Config.ExposedPorts
                    this.created = new Date(item.Created).getTime()
                    if (item.Config && item.Config['com.docker.compose.project']) this.project = item.Config.Labels['com.docker.compose.project']
                }
            }).catch(err => {
                console.log(err)
            })
        },
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped>
.sub-title{
    margin-bottom: 3px;
}
</style>
  
