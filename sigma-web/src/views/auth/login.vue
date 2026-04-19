<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const handleLogin = () => {
  // Nanti di sini kita panggil API Go kamu
  // Untuk sekarang, langsung lempar ke portal
  router.push('/portal')
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-950 text-white">
    <div class="glass p-8 rounded-[2rem] border border-slate-800 w-full max-w-md text-center">
      <img src="/logo/SIGMA.svg" class="w-20 h-20 mx-auto mb-6 drop-shadow-[0_0_15px_rgba(16,185,129,0.3)]" alt="SIGMA Logo" />
      <h1 class="text-3xl font-bold mb-6 italic tracking-tighter">SIGMA <span class="text-emerald-500">V2</span></h1>
      
      <div class="space-y-4 mb-6">
        <input type="text" placeholder="Username" class="w-full p-4 rounded-2xl bg-slate-900 border border-slate-800 focus:border-emerald-500 outline-none transition-all" />
        <input type="password" placeholder="Password" class="w-full p-4 rounded-2xl bg-slate-900 border border-slate-800 focus:border-emerald-500 outline-none transition-all" />
      </div>

      <button @click="handleLogin" class="w-full py-4 bg-emerald-500 rounded-2xl font-bold hover:bg-emerald-600 transition-all shadow-lg shadow-emerald-500/20">
        Masuk ke Sistem
      </button>
    </div>
  </div>
</template>