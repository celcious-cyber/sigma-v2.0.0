<template>
  <div class="min-h-screen flex items-center justify-center p-6 relative overflow-hidden">
    <!-- Background Accents -->
    <div class="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-indigo-600/20 blur-[120px] rounded-full"></div>
    <div class="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-sky-600/20 blur-[120px] rounded-full"></div>

    <div class="glass-card w-full max-w-md p-10 space-y-8 relative z-10 fade-in">
      <div class="text-center space-y-2">
        <h1 class="text-4xl font-black tracking-tight premium-gradient-text">SIGMA v2.0</h1>
        <p class="text-slate-400 text-sm">Masuk untuk mengelola ekosistem pondok</p>
      </div>

      <form @submit.prevent="handleLogin" class="space-y-6">
        <div class="space-y-2">
          <label class="text-xs font-semibold text-slate-400 uppercase tracking-wider ml-1">Email</label>
          <input 
            v-model="email"
            type="email" 
            placeholder="admin@sigma.com" 
            class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 outline-none focus:border-indigo-500 transition-colors"
            required
          />
        </div>

        <div class="space-y-2">
          <label class="text-xs font-semibold text-slate-400 uppercase tracking-wider ml-1">Kata Sandi</label>
          <input 
            v-model="password"
            type="password" 
            placeholder="••••••••" 
            class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 outline-none focus:border-indigo-500 transition-colors"
            required
          />
        </div>

        <button 
          type="submit" 
          :disabled="loading"
          class="w-full bg-gradient-to-r from-indigo-600 to-indigo-500 hover:from-indigo-500 hover:to-indigo-400 text-white font-bold py-3 rounded-xl shadow-lg shadow-indigo-600/20 transition-all hover-scale disabled:opacity-50"
        >
          {{ loading ? 'Mencoba Masuk...' : 'MASUK SEKARANG' }}
        </button>
      </form>

      <div class="pt-4 text-center">
        <a href="#" class="text-xs text-slate-500 hover:text-indigo-400 transition-colors">Lupa kata sandi? Hubungi admin pusat.</a>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const email = ref('')
const password = ref('')
const loading = ref(false)
const router = useRouter()

const handleLogin = async () => {
  loading.value = true
  try {
    const response = await axios.post('/api/v1/auth/admin/login', {
      identifier: email.value,
      password: password.value
    })
    
    // Simple state management for now
    localStorage.setItem('sigma_token', response.data.token)
    localStorage.setItem('sigma_user', JSON.stringify(response.data.user))
    
    router.push('/portal')
  } catch (error) {
    alert('Email atau password salah.')
  } finally {
    loading.value = false
  }
}
</script>
