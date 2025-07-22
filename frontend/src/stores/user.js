// src/stores/user.js
import { defineStore } from "pinia";

export const useUserStore = defineStore("user", {
    state: () => ({
        id: null,
        username: "",
        plan: "",
        token: "",
        avatarUrl: "",
    }),
    actions: {
        setUser({ id, username, plan, token, avatarUrl }) {
            this.id = id;
            this.username = username;
            this.plan = plan;
            this.token = token;
            this.avatarUrl = avatarUrl;
        },
        clearUser() {
            this.id = null;
            this.username = "";
            this.plan = "";
            this.token = "";
            this.avatarUrl = "";
        },
    },
    persist: true,
});
