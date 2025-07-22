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
                <span class="text-2xl font-bold text-blue-800">
                    Folder Contents
                </span>
                <span
                    v-if="folderName"
                    class="ml-4 text-lg text-gray-700 font-semibold"
                    >/ {{ folderName }}</span
                >
            </div>

            <!-- Upload File Button -->
            <div class="flex w-full justify-end">
                <button
                    @click="showUploadFileModal = true"
                    class="px-6 py-2 mb-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-semibold shadow"
                >
                    Upload File
                </button>
            </div>

            <table
                class="w-full border border-blue-200 rounded-xl overflow-hidden bg-white shadow"
            >
                <thead>
                    <tr class="bg-gray-100 text-gray-700">
                        <th class="text-left px-4 py-2 w-[60%]">Name</th>
                        <th class="text-left px-4 py-2 w-[20%]">Type</th>
                        <th class="text-left px-4 py-2 w-[20%]">Created At</th>
                        <th class="text-left px-4 py-2 w-[15%]">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <template v-if="folders.length">
                        <tr v-for="folder in folders" :key="folder.id">
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
                                {{ folder.name }}
                            </td>
                            <td class="px-4 py-2">Folder</td>
                            <td class="px-4 py-2">
                                {{ formatDate(folder.created_at) }}
                            </td>
                            <td class="px-4 py-2">
                                <button
                                    @click="goToFolder(folder.id)"
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
                                    @click="openFileModal(file)"
                                    class="text-2xl font-bold text-blue-700 hover:text-blue-900"
                                >
                                    …
                                </button>
                            </td>
                        </tr>
                    </template>
                    <tr v-if="!folders.length && !files.length">
                        <td colspan="4" class="py-8 text-center text-gray-400">
                            No files found in this folder.
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- Upload File Modal -->
        <div
            v-if="showUploadFileModal"
            class="fixed inset-0 z-50 flex items-center justify-center bg-white bg-opacity-60"
        >
            <div
                class="bg-blue-600 rounded-2xl shadow-2xl p-8 w-[90vw] max-w-md flex flex-col items-center gap-4"
            >
                <h2 class="text-xl font-semibold text-white mb-2">
                    Upload File to {{ folderName || "Folder" }}
                </h2>
                <input
                    type="file"
                    ref="uploadFileInput"
                    class="w-full px-4 py-2 rounded-xl border border-blue-300 bg-blue-50 text-blue-900"
                />
                <div class="flex gap-4 mt-4 w-full">
                    <button
                        @click="uploadFileToCurrentFolder"
                        class="flex-1 py-2 rounded-xl bg-white text-blue-700 font-semibold hover:bg-blue-50 hover:text-blue-800 transition"
                    >
                        Upload
                    </button>
                    <button
                        @click="showUploadFileModal = false"
                        class="flex-1 py-2 rounded-xl bg-blue-500 text-white font-semibold hover:bg-blue-700 transition"
                    >
                        Cancel
                    </button>
                </div>
            </div>
        </div>

        <!-- File Actions Modal -->
        <div
            v-if="showFileModal"
            class="fixed inset-0 z-50 flex items-center justify-center bg-white bg-opacity-50"
        >
            <div
                class="bg-blue-600 rounded-2xl shadow-2xl p-8 w-80 flex flex-col gap-4 items-center"
            >
                <!-- <button
                    @click="downloadFile(selectedFile.id)"
                    class="w-full border-white border-2 px-4 py-2 rounded-xl text-white text-lg font-semibold hover:bg-blue-700 transition cursor-pointer"
                >
                    Download
                </button> -->
                <button
                    @click="openRenameFileModal"
                    class="w-full border-white border-2 px-4 py-2 rounded-xl text-white text-lg font-semibold hover:bg-blue-700 transition cursor-pointer"
                >
                    Rename
                </button>
                <button
                    @click="createSettings"
                    class="w-full border-white border-2 px-4 py-2 rounded-xl text-white text-lg font-semibold hover:bg-blue-700 transition cursor-pointer"
                >
                    Create Settings
                </button>
                <button
                    @click="deleteFileFromModal"
                    class="w-full py-2 rounded-xl bg-red-500 text-white font-semibold hover:bg-red-700 transition"
                >
                    Delete
                </button>
                <button
                    @click="closeFileModal"
                    class="mt-2 w-full px-4 py-2 rounded-xl bg-white text-blue-600 font-semibold hover:bg-blue-50 hover:text-blue-700 transition cursor-pointer"
                >
                    Cancel
                </button>
            </div>
        </div>

        <!-- Rename File Modal -->
        <div
            v-if="showRenameFileModal"
            class="fixed inset-0 z-50 flex items-center justify-center bg-white bg-opacity-60"
        >
            <div
                class="bg-blue-600 rounded-2xl shadow-2xl p-8 w-96 flex flex-col items-center gap-4"
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
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import axiosInstance from "../axiosInstance";
import { showError, showSuccess } from "../utils/toast";

const route = useRoute();
const router = useRouter();

const folderId = route.params.id || "";

const folderName = ref("");
const folders = ref([]); // Subfolders
const files = ref([]); // Files

const showFileModal = ref(false);
const showRenameFileModal = ref(false);
const selectedFile = ref(null);
const renameFileName = ref("");

// Upload modal state
const showUploadFileModal = ref(false);
const uploadFileInput = ref(null);

// Load folder contents
function fetchFolderContents(id) {
    axiosInstance
        .get(`/api/folder/get/${id}`)
        .then((resp) => {
            folders.value = resp.data.folders || [];
            files.value = resp.data.files || [];
            folderName.value = resp.data.folder_name || "";
        })
        .catch((err) => {
            console.error(err.response?.data || err);
        });
}

function formatDate(dateStr) {
    if (!dateStr) return "";
    return new Date(dateStr).toLocaleDateString("en-CA");
}

function goBack() {
    router.back();
}

function goToFolder(id) {
    router.push({ name: "FolderContents", params: { id } });
}

// File modal logic
function openFileModal(file) {
    selectedFile.value = file;
    showFileModal.value = true;
}
function closeFileModal() {
    showFileModal.value = false;
    selectedFile.value = null;
}
function openRenameFileModal() {
    renameFileName.value = selectedFile.value?.name || "";
    showRenameFileModal.value = true;
    showFileModal.value = false;
}

function renameFile() {
    axiosInstance
        .put(`/api/file/rename/${selectedFile.value.id}`, {
            name: renameFileName.value,
        })
        .then(() => {
            showSuccess("File renamed successfully");
            showRenameFileModal.value = false;
            fetchFolderContents(folderId);
        })
        .catch((err) => {
            showError(err.response?.data || "Rename failed");
        });
}
function deleteFileFromModal() {
    axiosInstance
        .delete(`/api/file/delete/${selectedFile.value.id}`)
        .then(() => {
            showSuccess("File deleted successfully");
            showFileModal.value = false;
            fetchFolderContents(folderId);
        })
        .catch((err) => {
            showError(err.response?.data || "Delete failed");
        });
}
function createSettings() {
    // Redirect to settings page for file
    router.push(
        `/file/setting/create/${selectedFile.value.id}?folderId=${folderId}`
    );
    closeFileModal();
}

async function downloadFile(fileId) {
    try {
        const res = await axiosInstance.get(`/api/file/download/${fileId}`, {
            responseType: "blob",
        });

        // Extract file extension from content-type header
        const contentType = res.headers['content-type'] || '';
        let ext = contentType.split('/')[1] || 'bin'; // fallback to .bin

        // For content types like "application/pdf"
        if (ext.includes(';')) ext = ext.split(';')[0];

        const filename = `${fileId}.${ext}`;
        const url = URL.createObjectURL(res.data);
        const link = document.createElement("a");
        link.href = url;
        link.download = filename;
        document.body.appendChild(link);
        link.click();
        link.remove();
        URL.revokeObjectURL(url);

        showSuccess("Download started");
        closeFileModal();
    } catch (err) {
        showError("Download failed");
    }
}

// ---- Upload logic ----
function uploadFileToCurrentFolder() {
    const input = uploadFileInput.value;
    if (!input || !input.files.length) return;

    const file = input.files[0];
    const formData = new FormData();
    formData.append("file", file);
    formData.append("folder_id", folderId);

    axiosInstance
        .post(`/api/file/create`, formData, {
            headers: { "Content-Type": "multipart/form-data" },
        })
        .then(() => {
            showSuccess("File uploaded successfully");
            showUploadFileModal.value = false;
            fetchFolderContents(folderId);
        })
        .catch((err) => {
            showError(err.response?.data || "Upload failed");
        });
}

onMounted(() => {
    fetchFolderContents(folderId);
});
</script>
