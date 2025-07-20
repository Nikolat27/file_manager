<template>
    <div class="w-screen h-screen flex items-center justify-center bg-blue-50">
        <form
            class="w-full max-w-md bg-white rounded-2xl shadow-2xl p-8 flex flex-col gap-6 items-center"
            @submit.prevent="onLogin"
        >
            <h2 class="text-2xl font-bold text-blue-700 mb-2">Login</h2>
            <input
                v-model="username"
                type="text"
                placeholder="Username"
                class="w-full px-4 py-3 rounded-xl border border-blue-300 focus:outline-none focus:border-blue-500 bg-blue-50"
                required
            />
            <input
                v-model="password"
                type="password"
                placeholder="Password"
                class="w-full px-4 py-3 rounded-xl border border-blue-300 focus:outline-none focus:border-blue-500 bg-blue-50"
                required
            />
            <button
                type="submit"
                class="w-full py-3 bg-blue-600 hover:bg-blue-700 transition-colors text-white rounded-xl font-semibold text-lg shadow"
            >
                Login
            </button>
        </form>
    </div>
</template>

<script setup>
import { ref } from "vue";
import { showSuccess, showError } from "../utils/toast";
import { useUserStore } from "../stores/user";
import axiosInstance from "../axiosInstance";
import { useRouter } from "vue-router";

const username = ref("");
const password = ref("");

const router = useRouter();

const userStore = useUserStore();

async function login() {
    const payload = {
        username: username.value,
        password: password.value,
    };

    const headers = {
        "Content-Type": "application/json",
    };

    axiosInstance
        .post("/api/auth/login", payload, {
            headers: headers,
        })
        .then((resp) => {
            userStore.setUser({
                id: resp.data.id,
                username: resp.data.username,
                plan: resp.data.plan,
                token: resp.data.token,
            });

            showSuccess("User logged in successfully");
            router.push({ name: "home" });
        })
        .catch((err) => {
            showError(err.response.data);
        });
}

function onLogin() {
    if (!username || !password) {
        showError("username or password is missing");
        return;
    }

    login();
}
</script>
