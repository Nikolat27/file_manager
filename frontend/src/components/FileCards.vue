<script setup>
import { ref } from 'vue'

const fileInputRef = ref(null)

function triggerFileInput() {
  fileInputRef.value.click()
}

function handleFileChange(event) {
  const file = event.target.files[0]
  if (file) {
    console.log("Selected file:", file)
    // You can upload it here
  }
}
</script>
<template>
    <div class="flex flex-row items-center gap-x-8 text-[16px] font-semibold">
        <input
            type="file"
            ref="fileInputRef"
            class="hidden"
            @change="handleFileChange"
        />

        <div @click="triggerFileInput" class="cursor-pointer w-[140px] flex items-center justify-center h-[76px]
         rounded-2xl bg-black border-1 text-white border-gray-500 hover:opacity-70">
            <span>Upload</span>
        </div>
        <div class="cursor-pointer w-[140px] flex items-center justify-center h-[76px]
         rounded-2xl bg-white border-1 border-gray-500 hover:opacity-70">
            <span>Create Folder</span>
        </div>
    </div>
</template>