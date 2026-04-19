<template>
  <div class="space-y-8 fade-in">
    <!-- Welcome Section -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-black text-white">Ahlan wa Sahlan, Admin 👋</h1>
        <p class="text-slate-500 mt-1">Ini ringkasan aktivitas ekosistem SIGMA hari ini.</p>
      </div>
      <div class="flex space-x-3">
        <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-4 py-2 rounded-xl text-sm font-bold transition-all shadow-lg shadow-indigo-600/20">
          + DATA SANTRI
        </button>
        <button class="bg-white/5 hover:bg-white/10 text-white px-4 py-2 rounded-xl text-sm font-bold border border-white/10 transition-all">
          CETAK LAPORAN
        </button>
      </div>
    </div>

    <!-- Quick Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="glass-card p-6 border-l-4 border-indigo-500 hover-scale">
        <p class="text-xs font-bold text-slate-500 uppercase tracking-widest">Total Santri</p>
        <h3 class="text-3xl font-black mt-2">1,248</h3>
        <p class="text-[10px] text-indigo-400 mt-1 font-bold">+12 Santri Baru (Bulan ini)</p>
      </div>
      
      <div class="glass-card p-6 border-l-4 border-emerald-500 hover-scale">
        <p class="text-xs font-bold text-slate-500 uppercase tracking-widest">Keuangan (Bulan ini)</p>
        <h3 class="text-3xl font-black mt-2 text-emerald-400">Rp 84,2M</h3>
        <p class="text-[10px] text-slate-500 mt-1 font-bold">85% Target Terpenuhi</p>
      </div>

      <div class="glass-card p-6 border-l-4 border-amber-500 hover-scale">
        <p class="text-xs font-bold text-slate-500 uppercase tracking-widest">Pelanggaran</p>
        <h3 class="text-3xl font-black mt-2 text-amber-400">32</h3>
        <p class="text-[10px] text-slate-500 mt-1 font-bold">Menurun 15% dari kemarin</p>
      </div>

      <div class="glass-card p-6 border-l-4 border-sky-500 hover-scale">
        <p class="text-xs font-bold text-slate-500 uppercase tracking-widest">UKS / Sakit</p>
        <h3 class="text-3xl font-black mt-2 text-sky-400">5</h3>
        <p class="text-[10px] text-slate-500 mt-1 font-bold">Sedang dalam perawatan</p>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Recent Activity -->
      <div class="lg:col-span-2 glass-card overflow-hidden">
        <div class="p-6 border-b border-white/5 flex items-center justify-between">
          <h3 class="font-bold uppercase tracking-widest text-sm">Aktivitas Terkini</h3>
          <button class="text-xs text-indigo-400 font-bold hover:underline">LIHAT SEMUA</button>
        </div>
        <div class="divide-y divide-white/5">
          <div v-for="i in 5" :key="i" class="p-6 flex items-center justify-between hover:bg-white/5 transition-colors">
            <div class="flex items-center space-x-4">
              <div class="w-10 h-10 rounded-full bg-indigo-500/10 flex items-center justify-center text-indigo-400">
                <component :is="icons[i % icons.length]" :size="18" />
              </div>
              <div>
                <p class="text-sm font-bold">Aktivitas #{{ i }} - Baru</p>
                <p class="text-[10px] text-slate-500">2 jam yang lalu • Oleh Admin Utama</p>
              </div>
            </div>
            <div class="px-3 py-1 rounded-full bg-white/5 text-[10px] font-bold uppercase tracking-tighter">
              PENDING
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Shortcuts -->
      <div class="space-y-6">
        <div class="glass-card p-6 bg-gradient-to-br from-indigo-600/20 to-transparent">
          <h3 class="font-bold text-sm mb-4 tracking-widest uppercase">Pintas Cepat</h3>
          <div class="grid grid-cols-2 gap-4">
            <button class="shortcut-btn"><Zap :size="16" /><span>Izin Kilat</span></button>
            <button class="shortcut-btn"><FileText :size="16" /><span>SK Santri</span></button>
            <button class="shortcut-btn"><TrendingUp :size="16" /><span>Grafik Keu</span></button>
            <button class="shortcut-btn"><Archive :size="16" /><span>Arsip</span></button>
          </div>
        </div>

        <div class="glass-card p-6 border border-indigo-500/30">
          <div class="flex items-center space-x-2 mb-2">
            <div class="w-2 h-2 rounded-full bg-indigo-400 animate-pulse"></div>
            <p class="text-[10px] font-bold text-indigo-400 uppercase tracking-widest">Sistem Online</p>
          </div>
          <p class="text-xs text-slate-400">Node Go API terhubung di 127.0.0.1:3000 via SQLite WAL Mode.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Users, CreditCard, ShieldAlert, Zap, FileText, TrendingUp, Archive, GraduationCap } from 'lucide-vue-next'

const icons = [Users, CreditCard, ShieldAlert, GraduationCap, Zap]
</script>

<style scoped>
@reference "tailwindcss";

.shortcut-btn {
  @apply flex flex-col items-center justify-center space-y-2 p-4 rounded-xl bg-white/5 border border-white/5 hover:bg-white/10 hover:border-white/20 transition-all text-xs font-bold text-slate-300;
}
</style>
