<template>
    <div></div>
</template>
<script setup>
import { useMessage,useDialog,useNotification } from 'naive-ui'

window.$notify = useNotification()
window.$message = useMessage()
window.$dialog = useDialog()
</script>