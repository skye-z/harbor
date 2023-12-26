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
                            <n-button tertiary circle type="success"
                                :disabled="state.Status === 'exited' || state.Status === 'stop' || state.Status === 'dead' || state.Status === 'paused'"
                                @click="control('start')">
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
                            <n-button tertiary circle type="error"
                                :disabled="state.Status === 'exited' || state.Status === 'stop' || state.Status === 'dead'"
                                @click="control('stop')">
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
                            <n-button tertiary circle type="warning" :disabled="state.Status === 'paused'"
                                @click="control('pause')">
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
                            <n-button tertiary circle type="success" :disabled="state.Status != 'paused'"
                                @click="control('unpause')">
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
                            <n-button tertiary circle type="error"
                                :disabled="state.Status === 'exited' || state.Status === 'stop' || state.Status === 'dead'"
                                @click="control('kill')">
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
        <div class="loading flex align-center justify-center" v-if="loading">
            <n-spin />
        </div>
        <div v-else class="sub-body">
            <n-grid :x-gap="10" :y-gap="10" cols="1 800:2">
                <n-grid-item>
                    <div class="card pa-10">
                        <div class="card-title">容器信息</div>
                        <div class="flex align-center">
                            <div class="item-label">编号</div>
                            <div class="item-info line1">{{ id }}</div>
                        </div>
                        <div class="flex align-center">
                            <div class="item-label">目录</div>
                            <div class="item-info line1">{{ path }}</div>
                        </div>
                        <div class="flex align-center">
                            <div class="item-label">状态</div>
                            <div class="flex align-center">
                                <n-tag round :bordered="false" type="success" v-if="state.Status === 'running'">
                                    {{ state.State }}
                                    <template #icon>
                                        <n-icon>
                                            <CheckmarkCircle12Filled />
                                        </n-icon>
                                    </template>
                                </n-tag>
                                <n-tag round :bordered="false" type="error"
                                    v-else-if="state.Status === 'exited' || state.Status === 'stop' || state.Status === 'removing' || state.Status === 'dead'">
                                    {{ state.State }}
                                    <template #icon>
                                        <n-icon>
                                            <DismissCircle12Filled />
                                        </n-icon>
                                    </template>
                                </n-tag>
                                <n-tag round :bordered="false" type="warning" v-else>
                                    {{ state.State }}
                                    <template #icon>
                                        <n-icon>
                                            <AddCircle16Filled v-if="state.Status === 'created'" />
                                            <PauseCircle24Filled v-else-if="state.Status === 'paused'" />
                                            <ArrowSyncCircle24Filled v-else-if="state.Status === 'restarting'" />
                                            <ErrorCircle12Filled v-else />
                                        </n-icon>
                                    </template>
                                </n-tag>
                                <span class="ml-10">PID:{{ state.Pid }}</span>
                            </div>
                        </div>
                        <div class="flex align-center">
                            <div class="item-label">平台</div>
                            <div>{{ platform }} - {{ driver }} - {{ runtime }}</div>
                        </div>
                        <div class="flex align-center">
                            <div class="item-label">创建时间</div>
                            <div>
                                <template v-if="created">
                                    <n-time :time="created" /> (<n-time :time="created" :to="now" type="relative" />)
                                </template>
                                <span v-else>-</span>
                            </div>
                        </div>
                        <div class="flex align-center">
                            <div class="item-label">启动时间</div>
                            <div>
                                <template v-if="state && state.StartTime">
                                    <n-time :time="state.StartTime" /> (<n-time :time="state.StartTime" :to="now"
                                        type="relative" />)
                                </template>
                                <span v-else>-</span>
                            </div>
                        </div>
                    </div>
                </n-grid-item>
                <n-grid-item>
                    <div class="card pa-10">
                        <div class="card-title">启动信息</div>
                        <div class="flex align-center">
                            <div class="item-label">入口点</div>
                            <div>{{ run.inlet }}</div>
                        </div>
                        <div class="flex align-center">
                            <div class="item-label">工作目录</div>
                            <div>{{ run.work }}</div>
                        </div>
                        <div class="flex align-center">
                            <div class="item-label">启动命令</div>
                            <div class="item-info line1" style="display: -webkit-box;">
                                <template v-for="item in run.cmd" class="mr-5">{{ item }}&nbsp;&nbsp;</template>
                            </div>
                        </div>
                        <div class="flex align-center">
                            <div class="item-label">重启次数</div>
                            <div>{{ run.count }}次</div>
                        </div>
                        <div class="text-small text-gray">*指启动容器中应用服务的次数</div>
                    </div>

                    <div class="tools">
                        <n-button-group class="full-width">
                            <n-button type="info" @click="jump('logs')" strong secondary round>
                                <template #icon>
                                    <n-icon>
                                        <DocumentBulletListClock20Regular />
                                    </n-icon>
                                </template>
                                日志
                            </n-button>
                            <n-button type="info" @click="jump('terminal')" strong secondary>
                                <template #icon>
                                    <n-icon>
                                        <Terminal />
                                    </n-icon>
                                </template>终端</n-button>
                            <n-button type="info" @click="jump('stat')" strong secondary>
                                <template #icon>
                                    <n-icon>
                                        <ChartMultiple24Regular />
                                    </n-icon>
                                </template>统计</n-button>
                            <n-button type="info" @click="jump('edit')" strong secondary>
                                <template #icon>
                                    <n-icon>
                                        <BoxEdit24Regular />
                                    </n-icon>
                                </template>编辑</n-button>
                            <n-button type="info" @click="jump('copy')" strong secondary>
                                <template #icon>
                                    <n-icon>
                                        <CopyAdd24Regular />
                                    </n-icon>
                                </template>克隆</n-button>
                            <n-button type="info" @click="jump('rebuild')" strong secondary round>
                                <template #icon>
                                    <n-icon>
                                        <BranchCompare24Regular />
                                    </n-icon>
                                </template>重建</n-button>
                        </n-button-group>
                    </div>
                </n-grid-item>
            </n-grid>

            <div class="card mt-10 pa-10">
                <div class="card-title">环境变量</div>
                <div class="flex align-center" v-for="(value, key) in env">
                    <div class="env-key line1">{{ key }}</div>
                    <div class="env-value line1">{{ value }}</div>
                </div>
            </div>
            <div class="card mt-10 pa-10" v-if="labels != null">
                <div class="card-title">容器标签</div>
                <div class="flex align-center" v-for="(value, key) in labels">
                    <div class="label-key line1">{{ key }}</div>
                    <div class="label-value line1">{{ value }}</div>
                </div>
            </div>
            <div class="card mt-10 pa-10">
                <div class="card-title">存储设置</div>
                <div class="flex align-center" v-for="item in volumes">
                    <div class="volume-from line1">{{ item.Source }}</div>
                    <div class="volume-to line1">{{ item.Destination }}</div>
                </div>
            </div>
            <div class="card mt-10 pa-10">
                <div class="card-title">网络设置</div>
                <n-grid :x-gap="10" :y-gap="0" cols="1 800:2">
                    <n-grid-item>
                        <div class="flex align-center" v-if="network.ipv4Gateway">
                            <div class="item-label">IPv4网关</div>
                            <div>{{ network.ipv4Gateway }}</div>
                        </div>
                        <div class="flex align-center" v-if="network.ipv4IPAddress">
                            <div class="item-label">IPv4地址</div>
                            <div>{{ network.ipv4IPAddress }}</div>
                        </div>
                        <div class="flex align-center" v-if="network.ipv6Gateway">
                            <div class="item-label">IPv6网关</div>
                            <div>{{ network.ipv6Gateway }}</div>
                        </div>
                        <div class="flex align-center" v-if="network.ipv6IPAddress">
                            <div class="item-label">IPv6地址</div>
                            <div>{{ network.ipv6IPAddress }}</div>
                        </div>
                        <div class="flex align-center">
                            <div class="item-label">Mac地址</div>
                            <div>{{ network.mac }}</div>
                        </div>
                    </n-grid-item>
                    <n-grid-item>
                        <div class="item-label">端口映射</div>
                        <template class="flex align-center" v-for="(value, key) in ports">
                            <n-button strong secondary type="info" class="mr-10">
                                {{ value }}:{{ key }}
                            </n-button>
                        </template>
                    </n-grid-item>
                </n-grid>
            </div>
        </div>
    </div>
</template>
  
<script>
import { container } from "../../plugins/api";
import { 
    AddCircle16Filled, Play12Regular, RecordStop12Regular, Pause12Regular, 
    Replay20Filled, Power24Filled, Delete16Regular, ArrowReset24Filled, 
    ArrowSync20Filled, CheckmarkCircle12Filled, ErrorCircle12Filled, 
    PauseCircle24Filled, DismissCircle12Filled, ArrowSyncCircle24Filled,
    DocumentBulletListClock20Regular,ChartMultiple24Regular, BoxEdit24Regular,
    CopyAdd24Regular, BranchCompare24Regular
 } from '@vicons/fluent';

import { Terminal } from '@vicons/tabler';

export default {
    name: "ContainerInfo",
    components: {
        AddCircle16Filled, Play12Regular, RecordStop12Regular, Pause12Regular, 
        Replay20Filled, Power24Filled, Delete16Regular, ArrowReset24Filled, 
        ArrowSync20Filled, CheckmarkCircle12Filled, ErrorCircle12Filled, 
        PauseCircle24Filled, DismissCircle12Filled, ArrowSyncCircle24Filled,
        DocumentBulletListClock20Regular,ChartMultiple24Regular, BoxEdit24Regular,
        CopyAdd24Regular, BranchCompare24Regular, Terminal
    },
    data: () => ({
        loading: true,
        id: '',
        now: 0,
        name: '正在加载中',
        path: '',
        image: {
            id: '',
            name: 'loading...'
        },
        state: {},
        ports: {},
        network: {},
        project: '',
        created: '',
        platform: '',
        runtime: '',
        driver: '',
        run: {
            inlet: '',
            work: '',
            count: 0,
            cmd: []
        },
        env: {},
        labels: null,
        volumes: []
    }),
    methods: {
        init() {
            this.id = this.$route.params.id;
            this.getInfo();
        },
        cleanInfo() {
            this.loading = true;
            this.name = '正在加载中';
            this.image = {
                id: '',
                name: 'loading...'
            };
            this.state = {};
            this.ports = {};
            this.network = {};
            this.run = {};
            this.env = {};
            this.labels = null;
            this.volumes = [];
        },
        getInfo() {
            this.cleanInfo();
            this.now = new Date().getTime();
            container.getInfo(this.id).then(res => {
                if (res.state) {
                    let item = res.data
                    if (item.Name != undefined && item.Name != null) this.name = item.Name.substring(1);
                    this.image = {
                        id: item.Image,
                        name: item.Config.Image,
                    }
                    if (item.HostnamePath != undefined && item.HostnamePath != null)
                        this.path = item.HostnamePath.substring(0, item.HostnamePath.lastIndexOf('/'))
                    this.state = item.State;
                    if (this.state.StartedAt) this.state.StartTime = new Date(this.state.StartedAt).getTime()
                    this.state.State = this.getState(this.state.Status)
                    if (item.NetworkSettings && item.NetworkSettings.Networks) {
                        for (let key in item.NetworkSettings.Networks) {
                            let net = item.NetworkSettings.Networks[key]
                            this.network = {
                                ipv4Gateway: net.Gateway,
                                ipv4IPAddress: net.IPAMConfig ? net.IPAMConfig.IPv4Address : net.IPAddress,
                                ipv6Gateway: net.IPv6Gateway,
                                ipv6IPAddress: net.GlobalIPv6Address,
                                mac: net.MacAddress
                            };
                            break;
                        }
                    }
                    if (item.NetworkSettings && item.NetworkSettings.Ports) {
                        let map = {}
                        for (let key in item.NetworkSettings.Ports) {
                            let port = item.NetworkSettings.Ports[key];
                            if (port) map[key.substring(0, key.indexOf('/'))] = port[0].HostPort
                            else map[key.substring(0, key.indexOf('/'))] = 'Privacy'
                        }
                        this.ports = map
                    }
                    this.created = new Date(item.Created).getTime();
                    if (item.Config && item.Config['com.docker.compose.project']) this.project = item.Config.Labels['com.docker.compose.project'];
                    this.platform = item.Platform;
                    this.runtime = item.HostConfig.Runtime;
                    this.driver = item.Driver;
                    this.run = {
                        inlet: item.Config.Entrypoint ? item.Config.Entrypoint[0] : item.Path,
                        work: item.Config.WorkingDir,
                        count: item.RestartCount,
                        cmd: item.Config.Cmd
                    }
                    if (item.Config.Env) {
                        for (let i in item.Config.Env) {
                            let env = item.Config.Env[i]
                            this.env[env.substring(0, env.indexOf('='))] = env.substring(env.indexOf('=') + 1)
                        }
                    }
                    if (item.Config.Labels && Object.keys(item.Config.Labels).length > 0) 
                        this.labels = item.Config.Labels
                    this.volumes = item.Mounts
                }
                setTimeout(() => {
                    this.loading = false;
                }, 500);
            }).catch(err => {
                console.log(err)
                setTimeout(() => {
                    this.loading = false;
                }, 500);
            })
        },
        getState(state) {
            switch (state) {
                case 'created': return '已创建';
                case 'exited': return '已退出';
                case 'running': return '运行中';
                case 'paused': return '已暂停';
                case 'restarting': return '重启中';
                case 'removing': return '移除中';
                case 'dead': return '已死亡';
            }
        },
        control(action) {
            switch (action) {
                case 'start':
                case 'stop':
                case 'restart':
                case 'kill':
                case 'pause':
                case 'unpause':
                    container[action](this.id).then(res => {
                        window.$notify[res.state ? 'success' : 'warning']({
                            meta: this.name,
                            content: res.message,
                            duration: 2500,
                            keepAliveOnHover: true
                        })
                        this.getInfo();
                    }).catch(() => {
                        window.$notify.error({
                            content: "容器控制出错",
                            meta: "容器服务暂不可用",
                            duration: 2500,
                            keepAliveOnHover: true
                        })
                        this.getInfo();
                    })
                    break;
                default:
                    window.$message.warning("无效的控制指令")
                    return false;
            }
        },
        jump(path){
            this.$router.push('/container/'+this.id+'/'+path)
        }
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped>
.sub-header {
    background-color: #e8e8e8;
    width: calc(100vw - 145px);
    position: fixed;
    z-index: 999;
}

.loading {
    height: 100vh;
}

.card {
    padding-bottom: 5px;
}

.sub-title {
    margin-bottom: 3px;
}

.sub-body {
    padding: 66px 10px 10px 10px;
}

.card-title {
    font-size: 18px;
    font-weight: bold;
}

.tools {
    padding: 10px 0 0 0;
}

.tools:deep(.n-button) {
    width: 16.6%;
}

.item-label {
    min-width: 80px;
    padding: 5px 0;
    width: 80px;
}

.item-info {
    max-width: calc(100vw - 270px);
}

.env-key {
    max-width: 250px;
    min-width: 250px;
    padding: 5px 0;
    width: 250px;
}

.label-key {
    max-width: 50%;
    min-width: 280px;
    padding: 5px 0;
    width: 280px;
}

.volume-from {
    max-width: 250px;
    min-width: 250px;
    padding: 5px 0;
    width: 250px;
}
</style>
  
