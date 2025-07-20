import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/Home.vue";
import LoginView from "../views/Login.vue";
import RegisterView from "../views/Register.vue";
import CreateFileSettingsView from "../views/CreateFileSettings.vue";
import FolderContents from "../views/FolderContents.vue";

const routes = [
    { path: "/", redirect: "/home" },
    { path: "/home", name: "home", component: HomeView },
    { path: "/login", name: "login", component: LoginView },
    { path: "/register", name: "register", component: RegisterView },
    {
        path: "/file/setting/create/:id",
        name: "FileSettingCreate",
        component: CreateFileSettingsView,
    },
    {
        path: "/folder/get/:id",
        name: "FolderContents",
        component: FolderContents,
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;
