<template>
    <aside
        class="w-[71px] h-full bg-[#f7f5f2] border-r border-r-gray-300 absolute left-0"
    >
        <div class="flex flex-col gap-y-1 items-center pt-6">
            <button
                class="w-10 h-10 rounded-xl bg-[#ebe9e6] flex items-center justify-center hover:cursor-pointer hover:bg-[#d5d3d1]"
            >
                <svg
                    viewBox="0 0 24 24"
                    fill="none"
                    class="dig-UIIcon dig-UIIcon--standard _railButtonIcon_151kt_46"
                    width="24"
                    height="24"
                    focusable="false"
                    aria-label="Home"
                >
                    <path
                        d="M20.594 9.415 12 2.54 3.406 9.415l.938 1.17L5.5 9.662V19H10v-4.75a.75.75 0 0 1 .75-.75h2.5a.75.75 0 0 1 .75.75V19h4.5V9.66l1.157.925.937-1.171Z"
                        fill="currentColor"
                        vector-effect="non-scaling-stroke"
                    ></path>
                </svg>
            </button>
            <span class="text-[13px]">Home</span>
        </div>
        <div class="flex flex-col gap-y-1 items-center pt-6">
            <button
                @click="handleLogout"
                class="w-10 h-10 rounded-xl bg-[#ebe9e6] flex items-center justify-center hover:cursor-pointer hover:bg-[#d5d3d1]"
            >
                <img
                    width="25px"
                    height="25px"
                    src="../assets/icons8-logout-64.png"
                    alt=""
                />
            </button>
            <span class="text-[13px]">Logout</span>
            <span class="text-[13px] underline">{{ userStore.username }}</span>
        </div>
    </aside>
    <aside
        class="w-[228px] h-full bg-[#f7f5f2] border-r border-r-gray-300 absolute left-18"
    >
        <div class="flex flex-col gap-y-1 pt-6 px-6">
            <h5 class="text-[20px] font-bold mb-4">Home</h5>

            <div
                @click="goToAllFiles"
                :class="[
                    'w-[196px] h-[36px] rounded-xl hover:cursor-pointer flex justify-start items-center pl-4',
                    route.name === 'home'
                        ? 'bg-black text-white'
                        : 'bg-[#f7f5f2] text-black hover:bg-[#ebe9e6]',
                ]"
            >
                <span class="text-[14px] font-medium">All files</span>
            </div>
            <div
                @click="goToSharedUrls"
                :class="[
                    'w-[196px] h-[36px] rounded-xl hover:cursor-pointer flex justify-start items-center pl-4',
                    route.name === 'sharedUrls'
                        ? 'bg-black text-white'
                        : 'bg-[#f7f5f2] text-black hover:bg-[#ebe9e6]',
                ]"
            >
                <span class="text-[14px] font-medium">Shared Urls</span>
            </div>
            <div
                @click="goToSendApprovals"
                :class="[
                    'w-[196px] h-[36px] rounded-xl hover:cursor-pointer flex justify-start items-center pl-4',
                    route.name === 'SentApprovals'
                        ? 'bg-black text-white'
                        : 'bg-[#f7f5f2] text-black hover:bg-[#ebe9e6]',
                ]"
            >
                <span class="text-[14px] font-medium">Sent Approvals</span>
            </div>
            <div
                @click="goToReceivedApprovals"
                :class="[
                    'w-[196px] h-[36px] rounded-xl hover:cursor-pointer flex justify-start items-center pl-4',
                    route.name === 'ReceivedApprovals'
                        ? 'bg-black text-white'
                        : 'bg-[#f7f5f2] text-black hover:bg-[#ebe9e6]',
                ]"
            >
                <span class="text-[14px] font-medium">Received Approvals</span>
            </div>
            <div
                @click="goToTeams"
                :class="[
                    'w-[196px] h-[36px] rounded-xl hover:cursor-pointer flex justify-start items-center pl-4',
                    route.name === 'teams'
                        ? 'bg-black text-white'
                        : 'bg-[#f7f5f2] text-black hover:bg-[#ebe9e6]',
                ]"
            >
                <span class="text-[14px] font-medium">Teams</span>
            </div>
            <div
                @click="goToPlans"
                :class="[
                    'w-[196px] h-[36px] rounded-xl hover:cursor-pointer flex justify-start items-center pl-4',
                    route.name === 'plans'
                        ? 'bg-black text-white'
                        : 'bg-[#f7f5f2] text-black hover:bg-[#ebe9e6]',
                ]"
            >
                <span class="text-[14px] font-medium">Plans</span>
            </div>
        </div>
    </aside>
</template>
<script setup>
import { useUserStore } from "../stores/user";
import { useRoute, useRouter } from "vue-router";

const userStore = useUserStore();

const route = useRoute();
const router = useRouter();

function handleLogout() {
    userStore.token = "";
    router.push("/login");
}

function goToAllFiles() {
    router.push({ name: "home" });
}

function goToSharedUrls() {
    router.push({ name: "sharedUrls" });
}

function goToSendApprovals() {
    router.push({ name: "SentApprovals" });
}

function goToReceivedApprovals() {
    router.push({ name: "ReceivedApprovals" });
}

function goToTeams() {
    router.push({ name: "teams" });
}

function goToPlans() {
    router.push({ name: "plans" });
}
</script>
