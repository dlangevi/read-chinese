module.exports = {
  '*.{js,jsx,vue}': 'eslint',
  '*.{vue,ts}': [
    () => 'vue-tsc --noEmit',
  ],
};
