<template>
    <div class="flex flex-col gap-y-6">
        <span class="text-[22px] font-bold">All files</span>

        <table
            class="w-[81%] border border-blue-200 rounded-xl overflow-hidden"
        >
            <thead>
                <tr class="bg-gray-100 text-gray-700">
                    <th class="text-left px-4 py-2 w-[70%]">Name</th>
                    <th class="text-left px-4 py-2 w-[15%]">Updated At</th>
                    <th class="text-left px-4 py-2 w-[15%]">Modified At</th>
                    <th class="text-left px-4 py-2 w-[15%]">More</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td class="px-4 py-2 flex flex-row items-center gap-x-2 z">
                        <span>
                            <svg
                                viewBox="0 0 40 40"
                                fill="none"
                                role="img"
                                focusable="false"
                                width="32"
                                height="32"
                                class="dig-ContentIcon brws-file-name-cell-icon dig-ContentIcon--small dig-ctz1wx2_5-7-0 dig-ctz1wx4_5-7-0"
                                data-testid="FolderBaseDefaultSmall"
                                data-campaigns-element-id="file-thumbnail"
                                data-thumbnail-testid="fallback-file-thumbnail"
                            >
                                <path
                                    d="M15.002 7.004c.552.018.993.252 1.295.7l.785 2.12c.145.298.363.576.561.779.252.257.633.4 1.156.4H35.5l-.002 18c-.027.976-.3 1.594-.836 2.142-.565.577-1.383.858-2.41.858H8.748c-1.026 0-1.844-.28-2.409-.858-.564-.577-.838-1.415-.838-2.465V7.003h9.502Z"
                                    fill="var(--dig-color__foldericon__shadow, #75aaff)"
                                ></path>
                                <path
                                    d="M15.002 7.001c.552.018.993.252 1.295.7l.785 2.12c.145.298.363.576.561.779.252.257.633.4 1.156.4H35.5l-.002 16.84c-.027.976-.3 1.754-.836 2.302-.565.577-1.383.858-2.41.858H8.748c-1.026 0-1.844-.28-2.409-.858-.564-.577-.838-1.415-.838-2.465V7l9.502.001Z"
                                    fill="var(--dig-color__foldericon__container, #a0c4ff)"
                                ></path>
                            </svg>
                        </span>
                        example.txt
                    </td>
                    <td class="px-4 py-2">2025-07-19</td>
                    <td class="px-4 py-2">2025-07-19</td>
                    <td
                        class="px-4 py-2 text-2xl font-bold cursor-pointer select-none pl-8 pb-4"
                    >
                        ...
                    </td>
                    <div class="absolute ">

                    </div>
                </tr>
            </tbody>
        </table>
    </div>
</template>
