<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="sub-title">应用容器</div>
            <div class="flex align-center">
                <n-button quaternary circle type="primary" class="mr-10" @click="getList">
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
                    <n-data-table striped size="small" :row-key="rowKey" :bordered="false" :columns="columns"
                        :data="list" />
                </n-scrollbar>
            </div>
        </div>
    </div>
</template>
  
<script>
import { container } from "../../plugins/api";
import { NIcon, NTag, NButton, NTime } from "naive-ui";
import { ArrowSync20Filled, CheckmarkCircle12Filled, ErrorCircle12Filled, DismissCircle12Filled, ArrowSyncCircle24Filled } from '@vicons/fluent';
export default {
    name: "Container",
    components: { ArrowSync20Filled },
    data: () => ({
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
                            text: true
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
                    let ip1 = row1.network.length == 0 ? '-':row1.network[0].ip;
                    let ip2 = row2.network.length == 0 ? '-':row2.network[0].ip;
                    if(ip1 == ip2) return 0;
                    return ip1 > ip2 ? 1:-1;
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
                                class: 'taxt-small'
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
        ],
        list: [],
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
                }
            }).catch(err => {
                console.log(err)
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
        }
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
