<template>
    <div class="pa-10">
        Container
    </div>
</template>
  
<script>
export default {
    name: "Container",
    data: () => ({
    }),
    methods: {
        init() {
            console.log('[Init] Container')
        },
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
