<template>
    <div class="w-screen h-screen flex items-center justify-center bg-blue-50">
        <form
            class="w-full max-w-md bg-white rounded-2xl shadow-2xl p-8 flex flex-col gap-6 items-center"
            @submit.prevent="onRegister"
        >
            <h2 class="text-2xl font-bold text-blue-700 mb-2">Register</h2>
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
            <input
                v-model="repeatPassword"
                type="password"
                placeholder="Repeat Password"
                class="w-full px-4 py-3 rounded-xl border border-blue-300 focus:outline-none focus:border-blue-500 bg-blue-50"
                required
            />
            <button
                type="submit"
                class="w-full py-3 bg-blue-600 hover:bg-blue-700 transition-colors text-white rounded-xl font-semibold text-lg shadow"
            >
                Register
            </button>
            <span
                @click="goToLogin"
                class="cursor-pointer text-blue-500 underline"
                >login</span
            >
        </form>
    </div>
</template>

<script setup>
import { ref } from "vue";
import axiosInstance from "../axiosInstance";
import { showError, showSuccess } from "../utils/toast";
import { useRouter } from "vue-router";

const username = ref("");
const password = ref("");
const repeatPassword = ref("");

const router = useRouter();

async function register() {
    const payload = {
        username: username.value,
        password: password.value,
    };

    const headers = {
        "Content-Type": "application/json",
    };

    axiosInstance
        .post("/api/auth/register", payload, {
            headers: headers,
        })
        .then(() => {
            showSuccess("User registered successfully");
            router.push({ name: "login" });
        })
        .catch((err) => {
            showError(err.response.data.error);
        });
}

function goToLogin() {
    router.push({ name: "login" });
}

function onRegister() {
    if (password.value !== repeatPassword.value) {
        alert("Passwords do not match!");
        return;
    }

    register();
}
</script>
