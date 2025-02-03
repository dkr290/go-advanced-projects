/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.templ", "./**/*.templ"],
  safelist: [],
  plugins: [require("daisyui")],
};
