<script setup lang="ts">
import Sidebar from './components/Sidebar.vue';
import Header from './components/Header.vue';

</script>

<template>
  <div class="w-screen h-screen relative">
    <Sidebar></Sidebar>
    <Header></Header>
    <router-view />
  </div>
</template>