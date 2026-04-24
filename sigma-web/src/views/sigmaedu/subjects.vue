<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { 
  Search, Plus, Pencil, Trash2, 
  BookOpen, Hash, ArrowLeft, Sun, Moon,
  Check, X
} from 'lucide-vue-next'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()


// State
const subjects = ref<any[]>([])
const isLoading = ref(true)
const searchQuery = ref('')
const isModalOpen = ref(false)
const isSubmitting = ref(false)
const editingId = ref<number | null>(null)

const form = ref({
  name: '',
  code: ''
})


const fetchSubjects = async () => {
  try {
    const response = await axios.get('/api/v1/admin/edu/subjects')
    subjects.value = response.data
  } catch (err) {
    console.error('Gagal mengambil data mapel:', err)
  } finally {
    isLoading.value = false
  }
}

const openModal = (subject: any = null) => {
  if (subject) {
    editingId.value = subject.ID
    form.value = { name: subject.name, code: subject.code }
  } else {
    editingId.value = null
    form.value = { name: '', code: '' }
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
      await axios.put(`/api/v1/admin/edu/subjects/${editingId.value}`, form.value)
    } else {
      await axios.post('/api/v1/admin/edu/subjects', form.value)
    }
    await fetchSubjects()
    closeModal()
  } catch (err) {
    alert('Gagal menyimpan data. Pastikan kode unik.')
  } finally {
    isSubmitting.value = false
  }
}

const deleteSubject = async (id: number) => {
  if (!confirm('Yakin ingin menghapus mata pelajaran ini?')) return
  try {
    await axios.delete(`/api/v1/admin/edu/subjects/${id}`)
    await fetchSubjects()
  } catch (err) {
    alert('Gagal menghapus data.')
  }
}

onMounted(() => {
  fetchSubjects()
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
          <h1 class="text-2xl font-black italic uppercase tracking-tight">Mata <span class="text-blue-500">Pelajaran</span></h1>
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
    <div class="flex flex-col md:flex-row gap-4 justify-between items-center bg-sigma-surface p-4 rounded-[2rem] border border-sigma-border shadow-sm">
      <div class="relative w-full md:w-96">
        <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-sigma-muted" />
        <input v-model="searchQuery" 
               type="text" 
               placeholder="Cari Mata Pelajaran..." 
               class="w-full pl-12 pr-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight transition-all" />
      </div>
      <button @click="openModal()" 
              class="w-full md:w-auto px-8 py-3.5 bg-blue-600 text-white rounded-2xl font-black uppercase tracking-widest text-[10px] hover:bg-blue-700 transition-all shadow-xl shadow-blue-500/20 flex items-center justify-center gap-3">
        <Plus class="w-4 h-4" /> Tambah Pelajaran
      </button>
    </div>

    <!-- Subjects Table -->
    <div class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] overflow-hidden shadow-sm">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="text-[10px] text-sigma-muted uppercase tracking-[0.3em] border-b border-sigma-border bg-sigma-surface-alt/20">
              <th class="p-8">Mata Pelajaran</th>
              <th class="p-8">Kode Mapel</th>
              <th class="p-8 text-right">Opsi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-sigma-border">
            <tr v-for="subject in subjects" :key="subject.ID" 
                class="hover:bg-sigma-surface-alt/50 transition-all group">
              <td class="p-8">
                <div class="flex items-center gap-4">
                  <div class="w-12 h-12 rounded-2xl bg-blue-500/10 border border-blue-500/10 flex items-center justify-center text-blue-500 group-hover:bg-blue-500 group-hover:text-white transition-all duration-500">
                    <BookOpen class="w-6 h-6" />
                  </div>
                  <div>
                    <div class="font-black text-sigma-text text-sm uppercase tracking-tight">{{ subject.name }}</div>
                    <div class="text-[10px] text-sigma-muted mt-0.5">Mata Pelajaran Aktif</div>
                  </div>
                </div>
              </td>
              <td class="p-8">
                <div class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-sigma-surface-alt border border-sigma-border w-fit">
                  <Hash class="w-3 h-3 text-blue-500" />
                  <span class="text-xs font-mono font-black text-sigma-text tracking-widest">{{ subject.code }}</span>
                </div>
              </td>
              <td class="p-8 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button @click="openModal(subject)" class="p-3 rounded-xl bg-amber-500/10 text-amber-500 hover:bg-amber-500 hover:text-white transition-all border border-amber-500/10">
                    <Pencil class="w-4 h-4" />
                  </button>
                  <button @click="deleteSubject(subject.ID)" class="p-3 rounded-xl bg-red-500/10 text-red-400 hover:bg-red-500 hover:text-white transition-all border border-red-500/10">
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="subjects.length === 0 && !isLoading">
              <td colspan="3" class="p-20 text-center">
                <div class="flex flex-col items-center gap-4 opacity-30">
                  <BookOpen class="w-16 h-16" />
                  <p class="font-black uppercase tracking-widest text-xs">Belum ada data mata pelajaran</p>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Add/Edit -->
    <div v-if="isModalOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-6 backdrop-blur-xl bg-black/40">
      <div class="bg-sigma-surface w-full max-w-lg rounded-[3rem] border border-sigma-border shadow-2xl overflow-hidden animate-in zoom-in duration-300">
        <div class="p-8 bg-gradient-to-r from-blue-600 to-blue-900 flex justify-between items-center border-b border-white/10">
          <div>
            <h3 class="text-white font-black uppercase tracking-widest text-xs">{{ editingId ? 'Update' : 'Tambah' }} Pelajaran</h3>
            <p class="text-blue-100/60 text-[10px] mt-1 font-medium">Sistem Kurikulum Sigmaedu</p>
          </div>
          <button @click="closeModal" class="p-2 rounded-xl bg-white/10 text-white hover:bg-white/20 transition-all">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="p-10 space-y-6">
          <div class="space-y-2">
            <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Nama Mata Pelajaran</label>
            <div class="relative group">
              <BookOpen class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-sigma-muted group-focus-within:text-blue-500 transition-all" />
              <input v-model="form.name" type="text" placeholder="Contoh: Matematika" 
                     class="w-full pl-12 pr-4 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight" />
            </div>
          </div>

          <div class="space-y-2">
            <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Kode Pelajaran</label>
            <div class="relative group">
              <Hash class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-sigma-muted group-focus-within:text-blue-500 transition-all" />
              <input v-model="form.code" type="text" placeholder="Contoh: MTK-01" 
                     class="w-full pl-12 pr-4 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight uppercase" />
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
