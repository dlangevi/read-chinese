module.exports = {
  mode: 'jit',
  important: true,
  purge: [
    './src/**/*.{ts,js,vue}',
  ],
  theme: {
    extend: {
      colors: {
        primary: '#5438DC',
        'primary-dark': '#4125D0',
        secondary: '#48a9a6',
        neutral: '#F3DE8A',
        warning: '#cd533b',
        success: '#415D43',
      },
    },
  },
};
