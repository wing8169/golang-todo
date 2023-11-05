/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.{html,js,templ,go}"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/forms"), require("@tailwindcss/typography")],
};
