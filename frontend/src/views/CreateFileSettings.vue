<template>
    <div v-if="!fileId" class="text-red-600 text-xl text-center mt-10">
        Error: File ID is required!
    </div>
    <div v-else>
        <div
            class="min-h-screen flex items-center justify-center bg-blue-50 z-20"
        >
            <form
                class="bg-white rounded-2xl shadow-xl p-8 w-full max-w-lg flex flex-col gap-6"
                @submit.prevent="onSave"
            >
                <h2 class="text-2xl font-bold text-blue-700 mb-2">
                    Create File Settings
                </h2>
                <div
                    class="mb-4 text-blue-800 font-semibold flex items-center gap-2"
                >
                    <span>Your Plan:</span>
                    <span
                        class="px-3 py-1 rounded-xl bg-blue-100 border border-blue-300 text-blue-700 text-base"
                    >
                        {{ plan.toUpperCase() }}
                    </span>
                </div>
                <div class="text-blue-500 mb-4">
                    <span class="font-semibold">File ID:</span> {{ fileId }}
                </div>

                <!-- Password -->
                <div>
                    <label class="block text-blue-800 font-medium mb-1"
                        >Password</label
                    >
                    <input
                        v-model="password"
                        type="password"
                        placeholder="Set a password (optional)"
                        class="w-full px-4 py-3 rounded-xl border border-blue-300 focus:outline-none focus:border-blue-500 bg-blue-50"
                    />
                </div>

                <!-- Approver Required (disabled if free) -->
                <div class="flex items-center gap-2">
                    <input
                        v-model="approvable"
                        type="checkbox"
                        id="approvable"
                        class="accent-blue-600"
                        :disabled="isPlanFree"
                    />
                    <label
                        for="approvable"
                        class="text-blue-800 font-medium"
                        :class="
                            isPlanFree ? 'opacity-60 cursor-not-allowed' : ''
                        "
                    >
                        Approver required
                    </label>
                </div>

                <!-- View Only (disabled if free) -->
                <div class="flex items-center gap-2">
                    <input
                        v-model="viewOnly"
                        type="checkbox"
                        id="viewOnly"
                        class="accent-blue-600"
                        :disabled="isPlanFree"
                        :class="
                            isPlanFree
                                ? 'opacity-60 backdrop-blur cursor-not-allowed'
                                : ''
                        "
                    />
                    <label
                        for="viewOnly"
                        class="text-blue-800 font-medium"
                        :class="
                            isPlanFree ? 'opacity-60 cursor-not-allowed' : ''
                        "
                    >
                        View Only
                    </label>
                </div>

                <!-- Expiration Date -->
                <div>
                    <label
                        class="block text-blue-800 font-medium mb-1"
                        :class="
                            isPlanFree ? 'opacity-60  cursor-not-allowed' : ''
                        "
                    >
                        Expiration Date
                    </label>
                    <input
                        v-model="expirationDate"
                        type="date"
                        :disabled="isPlanFree"
                        class="w-full px-4 py-3 rounded-xl border border-blue-300 focus:outline-none focus:border-blue-500 bg-blue-50 disabled:bg-gray-100"
                        :class="
                            isPlanFree
                                ? 'opacity-60 backdrop-blur cursor-not-allowed'
                                : ''
                        "
                    />
                    <p v-if="isPlanFree" class="text-sm text-blue-600 mt-2">
                        Default expiration is
                        <span class="font-semibold"
                            >7 days for free plan users</span
                        >.
                    </p>
                </div>

                <!-- Max Downloads (disabled if free) -->
                <div>
                    <label
                        class="block text-blue-800 font-medium mb-1"
                        :class="
                            isPlanFree ? 'opacity-60 cursor-not-allowed' : ''
                        "
                    >
                        Max Downloads
                    </label>
                    <input
                        v-model.number="maxDownloads"
                        type="number"
                        min="1"
                        :disabled="isPlanFree"
                        placeholder="Unlimited"
                        class="w-full px-4 py-3 rounded-xl border border-blue-300 focus:outline-none focus:border-blue-500 bg-blue-50 disabled:bg-gray-100"
                        :class="
                            isPlanFree
                                ? 'opacity-60 backdrop-blur cursor-not-allowed'
                                : ''
                        "
                    />
                </div>

                <button
                    type="submit"
                    class="w-full py-3 bg-blue-600 hover:bg-blue-700 transition-colors text-white rounded-xl font-semibold text-lg shadow"
                >
                    Create Settings
                </button>
            </form>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import axiosInstance from "../axiosInstance";
import { useUserStore } from "../stores/user";
import { showError, showSuccess } from "../utils/toast";

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();

const fileId = route.params.id;
const folderId = route.params.folderId || route.query.folderId || "";
const plan = userStore.plan || "free";

if (!fileId) {
    router.replace({ name: "NotFound" });
}

// Form fields
const password = ref("");
const approvable = ref(false);
const viewOnly = ref(false);
const expirationDate = ref("");
const maxDownloads = ref(null);

const isPlanFree = computed(() => plan === "free");

function onSave() {
    const formData = new FormData();
    formData.append("password", password.value);
    formData.append("approvable", !isPlanFree.value ? approvable.value : false);
    formData.append("view_only", !isPlanFree.value ? viewOnly.value : false);
    formData.append(
        "expiration_at",
        !isPlanFree.value && expirationDate.value
            ? new Date(expirationDate.value).toISOString()
            : ""
    );

    if (!maxDownloads.value || !isPlanFree.value) {
        formData.append("max_downloads", -1); // default (unlimited)
    } else {
        formData.append("max_downloads", maxDownloads.value);
    }

    // Add folder_id
    if (folderId) {
        formData.append("folder_id", folderId);
    }

    axiosInstance
        .post(`/api/file/settings/create/${fileId}`, formData, {
            headers: { "Content-Type": "multipart/form-data" },
        })
        .then(() => {
            showSuccess("Settings saved successfully");
            router.push({ name: "sharedUrls" });
        })
        .catch((err) => {
            showError(err.response.data);
        });
}
</script>
