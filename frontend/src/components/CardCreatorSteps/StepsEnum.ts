export const StepsEnum = {
  NONE: 'none',
  SENTENCE: 'sentence',
  IMAGE: 'image',
  ENGLISH: 'english',
  CHINESE: 'chinese',
} as const;

export type StepsEnum = typeof StepsEnum[keyof typeof StepsEnum]
export type StepsEnumType = typeof StepsEnum[keyof typeof StepsEnum]
