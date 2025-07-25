import { createRouter, createWebHistory } from "vue-router";
import HomeLayout from "../layout/HomeLayout.vue";
import HomeView from "../views/Home.vue";
import LoginView from "../views/Login.vue";
import RegisterView from "../views/Register.vue";
import CreateFileSettingsView from "../views/CreateFileSettings.vue";
import FolderContents from "../views/FolderContents.vue";
import ShortUrlsView from "../views/ShortUrls.vue";
import GetFileView from "../views/GetFile.vue";
import SentApprovals from "../views/SentApprovals.vue";
import ReceivedApprovals from "../views/ReceivedApprovals.vue";
import TeamsView from "../views/Teams.vue";
import TeamFilesView from "../views/TeamFiles.vue";
import PlansView from "../views/Plans.vue";
import UserSearchView from "../views/UserSearch.vue";

const routes = [
    { path: "/", redirect: "/home" },
    {
        path: "/home",
        component: HomeLayout,
        children: [
            {
                path: "",
                name: "home",
                component: HomeView,
            },
            {
                path: "shared-urls",
                name: "sharedUrls",
                component: ShortUrlsView,
            },
        ],
    },
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
    {
        path: "/file/get/:id",
        name: "GetFile",
        component: GetFileView,
    },
    {
        path: "/approvals/sent",
        name: "SentApprovals",
        component: SentApprovals,
    },
    {
        path: "/approvals/received",
        name: "ReceivedApprovals",
        component: ReceivedApprovals,
    },
    {
        path: "/teams",
        name: "teams",
        component: TeamsView,
    },
    {
        path: "/teams/:id",
        name: "TeamFiles",
        component: TeamFilesView,
    },
    {
        path: "/plans",
        name: "plans",
        component: PlansView,
    },
    {
        path: "/search",
        name: "UserSearch",
        component: UserSearchView,
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;
