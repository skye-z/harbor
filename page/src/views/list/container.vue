<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="sub-title">应用容器</div>
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
                <n-button quaternary circle type="primary" class="mr-10" :loading="loading" @click="getList">
                    <template #icon>
                        <n-icon>
                            <ArrowSync20Filled />
                        </n-icon>
                    </template>
                </n-button>
                <n-button strong secondary type="primary">创建容器</n-button>
            </div>
        </div>
        <div class="pa-10">
            <div class="card">
                <n-scrollbar x-scrollable
                    style="max-height: calc(100vh - 80px);width: calc(100vw - 165px);border-radius: 8px;">
                    <n-data-table :loading="loading" size="small" :row-key="rowKey" :bordered="false" :columns="columns"
                        :data="list" @update:checked-row-keys="selectCallback" style="min-height: calc(100vh - 80px);" />
                </n-scrollbar>
            </div>
        </div>
    </div>
</template>
  
<script>
import { useRouter } from 'vue-router'
import { container } from "../../plugins/api";
import { NIcon, NTag, NButton, NButtonGroup, NTime } from "naive-ui";
import { 
    Play12Regular, RecordStop12Regular, Pause12Regular, Replay20Filled, 
    Power24Filled, Delete16Regular, AddCircle16Filled, ArrowReset24Filled, 
    ArrowSync20Filled, CheckmarkCircle12Filled, ErrorCircle12Filled, PauseCircle24Filled, 
    DismissCircle12Filled, ArrowSyncCircle24Filled, DocumentBulletListClock20Regular,
    ChartMultiple24Regular
 } from '@vicons/fluent';

import { Terminal } from '@vicons/tabler';

export default {
    name: "Container",
    components: { Play12Regular, RecordStop12Regular, Pause12Regular, Replay20Filled, Power24Filled, Delete16Regular, ArrowReset24Filled, ArrowSync20Filled },
    setup() {
        const router = useRouter()

        const openItem = (id, path) => {
            router.push(`/container/${id}${path}`);
        };

        return {
            columns: [
                {
                    type: 'selection'
                },
                {
                    title: "名称",
                    key: "name",
                    width: '150px',
                    ellipsis: {
                        tooltip: true
                    },
                    sorter: 'default',
                    render(row) {
                        return h(
                            NButton,
                            {
                                text: true,
                                onClick: () => openItem(row.id, '')

                            },
                            {
                                default: () => {
                                    return row.name
                                }
                            }
                        );
                    }
                },
                {
                    title: "状态",
                    key: "state",
                    width: '100px',
                    sorter: 'default',
                    render(row) {
                        let type = 'warning'
                        let icon = ErrorCircle12Filled;
                        if (row.state === 'exited' || row.state === 'stop' || row.state === 'removing' || row.state === 'dead') {
                            type = 'error'
                            icon = DismissCircle12Filled
                        } else if (row.state === 'running') {
                            type = 'success'
                            icon = CheckmarkCircle12Filled
                        } else if (row.state === 'restarting') icon = ArrowSyncCircle24Filled
                        else if (row.state === 'created') icon = AddCircle16Filled
                        else if (row.state === 'paused') icon = PauseCircle24Filled
                        return h(
                            NTag,
                            {
                                style: {
                                    marginRight: '6px'
                                },
                                type: type,
                                bordered: false,
                                round: true
                            },
                            {
                                icon: () => h(NIcon, null, { default: () => h(icon) }),
                                default: () => row.stateZh
                            }
                        )
                    }
                },
                {
                    title: "快捷入口",
                    key: "action",
                    align: 'center',
                    width: '140px',
                    render(row) {
                        return h(
                            NButtonGroup,
                            {
                            },
                            {
                                default: () => {
                                    return [h(
                                        NButton,
                                        {
                                            tertiary: true,
                                            round: true,
                                            size: 'small',
                                            onClick: () => openItem(row.id, '/logs')
                                        },
                                        {
                                            icon: () => h(NIcon, null, { default: () => h(DocumentBulletListClock20Regular) })
                                        }
                                    ), h(
                                        NButton,
                                        {
                                            tertiary: true,
                                            size: 'small',
                                            onClick: () => openItem(row.id, '/terminal')
                                        },
                                        {
                                            icon: () => h(NIcon, null, { default: () => h(Terminal) })
                                        }
                                    ), h(
                                        NButton,
                                        {
                                            tertiary: true,
                                            round: true,
                                            size: 'small',
                                            onClick: () => openItem(row.id, '/stat')
                                        },
                                        {
                                            icon: () => h(NIcon, null, { default: () => h(ChartMultiple24Regular) })
                                        }
                                    )]
                                }
                            }
                        );
                    }
                },
                {
                    title: "项目",
                    key: "project",
                    width: '100px',
                    ellipsis: {
                        tooltip: true
                    },
                    sorter: 'default',
                    render(row) {
                        return h(
                            NButton,
                            {
                                text: true
                            },
                            {
                                default: () => {
                                    return row.project
                                }
                            }
                        );
                    }
                },
                {
                    title: "镜像",
                    key: "image",
                    ellipsis: {
                        tooltip: true
                    },
                    minWidth: '150px',
                    sorter: 'default',
                    render(row) {
                        return h(
                            NButton,
                            {
                                text: true
                            },
                            {
                                default: () => {
                                    return row.image
                                }
                            }
                        );
                    }
                },
                {
                    title: "创建时间",
                    key: "created",
                    width: '170px',
                    sorter: 'default',
                    render(row) {
                        return h(
                            NTime,
                            {
                                time: row.created
                            }, null
                        );
                    }
                },
                {
                    title: "网络",
                    key: "network",
                    width: '140px',
                    sorter: (row1, row2) => {
                        let ip1 = row1.network.length == 0 ? '-' : row1.network[0].ip;
                        let ip2 = row2.network.length == 0 ? '-' : row2.network[0].ip;
                        if (ip1 == ip2) return 0;
                        return ip1 > ip2 ? 1 : -1;
                    },
                    render(row) {
                        if (row.network.length == 0) return "-"
                        return row.network[0].ip;
                    }
                },
                {
                    title: "端口",
                    key: "ports",
                    minWidth: '140px',
                    render(row) {
                        let list = [];
                        for (let key in row.ports) {
                            list.push(h(
                                NButton,
                                {
                                    strong: true,
                                    secondary: true,
                                    size: "small",
                                    style: {
                                        marginRight: "5px"
                                    },
                                    type: "info",
                                    class: 'taxt-small',
                                    onClick: () => {
                                        window.open(window.location.hostname + ":" + key)
                                    }
                                },
                                {
                                    default: () => {
                                        return key + ":" + row.ports[key]
                                    }
                                }
                            ))
                        }
                        return list;
                    }
                },
            ]
        }
    },
    data: () => ({
        list: [],
        select: [],
        loading: true,
        rowKey(row) {
            return row.id
        },
    }),
    methods: {
        init() {
            console.log('[Init] Container')
            this.getList();
        },
        getList() {
            this.loading = true
            container.getList().then(res => {
                if (res.state) {
                    let list = []
                    for (let i in res.data) {
                        let item = res.data[i]
                        list.push({
                            id: item.Id,
                            name: this.getName(item.Names),
                            image: item.Image,
                            state: item.State,
                            status: item.Status,
                            stateZh: this.getState(item.State),
                            ports: this.getPorts(item.Ports),
                            network: this.getIpAddress(item.NetworkSettings),
                            project: this.getProject(item.Labels),
                            created: item.Created * 1000
                        })
                    }
                    this.list = list
                    this.loading = false
                }
            }).catch(err => {
                console.log(err)
                this.loading = false
            })
        },
        getName(names) {
            if (names == undefined || names == null || names.length == 0) return ""
            return names[0].substring(1)
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
        getPorts(ports) {
            let map = {}
            for (let i in ports) {
                let item = ports[i];
                if (!item.PublicPort || map[item.PublicPort]) continue;
                map[item.PublicPort] = item.PrivatePort
            }
            return map;
        },
        getIpAddress(setting) {
            let list = []
            if (setting && setting.Networks) {
                for (let key in setting.Networks) {
                    let item = setting.Networks[key];
                    list.push({
                        name: key,
                        id: item.NetworkID,
                        ip: item.IPAddress,
                        mac: item.MacAddress
                    })
                }
            }
            return list;
        },
        getProject(tags) {
            if (tags && tags['com.docker.compose.project']) return tags['com.docker.compose.project'];
            return ''
        },
        selectCallback(e) {
            this.select = e
        },
        id2Name(id) {
            for (let i in this.list) {
                if (this.list[i].id === id) return this.list[i].name;
            }
            return "未知容器"
        },
        control(action) {
            let num = this.select.length;
            if (num == 0) {
                window.$message.info("请先选择要控制的容器")
                return false;
            }
            let now = 0
            switch (action) {
                case 'start':
                case 'stop':
                case 'restart':
                case 'kill':
                case 'pause':
                case 'unpause':
                    for (let i in this.select) {
                        let id = this.select[i];
                        container[action](id).then(res => {
                            window.$notify[res.state ? 'success' : 'warning']({
                                meta: this.id2Name(id),
                                content: res.message,
                                duration: 2500,
                                keepAliveOnHover: true
                            })
                            now++;
                        }).catch(() => {
                            window.$notify.error({
                                content: "容器控制出错",
                                meta: "容器服务暂不可用",
                                duration: 2500,
                                keepAliveOnHover: true
                            })
                            now++;
                        })
                    }
                    break;
                default:
                    window.$message.warning("无效的控制指令")
                    return false;
            }
            let timer = setInterval(() => {
                if (num == now) {
                    this.getList()
                    clearInterval(timer)
                }
            }, 1000)
        }
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
