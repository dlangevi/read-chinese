<template>
    <div class="bg-gray-800 mx-auto px-8">
      <div class="block flex items-center justify-between h-16">
        <img class="block h-32 w-auto"
             src="../assets/logo_transparent.png" alt="Read More" />
        <div class="flex-1 flex items-center justify-center sm:items-stretch sm:justify-start">
          <div class="hidden sm:block sm:ml-6">
            <div class="flex space-x-4">
              <router-link v-for="item in navigation" :key="item.name" :to="item.href"
                :class="[item.current ? 'bg-gray-900 text-white' :
                                        'text-gray-300 hover:bg-gray-700 hover:text-white',
                                        'px-3 py-2 rounded-md text-sm font-medium']"
                :aria-current="item.current ? 'page' : undefined">{{ item.name }}
                </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
</template>

<script setup>
// TODO get the image working better

// TODO have the current somehow match the current page
const navigation = [
  { name: 'BookLibrary', href: '/BookLibrary', current: true },
  { name: 'Manage FlashCards', href: '/FlashCards', current: true },
  { name: 'Manage Wordlist', href: '/WordLists', current: true },
  { name: 'About', href: '/About', current: true },
];
</script>
