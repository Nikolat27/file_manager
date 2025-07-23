<template>
    <div class="w-full max-w-4xl mx-auto py-12 pt-24">
        <h1 class="text-2xl font-bold mb-6">My Teams</h1>
        <div v-if="teams?.length === 0" class="text-gray-500">
            No teams found.
        </div>

        <div
            v-for="team in teams"
            :key="team.id"
            class="flex flex-row items-center bg-white rounded-2xl shadow-md mb-5 px-6 py-4 gap-6"
        >
            <!-- Avatar -->
            <img
                :src="getTeamAvatarUrl(team.avatar_url) || defaultAvatar"
                alt="avatar"
                class="w-20 h-20 rounded-full object-cover border-2 border-gray-300 bg-gray-100"
            />

            <div class="flex-1 flex flex-col gap-1">
                <div class="flex items-center gap-3">
                    <span class="text-lg font-bold">{{ team.name }}</span>
                    <span
                        class="px-3 py-0.5 rounded-lg text-xs font-semibold"
                        :class="
                            team.plan === 'premium'
                                ? 'bg-yellow-300 text-yellow-900'
                                : 'bg-gray-200 text-gray-700'
                        "
                    >
                        {{
                            team.plan.charAt(0).toUpperCase() +
                            team.plan.slice(1)
                        }}
                    </span>
                </div>
                <div class="text-gray-700 text-sm">{{ team.description }}</div>
                <div
                    class="flex flex-wrap gap-4 mt-1 text-[13px] text-gray-600"
                >
                    <span><b>Users:</b> {{ team.users?.length }}</span>
                    <span
                        ><b>Storage:</b>
                        {{ formatMB(team.storage_used) }} MB</span
                    >
                    <span
                        ><b>Updated:</b> {{ formatDate(team.updated_at) }}</span
                    >
                    <span
                        ><b>Created:</b> {{ formatDate(team.created_at) }}</span
                    >
                </div>
            </div>

            <div
                class="flex flex-row gap-2 justify-between items-end ml-4 h-full"
            >
                <button
                    @click="goToTeamFiles(team.id)"
                    class="cursor-pointer px-4 py-2 rounded-xl bg-blue-600 text-white font-semibold hover:bg-blue-700"
                >
                    View
                </button>
                <button
                    v-if="team.owner_id === userStore.id"
                    @click="deleteTeam(team.id)"
                    class="cursor-pointer px-4 py-2 rounded-xl bg-red-600 text-white font-semibold hover:bg-red-700"
                >
                    Delete
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import axiosInstance from "../axiosInstance";
import { useUserStore } from "../stores/user";
import { showSuccess, showError } from "../utils/toast";
import defaultAvatar from "../assets/images/Avatar-16-512.webp";

const userStore = useUserStore();

const router = useRouter();

function goToTeamFiles(teamId) {
    router.push({ name: "TeamFiles", params: { id: teamId } });
}

// Sample data (replace with API call in real usage)
const teams = ref([]);

function fetchUserTeams() {
    axiosInstance
        .get("/api/team/get")
        .then((resp) => {
            teams.value = resp.data;
            console.log(resp.data);
        })
        .catch((err) => {
            console.error(err);
        });
}

function deleteTeam(teamId) {
    axiosInstance
        .delete(`/api/team/delete/${teamId}`)
        .then(() => {
            showSuccess("deleted team successfully");
            fetchUserTeams();
        })
        .catch((err) => {
            showError(err.response.data.error);
        });
}

// Convert bytes to MB with 2 decimal digits
function formatMB(bytes) {
    if (!bytes) return "0";
    return (bytes / (1024 * 1024)).toFixed(2);
}

// Format date to readable string (YYYY-MM-DD)
function formatDate(dateStr) {
    if (!dateStr) return "";
    const d = new Date(dateStr);
    return d.toISOString().split("T")[0];
}

function getTeamAvatarUrl(avatarUrl) {
    if (!avatarUrl) {
        return null;
    }

    const VITE_BACKEND_BASE_URL = import.meta.env.VITE_BACKEND_BASE_URL || "http://localhost:8000";
    const staticUrl = VITE_BACKEND_BASE_URL + "/static/";

    return staticUrl + avatarUrl;
}

onMounted(() => {
    fetchUserTeams();
});
</script>

<style scoped>
div[role="team-row"]:hover {
    box-shadow: 0 2px 16px 0 rgba(0, 0, 0, 0.08);
    background: #f6f8fc;
}
</style>
