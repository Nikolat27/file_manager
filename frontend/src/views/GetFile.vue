<template>
    <div class="pt-32 flex flex-col items-center justify-center min-h-[400px]">
        <div v-if="loading" class="text-blue-500 text-lg font-semibold">
            Loading file info...
        </div>

        <div
            v-else-if="error && !passwordModal && !approvalModal"
            class="text-red-600 text-[14px] font-semibold"
        >
            {{ error }}
        </div>

        <div v-else-if="fileReady">
            <div
                class="bg-gray-50 p-6 rounded-xl shadow-md max-w-lg w-full flex flex-col items-center gap-6"
            >
                <h2 class="text-xl font-bold mb-4">{{ fileName }}</h2>

                <!-- Image Preview -->
                <div
                    v-if="['png', 'jpg', 'jpeg'].includes(fileFormat)"
                    class="w-full flex flex-col items-center gap-2"
                >
                    <img
                        :src="fileUrl"
                        class="max-w-full max-h-96 rounded shadow"
                        :alt="fileName"
                    />
                    <a
                        :href="fileUrl"
                        :download="fileName"
                        class="mt-2 bg-blue-600 hover:bg-blue-800 text-white px-4 py-2 rounded"
                        >Download Image</a
                    >
                </div>

                <!-- PDF Download Only -->
                <div
                    v-else-if="fileFormat === 'pdf'"
                    class="flex flex-col items-center gap-2"
                >
                    <svg
                        class="w-16 h-16 text-red-400"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <rect
                            width="20"
                            height="24"
                            x="2"
                            y="0"
                            rx="4"
                            fill="#fee2e2"
                        />
                        <text
                            x="12"
                            y="16"
                            text-anchor="middle"
                            font-size="10"
                            fill="#b91c1c"
                        >
                            PDF
                        </text>
                    </svg>
                    <span class="font-mono">{{ fileName }}</span>
                    <a
                        :href="fileUrl"
                        :download="fileName"
                        class="mt-2 bg-blue-600 hover:bg-blue-800 text-white px-4 py-2 rounded"
                        >Download PDF</a
                    >
                </div>

                <!-- ZIP Download -->
                <div
                    v-else-if="fileFormat === 'zip'"
                    class="flex flex-col items-center gap-2"
                >
                    <svg
                        class="w-16 h-16 text-gray-400"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <rect
                            width="20"
                            height="24"
                            x="2"
                            y="0"
                            rx="4"
                            fill="#eee"
                        />
                        <path
                            d="M8 2h8v2H8zM8 6h8v2H8zM8 10h8v2H8zM8 14h8v2H8z"
                            fill="#ccc"
                        />
                    </svg>
                    <span class="font-mono">{{ fileName }}</span>
                    <a
                        :href="fileUrl"
                        :download="fileName"
                        class="mt-2 bg-blue-600 hover:bg-blue-800 text-white px-4 py-2 rounded"
                        >Download ZIP</a
                    >
                </div>

                <!-- Fallback -->
                <div v-else class="flex flex-col items-center gap-2">
                    <span class="text-gray-500">Preview not available.</span>
                    <a
                        :href="fileUrl"
                        :download="fileName"
                        class="bg-blue-600 hover:bg-blue-800 text-white px-4 py-2 rounded"
                        >Download {{ fileName }}</a
                    >
                </div>
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

        <!-- Approval Required Modal -->
        <div
            v-if="approvalModal"
            class="fixed inset-0 flex items-center justify-center bg-gray-200 bg-opacity-50 z-50"
        >
            <div class="bg-white p-8 rounded-xl shadow-lg max-w-sm w-full">
                <h2 class="text-lg font-bold mb-4 text-red-700">
                    Access Requires Approval
                </h2>
                <p class="text-gray-600 mb-4 text-sm">
                    This file requires the owner's approval. Please explain your
                    reason.
                </p>
                <textarea
                    v-model="approvalReason"
                    placeholder="Reason for request"
                    rows="3"
                    class="w-full border rounded px-3 py-2 mb-4"
                ></textarea>
                <div class="flex justify-end gap-2">
                    <button
                        class="bg-gray-400 text-white px-4 py-2 rounded"
                        @click="approvalModal = false"
                    >
                        Cancel
                    </button>
                    <button
                        class="bg-blue-600 hover:bg-blue-800 text-white px-4 py-2 rounded"
                        @click="sendApprovalRequest"
                    >
                        Send Request
                    </button>
                </div>
                <div
                    v-if="approvalMessage"
                    class="mt-4 text-sm font-semibold text-green-700"
                >
                    {{ approvalMessage }}
                </div>
            </div>
        </div>

        <!-- Status-only Modal -->
        <div
            v-if="showStatusModal"
            class="fixed inset-0 flex items-center justify-center bg-gray-200 bg-opacity-50 z-50"
        >
            <div
                class="bg-white p-6 rounded-xl shadow max-w-sm w-full text-center"
            >
                <h2 class="text-lg font-bold mb-4 text-blue-700">
                    Approval Info
                </h2>
                <p class="text-gray-700">{{ approvalMessage }}</p>
                <button
                    class="mt-4 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-800"
                    @click="showStatusModal = false"
                >
                    Close
                </button>
            </div>
        </div>
    </div>
</template>
<script setup>
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import axiosInstance from "../axiosInstance";

const route = useRoute();
const loading = ref(true);
const error = ref("");
const passwordModal = ref(false);
const approvalModal = ref(false);
const approvalMessage = ref("");
const approvalReason = ref("");
const password = ref("");

const shortId = route.params.id;
const realFileId = ref(route.query.file_id || "");

const fileUrl = ref("");
const fileFormat = ref("");
const fileName = ref("");
const fileReady = ref(false);

const showStatusModal = ref(false);

async function fetchFile() {
    loading.value = true;
    error.value = "";
    try {
        const resp = await axiosInstance.get(`/api/file/get/${shortId}`);
        showFile(resp.data);
    } catch (err) {
        const status = err.response?.status;
        const data = err.response?.data || {};

        if (status === 428) {
            const type = data.type;

            if (type === "not_requested") {
                // show approval modal
                approvalModal.value = true;
            } else if (type === "pending") {
                approvalMessage.value =
                    "Your approval is pending. Please wait.";
                showStatusModal.value = true;
            } else if (type === "rejected") {
                approvalMessage.value = "Your approval was rejected.";
                showStatusModal.value = true;
            } else {
                approvalMessage.value =
                    "Your approval status is invalid. Please contact our admins";
                showStatusModal.value = true;
            }
        }
    } finally {
        loading.value = false;
    }
}

async function submitPassword() {
    loading.value = true;
    error.value = "";
    try {
        const resp = await axiosInstance.post(`/api/file/get/${shortId}`, {
            password: password.value,
        });
        passwordModal.value = false;
        showFile(resp.data);
    } catch (err) {
        const status = err.response?.status;
        const data = err.response?.data || {};
        if (status === 428) {
            passwordModal.value = false;
            approvalModal.value = true;
            if (!realFileId.value) realFileId.value = data.file_id;
        } else {
            error.value = data.message || "Incorrect password. Try again";
        }
    } finally {
        loading.value = false;
        password.value = "";
    }
}

async function sendApprovalRequest() {
    console.log(realFileId.value);

    if (!approvalReason.value.trim()) {
        approvalMessage.value = "Please enter a reason.";
        return;
    }

    try {
        await axiosInstance.post(`/api/approval/create`, {
            file_id: realFileId.value,
            reason: approvalReason.value.trim(),
        });
        approvalMessage.value =
            "Approval request sent successfully. Please wait for confirmation.";
    } catch (err) {
        approvalMessage.value =
            err.response?.data || "Failed to send approval request.";
    }
}

function showFile(data) {
    const backendBaseUrl =
        import.meta.env.backendBaseUrl || "http://localhost:8000";
    const staticUrl = backendBaseUrl + "/static/";

    fileUrl.value = staticUrl + data.file_address;
    fileName.value = data.name || fileUrl.value.split("/").pop();
    fileFormat.value = getFileFormat(fileUrl.value);
    fileReady.value = true;
}

function getFileFormat(fileUrl) {
    const pathname = fileUrl.split("?")[0].split("#")[0];
    const parts = pathname.split("/");
    const filename = parts.pop() || "";
    const dotIndex = filename.lastIndexOf(".");
    if (dotIndex === -1) return "";
    return filename.slice(dotIndex + 1).toLowerCase();
}

onMounted(fetchFile);
</script>
