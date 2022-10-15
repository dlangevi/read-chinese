import { inject } from 'vue';
import { UserSettingsKey } from '@/shared/types';
import type { UserSettingsType } from '@/shared/types';

export function getUserSettings():UserSettingsType {
  return inject(UserSettingsKey) as UserSettingsType;
}
