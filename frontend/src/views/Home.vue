<template>
    <div
        class="flex flex-col w-screen h-screen absolute top-22 left-82 gap-y-10"
    >
        <div>
            <FileCards></FileCards>
        </div>
        <div>
            <ShowFiles></ShowFiles>
        </div>
    </div>
</template>

<script setup lang="ts">
import FileCards from "../components/FileCards.vue";
import ShowFiles from "../components/ShowFiles.vue";
import { useUserStore } from "../stores/user";

const userStore = useUserStore();

console.log(userStore.token);
</script>
