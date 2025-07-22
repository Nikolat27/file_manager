<template>
    <div class="min-h-screen bg-gray-50 py-12 pt-20 px-4 flex flex-col items-center">
        <h1 class="text-3xl font-bold mb-10">Plans & Features</h1>
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
    </div>
</template>

<script setup>
const userPlans = [
    {
        name: "Free",
        storage: "2 GB",
        maxFileSize: "100 MB",
        maxFiles: "1,000",
        maxLinks: "10",
        expiration: "30 days",
        features: "Basic file upload/download, basic sharing",
    },
    {
        name: "Plus",
        storage: "100 GB",
        maxFileSize: "2 GB",
        maxFiles: "Unlimited",
        maxLinks: "100",
        expiration: "No expiration",
        features:
            "Password-protected sharing, extended history, advanced sharing",
    },
    {
        name: "Premium",
        storage: "1 TB",
        maxFileSize: "10 GB",
        maxFiles: "Unlimited",
        maxLinks: "Unlimited",
        expiration: "No expiration",
        features:
            "All features, priority support, versioning, advanced sharing, no ads",
    },
];

const teamPlans = [
    {
        name: "Free",
        storage: "10 GB",
        maxFileSize: "100 MB",
        maxMembers: "5",
        sharedFolders: "Limited",
        expiration: "30 days",
        features: "Basic collaboration, file sharing, team chat",
    },
    {
        name: "Premium",
        storage: "10 TB",
        maxFileSize: "10 GB",
        maxMembers: "Unlimited",
        sharedFolders: "Unlimited",
        expiration: "No expiration",
        features:
            "All collaboration features, admin controls, priority support, file approval workflows",
    },
];
</script>
