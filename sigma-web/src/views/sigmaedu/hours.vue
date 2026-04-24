<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { 
  Plus, Clock, ArrowLeft, Sun, Moon,
  Check, X, Coffee
} from 'lucide-vue-next'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()


// State
const hours = ref<any[]>([])
const isLoading = ref(true)
const isModalOpen = ref(false)
const isSubmitting = ref(false)
const editingId = ref<number | null>(null)

const form = ref({
  name: '',
  start_time: '',
  end_time: '',
  is_break: false
})


const fetchHours = async () => {
  try {
    const response = await axios.get('/api/v1/admin/edu/hours')
    hours.value = response.data
  } catch (err) {
    console.error('Gagal mengambil data jam pelajaran:', err)
  } finally {
    isLoading.value = false
  }
}

const openModal = (hour: any = null) => {
  if (hour) {
    editingId.value = hour.ID
    form.value = { 
      name: hour.name, 
      start_time: hour.start_time, 
      end_time: hour.end_time,
      is_break: hour.is_break 
    }
  } else {
    editingId.value = null
    form.value = { name: '', start_time: '', end_time: '', is_break: false }
  }
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
}

const handleSubmit = async () => {
  isSubmitting.value = true
  try {
    if (editingId.value) {
      await axios.put(`/api/v1/admin/edu/hours/${editingId.value}`, form.value)
    } else {
      await axios.post('/api/v1/admin/edu/hours', form.value)
    }
    await fetchHours()
    closeModal()
  } catch (err) {
    alert('Gagal menyimpan data.')
  } finally {
    isSubmitting.value = false
  }
}

const deleteHour = async (id: number) => {
  if (!confirm('Yakin ingin menghapus jam pelajaran ini?')) return
  try {
    await axios.delete(`/api/v1/admin/edu/hours/${id}`)
    await fetchHours()
  } catch (err) {
    alert('Gagal menghapus data.')
  }
}

onMounted(() => {
  fetchHours()
})
</script>

<template>
  <div class="p-8 space-y-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <div class="h-10 w-1 bg-blue-500 rounded-full"></div>
        <div>
          <h2 class="text-xs font-black uppercase tracking-[0.3em] text-sigma-muted mb-1">Kurikulum</h2>
          <h1 class="text-2xl font-black italic uppercase tracking-tight">Jam <span class="text-blue-500">Pelajaran</span></h1>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <button @click="router.push('/portal')" 
                class="flex items-center gap-3 px-5 py-2.5 rounded-2xl bg-sigma-surface border border-sigma-border text-sigma-muted hover:text-sigma-text hover:border-blue-500/30 transition-all group shadow-sm">
          <ArrowLeft class="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
          <span class="text-[10px] font-black uppercase tracking-widest">Portal</span>
        </button>
      </div>
    </div>

    <!-- Action Bar -->
    <div class="flex flex-col md:flex-row gap-4 justify-end items-center bg-sigma-surface p-4 rounded-[2rem] border border-sigma-border shadow-sm">
      <button @click="openModal()" 
              class="w-full md:w-auto px-8 py-3.5 bg-blue-600 text-white rounded-2xl font-black uppercase tracking-widest text-[10px] hover:bg-blue-700 transition-all shadow-xl shadow-blue-500/20 flex items-center justify-center gap-3">
        <Plus class="w-4 h-4" /> Tambah Sesi Waktu
      </button>
    </div>

    <!-- Hours Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div v-for="hour in hours" :key="hour.ID" 
           class="bg-sigma-surface border border-sigma-border p-6 rounded-[2.5rem] relative overflow-hidden group hover:border-blue-500/30 transition-all">
        
        <div v-if="hour.is_break" class="absolute top-0 right-0 p-3">
          <span class="px-3 py-1 rounded-full bg-amber-500/10 text-amber-500 text-[8px] font-black uppercase tracking-widest border border-amber-500/20">ISTIRAHAT</span>
        </div>

        <div class="flex flex-col gap-6">
          <div class="w-12 h-12 rounded-2xl flex items-center justify-center transition-all"
               :class="hour.is_break ? 'bg-amber-500/10 text-amber-500' : 'bg-blue-500/10 text-blue-500 group-hover:bg-blue-500 group-hover:text-white'">
            <component :is="hour.is_break ? Coffee : Clock" class="w-6 h-6" />
          </div>

          <div>
            <h3 class="text-lg font-black uppercase italic tracking-tighter">{{ hour.name }}</h3>
            <div class="flex items-center gap-2 mt-2">
              <span class="text-2xl font-black text-sigma-text">{{ hour.start_time }}</span>
              <div class="w-2 h-0.5 bg-sigma-muted"></div>
              <span class="text-2xl font-black text-sigma-muted">{{ hour.end_time }}</span>
            </div>
          </div>

          <div class="flex gap-2 pt-4 border-t border-sigma-border">
            <button @click="openModal(hour)" class="flex-1 py-2 rounded-xl bg-sigma-surface-alt text-sigma-muted hover:text-blue-500 hover:bg-blue-500/5 transition-all text-[10px] font-black uppercase tracking-widest border border-sigma-border">Edit</button>
            <button @click="deleteHour(hour.ID)" class="flex-1 py-2 rounded-xl bg-sigma-surface-alt text-sigma-muted hover:text-red-500 hover:bg-red-500/5 transition-all text-[10px] font-black uppercase tracking-widest border border-sigma-border">Hapus</button>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="hours.length === 0 && !isLoading" 
           class="md:col-span-4 py-20 bg-sigma-surface border border-sigma-border border-dashed rounded-[3rem] flex flex-col items-center justify-center gap-4 opacity-30">
        <Clock class="w-16 h-16" />
        <p class="font-black uppercase tracking-widest text-xs">Belum ada pengaturan jam pelajaran</p>
      </div>
    </div>

    <!-- Modal Add/Edit -->
    <div v-if="isModalOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-6 backdrop-blur-xl bg-black/40">
      <div class="bg-sigma-surface w-full max-w-lg rounded-[3rem] border border-sigma-border shadow-2xl overflow-hidden animate-in zoom-in duration-300">
        <div class="p-8 bg-gradient-to-r from-blue-600 to-blue-900 flex justify-between items-center border-b border-white/10">
          <div>
            <h3 class="text-white font-black uppercase tracking-widest text-xs">{{ editingId ? 'Update' : 'Tambah' }} Sesi Waktu</h3>
            <p class="text-blue-100/60 text-[10px] mt-1 font-medium">Pengaturan Jam KBM Sigmaedu</p>
          </div>
          <button @click="closeModal" class="p-2 rounded-xl bg-white/10 text-white hover:bg-white/20 transition-all">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="p-10 space-y-6">
          <div class="space-y-2">
            <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Nama Sesi (e.g. Jam Ke-1)</label>
            <input v-model="form.name" type="text" placeholder="Masukkan nama sesi..." 
                   class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight" />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Mulai</label>
              <input v-model="form.start_time" type="time" 
                     class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight" />
            </div>
            <div class="space-y-2">
              <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Selesai</label>
              <input v-model="form.end_time" type="time" 
                     class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight" />
            </div>
          </div>

          <div class="flex items-center gap-4 p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border">
            <button @click="form.is_break = !form.is_break" 
                    class="w-12 h-6 rounded-full relative transition-all"
                    :class="form.is_break ? 'bg-blue-500' : 'bg-sigma-border'">
              <div class="absolute top-1 w-4 h-4 bg-white rounded-full transition-all"
                   :class="form.is_break ? 'left-7' : 'left-1'"></div>
            </button>
            <div class="flex flex-col">
              <span class="text-xs font-black uppercase tracking-widest text-sigma-text">Waktu Istirahat</span>
              <span class="text-[10px] text-sigma-muted">Tandai jika sesi ini bukan jam pelajaran (Istirahat/Sholat).</span>
            </div>
          </div>

          <div class="flex gap-4 pt-4">
            <button @click="closeModal" class="flex-1 px-8 py-4 bg-sigma-surface-alt border border-sigma-border text-sigma-text rounded-2xl font-black uppercase tracking-widest text-[10px] hover:bg-sigma-border transition-all">
              Batal
            </button>
            <button @click="handleSubmit" :disabled="isSubmitting"
                    class="flex-1 px-8 py-4 bg-blue-600 text-white rounded-2xl font-black uppercase tracking-widest text-[10px] hover:bg-blue-700 transition-all shadow-xl shadow-blue-500/20 disabled:opacity-50 flex items-center justify-center gap-2">
              <Check v-if="!isSubmitting" class="w-4 h-4" />
              <div v-else class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
              {{ editingId ? 'Simpan Perubahan' : 'Tambah Data' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>


<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 130, 246, 0.1); border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(59, 130, 246, 0.2); }
</style>
