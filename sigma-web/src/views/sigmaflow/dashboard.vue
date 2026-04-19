<script setup lang="ts">
import { ref } from 'vue'
import { CircleDollarSign, ArrowLeft, Activity } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const router = useRouter()

const stats = ref([
  { label: 'Total Invoices', value: '42', color: 'indigo' },
  { label: 'Pending Payment', value: 'Rp 12.500.000', color: 'amber' },
  { label: 'Revenue (MTD)', value: 'Rp 85.200.000', color: 'emerald' }
])
</script>

<template>
  <div class="min-h-screen bg-slate-950 text-white p-8">
    <div class="max-w-6xl mx-auto">
      <div class="flex items-center justify-between mb-12">
        <div class="flex items-center gap-4">
          <div class="p-3 bg-indigo-500/10 text-indigo-400 rounded-2xl border border-indigo-500/20">
            <CircleDollarSign class="w-8 h-8" />
          </div>
          <div>
            <h1 class="text-3xl font-bold tracking-tight">Sigmaflow</h1>
            <p class="text-slate-400">Cashflow & Invoice Management</p>
          </div>
        </div>
        
        <button @click="router.push('/dashboard')" 
                class="flex items-center gap-2 px-6 py-3 bg-slate-900 hover:bg-slate-800 rounded-xl border border-slate-800 transition-all group">
          <ArrowLeft class="w-5 h-5 group-hover:-translate-x-1 transition-transform" />
          Kembali ke Dashboard
        </button>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
        <div v-for="stat in stats" :key="stat.label" 
             class="p-6 rounded-3xl bg-slate-900 border border-slate-800 hover:border-indigo-500/50 transition-colors">
          <p class="text-slate-400 text-sm mb-1">{{ stat.label }}</p>
          <p class="text-2xl font-bold">{{ stat.value }}</p>
        </div>
      </div>

      <div class="bg-slate-900/50 border border-slate-800 rounded-[2.5rem] p-8 text-center">
        <div class="w-20 h-20 bg-slate-800 rounded-full flex items-center justify-center mx-auto mb-6">
          <Activity class="w-10 h-10 text-slate-500" />
        </div>
        <h2 class="text-2xl font-bold mb-2">Module Under Development</h2>
        <p class="text-slate-400 max-w-md mx-auto">
          Sistem manajemen arus kas dan faktur sedang dalam tahap sinkronisasi data. 
          Silakan hubungi tim keuangan untuk informasi lebih lanjut.
        </p>
      </div>
    </div>
  </div>
</template>
