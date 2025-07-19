<template>
    <aside class="w-full h-[72px] flex flex-row items-center absolute top-0 left-[300px] px-6">
        <input type="text" class="w-[922px] h-[40px] hover:outline-2  outline-gray-400 outline-1
         hover:outline-black rounded-2xl pl-4" placeholder="Search for file names">

         <div class="h-[40px] absolute right-[300px] flex flex-row items-center gap-x-4 mr-6">
            <button class="rounded-full w-[32px] h-[32px] cursor-pointer">
               <svg viewBox="0 0 24 24" fill="none" class="dig-UIIcon" width="32" height="32" role="presentation" focusable="false" color="inherit"><path d="m17.608 12.971-.329-.219a1.747 1.747 0 0 1-.779-1.457v-1.67c0-1.094 0-2.332-.563-3.336C15.261 5.084 13.972 4.5 12 4.5c-1.972 0-3.26.585-3.937 1.787C7.5 7.292 7.5 8.531 7.5 9.624v1.672a1.747 1.747 0 0 1-.78 1.454l-.327.219A4.241 4.241 0 0 0 4.5 16.507v.993H10a1.857 1.857 0 0 0 2 2 1.857 1.857 0 0 0 2-2h5.5v-.993a4.242 4.242 0 0 0-1.892-3.536ZM6.047 16a2.743 2.743 0 0 1 1.178-1.781L7.553 14A3.244 3.244 0 0 0 9 11.296V9.622c0-.953 0-1.938.37-2.6C9.618 6.584 10.16 6 12 6c1.841 0 2.383.584 2.63 1.023.371.662.371 1.646.37 2.6v1.674A3.244 3.244 0 0 0 16.447 14l.329.219A2.744 2.744 0 0 1 17.953 16H6.047Z" fill="currentColor" vector-effect="non-scaling-stroke"></path></svg>
            </button>
            <button class="rounded-full w-[32px] h-[32px] cursor-pointer">
                <img class="w-full h-full" src="../assets/images.png" alt="">
            </button>
        </div>
    </aside>
</template>