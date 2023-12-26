<template>
    <div>
        <div class="sub-header pa-10 border-bottom flex align-center justify-between">
            <div class="flex align-center">
                <n-button quaternary circle class="mr-10" @click="back">
                    <template #icon>
                        <n-icon size="24" color="#999">
                            <ChevronLeft24Filled />
                        </n-icon>
                    </template>
                </n-button>
                <div>
                    <div class="sub-title">{{ name }}</div>
                    <div class="text-small">{{ image.name }}</div>
                </div>
            </div>
            <div class="flex align-center">
                <n-select style="width: 100px;" class="mr-10" v-model:value="line" :options="options.line" />
                <n-select style="width: 80px;" class="mr-10" v-model:value="size" :options="options.size" />
                <n-select style="width: 85px;" class="mr-10" v-model:value="height" :options="options.height" />
                <n-button quaternary circle type="primary" class="mr-10" :loading="loading" @click="getLogs">
                    <template #icon>
                        <n-icon>
                            <ArrowSync20Filled />
                        </n-icon>
                    </template>
                </n-button>
            </div>
        </div>
        <div class="pa-10">
            <div class="card">
                <n-log class="log" ref="logRef" language="verilog" :hljs="highlight" :font-size="size" :line-height="height"
                    :lines="logs" :loading="loading"></n-log>
            </div>
        </div>
    </div>
</template>
  
<script>
import { container } from "../../plugins/api";
import { ArrowSync20Filled, ChevronLeft24Filled } from '@vicons/fluent';
import hljs from 'highlight.js/lib/core'

import verilog from 'highlight.js/lib/languages/verilog'

export default {
    name: "ContainerLog",
    components: { ArrowSync20Filled, ChevronLeft24Filled },
    setup() {
        hljs.registerLanguage('verilog', verilog);
        return { highlight: hljs };
    },
    data: () => ({
        loading: true,
        id: '',
        name: '',
        image: {
            id: '',
            name: ''
        },
        logs: [],
        options: {
            line: [
                {
                    label: "100行",
                    value: 100
                },
                {
                    label: "200行",
                    value: 200
                },
                {
                    label: "400行",
                    value: 400
                },
                {
                    label: "800行",
                    value: 800
                },
                {
                    label: "1600行",
                    value: 1600
                },
                {
                    label: "3200行",
                    value: 3200
                },
            ],
            height: [
                {
                    label: "1.25倍",
                    value: 1.25
                },
                {
                    label: "1.5倍",
                    value: 1.5
                },
                {
                    label: "1.75倍",
                    value: 1.75
                },
                {
                    label: "2.0倍",
                    value: 2
                }
            ],
            size: [
                {
                    label: "12px",
                    value: 12
                },
                {
                    label: "14px",
                    value: 14
                },
                {
                    label: "18px",
                    value: 18
                },
                {
                    label: "20px",
                    value: 20
                },
                {
                    label: "24px",
                    value: 24
                },
                {
                    label: "28px",
                    value: 28
                },
                {
                    label: "32px",
                    value: 32
                },
            ]
        },
        height: 1.25,
        line: 100,
        size: 14
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
            this.loading = true;
            container.getLogs(this.id, this.line).then(res => {
                if (res.state) {
                    this.logs = res.data
                }
                setTimeout(() => {
                    this.scrollTo('bottom', false)
                }, 200);
                setTimeout(() => {
                    this.loading = false;
                }, 500);
            }).catch(err => {
                console.log(err)
                this.loading = false;
            })
        },
        scrollTo(position, silent) {
            this.$refs.logRef.scrollTo({ position, silent })
        },
        back() {
            this.$router.push('/container/' + this.id)
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

.log {
    height: calc(100vh - 80px) !important;
}

.log:deep(.n-scrollbar-content) {
    padding: 10px;
}
</style>
  
