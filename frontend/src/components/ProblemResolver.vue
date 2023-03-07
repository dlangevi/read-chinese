<template>
  <div class="flex gap-2">
    <button
      v-for="b, i in buttons"
      :key="b.text"
      class="btn-secondary btn-xs btn"
      @click="processAction(b, i)"
    >
      {{ b.text }}
    </button>
  </div>
</template>

<script lang="ts" setup>
import { onBeforeMount, ref, watch } from 'vue';
import { useCardQueue } from '@/stores/CardQueue';
import type { ICellRendererParams } from 'ag-grid-community';
import { backend } from '@wailsjs/models';
import {
  UpdateSentenceAudio,
  UpdateWordAudio,
  UpdatePinyin,
} from '@wailsjs/backend/ankiInterface';
import { StepsEnum } from '@/components/CardCreatorSteps/StepsEnum';

const props = defineProps<{
  params: ICellRendererParams<backend.ProblemCard>,
}>();

type ButtonInfo = {
  text: string,
  action:() => Promise<Error | void>,
}

const buttons = ref<ButtonInfo[]>([]);
let finishedLoad = false;

watch(buttons, () => {
  if (!finishedLoad) {
    return;
  }
  if (buttons.value.length === 0) {
    const rowData = props.params.data;
    if (rowData) {
      props.params.api.applyTransaction({
        remove: [rowData],
      });
    }
  }
}, {
  deep: true,
});

async function processAction(button: ButtonInfo, index : number) {
  await button.action();
  buttons.value.splice(index, 1);
}

onBeforeMount(() => {
  const problemCard = props.params.data as backend.ProblemCard;
  const problems : backend.Problems = problemCard.problems;

  // We always want to give the option to
  // open the card editor. This covers both
  // problems.MissingImage and problems.MissingSentence
  buttons.value.push({
    text: 'Open Editor',
    action: addToQueue,
  });

  if (problems.MissingPinyin) {
    // call the backend to add pinyin to the card
    buttons.value.push({
      text: 'Generate Pinyin',
      action: () =>
        UpdatePinyin(problemCard.noteId),
    });
  }
  if (problems.MissingSentenceAudio) {
    // call the backend to add audio to the card
    buttons.value.push({
      text: 'Generate Sentence Audio',
      action: () =>
        UpdateSentenceAudio(problemCard.noteId),
    });
  }
  if (problems.MissingWordAudio) {
    // call the backend to add audio to the card
    buttons.value.push({
      text: 'Generate Word Audio',
      action: () =>
        UpdateWordAudio(problemCard.noteId),
    });
  }
  finishedLoad = true;
});

const store = useCardQueue();
async function addToQueue() : Promise<void> {
  let editComplete = false;

  const problemCard = props.params.data as backend.ProblemCard;
  const problems : backend.Problems = problemCard.problems;

  const keySteps : StepsEnum[] = [];
  if (problems.MissingImage) {
    keySteps.push(StepsEnum.IMAGE);
  }
  if (problems.MissingSentence) {
    keySteps.push(StepsEnum.SENTENCE);
  }

  store.addWord({
    word: problemCard.word,
    sourceCardId: problemCard.noteId,
    keySteps,
  }, () => {
    editComplete = true;
  });

  // TODO this is hacky, but the chain of callbacks gets lost currently
  return new Promise((resolve, _reject) => {
    const checkVariable = () => {
      if (editComplete) {
        resolve(undefined);
      } else {
        setTimeout(checkVariable, 100);
      }
    };
    checkVariable();
  });
}

</script>
