<template>
    <div class="pa-10">
        Console
    </div>
</template>
  
<script>
export default {
    name: "Console",
    data: () => ({
    }),
    methods: {
        init() {
            console.log('[Init] Console')
        },
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
