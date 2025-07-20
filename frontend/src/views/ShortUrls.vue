<template>
    <div class="max-w-5xl">
        <h1 class="text-2xl font-bold mb-6">Shared File URLs</h1>
        <table
            class="min-w-full border border-gray-200 rounded-xl overflow-hidden"
        >
            <thead>
                <tr class="bg-gray-100">
                    <th class="px-4 py-2">Short URL</th>
                    <th class="px-4 py-2">Password?</th>
                    <th class="px-4 py-2">Approvable?</th>
                    <th class="px-4 py-2">Max Downloads</th>
                    <th class="px-4 py-2">Created At</th>
                    <th class="px-4 py-2">Expires At</th>
                    <th class="px-4 py-2">View Only</th>
                    <th class="px-4 py-2">Actions</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="url in sharedUrls" :key="url.id">
                    <td class="px-8 py-2">
                        <a
                            @click="redirectToGetFile(url.short_url)"
                            target="_blank"
                            class="text-blue-600 underline font-semibold cursor-pointer"
                            >click</a
                        >
                    </td>
                    <td class="px-4 py-2 text-center">
                        <span v-if="url.hashed_password">✅</span>
                        <span v-else>❌</span>
                    </td>
                    <td class="px-4 py-2 text-center">
                        <span v-if="url.approvable">✅</span>
                        <span v-else>❌</span>
                    </td>
                    <td class="px-4 py-2 text-center">
                        {{ url.max_downloads === -1 ? "-" : url.max_downloads }}
                    </td>
                    <td class="px-4 py-2 text-sm">
                        {{ formatDate(url.created_at) }}
                    </td>
                    <td class="px-4 py-2 text-sm">
                        {{ formatDate(url.expiration_at) }}
                    </td>
                    <td class="px-4 py-2 text-center">
                        <span title="Only Read" v-if="url.view_only">👁️</span>
                        <span title="Read-Write" v-else>✏️</span>
                    </td>
                    <td class="px-4 py-2 flex items-center gap-2">
                        <button
                            class="bg-red-500 hover:bg-red-700 text-white px-3 py-1 rounded"
                            @click="deleteUrl(url.id)"
                        >
                            Delete
                        </button>
                        <button
                            class="bg-blue-500 hover:bg-blue-700 text-white px-3 py-1 rounded"
                            @click="editUrl(url)"
                        >
                            Modify
                        </button>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>
<script setup>
import { ref, onMounted } from "vue";
import { useUserStore } from "../stores/user";
import axiosInstance from "../axiosInstance";
import { useRouter } from "vue-router";

const router = useRouter();

const userStore = useUserStore();
const sharedUrls = ref([]);

onMounted(async () => {
    sharedUrls.value = await fetchSharedUrls();
});

async function fetchSharedUrls() {
    axiosInstance.get("/api/file/settings/get").then((resp) => {
        sharedUrls.value = resp.data.sharedUrls;
    });
}

function formatDate(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleString();
}

function deleteUrl(id) {
    // Call API to delete and then update sharedUrls
    // await deleteSharedUrl(id);
    sharedUrls.value = sharedUrls.value.filter((url) => url.id !== id);
}

function redirectToGetFile(shortUrl) {
    router.push(`/file/get/${shortUrl}`);
}

function editUrl(url) {
    if (userStore.plan === "free") {
        alert("Free plan can only edit password protection.");
    } else {
        // Show full edit UI/modal
    }
}
</script>
