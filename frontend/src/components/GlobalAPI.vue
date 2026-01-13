<script setup lang="ts">
import { useMessage, useDialog, useNotification, useLoadingBar } from 'naive-ui'

(window as any).$notify = useNotification();
(window as any).$message = useMessage();
(window as any).$dialog = useDialog();
(window as any).$loading = useLoadingBar();
</script>
<template>
    <div></div>
</template>