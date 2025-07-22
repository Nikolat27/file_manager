<template>
    <div
        class="min-h-screen bg-gray-50 py-12 pt-20 px-4 flex flex-col items-center"
    >
        <h1 class="text-3xl font-bold mb-10">Plans & Features</h1>
        <button
            @click="openUpgradeModal"
            class="mb-8 bg-blue-700 hover:bg-blue-800 text-white font-bold px-8 py-3 rounded-2xl shadow transition"
        >
            Upgrade Plan
        </button>
        <div class="flex flex-col md:flex-row gap-8 w-full max-w-5xl">
            <!-- User Plans -->
            <div class="flex-1">
                <h2 class="text-xl font-semibold mb-4 text-blue-700">
                    User Plans
                </h2>
                <div
                    v-for="plan in userPlans"
                    :key="plan.name"
                    class="bg-white rounded-2xl shadow-md p-6 mb-6 border-t-4"
                    :class="
                        plan.name === 'Premium'
                            ? 'border-blue-700'
                            : plan.name === 'Plus'
                            ? 'border-blue-500'
                            : 'border-gray-300'
                    "
                >
                    <h3 class="text-lg font-bold mb-2">{{ plan.name }}</h3>
                    <ul class="text-gray-700 space-y-1">
                        <li><b>Storage:</b> {{ plan.storage }}</li>
                        <li><b>Max File Size:</b> {{ plan.maxFileSize }}</li>
                        <li><b>Max Files:</b> {{ plan.maxFiles }}</li>
                        <li><b>Max Shared Links:</b> {{ plan.maxLinks }}</li>
                        <li><b>File Expiration:</b> {{ plan.expiration }}</li>
                        <li><b>Features:</b> {{ plan.features }}</li>
                    </ul>
                </div>
            </div>

            <!-- Team Plans -->
            <div class="flex-1">
                <h2 class="text-xl font-semibold mb-4 text-green-700">
                    Team Plans
                </h2>
                <div
                    v-for="plan in teamPlans"
                    :key="plan.name"
                    class="bg-white rounded-2xl shadow-md p-6 mb-6 border-t-4"
                    :class="
                        plan.name === 'Premium'
                            ? 'border-green-700'
                            : 'border-gray-300'
                    "
                >
                    <h3 class="text-lg font-bold mb-2">{{ plan.name }}</h3>
                    <ul class="text-gray-700 space-y-1">
                        <li><b>Team Storage:</b> {{ plan.storage }}</li>
                        <li><b>Max File Size:</b> {{ plan.maxFileSize }}</li>
                        <li><b>Max Members:</b> {{ plan.maxMembers }}</li>
                        <li><b>Shared Folders:</b> {{ plan.sharedFolders }}</li>
                        <li><b>File Expiration:</b> {{ plan.expiration }}</li>
                        <li><b>Features:</b> {{ plan.features }}</li>
                    </ul>
                </div>
            </div>
        </div>

        <!-- Upgrade Plan Modal -->
        <div
            v-if="showUpgradeModal"
            class="fixed inset-0 z-50 bg-black bg-opacity-30 flex items-center justify-center"
        >
            <div
                class="bg-white rounded-2xl shadow-xl p-8 min-w-[360px] flex flex-col gap-4 relative"
            >
                <button
                    class="absolute top-2 right-3 text-2xl text-gray-400 hover:text-gray-600"
                    @click="closeUpgradeModal"
                >
                    &times;
                </button>
                <h2 class="text-xl font-bold mb-2 text-blue-700">
                    Upgrade Plan
                </h2>

                <!-- Step 1: Choose type -->
                <label class="font-semibold mb-2">Choose Account Type:</label>
                <div class="flex gap-3 mb-3">
                    <button
                        :class="
                            upgradeType === 'user'
                                ? 'bg-blue-600 text-white'
                                : 'bg-gray-200 text-gray-800'
                        "
                        class="rounded-xl px-4 py-2 font-semibold transition"
                        @click="selectUpgradeType('user')"
                    >
                        User
                    </button>
                    <button
                        :class="
                            upgradeType === 'team'
                                ? 'bg-green-700 text-white'
                                : 'bg-gray-200 text-gray-800'
                        "
                        class="rounded-xl px-4 py-2 font-semibold transition"
                        @click="selectUpgradeType('team')"
                    >
                        Team
                    </button>
                </div>

                <!-- Step 2: Choose plan -->
                <div v-if="upgradeType" class="mb-3">
                    <label class="font-semibold mb-2">
                        Choose
                        {{ upgradeType === "user" ? "User" : "Team" }} Plan:
                    </label>
                    <div class="flex flex-col gap-2 mt-1">
                        <button
                            v-for="plan in upgradeType === 'user'
                                ? userPlans
                                : teamPlans"
                            :key="plan.name"
                            :class="
                                selectedPlan === plan.name
                                    ? upgradeType === 'user'
                                        ? 'bg-blue-600 text-white'
                                        : 'bg-green-700 text-white'
                                    : 'bg-gray-100 text-gray-800 border border-gray-200'
                            "
                            class="rounded-xl px-4 py-2 font-semibold text-left transition"
                            @click="selectedPlan = plan.name"
                        >
                            {{ plan.name }}
                            <span class="text-xs text-gray-300 ml-2"
                                >({{
                                    upgradeType === "user"
                                        ? plan.storage +
                                          ", " +
                                          plan.maxFileSize +
                                          " file"
                                        : plan.storage +
                                          ", up to " +
                                          plan.maxMembers +
                                          " members"
                                }})</span
                            >
                        </button>
                    </div>
                </div>

                <!-- Step 3: Choose Team if team upgrade -->
                <div v-if="upgradeType === 'team' && selectedPlan" class="mb-3">
                    <label class="font-semibold mb-2">Select Your Team:</label>
                    <div v-if="teamLoading" class="text-gray-500">
                        Loading teams...
                    </div>
                    <div v-else-if="teamList.length === 0" class="text-red-400">
                        No teams found
                    </div>
                    <div v-else class="flex flex-col gap-1">
                        <label
                            v-for="team in teamList"
                            :key="team.id"
                            class="flex items-center gap-2 cursor-pointer"
                        >
                            <input
                                type="radio"
                                name="selectedTeam"
                                :value="team.id"
                                v-model="selectedTeam"
                                class="accent-blue-600"
                            />
                            <span>{{ team.name }}</span>
                        </label>
                    </div>
                </div>

                <!-- Apply Button (always last step) -->
                <button
                    class="w-full bg-blue-700 hover:bg-blue-800 text-white font-bold px-6 py-2 rounded-xl mt-2 disabled:opacity-60"
                    :disabled="!canApply || isUpgrading"
                    @click="upgradePlan"
                >
                    <span v-if="!isUpgrading">Apply</span>
                    <span v-else>Applying...</span>
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from "vue";
import axiosInstance from "../axiosInstance";
import { showSuccess, showError } from "../utils/toast";
import { useUserStore } from "../stores/user";

const userStore = useUserStore();

const userPlans = [
    {
        name: "free",
        storage: "2 GB",
        maxFileSize: "100 MB",
        maxFiles: "1,000",
        maxLinks: "10",
        expiration: "7 days",
        features: "Basic file upload/download, basic sharing",
    },
    {
        name: "plus",
        storage: "100 GB",
        maxFileSize: "2 GB",
        maxFiles: "Unlimited",
        maxLinks: "100",
        expiration: "Max 30 days (custom)",
        features: "Password-protected sharing, advanced sharing",
    },
    {
        name: "premium",
        storage: "1 TB",
        maxFileSize: "10 GB",
        maxFiles: "Unlimited",
        maxLinks: "Unlimited",
        expiration: "Max 180 days (custom)",
        features: "All features, priority support, advanced sharing, no ads",
    },
];

const teamPlans = [
    {
        name: "free",
        storage: "10 GB",
        maxFileSize: "100 MB",
        maxMembers: "5",
        sharedFolders: "Limited",
        expiration: "14 days",
        features: "Basic collaboration, file sharing, team chat",
    },
    {
        name: "premium",
        storage: "1 TB",
        maxFileSize: "10 GB",
        maxMembers: "Unlimited",
        sharedFolders: "Unlimited",
        expiration: "120 days (custom)",
        features:
            "All collaboration features, admin controls, priority support, file approval workflows",
    },
];

// Modal state
const showUpgradeModal = ref(false);
const upgradeType = ref(""); // "user" or "team"
const selectedPlan = ref("");
const isUpgrading = ref(false);

// For team selection
const teamList = ref([]);
const teamLoading = ref(false);
const selectedTeam = ref("");

// When to enable Apply button
const canApply = computed(() =>
    upgradeType.value === "user"
        ? !!selectedPlan.value
        : !!selectedPlan.value && !!selectedTeam.value
);

// Modal logic
function openUpgradeModal() {
    upgradeType.value = "";
    selectedPlan.value = "";
    showUpgradeModal.value = true;
    teamList.value = [];
    selectedTeam.value = "";
}

function closeUpgradeModal() {
    showUpgradeModal.value = false;
    upgradeType.value = "";
    selectedPlan.value = "";
    isUpgrading.value = false;
    teamList.value = [];
    selectedTeam.value = "";
}

// Handle type selection and fetch team list if needed
async function selectUpgradeType(type) {
    upgradeType.value = type;
    selectedPlan.value = "";
    if (type === "team") {
        teamLoading.value = true;
        selectedTeam.value = "";
        try {
            const resp = await axiosInstance.get("/api/team/get");
            teamList.value = resp.data || [];
        } catch (e) {
            teamList.value = [];
            showError("Failed to load teams");
        }
        teamLoading.value = false;
    }
}

// Apply/Upgrade action
async function upgradePlan() {
    if (!canApply.value) {
        showError("Please fill all fields.");
        return;
    }
    isUpgrading.value = true;
    try {
        if (upgradeType.value === "user") {
            await axiosInstance.put("/api/user/plan/change", {
                plan: selectedPlan.value,
            });

            userStore.plan = selectedPlan.value;
        } else if (upgradeType.value === "team") {
            await axiosInstance.put(
                `/api/team/plan/update/${selectedTeam.value}`,
                {
                    plan: selectedPlan.value,
                }
            );
        }
        showSuccess("Plan upgrade requested!");
        closeUpgradeModal();
    } catch (err) {
        showError(
            err?.response?.data?.error ||
                "Failed to upgrade plan. Please try again."
        );
        isUpgrading.value = false;
    }
}
</script>
