import { beforeEach, expect, test, describe } from 'vitest';
import { transformDefinition, useCardManager } from '../src/stores/CardManager';
import { setActivePinia, createPinia } from 'pinia';
import {
  backend,
} from '@wailsjs/models';

describe('CardManager unit tests', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  test('definition transformation', () => {
    const definitionA = '[pinyin] The text of pinyin';
    expect(transformDefinition(definitionA)).toStrictEqual([{
      definition: 'The text of pinyin',
      pronunciation: 'pinyin',
    }]);

    const anOldDefinition =
      '【答应】【答應】 dā ying ★★★★★<br>to answer/to respond/' +
      'to answer positively/to agree/to accept/to promise';
    expect(transformDefinition(anOldDefinition)).toStrictEqual([{
      definition: anOldDefinition,
      pronunciation: undefined,
    }]);

    // TODO transform old pinyin
    const jiaoLong =
      '[jiao1long2] a legendary dragon with the ability to control ' +
      'rain and floods/see also 蛟龍|蛟龙[jiao1 long2]';
    expect(transformDefinition(jiaoLong)).toStrictEqual([{
      definition: 'a legendary dragon with the ability to control ' +
      'rain and floods/see also 蛟龍|蛟龙[jiao1 long2]',
      pronunciation: 'jiao1long2',
    }]);

    const lu = '[lǚ] surname Lü<br>[lǚ] pitchpipe, pitch standard' +
      ', one of the twelve semitones in the traditional tone system';
    expect(transformDefinition(lu)).toStrictEqual([{
      definition: 'surname Lü',
      pronunciation: 'lǚ',
    },
    {
      definition: 'pitchpipe, pitch standard' +
      ', one of the twelve semitones in the traditional tone system',
      pronunciation: 'lǚ',
    }]);

    // TODO clean these types up?
    const fullOfHtml = '<br>[lǐzhì] reason/rational<br><br><br><div></div>';
    expect(transformDefinition(fullOfHtml)).toStrictEqual([{
      definition: fullOfHtml,
      pronunciation: undefined,
    }]);
    const multiSylable = '[wèicǐ] for this reason/with regards to this';
    expect(transformDefinition(multiSylable)).toStrictEqual([{
      definition: 'for this reason/with regards to this',
      pronunciation: 'wèicǐ',
    }]);
  });

  test('Update Sentence', () => {
    const cardManager = useCardManager();
    cardManager.loadCard({
      word: '你好',
      sourceCard: backend.RawAnkiNote.createFrom({
        noteId: 34,
        fields: backend.Fields.createFrom({
          word: '你好',
          sentence: '你好我是大卫',
          englishDefn: '[ni3hao3] Hello',
          chineseDefn: '',
          pinyin: 'ni3hao3',
          images: [],
        }),
      }),
    });
    expect(cardManager.sentence).toBe('你好我是大卫');
    expect(cardManager.sentenceSource).toBe(undefined);
  });
});
