module.exports = {
  '*.{js,ts,vue}': 'eslint',
  '*.{vue,ts}': [
    () => 'vue-tsc --noEmit',
  ],
};
