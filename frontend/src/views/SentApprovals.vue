<template>
    <div class="pt-30 px-2 w-[65%] mx-auto">
        <h2 class="text-xl font-bold mb-4">My Sent Approval Requests</h2>
        <table class="w-full bg-white rounded-xl shadow overflow-hidden">
            <thead>
                <tr class="bg-blue-50">
                    <th class="py-2 px-2 text-left">#</th>
                    <th class="py-2 px-2 text-left">File Name</th>
                    <th class="py-2 px-2 text-left">Status</th>
                    <th class="py-2 px-2 text-left">Reason</th>
                    <th class="py-2 px-2 text-left">Created At</th>
                    <th class="py-2 px-2 text-left">Reviewed At</th>
                    <th class="py-2 px-2 text-left">Delete</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(approval, idx) in approvals" :key="approval._id">
                    <td class="py-2 px-2">{{ idx + 1 }}</td>
                    <td class="py-2 px-2">{{ approval.file_name }}</td>
                    <td class="py-2 px-2 capitalize">
                        <span :class="statusClass(approval.status)">
                            {{ approval.status }}
                        </span>
                    </td>
                    <td
                        class="py-2 px-2 truncate max-w-[180px]"
                        :title="approval.reason"
                    >
                        {{ approval.reason || "-" }}
                    </td>
                    <td class="py-2 px-2 text-xs">
                        {{ formatDate(approval.created_at) }}
                    </td>
                    <td class="py-2 px-2 text-xs">
                        {{
                            approval.reviewed_at
                                ? formatDate(approval.reviewed_at)
                                : "-"
                        }}
                    </td>
                    <td class="py-2 px-2 text-center">
                        <button
                            @click="deleteApproval(approval.id)"
                            class="cursor-pointer text-red-500 hover:bg-red-100 rounded px-3 py-1"
                        >
                            🗑️
                        </button>
                    </td>
                </tr>
                <tr v-if="approvals.length === 0">
                    <td colspan="7" class="text-center py-6 text-gray-400">
                        No approvals sent.
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import axiosInstance from "../axiosInstance";
import { showError, showSuccess } from "../utils/toast";
const approvals = ref([]);

function statusClass(status) {
    if (status === "pending") return "text-yellow-500 font-bold";
    if (status === "approved") return "text-green-600 font-bold";
    if (status === "rejected") return "text-red-500 font-bold";
    return "";
}

function formatDate(str) {
    return str ? new Date(str).toLocaleString() : "-";
}

function deleteApproval(id) {
    axiosInstance
        .delete(`/api/approval/delete/${id}`)
        .then((resp) => {
            console.log(resp.data);
            showSuccess("approve request deleted successfully");
            approvals.value = approvals.value.filter((a) => a.id !== id);
        })
        .catch((err) => {
            showError(err.response.data.error);
        });
}

function fetchSentApprovals() {
    axiosInstance
        .get("/api/approval/get")
        .then((resp) => {
            if (!resp.data.approvals) {
                approvals.value = [];
            } else {
                approvals.value = resp.data.approvals;
                console.log(approvals.value);
            }
        })
        .catch((err) => {
            console.error(err.response.data);
        });
}

onMounted(() => {
    fetchSentApprovals();
});
</script>
