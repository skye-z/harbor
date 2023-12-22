<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="sub-title">存储卷</div>
            <div class="flex align-center">
                <n-button quaternary circle type="primary" class="mr-10" :loading="loading" @click="getList">
                    <template #icon>
                        <n-icon>
                            <ArrowSync20Filled />
                        </n-icon>
                    </template>
                </n-button>
                <n-button strong secondary type="primary">创建卷</n-button>
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
import { volume } from "../../plugins/api";
import { NTag, NButton, NTime } from "naive-ui";
import { ArrowSync20Filled } from '@vicons/fluent';

export default {
    name: "Volume",
    components: { ArrowSync20Filled },
    data: () => ({
        columns: [
            {
                type: 'selection'
            },
            {
                title: "名称",
                key: "name",
                width: '200px',
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
                                return row.project
                            }
                        }
                    );
                }
            },
            {
                title: "挂载点",
                key: "mount",
                ellipsis: {
                    tooltip: true
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
            console.log('[Init] Volume')
            this.getList();
        },
        getList() {
            this.loading = true
            volume.getList().then(res => {
                if (res.state) {
                    let list = []
                    for (let i in res.data) {
                        let item = res.data[i]
                        list.push({
                            id: item.Id,
                            name: item.Name,
                            mount: item.Mountpoint,
                            project: this.getProject(item.Labels),
                            created: new Date(item.CreatedAt).getTime()
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
  
