<template>
    <div class="pa-10">
        Logs
    </div>
</template>
  
<script>
export default {
    name: "Logs",
    data: () => ({
    }),
    methods: {
        init() {
            console.log('[Init] Logs')
        },
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
