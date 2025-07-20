<script setup lang="ts">
import Sidebar from "./components/Sidebar.vue";
import Header from "./components/Header.vue";
import { useRoute } from "vue-router";
import { computed } from "vue";

const route = useRoute();

const isAuthPage = computed(() => {
    return route.name === "login" || route.name === "register";
});
</script>

<template>
    <div class="w-screen h-screen relative">
        <Sidebar v-if="!isAuthPage"></Sidebar>
        <Header v-if="!isAuthPage"></Header>
        <router-view />
    </div>
</template>
