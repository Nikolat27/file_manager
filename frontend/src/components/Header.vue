<template>
    <aside
        class="w-full h-[72px] flex flex-row items-center absolute top-0 left-[300px] px-6"
    >
        <input
            v-model="searchText"
            type="text"
            class="w-[922px] h-[40px] hover:outline-2 outline-gray-400 outline-1 hover:outline-black rounded-2xl pl-4"
            placeholder="Search for file names"
        />
        <button
            :class="
                searchText?.length === 0
                    ? 'cursor-default opacity-60'
                    : 'cursor-pointer opacity-100'
            "
            class="ml-4 text-white bg-blue-600 rounded-full w-20 h-10 font-semibold"
            @click="handleSearch"
        >
            search
        </button>

        <div
            class="h-[40px] absolute right-[300px] flex flex-row items-center gap-x-4 mr-6"
        >
            <button
                @click="openAvatarModal"
                class="rounded-full w-[32px] h-[32px] cursor-pointer"
            >
                <img
                    class="w-full h-full"
                    :src="userStore.avatarUrl || defaultAvatarUrl"
                    alt=""
                />
            </button>
        </div>
    </aside>
    <!-- User Avatar Upload Modal -->
    <div
        v-if="showAvatarModal"
        class="fixed inset-0 bg-white bg-opacity-20 flex items-center justify-center z-50"
    >
        <div
            class="bg-black text-white rounded-xl p-8 flex flex-col items-center gap-4 w-[90vw] max-w-md shadow-2xl"
        >
            <span class="font-semibold"
                >Your user id:
                <span class="underline">{{ userStore.id }}</span></span
            >
            <h2 class="text-xl font-semibold mb-2">Change Profile Avatar</h2>

            <!-- Avatar Image Preview -->
            <div v-if="avatarPreview" class="mb-2 rounded-full">
                <img
                    :src="avatarPreview"
                    alt="avatar preview"
                    class="rounded-full w-24 h-24 object-cover border-2 border-gray-500"
                />
            </div>

            <input
                type="file"
                ref="avatarInputRef"
                class="hidden"
                accept="image/*"
                @change="handleAvatarChange"
            />
            <button
                @click="triggerAvatarInput"
                class="bg-blue-700 text-white px-4 py-1 rounded-xl font-semibold"
            >
                Choose Image
            </button>
            <span v-if="avatarFileName" class="text-xs text-gray-200">{{
                avatarFileName
            }}</span>

            <div class="flex gap-4 mt-4">
                <button
                    @click="uploadAvatar"
                    :disabled="!avatarFile"
                    class="cursor-pointer px-6 py-2 rounded-xl bg-blue-600 text-white font-semibold hover:bg-blue-700 disabled:opacity-60"
                >
                    Upload
                </button>
                <button
                    @click="closeAvatarModal"
                    class="cursor-pointer px-6 py-2 rounded-xl bg-gray-200 text-gray-800 font-semibold hover:bg-gray-300"
                >
                    Cancel
                </button>
            </div>
        </div>
    </div>
</template>
<script setup>
import { showError, showSuccess } from "../utils/toast";
import { ref } from "vue";
import axiosInstance from "../axiosInstance";
import { useUserStore } from "../stores/user";
import defaultAvatarUrl from "../assets/images/images.png";

import { useRouter } from "vue-router";

const searchText = ref("");
const router = useRouter();

function handleSearch() {
    if (!searchText.value) {
        return;
    }

    router.push({ name: "UserSearch", query: { q: searchText.value } });
    return;
}

const userStore = useUserStore();

const showAvatarModal = ref(false);
const avatarFile = ref(null);
const avatarPreview = ref("");
const avatarFileName = ref("");
const avatarInputRef = ref(null);

function openAvatarModal() {
    showAvatarModal.value = true;
    avatarFile.value = null;
    avatarFileName.value = "";
    avatarPreview.value = "";
}

function closeAvatarModal() {
    showAvatarModal.value = false;
    avatarFile.value = null;
    avatarFileName.value = "";
    avatarPreview.value = "";
}

function triggerAvatarInput() {
    avatarInputRef.value.click();
}

function handleAvatarChange(e) {
    const file = e.target.files[0];
    if (file) {
        avatarFile.value = file;
        avatarFileName.value = file.name;

        // Show image preview
        const reader = new FileReader();
        reader.onload = (event) => {
            avatarPreview.value = event.target.result;
        };
        reader.readAsDataURL(file);
    }
}

async function uploadAvatar() {
    if (!avatarFile.value) return;
    const formData = new FormData();
    formData.append("file", avatarFile.value);

    try {
        await axiosInstance.post("/api/user/avatar/upload", formData);
        showSuccess("Avatar updated!");
        closeAvatarModal();
        getUserData();
    } catch (err) {
        showError(err.response.data.error || "Failed to upload avatar");
    }
}

async function getUserData() {
    axiosInstance.get("/api/user/get").then((resp) => {
        const VITE_BACKEND_BASE_URL =
            import.meta.env.VITE_BACKEND_BASE_URL || "http://localhost:8000";

        const staticUrl = VITE_BACKEND_BASE_URL + "/static/";
        let avatarUrl = staticUrl + resp.data.avatar_url;
        if (resp.data.avatarUrl) {
            avatarUrl = null;
        }

        // only update the avatar url
        userStore.avatarUrl = avatarUrl;
    });
}
</script>
