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
                    <button
                        @click="downloadFile()"
                        class="cursor-pointer mt-2 bg-blue-600 hover:bg-blue-800 text-white px-4 py-2 rounded"
                    >
                        Download Image
                    </button>
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
                    >
                        Download PDF
                    </a>
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
                    >
                        Download ZIP
                    </a>
                </div>
                <!-- Fallback -->
                <div v-else class="flex flex-col items-center gap-2">
                    <span class="text-gray-500">Preview not available.</span>
                    <a
                        :href="fileUrl"
                        :download="fileName"
                        class="bg-blue-600 hover:bg-blue-800 text-white px-4 py-2 rounded"
                    >
                        Download {{ fileName }}
                    </a>
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
                <div
                    v-if="error"
                    class="mt-2 text-red-600 text-lg font-semibold"
                >
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
                        :disabled="!approvalReason.trim() || loading"
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

        <!-- Approval/Status Modal -->
        <div
            v-if="showStatusModal"
            class="fixed inset-0 flex items-center justify-center bg-gray-200 bg-opacity-50 z-50"
        >
            <div
                class="bg-white p-8 rounded-xl shadow-lg max-w-sm w-full text-center"
            >
                <h2 class="text-lg font-bold mb-4 text-blue-700">
                    Approval Info
                </h2>
                <p class="text-gray-700 mb-6">{{ approvalMessage }}</p>

                <!-- Show 'Continue' button only if approved -->
                <button
                    v-if="approvalStatus === 'approved'"
                    @click="
                        () => {
                            showStatusModal = false;
                            handleFileAccess(shortUrl);
                        }
                    "
                    class="bg-blue-600 hover:bg-blue-800 text-white px-4 py-2 rounded"
                    autofocus
                >
                    Continue
                </button>
                <!-- Show only a Close button otherwise -->
                <button
                    v-else
                    @click="showStatusModal = false"
                    class="bg-gray-400 text-white px-4 py-2 rounded"
                >
                    Close
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { showSuccess, showError, showInfo } from "../utils/toast";
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
const approvalStatus = ref(""); // pending, approved, rejected, not_requested

const shortUrl = route.params.id;

const fileUrl = ref("");
const fileFormat = ref("");
const fileName = ref("");
const fileReady = ref(false);

const showStatusModal = ref(false);

function resetModals() {
    passwordModal.value = false;
    approvalModal.value = false;
    showStatusModal.value = false;
    error.value = "";
    approvalMessage.value = "";
    approvalReason.value = "";
    approvalStatus.value = "";
}

function parseError(err) {
    let data = err?.response?.data;
    if (typeof data === "string") {
        try {
            data = JSON.parse(data);
        } catch {
            data = { error: data };
        }
    }
    if (!data || typeof data !== "object") data = { error: "Unknown error" };
    return data;
}

function showFile(data) {
    const backendBaseUrl =
        import.meta.env.backendBaseUrl || "http://localhost:8000";
    const staticUrl = backendBaseUrl + "/static/";

    fileUrl.value = staticUrl + data.file_address;
    fileName.value = data.file_address.split("/").pop() || "file";
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

async function handleFileAccess(fileShortUrl) {
    loading.value = true;
    error.value = "";
    resetModals();

    try {
        let resp;
        if (password.value) {
            // If we have password, try POSTing with it
            resp = await axiosInstance.post(`/api/file/get/${fileShortUrl}`, {
                password: password.value,
            });
        } else {
            // Else, try GET (public or approval only)
            resp = await axiosInstance.get(`/api/file/get/${fileShortUrl}`);
        }
        showFile(resp.data);
    } catch (err) {
        const data = parseError(err);

        if (
            data.error === "password is required" ||
            data.error === "incorrect password"
        ) {
            passwordModal.value = true;
            error.value = data.error || "Password is required.";
            // Clear wrong password on error
            password.value = "";
        } else if (err.response?.status === 428) {
            await checkApprovalStatus(fileShortUrl);
        } else {
            error.value = data.error || "Failed to fetch file data.";
            showError(error.value);
        }
    } finally {
        loading.value = false;
    }
}

// Step 2: Inquiry about approval status, then act accordingly
async function checkApprovalStatus(fileShortUrl) {
    try {
        const resp = await axiosInstance.get(
            `/api/approval/check/${fileShortUrl}`
        );
        // If approved (owner or approval granted)
        approvalStatus.value = "approved";
        approvalMessage.value =
            resp.data?.message ||
            "Approval granted. Please continue to access the file.";
        showStatusModal.value = true; // Show modal with "Continue" button
    } catch (err) {
        // Approval not granted yet
        const data = parseError(err);
        approvalStatus.value = data.error || "unknown";
        switch (approvalStatus.value) {
            case "pending":
                showStatusModal.value = true;
                approvalMessage.value =
                    "Your approval is pending. Please wait for the owner to approve.";
                break;
            case "rejected":
                showStatusModal.value = true;
                approvalMessage.value =
                    "Your approval was rejected. You may contact the owner or try again later.";
                break;
            case "not_requested":
                approvalModal.value = true; // Show the approval request form
                break;
            default:
                error.value = "Unknown approval status.";
                showError(error.value);
        }
    }
}

async function submitPassword() {
    loading.value = true;
    error.value = "";
    try {
        const resp = await axiosInstance.post(`/api/file/get/${shortUrl}`, {
            password: password.value,
        });
        passwordModal.value = false;
        showFile(resp.data);
    } catch (err) {
        if (err.response?.status === 428) {
            await checkApprovalStatus(shortUrl);
        }
        const data = parseError(err);
        error.value = data.error || "Incorrect password. Try again";
    } finally {
        loading.value = false;
    }
}

async function sendApprovalRequest() {
    loading.value = true;
    try {
        await axiosInstance.post(`/api/approval/create`, {
            short_url: shortUrl,
            reason: approvalReason.value.trim(),
        });
        showInfo("Approval request sent. Please wait for confirmation.");
        approvalModal.value = false;
        approvalMessage.value =
            "Approval request sent. Please wait for confirmation.";
        showStatusModal.value = true;
        approvalStatus.value = "pending";
    } catch (err) {
        const data = parseError(err);
        approvalMessage.value =
            data.error || "Failed to send approval request.";
        showError(approvalMessage.value);
    } finally {
        loading.value = false;
    }
}

async function downloadFile() {
    try {
        const res = await axiosInstance.get(`/api/file/download/${shortUrl}`, {
            responseType: "blob",
        });
        const url = URL.createObjectURL(res.data);
        const link = document.createElement("a");
        link.href = url;
        link.download = fileName.value || `${shortUrl}.bin`;
        document.body.appendChild(link);
        link.click();
        link.remove();
        URL.revokeObjectURL(url);
        showSuccess("Download started");
    } catch (err) {
        showError("Download failed");
    }
}

onMounted(() => handleFileAccess(shortUrl));
</script>
