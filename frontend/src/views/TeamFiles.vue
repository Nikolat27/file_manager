<template>
    <div class="flex flex-col gap-y-6 pl-120 pt-30">
        <span class="text-[22px] font-bold mb-2">All files</span>

        <!-- Top action buttons -->
        <div class="flex gap-4 mb-2">
            <button
                @click="showCreateFolderModal = true"
                class="px-4 py-2 rounded-xl bg-blue-600 text-white font-semibold hover:bg-blue-700"
            >
                Create Folder
            </button>
            <button
                @click="showUploadFileModal = true"
                class="px-4 py-2 rounded-xl bg-gray-200 text-blue-700 font-semibold hover:bg-gray-300"
            >
                Upload File
            </button>
            <button
                title="only admins can add users"
                class="bg-blue-600 text-white px-4 py-2 rounded-xl font-semibold hover:bg-blue-700"
                @click="showAddUserModal = true"
            >
                Add User
            </button>
        </div>

        <!-- Table -->
        <table
            class="w-[81%] h-auto border border-blue-200 rounded-xl overflow-hidden"
        >
            <thead>
                <tr class="bg-gray-100 text-gray-700">
                    <th class="text-left px-4 py-2 w-[70%]">Name</th>
                    <th class="text-left px-4 py-2 w-[15%]">Created At</th>
                    <th class="text-left px-4 py-2 w-[15%]">Expire At</th>
                    <th class="text-left px-4 py-2 w-[15%]">View</th>
                    <th class="text-left px-4 py-2 w-[15%]">More</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="folder in folders" :key="'folder-' + folder.id">
                    <td class="px-4 py-2 flex flex-row items-center gap-x-2 z">
                        <span title="folder">
                            <!-- Folder Icon -->
                            <svg
                                viewBox="0 0 40 40"
                                fill="none"
                                width="32"
                                height="32"
                            >
                                <path
                                    d="M15.002 7.004c.552.018.993.252 1.295.7l.785 2.12c.145.298.363.576.561.779.252.257.633.4 1.156.4H35.5l-.002 18c-.027.976-.3 1.594-.836 2.142-.565.577-1.383.858-2.41.858H8.748c-1.026 0-1.844-.28-2.409-.858-.564-.577-.838-1.415-.838-2.465V7.003h9.502Z"
                                    fill="#75aaff"
                                />
                                <path
                                    d="M15.002 7.001c.552.018.993.252 1.295.7l.785 2.12c.145.298.363.576.561.779.252.257.633.4 1.156.4H35.5l-.002 16.84c-.027.976-.3 1.754-.836 2.302-.565.577-1.383.858-2.41.858H8.748c-1.026 0-1.844-.28-2.409-.858-.564-.577-.838-1.415-.838-2.465V7l9.502.001Z"
                                    fill="#a0c4ff"
                                />
                            </svg>
                        </span>
                        <span class="cursor-pointer hover:underline">
                            {{ folder.name }}
                        </span>
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
                            <!-- File Icon -->
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
                    <td class="px-4 py-2">
                        {{ file.expire_at ? formatDate(file.expire_at) : "-" }}
                    </td>
                    <td
                        class="relative px-4 py-2 text-xl font-semibold cursor-pointer select-none pl-8 pb-4"
                    >
                        <span
                            @click="goToFile(file.id)"
                            class="cursor-pointer hover:text-blue-600 !text-[15px]"
                            >Click</span
                        >
                    </td>
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

        <!-- Folder/File Action Modal -->
        <div
            v-if="showModal"
            class="fixed inset-0 z-50 flex items-center justify-center bg-white bg-opacity-50"
        >
            <div
                class="bg-blue-600 rounded-2xl shadow-2xl p-6 w-72 flex flex-col gap-2 items-center"
            >
                <template v-if="isFolderModal">
                    <!-- Upload File Button (in folder modal) -->
                    <button
                        @click="openUploadFileToFolderModal"
                        class="w-full border-white border-2 mb-3 px-4 py-2 rounded-xl text-white text-lg font-semibold hover:bg-blue-700 transition cursor-pointer"
                    >
                        Upload File
                    </button>
                    <button
                        @click="
                            showRenameModal = true;
                            showModal = false;
                        "
                        class="w-full border-white border-2 mb-3 px-4 py-2 rounded-xl text-white text-lg font-semibold hover:bg-blue-700 transition cursor-pointer"
                    >
                        Rename
                    </button>
                    <button
                        @click="
                            deleteFolderFromModal();
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
                        class="w-full border-white border-2 px-4 py-2 rounded-xl text-white text-lg font-semibold hover:bg-blue-700 transition cursor-pointer"
                    >
                        Create Short Url
                    </button>
                    <button
                        @click="downloadFile()"
                        class="w-full border-white border-2 px-4 py-2 rounded-xl text-white text-lg font-semibold hover:bg-blue-700 transition cursor-pointer"
                    >
                        Download
                    </button>
                    <button
                        @click="openRenameFileModal"
                        class="w-full border-white border-2 px-4 py-2 rounded-xl text-white text-lg font-semibold hover:bg-blue-700 transition cursor-pointer"
                    >
                        Rename
                    </button>
                    <button
                        @click="handleDelete"
                        class="w-full border-red border-2 px-4 py-2 rounded-xl text-red-300 text-lg font-semibold hover:bg-red-600 hover:text-white transition cursor-pointer"
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

        <!-- Upload File Modal (root) -->
        <div
            v-if="showUploadFileModal"
            class="fixed inset-0 z-50 flex items-center justify-center bg-white bg-opacity-60"
        >
            <div
                class="bg-blue-600 rounded-2xl shadow-2xl p-8 w-[90vw] max-w-md flex flex-col items-center gap-4"
            >
                <h2 class="text-xl font-semibold text-white mb-2">
                    Upload File
                </h2>
                <input
                    v-model="customFileName"
                    type="text"
                    placeholder="File name (optional)"
                    class="w-full px-4 py-2 rounded-xl border border-blue-300 bg-blue-50 text-blue-900 font-semibold"
                />
                <input
                    type="file"
                    ref="uploadFileInput"
                    class="w-full px-4 py-2 rounded-xl border border-blue-300 bg-blue-50 text-blue-900"
                />
                <div class="flex gap-4 mt-4 w-full">
                    <button
                        @click="uploadFile"
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

        <!-- Upload File To Folder Modal -->
        <div
            v-if="showUploadFileToFolderModal"
            class="fixed inset-0 z-50 flex items-center justify-center bg-white bg-opacity-60"
        >
            <div
                class="bg-blue-600 rounded-2xl shadow-2xl p-8 w-[90vw] max-w-md flex flex-col items-center gap-4"
            >
                <h2 class="text-xl font-semibold text-white mb-2">
                    Upload File to {{ currentItem?.name }}
                </h2>
                <input
                    v-model="customFileName"
                    type="text"
                    placeholder="File name (optional)"
                    class="w-full px-4 py-2 rounded-xl border border-blue-300 bg-blue-50 text-blue-900 font-semibold"
                />
                <input
                    type="file"
                    ref="uploadFileToFolderInput"
                    class="w-full px-4 py-2 rounded-xl border border-blue-300 bg-blue-50 text-blue-900"
                />
                <div class="flex gap-4 mt-4 w-full">
                    <button
                        @click="uploadFileToFolder"
                        class="flex-1 py-2 rounded-xl bg-white text-blue-700 font-semibold hover:bg-blue-50 hover:text-blue-800 transition"
                    >
                        Upload
                    </button>
                    <button
                        @click="showUploadFileToFolderModal = false"
                        class="flex-1 py-2 rounded-xl bg-blue-500 text-white font-semibold hover:bg-blue-700 transition"
                    >
                        Cancel
                    </button>
                </div>
            </div>
        </div>

        <!-- Create Folder Modal -->
        <div
            v-if="showCreateFolderModal"
            class="fixed inset-0 z-50 flex items-center justify-center bg-white bg-opacity-60"
        >
            <div
                class="bg-blue-600 rounded-2xl shadow-2xl p-8 w-[90vw] max-w-md flex flex-col items-center gap-4"
            >
                <h2 class="text-xl font-semibold text-white mb-2">
                    Create Folder
                </h2>
                <input
                    v-model="newFolderName"
                    type="text"
                    placeholder="Enter folder name"
                    class="w-full px-4 py-2 rounded-xl border border-blue-300 bg-blue-50 text-blue-900 font-semibold"
                />
                <div class="flex gap-4 mt-4 w-full">
                    <button
                        @click="createFolder"
                        class="flex-1 py-2 rounded-xl bg-white text-blue-700 font-semibold hover:bg-blue-50 hover:text-blue-800 transition"
                    >
                        Create
                    </button>
                    <button
                        @click="showCreateFolderModal = false"
                        class="flex-1 py-2 rounded-xl bg-blue-500 text-white font-semibold hover:bg-blue-700 transition"
                    >
                        Cancel
                    </button>
                </div>
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
                </div>
            </div>
        </div>

        <!-- Add User Modal -->
        <div
            v-if="showAddUserModal"
            class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50"
        >
            <div
                class="bg-white rounded-2xl shadow-xl p-8 w-[350px] flex flex-col gap-4"
            >
                <h2 class="text-xl font-bold mb-2">Add User to Team</h2>
                <input
                    v-model="newUserId"
                    placeholder="Enter User ID"
                    class="border px-3 py-2 rounded-lg w-full mb-2"
                />
                <div class="flex gap-3">
                    <button
                        @click="addUser"
                        class="bg-blue-600 text-white px-4 py-2 rounded-xl hover:bg-blue-700"
                    >
                        Add
                    </button>
                    <button
                        @click="showAddUserModal = false"
                        class="bg-gray-300 text-gray-800 px-4 py-2 rounded-xl hover:bg-gray-400"
                    >
                        Cancel
                    </button>
                </div>
                <div v-if="errorMsg" class="text-red-600 text-sm font-semibold">
                    {{ errorMsg }}
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, watch, reactive } from "vue";
import { useRoute, useRouter } from "vue-router";
import axiosInstance from "../axiosInstance";
import { showError, showSuccess } from "../utils/toast";

const route = useRoute();
const router = useRouter();

const showModal = ref(false);
const isFolderModal = ref(false);
const showRenameModal = ref(false);
const showRenameFileModal = ref(false);
const showUploadFileModal = ref(false);
const showUploadFileToFolderModal = ref(false);
const showCreateFolderModal = ref(false);

const uploadFileInput = ref(null);
const uploadFileToFolderInput = ref(null);

const currentItem = ref(null);
const renameFolderName = ref("");
const renameFileName = ref("");
const newFolderName = ref("");

const files = ref([]);
const folders = ref([]);

const showAddUserModal = ref(false);
const newUserId = ref("");
const errorMsg = ref("");

const customFileName = ref("");
watch([showUploadFileModal, showUploadFileToFolderModal], ([valA, valB]) => {
    if (!valA && !valB) customFileName.value = "";
});

async function addUser() {
    if (!newUserId.value) {
        errorMsg.value = "User ID cannot be empty";
        return;
    }
    errorMsg.value = "";

    try {
        // Adjust endpoint as needed!
        await axiosInstance.post(`/api/team/user/add/${route.params.id}`, {
            user_id: newUserId.value,
        });
        showAddUserModal.value = false;
        newUserId.value = "";
        showSuccess("User added successfully");
    } catch (err) {
        errorMsg.value = err.response?.data?.error || "Failed to add user";
    }
}

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

function openUploadFileToFolderModal() {
    showUploadFileToFolderModal.value = true;
    showModal.value = false;
}

async function downloadFile() {
    try {
        const fileId = currentItem.value.id;
        const res = await axiosInstance.get(`/api/file/download/${fileId}`, {
            responseType: "blob",
        });
        let ext =
            res.headers["content-type"]?.split("/")[1]?.split(";")[0] || "bin";
        const filename = `${fileId}.${ext}`;
        const url = URL.createObjectURL(res.data);
        const link = Object.assign(document.createElement("a"), {
            href: url,
            download: filename,
        });
        document.body.appendChild(link);
        link.click();
        link.remove();
        URL.revokeObjectURL(url);

        showSuccess("Download started");
        closeModal();
    } catch {
        showError("Download failed");
    }
}

async function uploadFile() {
    const fileInputEl = uploadFileInput.value;
    if (!fileInputEl || !fileInputEl.files.length) return;

    const file = fileInputEl.files[0];
    const formData = new FormData();
    formData.append("file", file);
    formData.append("team_id", route.params.id);
    if (customFileName.value.trim()) {
        formData.append("file_name", customFileName.value.trim());
    }

    try {
        await axiosInstance.post(
            `/api/team/file/upload/${route.params.id}`,
            formData,
            {
                headers: { "Content-Type": "multipart/form-data" },
            }
        );
        showSuccess("File uploaded successfully");
        showUploadFileModal.value = false;
        getFiles();
    } catch (err) {
        showError(err.response?.data.error || "Upload failed");
    }
}

function uploadFileToFolder() {
    const fileInputEl = uploadFileToFolderInput.value;
    if (!fileInputEl || !fileInputEl.files.length) return;

    const file = fileInputEl.files[0];
    const formData = new FormData();
    formData.append("file", file);
    formData.append("folder_id", currentItem.value.id);
    formData.append("team_id", route.params.id);
    if (customFileName.value.trim()) {
        formData.append("file_name", customFileName.value.trim());
    }

    axiosInstance
        .post(`/api/file/create`, formData, {
            headers: { "Content-Type": "multipart/form-data" },
        })
        .then(() => {
            showSuccess("File uploaded successfully");
            showUploadFileToFolderModal.value = false;
            getFiles();
        })
        .catch((err) => {
            showError(err.response?.data.error || "Upload failed");
        });
}

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
            team_id: route.params.id,
        })
        .then(() => {
            showSuccess("File renamed successfully");
            currentItem.value.name = renameFileName.value;
            showRenameFileModal.value = false;
            getFiles();
        })
        .catch((err) => {
            closeModal();
            showError(err.response?.data.error || "Rename failed");
        });
}

function deleteFile(fileId) {
    axiosInstance
        .delete(`/api/file/delete/${fileId}?team_id=${route.params.id}`)
        .then(() => {
            showSuccess("file deleted successfully");
            getFiles();
        })
        .catch((err) => {
            showError(err.response?.data.error || "Delete failed");
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

function renameFolder() {
    axiosInstance
        .put(`/api/folder/rename/${currentItem.value.id}`, {
            name: renameFolderName.value,
            team_id: route.params.id,
        })
        .then(() => {
            showSuccess("Folder renamed successfully");
            currentItem.value.name = renameFolderName.value;
            showRenameModal.value = false;
            getFolders();
        })
        .catch((err) => {
            showError(err.response?.data.error || "Rename failed");
            closeModal();
        });
}

function deleteFolder(folderId) {
    axiosInstance
        .delete(`/api/folder/delete/${folderId}?team_id=${route.params.id}`)
        .then(() => {
            showSuccess("folder deleted successfully");
            getFolders();
        })
        .catch((err) => {
            showError(err.response?.data.error || "Delete failed");
        });
}

function deleteFolderFromModal() {
    showRenameModal.value = false;
    deleteFolder(currentItem.value.id);
}

async function createFolder() {
    if (!newFolderName.value.trim()) {
        showError("Folder name required");
        return;
    }
    try {
        await axiosInstance.post("/api/folder/create", {
            name: newFolderName.value,
            team_id: route.params.id,
        });
        showSuccess("Folder created");
        showCreateFolderModal.value = false;
        newFolderName.value = "";
        getFolders();
    } catch (err) {
        showError(err.response?.data.error || "Create failed");
    }
}

function formatDate(dateStr) {
    if (!dateStr) return "";
    return new Date(dateStr).toLocaleDateString("en-CA");
}

function goToFile(id) {
    const shortUrl = filesShortUrls[id];
    router.push({ name: "GetFile", params: { id: shortUrl } });
}

const filesShortUrls = reactive([]);

function getFiles() {
    axiosInstance
        .get(`/api/file/get?team_id=${route.params.id}`)
        .then((resp) => {
            files.value = resp.data.files;
            Object.assign(filesShortUrls, resp.data.shortUrls);
        });
}

function getFolders() {
    axiosInstance
        .get(`/api/folder/get?team_id=${route.params.id}`)
        .then((resp) => {
            folders.value = resp.data;
        });
}

function goToFolder(id) {
    router.push({ name: "FolderContents", params: { id } });
}

onMounted(() => {
    getFiles();
    getFolders();
});
</script>
