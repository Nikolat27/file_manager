<template>
    <div class="flex flex-col gap-y-6">
        <span class="text-[22px] font-bold">All files</span>

        <table class="w-[81%] border rounded-xl overflow-hidden">
            <thead>
            <tr class="bg-gray-100 text-gray-700">
                <th class="text-left px-4 py-2 w-[70%]">Name</th>
                <th class="text-left px-4 py-2 w-[15%]">Updated At</th>
                <th class="text-left px-4 py-2 w-[15%]">Modified At</th>
            </tr>
            </thead>
            <tbody>
            <tr>
                <td class="px-4 py-2">example.txt</td>
                <td class="px-4 py-2">2025-07-19</td>
                <td class="px-4 py-2">2025-07-19</td>
            </tr>
            </tbody>
        </table>
    </div>
</template>