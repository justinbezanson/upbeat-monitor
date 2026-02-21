const tailwindcssPlugin = require('@tailwindcss/postcss');
const autoprefixerPlugin = require('autoprefixer');

module.exports = {
  plugins: {
    '@tailwindcss/postcss': {},
    autoprefixer: autoprefixerPlugin(),
  },
};
