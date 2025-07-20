import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/Home.vue";
import LoginView from "../views/Login.vue";
import RegisterView from "../views/Register.vue";

const routes = [
    { path: "/", redirect: "/home" },
    { path: "/home", name: "home", component: HomeView },
    { path: "/login", name: "login", component: LoginView },
    { path: "/register", name: "register", component: RegisterView },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;
