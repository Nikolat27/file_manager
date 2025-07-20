// src/stores/user.js
import { defineStore } from "pinia";

export const useUserStore = defineStore("user", {
    state: () => ({
        id: null,
        username: "",
        plan: "",
        token: "",
    }),
    actions: {
        setUser({ id, username, plan, token }) {
            this.id = id;
            this.username = username;
            this.plan = plan;
            this.token = token;
        },
        clearUser() {
            this.id = null;
            this.username = "";
            this.plan = "";
            this.token = "";
        },
    },
    persist: true,
});
