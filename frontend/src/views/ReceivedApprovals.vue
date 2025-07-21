<template>
    <div class="px-2 w-[65%] mx-auto pt-30">
        <h2 class="text-xl font-bold mb-4">Approval Requests To My Files</h2>
        <table class="w-full bg-white rounded-xl shadow overflow-hidden">
            <thead>
                <tr class="bg-blue-50">
                    <th class="py-2 px-2 text-left">#</th>
                    <th class="py-2 px-2 text-left">From User</th>
                    <th class="py-2 px-2 text-left">For File</th>
                    <th class="py-2 px-2 text-left">Reason</th>
                    <th class="py-2 px-2 text-left">Created At</th>
                    <th class="py-2 px-2 text-left">Status</th>
                    <th class="py-2 px-2 text-left">Action</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(approval, idx) in requests" :key="approval._id">
                    <td class="py-2 px-2">{{ idx + 1 }}</td>
                    <td class="py-2 px-2">{{ approval.requester_name }}</td>
                    <td class="py-2 px-2">{{ approval.file_name }}</td>
                    <td
                        class="py-2 px-2 truncate max-w-[180px]"
                        :title="approval.reason"
                    >
                        {{ approval.reason || "-" }}
                    </td>
                    <td class="py-2 px-2 text-xs">
                        {{ formatDate(approval.created_at) }}
                    </td>
                    <td class="py-2 px-2 capitalize">
                        <span :class="statusClass(approval.status)">
                            {{ approval.status }}
                        </span>
                    </td>
                    <td class="py-2 px-2 flex gap-2">
                        <template v-if="approval.status === 'pending'">
                            <button
                                @click="approve(approval._id)"
                                class="cursor-pointer text-green-600 hover:bg-green-100 rounded px-2 py-1"
                            >
                                Approve
                            </button>
                            <button
                                @click="reject(approval._id)"
                                class="cursor-pointer text-red-500 hover:bg-red-100 rounded px-2 py-1"
                            >
                                Reject
                            </button>
                        </template>
                        <template v-else>
                            <span class="text-gray-400 italic">Reviewed</span>
                        </template>
                    </td>
                </tr>
                <tr v-if="requests.length === 0">
                    <td colspan="7" class="text-center py-6 text-gray-400">
                        No incoming requests.
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
const requests = ref([]);

function statusClass(status) {
    if (status === "pending") return "text-yellow-500 font-bold";
    if (status === "approved") return "text-green-600 font-bold";
    if (status === "rejected") return "text-red-500 font-bold";
    return "";
}
function formatDate(str) {
    return str ? new Date(str).toLocaleString() : "-";
}
function approve(id) {
    // call backend, set status to approved
}
function reject(id) {
    // call backend, set status to rejected
}
onMounted(() => {
    requests.value = [
        {
            _id: "2001",
            requester_name: "Alex N.",
            file_name: "Marketing_Strategy_Notes.docx",
            reason: "To add feedback from sales.",
            created_at: "2025-07-19T11:20:00Z",
            status: "pending",
        },
        {
            _id: "2002",
            requester_name: "Fatemeh S.",
            file_name: "Meeting_Recording.mp4",
            reason: "",
            created_at: "2025-07-18T16:45:00Z",
            status: "approved",
        },
        {
            _id: "2003",
            requester_name: "Ali M.",
            file_name: "Demo_Prototype.zip",
            reason: "Need to download and test.",
            created_at: "2025-07-17T08:10:00Z",
            status: "rejected",
        },
    ];
});
</script>
