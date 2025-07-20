<template>
    <div class="flex flex-col gap-y-6">
        <span class="text-[22px] font-bold">All files</span>

        <table
            class="w-[81%] h-auto border border-blue-200 rounded-xl overflow-hidden"
        >
            <thead>
                <tr class="bg-gray-100 text-gray-700">
                    <th class="text-left px-4 py-2 w-[70%]">Name</th>
                    <th class="text-left px-4 py-2 w-[15%]">Updated At</th>
                    <th class="text-left px-4 py-2 w-[15%]">Expire At</th>
                    <th class="text-left px-4 py-2 w-[15%]">More</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="folder in folders">
                    <td class="px-4 py-2 flex flex-row items-center gap-x-2 z">
                        <span title="folder">
                            <svg
                                viewBox="0 0 40 40"
                                fill="none"
                                role="img"
                                focusable="false"
                                width="32"
                                height="32"
                                class="dig-ContentIcon brws-file-name-cell-icon dig-ContentIcon--small dig-ctz1wx2_5-7-0 dig-ctz1wx4_5-7-0"
                                data-testid="FolderBaseDefaultSmall"
                                data-campaigns-element-id="file-thumbnail"
                                data-thumbnail-testid="fallback-file-thumbnail"
                            >
                                <path
                                    d="M15.002 7.004c.552.018.993.252 1.295.7l.785 2.12c.145.298.363.576.561.779.252.257.633.4 1.156.4H35.5l-.002 18c-.027.976-.3 1.594-.836 2.142-.565.577-1.383.858-2.41.858H8.748c-1.026 0-1.844-.28-2.409-.858-.564-.577-.838-1.415-.838-2.465V7.003h9.502Z"
                                    fill="var(--dig-color__foldericon__shadow, #75aaff)"
                                ></path>
                                <path
                                    d="M15.002 7.001c.552.018.993.252 1.295.7l.785 2.12c.145.298.363.576.561.779.252.257.633.4 1.156.4H35.5l-.002 16.84c-.027.976-.3 1.754-.836 2.302-.565.577-1.383.858-2.41.858H8.748c-1.026 0-1.844-.28-2.409-.858-.564-.577-.838-1.415-.838-2.465V7l9.502.001Z"
                                    fill="var(--dig-color__foldericon__container, #a0c4ff)"
                                ></path>
                            </svg>
                        </span>
                        {{ folder.name }}
                    </td>
                    <td class="px-4 py-2">
                        {{ formatDate(folder.created_at) }}
                    </td>
                    <td class="px-4 py-2">
                        {{ formatDate(folder.expire_at) }}
                    </td>
                    <td
                        class="relative px-4 py-2 text-2xl font-bold cursor-pointer select-none pl-8 pb-4"
                    >
                        ...
                    </td>
                </tr>
                <tr v-for="file in files">
                    <td class="px-4 py-2 flex flex-row items-center gap-x-2">
                        <span title="file">
                            <svg
                                class="mr-2"
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                                width="24"
                                height="24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M7 3h6l5 5v13a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2zm6 0v5h5"
                                />
                            </svg>
                        </span>
                        {{ file.name }}
                    </td>
                    <td class="px-4 py-2">{{ formatDate(file.created_at) }}</td>
                    <td class="px-4 py-2">{{ formatDate(file.expire_at) }}</td>
                    <td
                        class="relative px-4 py-2 text-xl font-semibold cursor-pointer select-none pl-8 pb-4"
                    >
                        <span
                            @click="openModal(file.id)"
                            class="hover:text-blue-600 text-2xl"
                            >…</span
                        >
                    </td>

                    <!-- Modal Action Menu -->
                    <div
                        v-if="showModal"
                        class="fixed inset-0 z-50 flex items-center justify-center bg-white bg-opacity-50"
                    >
                        <div
                            class="bg-blue-600 rounded-2xl shadow-2xl p-6 w-72 flex flex-col gap-2 items-center"
                        >
                            <button
                                @click="handleCreate"
                                class="w-full px-4 py-2 rounded-xl text-white text-lg font-semibold hover:bg-blue-700 transition cursor-pointer"
                            >
                                Create
                            </button>
                            <button
                                @click="handleEdit"
                                class="w-full px-4 py-2 rounded-xl text-white text-lg font-semibold hover:bg-blue-700 transition cursor-pointer"
                            >
                                Edit
                            </button>
                            <button
                                @click="handleDelete"
                                class="w-full px-4 py-2 rounded-xl text-red-200 text-lg font-semibold hover:bg-red-600 hover:text-white transition cursor-pointer"
                            >
                                Delete
                            </button>
                            <button
                                @click="closeModal"
                                class="mt-3 w-full px-4 py-2 rounded-xl bg-white text-blue-600 font-semibold hover:bg-blue-50 hover:text-blue-700 transition cursor-pointer"
                            >
                                Cancel
                            </button>
                        </div>
                    </div>
                </tr>
            </tbody>
        </table>
    </div>
</template>
<script setup>
import { ref, onMounted } from "vue";
import axiosInstance from "../axiosInstance";
import { useRouter } from "vue-router";

const showModal = ref(false);
const currentFileId = ref(null);

function openModal(fileId) {
    showModal.value = true;
    currentFileId.value = fileId;
}
function closeModal() {
    showModal.value = false;
    currentFileId.value = null;
}

// Action handlers
function handleCreate() {
    router.push(`/file/setting/create/${currentFileId.value}`)
}

function handleEdit() {
    alert("Edit for file " + currentFileId.value);
    closeModal();
}
function handleDelete() {
    alert("Delete for file " + currentFileId.value);
    closeModal();
}

const router = useRouter();

const files = ref([]);
const folders = ref([]);

function formatDate(dateStr) {
    if (!dateStr) return "";
    return new Date(dateStr).toLocaleDateString("en-CA");
}

function getFiles() {
    axiosInstance.get("/api/file/get").then((resp) => {
        files.value = resp.data;
    });
}

function getFolders() {
    axiosInstance.get("/api/folder/get").then((resp) => {
        folders.value = resp.data;
    });
}

function goToEdit(id) {
    if (id) {
        router.push(`/file/edit/${id}`);
    } else {
        alert("No file id provided!");
    }
}

onMounted(async () => {
    getFiles();
    getFolders();
});
</script>
