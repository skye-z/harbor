<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="sub-title">镜像仓库</div>
            <div class="flex align-center">
                <n-button quaternary circle type="primary" class="mr-10" :loading="loading" @click="getList">
                    <template #icon>
                        <n-icon>
                            <ArrowSync20Filled />
                        </n-icon>
                    </template>
                </n-button>
                <n-button strong secondary type="primary" class="mr-10">创建镜像</n-button>
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
import { image } from "../../plugins/api";
import { NTag, NButton, NTime } from "naive-ui";
import { ArrowSync20Filled } from '@vicons/fluent';

export default {
    name: "Image",
    components: { ArrowSync20Filled },
    data: () => ({
        columns: [
            {
                type: 'selection'
            },
            {
                title: "编号",
                key: "id",
                width: '200px',
                ellipsis: true,
                sorter: 'default',
                render(row) {
                    return h(
                        NButton,
                        {
                            text: true
                        },
                        {
                            default: () => {
                                return row.id.substring(0, 20) + '...'
                            }
                        }
                    );
                }
            },
            {
                title: "标签",
                key: "tags",
                render(row) {
                    let list = [];
                    for (let key in row.tags) {
                        list.push(h(
                            NTag,
                            {
                                style: {
                                    marginRight: '6px'
                                },
                                type: "info",
                                bordered: false,
                                round: true
                            },
                            {
                                default: () => row.tags[key]
                            }
                        ));
                    }
                    return list;
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
            console.log('[Init] Image')
            this.getList();
        },
        getList() {
            this.loading = true
            image.getList().then(res => {
                if (res.state) {
                    let list = []
                    for (let i in res.data) {
                        let item = res.data[i]
                        list.push({
                            id: item.Id,
                            tags: item.RepoTags,
                            size: item.Size,
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
        }
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
