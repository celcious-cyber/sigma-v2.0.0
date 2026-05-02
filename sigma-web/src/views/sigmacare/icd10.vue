<script setup lang="ts">
import { ref } from 'vue'
import { 
  Database, Search, Filter, 
  ChevronRight, 
  BookOpen
} from 'lucide-vue-next'

const diseases = ref([
  { code: 'J00', name: 'Nasofaringitis Akut (Common Cold)', category: 'Pernapasan' },
  { code: 'K29.7', name: 'Gastritis, tidak spesifik', category: 'Pencernaan' },
  { code: 'L20', name: 'Dermatitis Atopik', category: 'Kulit' },
  { code: 'A09', name: 'Diare dan Gastroenteritis', category: 'Infeksi' },
  { code: 'G44.2', name: 'Tension-type Headache', category: 'Saraf' },
  { code: 'H10.9', name: 'Konjungtivitis, tidak spesifik', category: 'Mata' }
])
</script>

<template>
  <div class="p-8 space-y-8 min-h-screen bg-slate-50/50">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
      <div>
        <h1 class="text-3xl font-black tracking-tight text-slate-800 italic uppercase">Katalog <span class="text-rose-600">Patologi</span></h1>
        <p class="text-slate-500 font-bold mt-1 tracking-wide uppercase text-xs flex items-center gap-2">
          <Database class="w-4 h-4" /> Referensi Standar ICD-10 (Klasifikasi Penyakit Internasional)
        </p>
      </div>
      
      <div class="relative">
        <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
        <input 
          type="text" 
          placeholder="Cari Kode atau Nama Penyakit..." 
          class="pl-11 pr-6 py-3.5 bg-white border border-slate-200 rounded-2xl text-sm font-bold focus:outline-none focus:ring-4 focus:ring-rose-500/10 focus:border-rose-500 transition-all w-64 md:w-80 shadow-sm"
        />
      </div>
    </div>

    <!-- Info Banner -->
    <div class="bg-slate-800 rounded-[2.5rem] p-10 text-white relative overflow-hidden shadow-2xl">
      <div class="absolute -right-10 -bottom-10 opacity-10">
        <BookOpen class="w-48 h-48 rotate-12" />
      </div>
      <div class="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-8">
        <div class="space-y-4">
          <h2 class="text-2xl font-black italic uppercase tracking-tighter flex items-center gap-4">
            <span class="w-12 h-12 bg-rose-600 rounded-2xl flex items-center justify-center text-white italic">i</span>
            Standar Pengodean Klinis
          </h2>
          <p class="text-slate-400 font-bold text-sm max-w-xl leading-relaxed uppercase tracking-wider">
            Gunakan katalog ini untuk memastikan pengodean rekam medis santri sesuai dengan standar kesehatan internasional guna sinkronisasi data yang akurat.
          </p>
        </div>
        <div class="flex gap-4">
          <div class="bg-white/5 p-6 rounded-3xl backdrop-blur-md border border-white/10 text-center min-w-[120px]">
            <p class="text-[10px] font-black text-white/40 uppercase tracking-widest mb-1">Total Kode</p>
            <p class="text-2xl font-black italic tracking-tighter">14,400+</p>
          </div>
          <div class="bg-white/5 p-6 rounded-3xl backdrop-blur-md border border-white/10 text-center min-w-[120px]">
            <p class="text-[10px] font-black text-white/40 uppercase tracking-widest mb-1">Edisi</p>
            <p class="text-2xl font-black italic tracking-tighter">ICD-10</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Disease Table -->
    <div class="bg-white rounded-[2.5rem] border border-slate-100 shadow-sm overflow-hidden">
      <div class="p-8 border-b border-slate-50 bg-slate-50/30 flex items-center justify-between">
        <h3 class="text-lg font-black text-slate-800 uppercase italic">Katalog Penyakit UKS</h3>
        <button class="p-2.5 bg-white border border-slate-200 text-slate-400 rounded-xl hover:text-rose-600 transition-colors">
          <Filter class="w-4 h-4" />
        </button>
      </div>
      
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead>
            <tr class="bg-slate-50/50">
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-[0.2em]">Kode ICD-10</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-[0.2em]">Nama Penyakit / Keluhan</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-[0.2em]">Kategori Sistem</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-[0.2em] text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-50">
            <tr v-for="d in diseases" :key="d.code" class="hover:bg-rose-50/20 transition-all group">
              <td class="px-8 py-6">
                <div class="flex items-center gap-4">
                  <div class="w-10 h-10 bg-slate-50 text-slate-400 rounded-xl flex items-center justify-center font-black text-xs group-hover:bg-rose-600 group-hover:text-white transition-all">
                    {{ d.code }}
                  </div>
                </div>
              </td>
              <td class="px-8 py-6">
                <p class="font-black text-slate-800 uppercase tracking-tight group-hover:text-rose-600 transition-colors">{{ d.name }}</p>
              </td>
              <td class="px-8 py-6">
                <span class="px-3 py-1 bg-slate-100 text-slate-500 rounded-lg text-[9px] font-black uppercase tracking-widest italic">{{ d.category }}</span>
              </td>
              <td class="px-8 py-6 text-right">
                <button class="p-2 text-slate-300 hover:text-slate-600 transition-colors">
                  <ChevronRight class="w-5 h-5" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
