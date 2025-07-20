<script setup>
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import axiosInstance from "../axiosInstance";

const route = useRoute();
const fileData = ref(null);
const loading = ref(true);
const error = ref("");
const passwordModal = ref(false);
const password = ref("");
const id = route.params.id;

async function fetchFile() {
    loading.value = true;
    error.value = "";
    try {
        const resp = await axiosInstance.get(`/api/file/get/${id}`);
        fileData.value = resp.data;
    } catch (err) {
        if (err.response?.status === 406) {
            // Password required
            passwordModal.value = true;
            error.value = "";
        } else {
            error.value = err.response?.data?.message || "Unknown error";
        }
    } finally {
        loading.value = false;
    }
}

async function submitPassword() {
    loading.value = true;
    error.value = "";
    try {
        const resp = await axiosInstance.post(`/api/file/get/${id}`, {
            password: password.value,
        });
        fileData.value = resp.data;
        passwordModal.value = false;
    } catch (err) {
        error.value =
            err.response?.data?.message || "Incorrect password. Try again";
    } finally {
        loading.value = false;
        password.value = "";
    }
}

onMounted(fetchFile);
</script>

<template>
    <div class="pt-32 flex flex-col items-center justify-center min-h-[400px]">
        <div v-if="loading" class="text-blue-500 text-lg font-semibold">
            Loading file info...
        </div>
        <div
            v-else-if="error && !passwordModal"
            class="text-red-600 text-[14px] font-semibold"
        >
            {{ error }}
        </div>
        <div v-else-if="fileData">
            <div class="bg-gray-50 p-6 rounded-xl shadow-md max-w-lg w-full">
                <h2 class="text-xl font-bold mb-4">File Info</h2>
                <pre class="whitespace-pre-wrap text-sm">{{ fileData }}</pre>
            </div>
        </div>

        <!-- Password Modal -->
        <div
            v-if="passwordModal"
            class="fixed inset-0 flex items-center justify-center bg-gray-200 bg-opacity-50 z-50"
        >
            <div class="bg-white p-8 rounded-xl shadow-lg max-w-sm w-full">
                <h2 class="text-lg font-bold mb-4">Password Required</h2>
                <input
                    v-model="password"
                    type="password"
                    class="w-full border rounded px-3 py-2 mb-4"
                    placeholder="Enter file password"
                    @keyup.enter="submitPassword"
                    autofocus
                />
                <div class="flex justify-end gap-2">
                    <button
                        class="bg-gray-400 text-white px-4 py-2 rounded"
                        @click="passwordModal = false"
                        :disabled="loading"
                    >
                        Cancel
                    </button>
                    <button
                        class="bg-blue-600 hover:bg-blue-800 text-white px-4 py-2 rounded"
                        @click="submitPassword"
                        :disabled="!password || loading"
                    >
                        Submit
                    </button>
                </div>
                <div v-if="error" class="mt-2 text-red-600 text-sm">
                    {{ error }}
                </div>
            </div>
        </div>
    </div>
</template>
