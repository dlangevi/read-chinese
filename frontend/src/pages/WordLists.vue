<template>
  <with-sidebar>
    <template #sidebar>
      <button class="btn-primary btn" @click="importAnki">
        Sync from Anki
      </button>
      <n-cascader
        :options="options"
        placeholder="Load HSK Words"
        :show-path="true"
        expand-trigger="click"
        check-strategy="child"
        @update:value="loadHsk"
      />
      <button class="btn-primary btn" @click="makeCards">
        Make Cards
      </button>
      <div class="border-2 p-2 text-center">
        {{ words.length }} remaining words
      </div>
    </template>
    <div class="flex h-full w-full flex-col">
      <div class="shrink-0 text-center">
        <h1 class="mt-5 text-center text-xl">
          Maybe you know these words?
        </h1>
        <p>For now lets just mark learned words you already know 好不好?</p>
      </div>
      <unknown-words
        class="m-4 mx-auto h-full w-5/6 grow"
        show-definitions
        :words="words"
      />
    </div>
  </with-sidebar>
</template>

<script lang="ts" setup>
import WithSidebar from '@/layouts/WithSidebar.vue';
import { NCascader } from 'naive-ui';
import type { CascaderOption } from 'naive-ui';
import { ref, onBeforeMount } from 'vue';
import { backend } from '@wailsjs/models';
import { LearningTarget } from '@wailsjs/backend/bookLibrary';
import { GetUnknownHskWords } from '@wailsjs/backend/KnownWords';
import { ImportAnkiKeywords } from '@wailsjs/backend/AnkiInterface';
import UnknownWords from '../components/UnknownWords.vue';

import { useCardQueue, ActionsEnum } from '@/stores/CardQueue';

type HskVersion = '2.0' | '3.0';
type HskLevel = 1 | 2 | 3 | 4 | 5 | 6 | 7;

interface HskCascaderOption extends CascaderOption {
  level: HskLevel,
  version: HskVersion,
}
const options:CascaderOption[] = ['2.0', '3.0'].map((version) => ({
  label: `HSK ${version}`,
  value: version,
  children: (version === '2.0'
    ? [1, 2, 3, 4, 5, 6]
    : [1, 2, 3, 4, 5, 6, 7]).map((level) => ({
    value: `${version}-${level}`,
    version,
    level,
    label: `Level ${level}`,
  })),
}));

const words = ref<backend.UnknownWordEntry[]>([]);
onBeforeMount(async () => {
  words.value = await LearningTarget();
});

function importAnki() {
  ImportAnkiKeywords();
}

async function loadHsk(_: string, option: HskCascaderOption) {
  words.value = await GetUnknownHskWords(option.version, option.level);
}

const store = useCardQueue();
async function makeCards() {
  const sorted = words.value.slice();
  sorted.sort((a, b) => {
    if (a.occurance === undefined || b.occurance === undefined) {
      return 0;
    }
    if (a.occurance > b.occurance) {
      return -1;
    }
    return 1;
  });
  const topWords = sorted.slice(0, 50);
  topWords.forEach((word) => {
    store.addWord(word.word, ActionsEnum.CREATE);
  });
}

</script>
