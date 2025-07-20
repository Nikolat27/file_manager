<template>
    <div class="flex flex-row items-center gap-x-8 text-[16px] font-semibold">
        <!-- Upload File Button -->
        <input
            type="file"
            ref="fileInputRef"
            class="hidden"
            @change="handleFileChange"
        />
        <div
            @click="triggerFileInput"
            class="cursor-pointer w-[140px] flex items-center justify-center h-[76px] rounded-2xl bg-black border-1 text-white border-gray-500 hover:opacity-70"
        >
            <span>Upload</span>
        </div>

        <!-- Create Folder Button -->
        <div
            @click="openFolderModal"
            class="cursor-pointer w-[140px] flex items-center justify-center h-[76px] rounded-2xl bg-white border-1 border-gray-500 hover:opacity-70"
        >
            <span>Create Folder</span>
        </div>
    </div>

    <!-- Modal Popup for File Upload -->
    <div
        v-if="showFileModal"
        class="fixed inset-0 bg-white bg-opacity-20 flex items-center justify-center z-50"
    >
        <div
            class="bg-black text-white rounded-xl p-8 flex flex-col items-center gap-4 w-[90vw] max-w-md shadow-2xl"
        >
            <h2 class="text-xl font-semibold mb-2">Optional: Name Your File</h2>
            <input
                v-model="customFileName"
                type="text"
                placeholder="File name (optional)"
                class="w-full px-4 py-2 text-white rounded border border-gray-300 focus:outline-none font-semibold"
            />
            <div class="flex gap-4 mt-4">
                <button
                    @click="handleUpload"
                    class="px-6 py-2 rounded-xl bg-blue-600 text-white font-semibold hover:bg-blue-700"
                >
                    Upload
                </button>
                <button
                    @click="closeFileModal"
                    class="px-6 py-2 rounded-xl bg-gray-200 text-gray-800 font-semibold hover:bg-gray-300"
                >
                    Cancel
                </button>
            </div>
        </div>
    </div>

    <!-- Modal Popup for Folder Creation -->
    <div
        v-if="showFolderModal"
        class="fixed inset-0 bg-white bg-opacity-20 flex items-center justify-center z-50"
    >
        <div
            class="bg-black text-white rounded-xl p-8 flex flex-col items-center gap-4 w-[90vw] max-w-md shadow-2xl"
        >
            <h2 class="text-xl font-semibold mb-2">Name Your Folder</h2>
            <input
                v-model="folderName"
                type="text"
                placeholder="Folder name"
                class="w-full px-4 py-2 text-white rounded border border-gray-300 focus:outline-none font-semibold"
            />
            <div class="flex gap-4 mt-4">
                <button
                    @click="createFolder"
                    class="px-6 py-2 rounded-xl bg-blue-600 text-white font-semibold hover:bg-blue-700"
                >
                    Create
                </button>
                <button
                    @click="closeFolderModal"
                    class="px-6 py-2 rounded-xl bg-gray-200 text-gray-800 font-semibold hover:bg-gray-300"
                >
                    Cancel
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from "vue";
import axiosInstance from "../axiosInstance";
import { useUserStore } from "../stores/user";
import { showSuccess } from "../utils/toast";

const fileInputRef = ref(null);
const showFileModal = ref(false);
const showFolderModal = ref(false);

const selectedFile = ref(null);
const customFileName = ref("");
const folderName = ref("");

const userStore = useUserStore();

// File upload logic
function triggerFileInput() {
    fileInputRef.value.click();
}

function handleFileChange(event) {
    const file = event.target.files[0];
    if (file) {
        selectedFile.value = file;
        customFileName.value = "";
        showFileModal.value = true;
    }
}

function handleUpload() {
    if (!selectedFile.value) return;

    const formData = new FormData();
    formData.append("file", selectedFile.value);
    if (customFileName.value) {
        formData.append("name", customFileName.value);
    }

    axiosInstance
        .post("/api/file/create", formData, {
            headers: {
                "Content-Type": "multipart/form-data",
                Authorization: userStore.token,
            },
        })
        .then((res) => {
            console.log("Upload successful:", res);
            showFileModal.value = false;
            selectedFile.value = null;
            customFileName.value = "";
            // Show a toast or refresh file list if needed
            showSuccess("file uploaded successfully");
        })
        .catch((err) => {
            console.error("Upload failed:", err);
            // Show error toast if you want
        });
}

function closeFileModal() {
    showFileModal.value = false;
    selectedFile.value = null;
    customFileName.value = "";
}

// Folder creation logic
function openFolderModal() {
    folderName.value = "";
    showFolderModal.value = true;
}

function closeFolderModal() {
    showFolderModal.value = false;
    folderName.value = "";
}

function createFolder() {
    axiosInstance
        .post(
            "/api/folder/create",
            { name: folderName.value },
            {
                headers: {
                    Authorization: userStore.token,
                },
            }
        )
        .then((res) => {
            console.log("Folder created:", res);
            showFolderModal.value = false;
            folderName.value = "";
            showSuccess("folder created successfully");
        })
        .catch((err) => {
            console.error("Folder creation failed:", err);
        });
}
</script>
