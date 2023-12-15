<template>
    <div class="pa-10">
        User
    </div>
</template>
  
<script>
export default {
    name: "User",
    data: () => ({
    }),
    methods: {
        init() {
            console.log('[Init] User')
        },
    },
    mounted() {
        this.init()
    }
};
</script>
  
<style scoped></style>
  
