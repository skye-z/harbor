<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="sub-title">内部网络</div>
            <div class="flex align-center">
                <n-button quaternary circle type="primary" class="mr-10" :loading="loading" @click="getList">
                    <template #icon>
                        <n-icon>
                            <ArrowSync20Filled />
                        </n-icon>
                    </template>
                </n-button>
                <n-button strong secondary type="primary">创建网络</n-button>
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
import { network } from "../../plugins/api";
import { NTag, NButton, NTime } from "naive-ui";
import { ArrowSync20Filled } from '@vicons/fluent';

export default {
    name: "Network",
    components: { ArrowSync20Filled },
    data: () => ({
        columns: [
            {
                type: 'selection'
            },
            {
                title: "名称",
                key: "name",
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
                title: "IPv4网段",
                key: "ipv4subnet",
                width: '150px',
                ellipsis: {
                    tooltip: true
                },
                render(row) {
                    if (row.ipam && row.ipam.Config && row.ipam.Config[0])
                        return row.ipam.Config[0].Subnet;
                    else return '';
                }
            },
            {
                title: "IPv4网关",
                key: "ipv4gateway",
                width: '130px',
                ellipsis: {
                    tooltip: true
                },
                render(row) {
                    if (row.ipam && row.ipam.Config && row.ipam.Config[0])
                        return row.ipam.Config[0].Gateway;
                    else return '';
                }
            },
            {
                title: "IPv6网段",
                key: "ipv6subnet",
                width: '150px',
                ellipsis: {
                    tooltip: true
                },
                render(row) {
                    if (row.ipam && row.ipam.Config && row.ipam.Config[1])
                        return row.ipam.Config[1].Subnet;
                    else return '';
                }
            },
            {
                title: "IPv6网关",
                key: "ipv6gateway",
                width: '130px',
                ellipsis: {
                    tooltip: true
                },
                render(row) {
                    if (row.ipam && row.ipam.Config && row.ipam.Config[1])
                        return row.ipam.Config[1].Gateway;
                    else return '';
                }
            },
            {
                title: "驱动",
                key: "driver",
                width: '70px',
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
                                return row.driver
                            }
                        }
                    );
                }
            },
            {
                title: "创建时间",
                key: "created",
                width: '180px',
                sorter: 'default',
                render(row) {
                    return h(
                        NTime,
                        {
                            time: row.created
                        }, null
                    );
                }
            }
        ],
        list: [],
        select: [],
        loading: true,
        rowKey(row) {
            return row.id
        },
    }),
    methods: {
        init() {
            console.log('[Init] Network')
            this.getList();
        },
        getList() {
            this.loading = true
            network.getList().then(res => {
                if (res.state) {
                    let list = []
                    for (let i in res.data) {
                        let item = res.data[i]
                        list.push({
                            id: item.Id,
                            name: item.Name,
                            ipam: item.IPAM,
                            scope: item.Scope,
                            driver: item.Driver,
                            ipv6: item.EnableIPv6,
                            internal: item.Internal,
                            project: this.getProject(item.Labels),
                            created: new Date(item.Created).getTime()
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
  
