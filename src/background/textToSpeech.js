// pull in the required packages.
import {
  SpeechConfig, AudioConfig, SpeechSynthesizer, ResultReason,
} from 'microsoft-cognitiveservices-speech-sdk';
import tmp from 'tmp';
import { getOptionValue } from './database';

// TODO optionally run this code
const subscriptionKey = getOptionValue('AzureApiKey');
const serviceRegion = 'eastus';

const speechConfig = SpeechConfig.fromSubscription(
  subscriptionKey,
  serviceRegion,
);
speechConfig.speechRecognitionLanguage = 'zh-CN';
speechConfig.speechSynthesisLanguage = 'zh-CN';

// Favorite three for now
const myVoices = [
  'zh-CN-YunxiNeural',
  'zh-CN-XiaochenNeural',
  'zh-CN-XiaoshuangNeural', // child
];
const nextVoice = (function nextVoice() {
  let next = 0;
  return () => {
    next = (next + 1) % 3;
    return next;
  };
}());

export async function synthesize(text) {
  const voice = myVoices[nextVoice()];
  const audioFile = tmp.fileSync({ postfix: '.wav' });
  const audioConfig = AudioConfig.fromAudioFileOutput(audioFile.name);
  speechConfig.speechSynthesisVoiceName = voice;
  let synthesizer = new SpeechSynthesizer(speechConfig, audioConfig);
  return new Promise((resolve, reject) => {
    synthesizer.speakTextAsync(
      text,
      (result) => {
        console.log(result);
        if (result.reason === ResultReason.SynthesizingAudioCompleted) {
          console.log('synthesis finished.');
          console.log(`Saved ${audioFile.name}`);
        } else {
          console.error(`Speech synthesis canceled, ${result.errorDetails
          }\nDid you update the subscription info?`);
        }
        synthesizer.close();
        synthesizer = undefined;
        resolve(audioFile.name);
      },
      (err) => {
        console.trace(`err - ${err}`);
        synthesizer.close();
        synthesizer = undefined;
        reject(err);
      },
    );
  });
}
