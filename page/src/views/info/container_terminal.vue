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
        </div>
        <div class="pa-10">
                <div ref="xterm" id="xterm" class="xterm"></div>
        </div>
    </div>
</template>
  
<script>
import { container } from "../../plugins/api";
import { ArrowSync20Filled, ChevronLeft24Filled } from '@vicons/fluent';

import "xterm/css/xterm.css";
import { Terminal } from "xterm";
import { FitAddon } from 'xterm-addon-fit';
import { AttachAddon } from 'xterm-addon-attach';
import { CanvasAddon } from 'xterm-addon-canvas';

export default {
    name: "ContainerTerminal",
    components: { ArrowSync20Filled, ChevronLeft24Filled },
    data: () => ({
        loading: true,
        id: '',
        name: '',
        image: {
            id: '',
            name: ''
        },
        term: null,
        socket: null,
        fitAddon: null,
        socketURI: 'ws://192.168.1.170:12800/api/container/terminal',
        connect: true,
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
                this.initConnect();
            }).catch(err => {
                console.log(err)
                setTimeout(() => {
                    this.loading = false;
                }, 500);
            })
        },
        initConnect() {
            this.term = new Terminal({
                // theme: {
                //     background,
                //     foreground
                // },
                fontSize: 14
            })
            // 加载插件
            this.addPlugins();
            setTimeout(() => {
                // 打开Dom元素
                this.term.open(this.$refs.xterm)
                // 自适应窗口大小
                this.fitAddon.fit()
            }, 100)
            // 创建连接
            this.addSocket()
            // 输入聚焦
            this.term.focus()
            // 加载大小变动事件
            this.addResizeEvent();
        },
        addPlugins() {
            // 加载Canvas渲染
            this.term.loadAddon(new CanvasAddon())
            // 加载窗口自适应插件
            this.fitAddon = new FitAddon()
            this.term.loadAddon(this.fitAddon)
        },
        addSocket() {
            // 创建WebSocket连接
            this.socket = new WebSocket(this.socketURI + '?cols=' + this.term.cols + '&rows' + this.term.rows + '&id=' + this.id)
            // 连接开启事件
            this.socket.onopen = () => {
                window.dispatchEvent(new CustomEvent("cache:connect", { detail: { id: this.id, connect: true } }))
            };
            // 连接关闭事件
            this.socket.onclose = () => {
                window.dispatchEvent(new CustomEvent("cache:connect", { detail: { id: this.id, connect: false } }))
            };
            this.socket.onerror = () => this.close()
            // 加载WebSocket插件
            this.term.loadAddon(new AttachAddon(this.socket))
        },
        addResizeEvent() {
            let timeout = 0
            window.addEventListener('resize', () => {
                this.fitAddon.fit();
                clearTimeout(timeout)
                timeout = setTimeout(() => {
                    this.socket.send("!~" + this.term.cols + ":" + this.term.rows)
                }, 500)
            });
        },
        close() {
            this.connect = false
            // 关闭连接
            try {
                if (this.socket) this.socket.close()
            } catch (err) {
                console.log(err)
            }
            // 销毁终端
            try {
                if (this.term) this.term.dispose()
            } catch (err) {
                console.log(err)
            }
            document.getElementById('xterm').innerHTML = "";
            window.dispatchEvent(new CustomEvent("cache:connect", { detail: { id: this.id, connect: false } }))
            console.log('Terminal Close')
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

.xterm {
    height: calc(100vh - 80px);
}

.xterm:deep(.terminal) {
    padding: 5px;
}
</style>
  
