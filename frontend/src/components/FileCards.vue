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

        <!-- Create Team Button -->
        <div
            @click="openTeamModal"
            class="cursor-pointer w-[140px] flex items-center justify-center h-[76px] rounded-2xl bg-black border-1 text-white border-gray-500 hover:opacity-70"
        >
            <span>Create Team</span>
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

    <!-- Add this just below your other modals -->
    <!-- Create Team Modal -->
    <div
        v-if="showTeamModal"
        class="fixed inset-0 bg-white bg-opacity-20 flex items-center justify-center z-50"
    >
        <div
            class="bg-black text-white rounded-xl p-8 flex flex-col items-center gap-4 w-[90vw] max-w-md shadow-2xl"
        >
            <h2 class="text-xl font-semibold mb-2">Create Team</h2>
            <input
                v-model="teamName"
                type="text"
                placeholder="Team name"
                class="w-full px-4 py-2 text-white rounded border border-gray-300 focus:outline-none font-semibold"
            />

            <!-- Avatar Upload -->
            <div class="w-full flex items-center gap-2">
                <input
                    type="file"
                    ref="teamAvatarInputRef"
                    class="hidden"
                    accept="image/*"
                    @change="handleTeamAvatarChange"
                />
                <button
                    @click="triggerTeamAvatarInput"
                    class="bg-blue-700 text-white px-4 py-1 rounded-xl font-semibold"
                >
                    Upload Avatar
                </button>
                <span v-if="teamAvatarName" class="text-xs text-gray-200">{{
                    teamAvatarName
                }}</span>
            </div>

            <textarea
                v-model="teamDescription"
                placeholder="Description (optional)"
                class="w-full px-4 py-2 text-white rounded border border-gray-300 focus:outline-none font-semibold resize-none"
                rows="3"
            ></textarea>

            <!-- Plan Select -->
            <select
                v-model="teamPlan"
                class="w-full px-4 py-2 text-white rounded border border-gray-300 bg-black focus:outline-none font-semibold"
            >
                <option value="free">Free</option>
                <option value="premium">Premium</option>
            </select>

            <div class="flex gap-4 mt-4">
                <button
                    @click="createTeam"
                    class="px-6 py-2 rounded-xl bg-blue-600 text-white font-semibold hover:bg-blue-700"
                >
                    Create
                </button>
                <button
                    @click="closeTeamModal"
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
import { showError, showSuccess } from "../utils/toast";

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
        formData.append("file_name", customFileName.value);
    }

    axiosInstance
        .post("/api/file/create", formData, {
            headers: {
                "Content-Type": "multipart/form-data",
                Authorization: userStore.token,
            },
        })
        .then(() => {
            showFileModal.value = false;
            selectedFile.value = null;
            customFileName.value = "";

            // Show a toast or refresh file list if needed
            showSuccess("file uploaded successfully");
            window.location.reload();
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

const showTeamModal = ref(false);
const teamName = ref("");
const teamAvatar = ref(null);
const teamAvatarName = ref(""); // for showing filename
const teamDescription = ref("");
const teamPlan = ref("free");

const teamAvatarInputRef = ref(null);

function openTeamModal() {
    teamName.value = "";
    teamAvatar.value = null;
    teamAvatarName.value = "";
    teamDescription.value = "";
    teamPlan.value = "free";
    showTeamModal.value = true;
}
function closeTeamModal() {
    showTeamModal.value = false;
    teamName.value = "";
    teamAvatar.value = null;
    teamAvatarName.value = "";
    teamDescription.value = "";
    teamPlan.value = "free";
}
function triggerTeamAvatarInput() {
    teamAvatarInputRef.value.click();
}
function handleTeamAvatarChange(e) {
    const file = e.target.files[0];
    if (file) {
        teamAvatar.value = file;
        teamAvatarName.value = file.name;
    }
}

async function createTeam() {
    if (!teamName.value) {
        showError("Team name is required!");
        return;
    }

    const formData = new FormData();
    formData.append("name", teamName.value);
    if (teamAvatar.value) formData.append("file", teamAvatar.value);
    if (teamDescription.value)
        formData.append("description", teamDescription.value);

    try {
        await axiosInstance.post("/api/team/create", formData, {
            headers: {
                "Content-Type": "multipart/form-data",
            },
        });
        showSuccess("team created successfully");
        closeTeamModal();
    } catch (err) {
        showError(err.response.data.error);
        console.error("Team creation failed:", err);
    }
}
</script>
