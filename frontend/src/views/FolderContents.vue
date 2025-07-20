<template>
    <div class="min-h-screen bg-blue-50 flex flex-col items-center py-8 pt-40">
        <div class="w-full max-w-5xl flex flex-col gap-y-6">
            <div class="flex items-center gap-3">
                <button
                    @click="goBack"
                    class="cursor-pointer text-blue-600 hover:text-blue-800 font-bold text-lg flex items-center"
                >
                    <svg width="20" height="20" fill="none" viewBox="0 0 20 20">
                        <path
                            d="M13 17l-5-5 5-5"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        />
                    </svg>
                    Back
                </button>
                <span class="text-2xl font-bold text-blue-800"
                    >Folder Contents</span
                >
                <span
                    v-if="folderName"
                    class="ml-4 text-lg text-gray-700 font-semibold"
                    >/ {{ folderName }}</span
                >
            </div>

            <table
                class="w-full border border-blue-200 rounded-xl overflow-hidden bg-white shadow"
            >
                <thead>
                    <tr class="bg-gray-100 text-gray-700">
                        <th class="text-left px-4 py-2 w-[60%]">Name</th>
                        <th class="text-left px-4 py-2 w-[20%]">Type</th>
                        <th class="text-left px-4 py-2 w-[20%]">Updated At</th>
                        <th class="text-left px-4 py-2 w-[15%]">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <template v-if="folders.length">
                        <tr v-for="sub in folders" :key="sub.id">
                            <td class="px-4 py-2 flex items-center gap-x-2">
                                <svg
                                    viewBox="0 0 40 40"
                                    fill="none"
                                    width="28"
                                    height="28"
                                >
                                    <path
                                        d="M15.002 7.004c.552.018.993.252 1.295.7l.785 2.12c.145.298.363.576.561.779.252.257.633.4 1.156.4H35.5l-.002 18c-.027.976-.3 1.594-.836 2.142-.565.577-1.383.858-2.41.858H8.748c-1.026 0-1.844-.28-2.409-.858-.564-.577-.838-1.415-.838-2.465V7.003h9.502Z"
                                        fill="#75aaff"
                                    ></path>
                                    <path
                                        d="M15.002 7.001c.552.018.993.252 1.295.7l.785 2.12c.145.298.363.576.561.779.252.257.633.4 1.156.4H35.5l-.002 16.84c-.027.976-.3 1.754-.836 2.302-.565.577-1.383.858-2.41.858H8.748c-1.026 0-1.844-.28-2.409-.858-.564-.577-.838-1.415-.838-2.465V7l9.502.001Z"
                                        fill="#a0c4ff"
                                    ></path>
                                </svg>
                                {{ sub.name }}
                            </td>
                            <td class="px-4 py-2">Folder</td>
                            <td class="px-4 py-2">
                                {{ formatDate(sub.created_at) }}
                            </td>
                            <td class="px-4 py-2">
                                <button
                                    @click="goToFolder(sub.id)"
                                    class="text-blue-600 hover:text-blue-800 font-semibold"
                                >
                                    Open
                                </button>
                            </td>
                        </tr>
                    </template>
                    <template v-if="files.length">
                        <tr v-for="file in files" :key="file.id">
                            <td class="px-4 py-2 flex items-center gap-x-2">
                                <svg
                                    class="mr-1"
                                    width="22"
                                    height="22"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M7 3h6l5 5v13a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2zm6 0v5h5"
                                    />
                                </svg>
                                {{ file.name }}
                            </td>
                            <td class="px-4 py-2">File</td>
                            <td class="px-4 py-2">
                                {{ formatDate(file.created_at) }}
                            </td>
                            <td class="px-4 py-2">
                                <button
                                    @click="downloadFile(file.id)"
                                    class="text-blue-600 hover:text-blue-800 font-semibold"
                                >
                                    Download
                                </button>
                            </td>
                        </tr>
                    </template>
                    <tr v-if="!folders.length && !files.length">
                        <td colspan="4" class="py-8 text-center text-gray-400">
                            No files or subfolders found in this folder.
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import axiosInstance from "../axiosInstance";

const route = useRoute();
const router = useRouter();

const folderId = route.params.id || ""; // <- gets id from URL (e.g. /folder/:id)

const folderName = ref("Folder Name"); // Will be filled after fetching
const folders = ref([]); // Subfolders
const files = ref([]); // Files

function goBack() {
    router.back();
}

function downloadFile(id) {
    // Implement file download logic
}
function formatDate(dateStr) {
    if (!dateStr) return "";
    return new Date(dateStr).toLocaleDateString("en-CA");
}

function fetchFolderContents(id) {
    axiosInstance
        .get(`/api/folder/get/${id}`)
        .then((resp) => {
            files.value = resp.data.files;
            console.log(resp);
        })
        .catch((err) => {
            console.error(err.response.data);
        });
}

onMounted(() => {
    fetchFolderContents(route.params.id);
});
</script>
