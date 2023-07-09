<template>
  <div class="flex h-screen flex-col">
    <div class="grow bg-gradient-to-b from-accent to-accent-focus">
      <div class="mx-auto max-w-7xl py-24 px-4 sm:px-6 lg:px-8">
        <h1
          :class="[
            'text-4xl',
            'font-extrabold',
            'text-accent-content',
            'sm:text-5xl md:text-6xl'
          ]"
        >
          Welcome to read-chinese
        </h1>
        <div class="flex">
          <div class="mt-6 w-2/3 text-accent-content">
            <p class="text-xl">
              Designed to help manage a dynamic flashcard library to aid in
              reading books, lets get started by making sure some beggining
              steps have been completed
            </p>
            <p class="mt-4 text-lg">
              Before getting started make sure the all the entries to the right
              are configured properly. After completing those, checkout the
              other items below for things you can do next to fully unlock the
              features of this application
            </p>
            <p class="mt-4 text-lg">
              To ask questions, suggest improvements, or get support in setup
              please join the community
              <a
                href="#"
                class="text-info hover:text-info-content"
                @click="
                  BrowserOpenURL('https://discord.gg/RVgCQjTVYq')"
              >
                discord
              </a>
            </p>
            <p class="mt-4 text-lg">
              For more information checkout:
            </p>
            <ul class="mt-4 list-inside list-disc text-lg">
              <li>
                <a
                  href="#"
                  class="text-info hover:text-info-content"
                  @click="
                    BrowserOpenURL('https://github.com/dlangevi/read-chinese')"
                >
                  The github project
                </a>
              </li>
              <li class="">
                <a
                  href="#"
                  class="text-info hover:text-info-content"
                  @click="BrowserOpenURL('https://youtu.be/dAqE2YquGlw')"
                >
                  A quick demo of the project
                </a>
              </li>
            </ul>
          </div>
          <health-check class="w-1/3" />
        </div>
      </div>
    </div>
    <div class="bg-gray-100">
      <div
        class="mx-auto flex max-w-7xl
      justify-between
      py-6 px-4 sm:px-6 lg:px-8"
      >
        <div>
          <h2 class="text-lg font-medium text-gray-900">Application Info</h2>
          <p class="mt-1 text-sm text-gray-600">
            Version 0.0.3
          </p>
        </div>
        <div>
          <h2 class="text-lg font-medium text-gray-900">Diagnostics</h2>
          <p class="mt-1 text-sm text-gray-600">
            Total Books: Lots
          </p>
          <p class="mt-1 text-sm text-gray-600">
            Feeling Good: I think so
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { BrowserOpenURL } from '@runtime/runtime';
import HealthCheck from '@/components/HealthCheck.vue';
</script>
