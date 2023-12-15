<template>
    <div id="auth-container" class="flex align-center justify-center no-select">
        <div id="auth-box">
            <div class="flex align-center justify-center">
                <img id="auth-logo" src="../assets/icon/icon-light@2x.png" />
                <div id="auth-title">Harbor</div>
            </div>
            <div id="auth-card" class="card">
                <n-form ref="formRef" :model="form" :disabled="loading" label-placement="top">
                    <n-form-item label="账户名称">
                        <n-input v-model:value="form.account" type="text" placeholder="请输入账户名称" clearable />
                    </n-form-item>
                    <n-form-item label="登陆密码">
                        <n-input v-model:value="form.voucher" type="password" placeholder="请输入登陆密码" clearable />
                    </n-form-item>
                </n-form>
                <n-button type="primary" class="full-width" :loading="loading" @click="login">登陆</n-button>
            </div>
            <div class="flex align-center justify-between">
                <div class="flex align-center">
                    <n-button quaternary circle class="text-gray" :disabled="loading">
                        <template #icon>
                            <n-icon>
                                <Carbon3rdPartyConnected />
                            </n-icon>
                        </template>
                    </n-button>
                </div>
                <n-button quaternary size="small" class="text-gray" :disabled="loading">忘记密码</n-button>
            </div>
        </div>
    </div>
</template>
  
<script>
import { Carbon3rdPartyConnected } from '@vicons/carbon'
export default {
    name: "Home",
    components: { Carbon3rdPartyConnected },
    data: () => ({
        loading: false,
        form: {
            account: '',
            voucher: ''
        }
    }),
    methods: {
        init() {
            console.log('[Init] Auth')
        },
        login() {
            this.loading = true;
            setTimeout(() => {
                window.$message.success("登陆成功, 欢迎回来.")
                localStorage.setItem("user:access:token", "test token")
            }, 1000)
            setTimeout(() => {
                this.$router.push('/')
            }, 2000)
        }
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped>
#auth-container {
    height: 100vh;
}

#auth-logo {
    width: 48px;
}

#auth-title {
    font-weight: bold;
    font-size: 42px;
    color: #17e3f0;
}

#auth-card {
    margin-bottom: 10px;
    padding: 20px;
    width: 300px;
}
</style>
  
