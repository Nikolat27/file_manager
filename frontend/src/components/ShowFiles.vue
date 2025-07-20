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
                <tr v-for="folder in folders" :key="'folder-' + folder.id">
                    <td class="px-4 py-2 flex flex-row items-center gap-x-2 z">
                        <span title="folder">
                            <!-- FOLDER ICON SVG -->
                            <svg
                                viewBox="0 0 40 40"
                                fill="none"
                                role="img"
                                focusable="false"
                                width="32"
                                height="32"
                                class="dig-ContentIcon brws-file-name-cell-icon dig-ContentIcon--small dig-ctz1wx2_5-7-0 dig-ctz1wx4_5-7-0"
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
                    <td class="px-4 py-2 text-sm">no expiration for folders</td>
                    <td
                        class="relative px-4 py-2 text-2xl font-bold cursor-pointer select-none pl-8 pb-4"
                    >
                        <span
                            @click="openModal(folder, true)"
                            class="hover:text-blue-600 text-2xl"
                            >…</span
                        >
                    </td>
                </tr>
                <tr v-for="file in files" :key="'file-' + file.id">
                    <td class="px-4 py-2 flex flex-row items-center gap-x-2">
                        <span title="file">
                            <!-- FILE ICON SVG -->
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
                            @click="openModal(file, false)"
                            class="hover:text-blue-600 text-2xl"
                            >…</span
                        >
                    </td>
                </tr>
            </tbody>
        </table>

        <!-- Modal Action Menu -->
        <div
            v-if="showModal"
            class="fixed inset-0 z-50 flex items-center justify-center bg-white bg-opacity-50"
        >
            <div
                class="bg-blue-600 rounded-2xl shadow-2xl p-6 w-72 flex flex-col gap-2 items-center"
            >
                <template v-if="isFolderModal">
                    <button
                        @click="
                            showRenameModal = true;
                            showModal = false;
                        "
                        class="w-full px-4 py-2 rounded-xl text-white text-lg font-semibold hover:bg-blue-700 transition cursor-pointer"
                    >
                        Rename
                    </button>
                    <button
                        @click="
                            deleteFolderFromModal;
                            showModal = false;
                        "
                        class="w-full py-2 rounded-xl bg-red-500 text-white font-semibold hover:bg-red-700 transition"
                    >
                        Delete
                    </button>
                </template>
                <template v-else>
                    <button
                        @click="handleCreate"
                        class="w-full px-4 py-2 rounded-xl text-white text-lg font-semibold hover:bg-blue-700 transition cursor-pointer"
                    >
                        Create
                    </button>
                    <button
                        @click="openRenameFileModal"
                        class="w-full px-4 py-2 rounded-xl text-white text-lg font-semibold hover:bg-blue-700 transition cursor-pointer"
                    >
                        Rename
                    </button>
                    <button
                        @click="handleDelete"
                        class="w-full px-4 py-2 rounded-xl text-red-200 text-lg font-semibold hover:bg-red-600 hover:text-white transition cursor-pointer"
                    >
                        Delete
                    </button>
                </template>
                <button
                    @click="closeModal"
                    class="mt-3 w-full px-4 py-2 rounded-xl bg-white text-blue-600 font-semibold hover:bg-blue-50 hover:text-blue-700 transition cursor-pointer"
                >
                    Cancel
                </button>
            </div>
        </div>

        <!-- Rename Folder Modal -->
        <div
            v-if="showRenameModal"
            class="fixed inset-0 z-50 flex items-center justify-center bg-white bg-opacity-60"
        >
            <div
                class="bg-blue-600 rounded-2xl shadow-2xl p-8 w-[90vw] max-w-md flex flex-col items-center gap-4"
            >
                <h2 class="text-xl font-semibold text-white mb-2">
                    Rename Folder
                </h2>
                <input
                    v-model="renameFolderName"
                    type="text"
                    placeholder="Enter new folder name"
                    class="w-full px-4 py-2 rounded-xl border border-blue-300 focus:outline-none focus:border-white bg-blue-50 text-blue-900 font-semibold"
                />
                <div class="flex gap-4 mt-4 w-full">
                    <button
                        @click="renameFolder"
                        class="flex-1 py-2 rounded-xl bg-white text-blue-700 font-semibold hover:bg-blue-50 hover:text-blue-800 transition"
                    >
                        Rename
                    </button>
                    <button
                        @click="showRenameModal = false"
                        class="flex-1 py-2 rounded-xl bg-blue-500 text-white font-semibold hover:bg-blue-700 transition"
                    >
                        Cancel
                    </button>
                </div>
            </div>
        </div>

        <!-- Rename File Modal -->
        <div
            v-if="showRenameFileModal"
            class="fixed inset-0 z-50 flex items-center justify-center bg-white bg-opacity-60"
        >
            <div
                class="bg-blue-600 rounded-2xl shadow-2xl p-8 w-[90vw] max-w-md flex flex-col items-center gap-4"
            >
                <h2 class="text-xl font-semibold text-white mb-2">
                    Rename File
                </h2>
                <input
                    v-model="renameFileName"
                    type="text"
                    placeholder="Enter new file name"
                    class="w-full px-4 py-2 rounded-xl border border-blue-300 focus:outline-none focus:border-white bg-blue-50 text-blue-900 font-semibold"
                />
                <div class="flex gap-4 mt-4 w-full">
                    <button
                        @click="renameFile"
                        class="flex-1 py-2 rounded-xl bg-white text-blue-700 font-semibold hover:bg-blue-50 hover:text-blue-800 transition"
                    >
                        Rename
                    </button>
                    <button
                        @click="showRenameFileModal = false"
                        class="flex-1 py-2 rounded-xl bg-blue-500 text-white font-semibold hover:bg-blue-700 transition"
                    >
                        Cancel
                    </button>
                    <button
                        @click="deleteFileFromModal"
                        class="flex-1 py-2 rounded-xl bg-red-500 text-white font-semibold hover:bg-red-700 transition"
                    >
                        Delete
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import axiosInstance from "../axiosInstance";
import { useRouter } from "vue-router";
import { showError, showSuccess } from "../utils/toast";

const showModal = ref(false);
const isFolderModal = ref(false);
const showRenameModal = ref(false);
const showRenameFileModal = ref(false);
const currentItem = ref(null);
const renameFolderName = ref("");
const renameFileName = ref("");

function openModal(item, isFolder) {
    showModal.value = true;
    isFolderModal.value = isFolder;
    currentItem.value = item;
    if (isFolder) {
        renameFolderName.value = item.name;
    } else {
        renameFileName.value = item.name;
    }
}
function closeModal() {
    showModal.value = false;
    currentItem.value = null;
    isFolderModal.value = false;
}

// --- File actions ---
const router = useRouter();
function handleCreate() {
    router.push(`/file/setting/create/${currentItem.value.id}`);
    closeModal();
}
function openRenameFileModal() {
    showRenameFileModal.value = true;
    showModal.value = false;
}
function renameFile() {
    axiosInstance
        .put(`/api/file/rename/${currentItem.value.id}`, {
            name: renameFileName.value,
        })
        .then(() => {
            showSuccess("File renamed successfully");
            currentItem.value.name = renameFileName.value;
            showRenameFileModal.value = false;
        })
        .catch((err) => {
            showError(err.response.data);
        });
}
function handleEdit() {
    openRenameFileModal();
}

function deleteFile(fileId) {
    axiosInstance
        .delete(`/api/file/delete/${fileId}`)
        .then(() => {
            showSuccess("file deleted successfully");
        })
        .catch((err) => {
            showError(err.response.data);
        });
}
function deleteFileFromModal() {
    showRenameFileModal.value = false;
    deleteFile(currentItem.value.id);
}

function handleDelete() {
    if (isFolderModal.value) {
        deleteFolder(currentItem.value.id);
    } else {
        deleteFile(currentItem.value.id);
    }
    closeModal();
}

// --- Folder actions ---
function renameFolder() {
    axiosInstance
        .put(`/api/folder/rename/${currentItem.value.id}`, {
            name: renameFolderName.value,
        })
        .then(() => {
            showSuccess("Folder renamed successfully");
            currentItem.value.name = renameFolderName.value;
            showRenameModal.value = false;
        })
        .catch((err) => {
            showError(err.response.data);
        });
}
function deleteFolder(folderId) {
    axiosInstance
        .delete(`/api/folder/delete/${folderId}`)
        .then(() => {
            showSuccess("folder deleted successfully");
        })
        .catch((err) => {
            showError(err.response.data);
        });
}
function deleteFolderFromModal() {
    showRenameModal.value = false;
    deleteFolder(currentItem.value.id);
}

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

onMounted(async () => {
    getFiles();
    getFolders();
});
</script>
