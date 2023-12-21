<template>
    <div class="no-select">
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="sub-title">控制台</div>
        </div>
        <div class="pa-10">
            <div v-if="info.name">
                <div class="card mb-10 pa-10">
                    <div class="server-name">{{ info.name }}</div>
                    <div class="text-gray text-small line1">{{ info.docker.id }}</div>
                    <div>{{ info.cpu }}核 {{ info.memory }}</div>
                </div>
                <div class="card mb-10 pa-10">
                    <n-grid :x-gap="10" :y-gap="10" cols="1 300:2 600:4" class="text-center">
                        <n-grid-item class="pt-10">
                            <n-progress type="dashboard" gap-position="bottom" :stroke-width="20" :percentage="80">
                                <div>
                                    <div class="text-big">80%</div>
                                    <div class="text-small text-gray">CPU</div>
                                </div>
                            </n-progress>
                        </n-grid-item>
                        <n-grid-item class="pt-10">
                            <n-progress type="dashboard" gap-position="bottom" :stroke-width="20" :percentage="80">
                                <div>
                                    <div class="text-big">80%</div>
                                    <div class="text-small text-gray">物理内存</div>
                                </div>
                            </n-progress>
                        </n-grid-item>
                        <n-grid-item class="pt-10">
                            <n-progress type="dashboard" gap-position="bottom" :stroke-width="20" :percentage="80">
                                <div>
                                    <div class="text-big">80%</div>
                                    <div class="text-small text-gray">虚拟内存</div>
                                </div>
                            </n-progress>
                        </n-grid-item>
                        <n-grid-item class="pt-10">
                            <n-progress type="dashboard" gap-position="bottom" :stroke-width="20" :percentage="80">
                                <div>
                                    <div class="text-big">80%</div>
                                    <div class="text-small text-gray">磁盘</div>
                                </div>
                            </n-progress>
                        </n-grid-item>
                    </n-grid>
                </div>
                <n-grid :x-gap="10" :y-gap="10" cols="1 600:3 1200:6">
                    <n-grid-item>
                        <div class="card pa-10">
                            <div class="item-title">引擎</div>
                            <div class="item-bar flex align-center">
                                <div class="mr-10">v{{ info.docker.version }}</div>
                                <div>{{ info.docker.runtime }}</div>
                            </div>
                        </div>
                    </n-grid-item>
                    <n-grid-item>
                        <div class="card pa-10">
                            <div class="item-number float-right">{{ info.containers.total }}</div>
                            <div class="item-title">容器</div>
                            <div class="item-bar flex align-center">
                                <div class="flex align-center mr-20">
                                    <n-icon size="14" color="#008000" class="mr-5">
                                        <Play32Filled />
                                    </n-icon>
                                    <div>{{ info.containers.running }}</div>
                                </div>
                                <div class="flex align-center mr-20">
                                    <n-icon size="14" color="#ffa500" class="mr-5">
                                        <Pause48Filled />
                                    </n-icon>
                                    <div>{{ info.containers.paused }}</div>
                                </div>
                                <div class="flex align-center">
                                    <n-icon size="14" color="#ff0000" class="mr-5">
                                        <RecordStop32Filled />
                                    </n-icon>
                                    <div>{{ info.containers.stopped }}</div>
                                </div>
                            </div>
                        </div>
                    </n-grid-item>
                    <n-grid-item>
                        <div class="card pa-10">
                            <div class="item-number float-right">{{ info.image.number }}</div>
                            <div class="item-title">镜像</div>
                            <div class="item-bar flex align-center">
                                <n-icon size="14" color="#999" class="mr-5">
                                    <LinkSquare24Filled />
                                </n-icon>
                                <div>{{ info.image.mirror[0] }}</div>
                            </div>
                        </div>
                    </n-grid-item>
                    <n-grid-item>
                        <div class="card pa-10">
                            <div class="item-number float-right">{{ info.network.number }}</div>
                            <div class="item-title">网络</div>
                            <div class="item-bar flex align-center">
                                <div class="mr-10">IPTables </div>
                                <div class="flex align-center mr-10">
                                    <n-icon size="14" color="#008000" class="mr-5" v-if="info.network.ipv4">
                                        <CheckmarkCircle12Filled />
                                    </n-icon>
                                    <n-icon size="14" color="#ff0000" class="mr-5" v-else>
                                        <DismissCircle12Filled />
                                    </n-icon>
                                    <div>IPv4</div>
                                </div>
                                <div class="flex align-center">
                                    <n-icon size="14" color="#008000" class="mr-5" v-if="info.network.ipv6">
                                        <CheckmarkCircle12Filled />
                                    </n-icon>
                                    <n-icon size="14" color="#ff0000" class="mr-5" v-else>
                                        <DismissCircle12Filled />
                                    </n-icon>
                                    <div>IPv6</div>
                                </div>
                            </div>
                        </div>
                    </n-grid-item>
                    <n-grid-item>
                        <div class="card pa-10">
                            <div class="item-number float-right">{{ info.volume.number }}</div>
                            <div class="item-title">存储</div>
                            <div class="item-bar text-gray">{{ info.volume.root }}</div>
                        </div>
                    </n-grid-item>
                    <n-grid-item>
                        <div class="card pa-10">
                            <div class="item-number float-right">{{ info.user.number }}</div>
                            <div class="item-title">用户</div>
                            <div class="item-bar flex align-center">
                                <div class="flex align-center mr-20">
                                    <n-icon size="14" color="#008000" class="mr-5" v-if="info.user.totp">
                                        <CheckmarkCircle12Filled />
                                    </n-icon>
                                    <n-icon size="14" color="#ff0000" class="mr-5" v-else>
                                        <DismissCircle12Filled />
                                    </n-icon>
                                    <div>TOTP</div>
                                </div>
                                <div class="flex align-center">
                                    <n-icon size="14" color="#008000" class="mr-5" v-if="info.user.oauth">
                                        <CheckmarkCircle12Filled />
                                    </n-icon>
                                    <n-icon size="14" color="#ff0000" class="mr-5" v-else>
                                        <DismissCircle12Filled />
                                    </n-icon>
                                    <div>OAuth2</div>
                                </div>
                            </div>
                        </div>
                    </n-grid-item>
                </n-grid>
            </div>
            <template v-else>
                <n-result v-if="error" style="margin-top: 24vh;" status="warning" title="无法连接到服务器"
                    description="请检查您的网络与防火墙">
                    <template #footer>
                        <n-button @click="init">重新连接</n-button>
                    </template>
                </n-result>
                <n-space v-else vertical>
                    <n-skeleton height="90px" :sharp="false" />
                    <n-skeleton height="150px" :sharp="false" />
                    <n-grid :x-gap="10" :y-gap="10" cols="1 600:3 1200:6">
                        <n-grid-item>
                            <n-skeleton height="69px" :sharp="false" />
                        </n-grid-item>
                        <n-grid-item>
                            <n-skeleton height="69px" :sharp="false" />
                        </n-grid-item>
                        <n-grid-item>
                            <n-skeleton height="69px" :sharp="false" />
                        </n-grid-item>
                        <n-grid-item>
                            <n-skeleton height="69px" :sharp="false" />
                        </n-grid-item>
                        <n-grid-item>
                            <n-skeleton height="69px" :sharp="false" />
                        </n-grid-item>
                        <n-grid-item>
                            <n-skeleton height="69px" :sharp="false" />
                        </n-grid-item>
                    </n-grid>
                </n-space>
            </template>
        </div>
    </div>
</template>
  
<script>
import { docker } from "../../plugins/api";
import {
    Server24Filled, Play32Filled, Pause48Filled, RecordStop32Filled, LinkSquare24Filled, CheckmarkCircle12Filled, DismissCircle12Filled
} from '@vicons/fluent';

export default {
    name: "Console",
    components: { Server24Filled, Play32Filled, Pause48Filled, RecordStop32Filled, LinkSquare24Filled, CheckmarkCircle12Filled, DismissCircle12Filled },
    data: () => ({
        info: {
            name: '',
            system: '',
            kernel: '',
            cpu: 0,
            memory: 0,
            containers: {
                total: 0,
                running: 0,
                paused: 0,
                stopped: 0
            },
            docker: {
                id: '',
                version: '',
                runtime: ''
            },
            image: {
                number: 0,
                mirror: []
            },
            network: {
                number: 0,
                ipv4: false,
                ipv6: false
            },
            volume: {
                number: 0,
                root: ''
            },
            user: {
                number: 0,
                totp: false,
                oauth: false
            }
        },
        error: false
    }),
    methods: {
        init() {
            console.log('[Init] Console')
            this.error = false;
            this.getDockerInfo();
        },
        getDockerInfo() {
            docker.getInfo().then(res => {
                if (res.state) {
                    let data = res.data
                    let mirror = [];
                    for (let i in data.RegistryConfig.Mirrors) {
                        let item = data.RegistryConfig.Mirrors[i];
                        if (item.indexOf('//') !== -1) item = item.substring(item.indexOf('//') + 2);
                        if (item.indexOf('/') !== -1) item = item.substring(0, item.indexOf('/'));
                        mirror.push(item)
                    }
                    this.info = {
                        name: data.Name,
                        system: data.OperatingSystem,
                        kernel: data.KernelVersion,
                        cpu: data.NCPU,
                        memory: this.getMemory(data.MemTotal),
                        containers: {
                            total: data.Containers,
                            running: data.ContainersRunning,
                            paused: data.ContainersPaused,
                            stopped: data.ContainersStopped
                        },
                        docker: {
                            id: data.ID,
                            version: data.ServerVersion,
                            runtime: data.DefaultRuntime.toUpperCase()
                        },
                        image: {
                            number: data.Images,
                            mirror: mirror
                        },
                        network: {
                            number: 0,
                            ipv4: data.BridgeNfIptables,
                            ipv6: data.BridgeNfIp6tables
                        },
                        volume: {
                            number: 0,
                            root: data.DockerRootDir
                        },
                        user: {
                            number: 0,
                            totp: false,
                            oauth: false
                        }
                    }
                }
            }).catch(err => {
                console.log(err)
                setTimeout(() => {
                    this.error = true
                }, 1000)
            })
        },
        getMemory(num) {
            let unit = 'KB'
            num = num / 1024
            if (num >= 1024) {
                num = num / 1024
                unit = 'MB'
            }
            if (num >= 1024) {
                num = num / 1024
                unit = 'GB'
            }
            return num.toFixed(2) + unit;
        }
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

.server-name {
    line-height: 26px;
    font-size: 24px;
}

.item-title {
    font-size: 20px;
    line-height: 20px;
    margin-bottom: 15px;
}

.item-number {
    border-radius: 0 8px 0 20px;
    padding: 8px 5px 10px 10px;
    background-color: #999;
    margin-right: -10px;
    text-align: center;
    margin-top: -10px;
    line-height: 20px;
    font-weight: 900;
    font-size: 32px;
    color: #fff;
    width: 65px;
}

.item-bar {
    line-height: 14px;
}
</style>
  
