/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{vue,js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                // Kita siapin warna khas SIGMA (Biru Modern)
                "sigma-blue": "#1e40af",
                "sigma-dark": "#0f172a",
            }
        },
    },
    plugins: [],
}