<template>
  <div class="">
    <div
      class="
        mx-auto
        flex
        h-16
        items-center
        bg-accent
        px-8"
    >
      <h2 class="mt-5 text-center text-3xl">
        Welcome to read-chinese
      </h2>
      <img
        class="h-32 w-auto"
        src="../assets/logo_transparent.png"
        alt="Read More"
      >
    </div>
    <div class="m-6 flex">
      <div class="w-2/3">
        <p class="mx-auto w-3/4">
          Before getting started make sure the all the entries to the right are
          configured properly. After completing those, checkout the other items
          below for things you can do next to fully unlock the features of this
          application
        </p>
        <p>
          Designed to help manage a dynamic flashcard
          library to aid in reading books, lets get started by making sure
          some beggining steps have been completed
        </p>
      </div>
      <div class="w-1/3">
        <recommended-actions />
        <health-check />
      </div>
    </div>
    <chat-health />
    <chat-list />
  </div>
</template>

<script lang="ts" setup>
import RecommendedActions from '@/components/RecommendedActions.vue';
import HealthCheck from '@/components/HealthCheck.vue';
import ChatHealth from '@/components/ChatHealth.vue';
import ChatList from '@/components/ChatList.vue';

</script>
