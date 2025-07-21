<template>
    <div class="pt-30 px-2 w-[65%] mx-auto">
        <h2 class="text-xl font-bold mb-4">My Sent Approval Requests</h2>
        <table class="w-full bg-white rounded-xl shadow overflow-hidden">
            <thead>
                <tr class="bg-blue-50">
                    <th class="py-2 px-2 text-left">#</th>
                    <th class="py-2 px-2 text-left">To File</th>
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
                            @click="deleteApproval(approval._id)"
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
    approvals.value = approvals.value.filter((a) => a._id !== id);
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

    approvals.value = [
        {
            _id: "1001",
            file_name: "Team_Plan_2025.pdf",
            status: "pending",
            reason: "Need to review latest updates for my task.",
            created_at: "2025-07-18T14:15:00Z",
            reviewed_at: null,
        },
        {
            _id: "1002",
            file_name: "Project_Schema.png",
            status: "approved",
            reason: "",
            created_at: "2025-07-16T10:40:00Z",
            reviewed_at: "2025-07-16T13:22:00Z",
        },
        {
            _id: "1003",
            file_name: "Financial_Report.xlsx",
            status: "rejected",
            reason: "Required for annual review.",
            created_at: "2025-07-13T09:00:00Z",
            reviewed_at: "2025-07-13T12:30:00Z",
        },
    ];
});
</script>
