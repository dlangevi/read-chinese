<template>
  <div class="container mx-auto px-4">
    <div class="text-center">
      <h2 class="text-center mt-5">Settings</h2>
      <p>Lets go</p>
      <div
        v-for="(contents, section, index) in sections" :key="index"
        class="grid gap-4 grid-cols-4 m-4 p-4 border-2">
        <div  class="border-4">
          {{section}}
        </div>
        <div
          class="col-span-3 grid gap-4 grid-cols-3">
            <component
              v-for="content in contents" :key="content"
              :is="content.type"
              v-bind="content.props"
              />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import UserSettings from '@/userSettings';

const sections = {
  CardCreationSettings: [
    UserSettings.PopulateEnglish.Render,
    UserSettings.AutoAdvanceSentence.Render,
    UserSettings.AutoAdvanceEnglish.Render,
    UserSettings.AutoAdvanceCard.Render,
    UserSettings.EnableChinese.Render,
    UserSettings.PopulateChinese.Render,
    UserSettings.GenerateTermAudio.Render,
    UserSettings.GenerateSentenceAudio.Render,
  ],
};
</script>
