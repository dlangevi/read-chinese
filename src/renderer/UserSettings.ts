import { inject } from 'vue';
import { UserSettingsKey, UserSettingsType } from '../shared/types';

export function getUserSettings():UserSettingsType {
  return inject(UserSettingsKey) as UserSettingsType;
}
