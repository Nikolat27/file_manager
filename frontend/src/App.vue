<script setup lang="ts">
import Sidebar from "./components/Sidebar.vue";
import Header from "./components/Header.vue";
import { useRoute, useRouter } from "vue-router";
import { computed } from "vue";
import { useUserStore } from "./stores/user";

const route = useRoute();
const router = useRouter();

const isAuthPage = computed(() => {
    return route.name === "login" || route.name === "register";
});

const userStore = useUserStore();

function checkUserAuthentication() {
    if (!userStore.token && route.name != "register") {
        router.push("/login");
    }
}

checkUserAuthentication();
</script>

<template>
    <div class="w-screen h-screen relative">
        <Sidebar v-if="!isAuthPage"></Sidebar>
        <Header v-if="!isAuthPage"></Header>
        <router-view />
    </div>
</template>
