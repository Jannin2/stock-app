// frontend/postcss.config.cjs
module.exports = {
  plugins: {
    // Para Tailwind CSS v4.x, el plugin PostCSS es ahora '@tailwindcss/postcss'
    '@tailwindcss/postcss': {}, // <--- ESTA ES LA CLAVE DEL CAMBIO
    'autoprefixer': {},
  },
};