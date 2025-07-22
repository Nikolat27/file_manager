<template>
    <aside
        class="w-full h-[72px] flex flex-row items-center absolute top-0 left-[300px] px-6"
    >
        <input
            type="text"
            class="w-[922px] h-[40px] hover:outline-2 outline-gray-400 outline-1 hover:outline-black rounded-2xl pl-4"
            placeholder="Search for file names"
        />

        <div
            class="h-[40px] absolute right-[300px] flex flex-row items-center gap-x-4 mr-6"
        >
            <button class="rounded-full w-[32px] h-[32px] cursor-pointer">
                <svg
                    viewBox="0 0 24 24"
                    fill="none"
                    class="dig-UIIcon"
                    width="32"
                    height="32"
                    role="presentation"
                    focusable="false"
                    color="inherit"
                >
                    <path
                        d="m17.608 12.971-.329-.219a1.747 1.747 0 0 1-.779-1.457v-1.67c0-1.094 0-2.332-.563-3.336C15.261 5.084 13.972 4.5 12 4.5c-1.972 0-3.26.585-3.937 1.787C7.5 7.292 7.5 8.531 7.5 9.624v1.672a1.747 1.747 0 0 1-.78 1.454l-.327.219A4.241 4.241 0 0 0 4.5 16.507v.993H10a1.857 1.857 0 0 0 2 2 1.857 1.857 0 0 0 2-2h5.5v-.993a4.242 4.242 0 0 0-1.892-3.536ZM6.047 16a2.743 2.743 0 0 1 1.178-1.781L7.553 14A3.244 3.244 0 0 0 9 11.296V9.622c0-.953 0-1.938.37-2.6C9.618 6.584 10.16 6 12 6c1.841 0 2.383.584 2.63 1.023.371.662.371 1.646.37 2.6v1.674A3.244 3.244 0 0 0 16.447 14l.329.219A2.744 2.744 0 0 1 17.953 16H6.047Z"
                        fill="currentColor"
                        vector-effect="non-scaling-stroke"
                    ></path>
                </svg>
            </button>
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
    axiosInstance.get("/api/user").then((resp) => {
        const backendUrl =
            import.meta.env.backendUrl || "http://localhost:8000";

        const staticUrl = backendUrl + "/static/";
        let avatarUrl = staticUrl + resp.data.avatar_url;
        if (resp.data.avatarUrl) {
            avatarUrl = null;
        }

        // only update the avatar url
        userStore.avatarUrl = avatarUrl;
    });
}
</script>
