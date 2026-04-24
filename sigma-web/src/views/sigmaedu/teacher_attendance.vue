<script setup lang="ts">
import { ref, onMounted, watch, reactive, computed } from 'vue'
import { 
  Users, Save, Calendar as CalendarIcon, 
  CheckCircle2, XCircle, Clock,
  Plus, Trash2, Edit, Search, RotateCw
} from 'lucide-vue-next'
import axios from 'axios'

// State
const journals = ref<any[]>([])
const classrooms = ref<any[]>([])
const subjects = ref<any[]>([])
const studyHours = ref<any[]>([])
const teachers = ref<any[]>([])
const isLoading = ref(false)
const isSaving = ref(false)
const showAddModal = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)
const searchQuery = ref('')

const filters = ref({
  date: new Date().toISOString().split('T')[0],
  classroom_id: null as number | null,
  subject_id: null as number | null
})

const filteredJournals = computed(() => {
  if (!searchQuery.value) return journals.value
  const query = searchQuery.value.toLowerCase()
  return journals.value.filter(j => 
    j.teacher?.name?.toLowerCase().includes(query) ||
    j.subject?.name?.toLowerCase().includes(query) ||
    j.classroom?.name?.toLowerCase().includes(query)
  )
})

const newJournal = reactive({
  teacher_id: null as number | null,
  classroom_id: null as number | null,
  subject_id: null as number | null,
  study_hour_id: null as number | null,
  status: 'Hadir',
  remarks: '',
  date: new Date().toISOString().split('T')[0]
})

const fetchData = async () => {
  try {
    const [clRes, subRes, hrRes, tRes] = await Promise.all([
      axios.get('/api/v1/admin/edu/classrooms'),
      axios.get('/api/v1/admin/edu/subjects'),
      axios.get('/api/v1/admin/edu/hours'),
      axios.get('/api/v1/admin/base/teachers')
    ])
    classrooms.value = clRes.data
    subjects.value = subRes.data
    studyHours.value = hrRes.data
    teachers.value = tRes.data
    
    await loadJournals()
  } catch (err) {
    console.error('Gagal mengambil data awal:', err)
  }
}

const loadJournals = async () => {
  isLoading.value = true
  try {
    const params: any = { date: filters.value.date }
    if (filters.value.classroom_id) params.classroom_id = filters.value.classroom_id
    if (filters.value.subject_id) params.subject_id = filters.value.subject_id
    
    const res = await axios.get('/api/v1/admin/edu/journals', { params })
    journals.value = res.data
  } catch (err) {
    console.error('Gagal memuat presensi:', err)
  } finally {
    isLoading.value = false
  }
}

const refreshData = () => loadJournals()

const showToast = ref(false)
const toastMessage = ref('')
const toastType = ref('success')

const openAddModal = () => {
  isEdit.value = false
  editingId.value = null
  newJournal.teacher_id = null
  newJournal.classroom_id = null
  newJournal.subject_id = null
  newJournal.study_hour_id = null
  newJournal.status = 'Hadir'
  newJournal.remarks = ''
  showAddModal.value = true
}

const openEditModal = (j: any) => {
  isEdit.value = true
  editingId.value = j.ID
  newJournal.teacher_id = j.teacher_id
  newJournal.classroom_id = j.classroom_id
  newJournal.subject_id = j.subject_id
  newJournal.study_hour_id = j.study_hour_id
  newJournal.status = j.status
  newJournal.remarks = j.remarks
  showAddModal.value = true
}

const handleSubmit = async () => {
  if (!newJournal.teacher_id || !newJournal.classroom_id || !newJournal.subject_id || !newJournal.study_hour_id) {
    toastMessage.value = 'Lengkapi semua data!'
    toastType.value = 'error'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
    return
  }

  isSaving.value = true
  try {
    const payload = {
      ...newJournal,
      date: new Date(filters.value.date).toISOString()
    }
    
    if (isEdit.value && editingId.value) {
      await axios.put(`/api/v1/admin/edu/journals/${editingId.value}`, payload)
      toastMessage.value = 'Presensi Berhasil Diperbarui!'
    } else {
      await axios.post('/api/v1/admin/edu/journals', payload)
      toastMessage.value = 'Presensi Berhasil Ditambahkan!'
    }
    
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
    
    showAddModal.value = false
    await loadJournals()
  } catch (err) {
    console.error('Gagal menyimpan presensi:', err)
  } finally {
    isSaving.value = false
  }
}

const handleDelete = async (id: any) => {
  const targetId = id
  if (!targetId) {
    console.error('ID tidak ditemukan:', id)
    return
  }

  if (!confirm('Hapus data presensi ini?')) return
  try {
    await axios.delete(`/api/v1/admin/edu/journals/${targetId}`)
    toastMessage.value = 'Presensi Berhasil Dihapus'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
    await loadJournals()
  } catch (err) {
    console.error('Gagal menghapus presensi:', err)
    toastMessage.value = 'Gagal menghapus data'
    toastType.value = 'error'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  }
}

watch(filters, loadJournals, { deep: true })
onMounted(fetchData)
</script>

<template>
  <div class="teacher-attendance-page">
    <div class="p-8 space-y-8">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="flex items-center gap-4">
          <div class="h-10 w-1 bg-rose-500 rounded-full shadow-[0_0_15px_rgba(244,63,94,0.5)]"></div>
          <div>
            <h2 class="text-xs font-black uppercase tracking-[0.3em] text-sigma-muted mb-1">Manajemen Guru / SDM</h2>
            <h1 class="text-2xl font-black italic uppercase tracking-tight">Presensi <span class="text-rose-500">Mengajar Guru</span></h1>
          </div>
        </div>

        <button @click="openAddModal"
                class="px-8 py-3 bg-rose-600 text-white rounded-xl font-black uppercase tracking-widest text-[10px] hover:bg-rose-700 transition-all shadow-lg shadow-rose-500/20 flex items-center gap-3 active:scale-95">
          <Plus class="w-4 h-4" />
          Input Presensi Baru
        </button>
      </div>

      <!-- Search & Filters -->
      <div class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] p-6 shadow-sm flex flex-wrap items-center justify-between gap-8">
        <div class="flex items-center gap-6">
          <div class="space-y-1 min-w-[200px]">
            <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Pilih Tanggal</label>
            <div class="relative">
              <CalendarIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-rose-500" />
              <input v-model="filters.date" type="date" 
                     class="w-full pl-12 pr-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-rose-500/50" />
            </div>
          </div>
          
          <div class="space-y-1 min-w-[150px]">
            <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Kelas</label>
            <select v-model="filters.classroom_id" 
                    class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-rose-500/50 appearance-none">
              <option :value="null">Semua Kelas</option>
              <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
            </select>
          </div>

          <div class="space-y-1 min-w-[150px]">
            <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Mapel</label>
            <select v-model="filters.subject_id" 
                    class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-rose-500/50 appearance-none">
              <option :value="null">Semua Mapel</option>
              <option v-for="sub in subjects" :key="sub.ID" :value="sub.ID">{{ sub.name }}</option>
            </select>
          </div>

          <div class="flex items-center gap-2">
            <div class="space-y-1 min-w-[250px]">
              <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Pencarian</label>
              <div class="relative">
                <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-rose-500" />
                <input v-model="searchQuery" type="text" placeholder="Cari Guru..."
                       class="w-full pl-12 pr-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-rose-500/50" />
              </div>
            </div>
            <div class="pt-4">
              <button @click="refreshData" 
                      class="p-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-rose-500 hover:bg-rose-500 hover:text-white transition-all shadow-sm"
                      title="Refresh Data">
                <RotateCw class="w-4 h-4" :class="{ 'animate-spin': isLoading }" />
              </button>
            </div>
          </div>
        </div>

        <div class="flex items-center gap-4 px-6 py-3 bg-rose-500/5 rounded-2xl border border-rose-500/10">
          <div class="text-right">
            <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted block">Total Records</span>
            <span class="text-sm font-black text-rose-500">{{ filteredJournals.length }} Entri</span>
          </div>
        </div>
      </div>

      <!-- Table View -->
      <div class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] overflow-hidden shadow-sm">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-sigma-surface-alt/50 border-b border-sigma-border">
              <th class="px-8 py-5 text-[8px] font-black uppercase tracking-widest text-sigma-muted">Jam</th>
              <th class="px-8 py-5 text-[8px] font-black uppercase tracking-widest text-sigma-muted">Guru Pengajar</th>
              <th class="px-8 py-5 text-[8px] font-black uppercase tracking-widest text-sigma-muted">Kelas & Mapel</th>
              <th class="px-8 py-5 text-[8px] font-black uppercase tracking-widest text-sigma-muted text-center">Status</th>
              <th class="px-8 py-5 text-[8px] font-black uppercase tracking-widest text-sigma-muted">Catatan Materi</th>
              <th class="px-8 py-5 text-[8px] font-black uppercase tracking-widest text-sigma-muted text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-sigma-border/50">
            <tr v-if="isLoading" class="animate-pulse">
              <td colspan="6" class="px-8 py-20 text-center text-[10px] font-black uppercase tracking-widest text-sigma-muted">
                Memproses Data...
              </td>
            </tr>
            <tr v-else-if="filteredJournals.length > 0" v-for="j in filteredJournals" :key="j.ID"
                class="hover:bg-rose-500/[0.02] transition-all group">
              <td class="px-8 py-5">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-lg bg-rose-500/10 flex items-center justify-center text-rose-500">
                    <Clock class="w-4 h-4" />
                  </div>
                  <span class="text-xs font-black">{{ j.study_hour?.name }}</span>
                </div>
              </td>
              <td class="px-8 py-5">
                <div class="flex flex-col">
                  <span class="text-xs font-bold text-sigma-text">{{ j.teacher?.name }}</span>
                  <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted">NIP: {{ j.teacher?.nip || '-' }}</span>
                </div>
              </td>
              <td class="px-8 py-5">
                <div class="flex flex-col">
                  <span class="text-xs font-bold text-sigma-text">{{ j.classroom?.name }}</span>
                  <span class="text-[9px] font-black text-rose-500 uppercase tracking-tighter">{{ j.subject?.name }}</span>
                </div>
              </td>
              <td class="px-8 py-5 text-center">
                <span class="px-3 py-1 rounded-full text-[8px] font-black uppercase tracking-widest inline-block"
                      :class="{
                        'bg-emerald-500/10 text-emerald-500': j.status === 'Hadir',
                        'bg-amber-500/10 text-amber-500': j.status === 'Izin',
                        'bg-blue-500/10 text-blue-500': j.status === 'Sakit',
                        'bg-rose-500/10 text-rose-500': j.status === 'Alpha'
                      }">
                  {{ j.status }}
                </span>
              </td>
              <td class="px-8 py-5 max-w-xs">
                <p class="text-[11px] text-sigma-text/70 line-clamp-2 italic leading-relaxed">
                  "{{ j.remarks || 'Tidak ada catatan.' }}"
                </p>
              </td>
              <td class="px-8 py-5 text-right whitespace-nowrap">
                <div class="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-all">
                  <button @click="openEditModal(j)" class="p-2 text-rose-500 hover:bg-rose-500/10 rounded-xl transition-all">
                    <Edit class="w-4 h-4" />
                  </button>
                  <button @click="handleDelete(j.ID || j.id)" class="p-2 text-rose-500 hover:bg-rose-500/10 rounded-xl transition-all">
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
            <tr v-else>
              <td colspan="6" class="px-8 py-20 text-center opacity-30">
                <Users class="w-12 h-12 mx-auto mb-4 text-sigma-muted" />
                <p class="font-black uppercase tracking-widest text-[10px]">Data tidak ditemukan</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>


      <!-- Add Modal -->
      <Transition name="modal">
        <div v-if="showAddModal" class="fixed inset-0 z-[100] flex items-center justify-center p-6">
          <div class="absolute inset-0 bg-sigma-bg/80 backdrop-blur-xl" @click="showAddModal = false"></div>
          <div class="relative bg-sigma-surface border border-sigma-border w-full max-w-2xl rounded-[3rem] shadow-2xl overflow-hidden animate-modal-in">
            <div class="p-8 border-b border-sigma-border bg-sigma-surface-alt/50 flex items-center justify-between">
              <div class="flex items-center gap-4">
                <div class="w-10 h-10 rounded-2xl bg-rose-500 flex items-center justify-center text-white shadow-lg shadow-rose-500/20">
                  <Plus v-if="!isEdit" class="w-6 h-6" />
                  <Edit v-else class="w-6 h-6" />
                </div>
                <h2 class="text-xl font-black italic uppercase tracking-tight">
                  {{ isEdit ? 'Edit' : 'Input' }} <span class="text-rose-500">Presensi Guru</span>
                </h2>
              </div>
              <button @click="showAddModal = false" class="w-10 h-10 rounded-full hover:bg-sigma-surface-alt flex items-center justify-center transition-all">
                <XCircle class="w-6 h-6 text-sigma-muted" />
              </button>
            </div>

            <div class="p-8 space-y-6">
              <div class="grid grid-cols-2 gap-6">
                <div class="space-y-1">
                  <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Pilih Guru</label>
                  <select v-model="newJournal.teacher_id" class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-rose-500/50">
                    <option :value="null">Pilih Guru...</option>
                    <option v-for="t in teachers" :key="t.ID" :value="t.ID">{{ t.name }}</option>
                  </select>
                </div>
                <div class="space-y-1">
                  <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Pilih Kelas</label>
                  <select v-model="newJournal.classroom_id" class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-rose-500/50">
                    <option :value="null">Pilih Kelas...</option>
                    <option v-for="c in classrooms" :key="c.ID" :value="c.ID">{{ c.name }}</option>
                  </select>
                </div>
                <div class="space-y-1">
                  <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Mata Pelajaran</label>
                  <select v-model="newJournal.subject_id" class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-rose-500/50">
                    <option :value="null">Pilih Mapel...</option>
                    <option v-for="s in subjects" :key="s.ID" :value="s.ID">{{ s.name }}</option>
                  </select>
                </div>
                <div class="space-y-1">
                  <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Jam Pelajaran</label>
                  <select v-model="newJournal.study_hour_id" class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-rose-500/50">
                    <option :value="null">Pilih Jam...</option>
                    <option v-for="h in studyHours" :key="h.ID" :value="h.ID">{{ h.name }} ({{ h.start_time }} - {{ h.end_time }})</option>
                  </select>
                </div>
                <div class="space-y-1">
                  <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Status Kehadiran</label>
                  <select v-model="newJournal.status" class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-rose-500/50">
                    <option value="Hadir">Hadir</option>
                    <option value="Izin">Izin</option>
                    <option value="Sakit">Sakit</option>
                    <option value="Alpha">Alpha</option>
                  </select>
                </div>
              </div>
              <div class="space-y-1">
                <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Catatan Pembelajaran</label>
                <textarea v-model="newJournal.remarks" rows="4" 
                          class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-rose-500/50 resize-none"
                          placeholder="Tuliskan catatan mengajar atau kendala di kelas..."></textarea>
              </div>
            </div>

            <div class="p-8 bg-sigma-surface-alt/30 border-t border-sigma-border flex justify-end gap-4">
              <button @click="showAddModal = false" class="px-6 py-3 rounded-xl font-black uppercase tracking-widest text-[10px] text-sigma-muted hover:text-sigma-text transition-all">Batal</button>
              <button @click="handleSubmit" :disabled="isSaving"
                      class="px-8 py-3 bg-rose-600 text-white rounded-xl font-black uppercase tracking-widest text-[10px] hover:bg-rose-700 transition-all shadow-lg shadow-rose-500/20 disabled:opacity-50 flex items-center gap-3">
                <Save v-if="!isSaving" class="w-4 h-4" />
                <div v-else class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                {{ isEdit ? 'Perbarui Presensi' : 'Simpan Presensi' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Notification Toast -->
      <Transition name="toast">
        <div v-if="showToast" 
             class="fixed bottom-10 left-1/2 -translate-x-1/2 z-[200] flex items-center gap-4 px-8 py-4 rounded-[2rem] border shadow-2xl backdrop-blur-2xl transition-all"
             :class="toastType === 'success' ? 'bg-rose-500/90 border-rose-400/50 text-white shadow-rose-500/20' : 'bg-rose-500/90 border-rose-400/50 text-white shadow-rose-500/20'">
          <div class="w-8 h-8 rounded-full bg-white/20 flex items-center justify-center">
            <CheckCircle2 v-if="toastType === 'success'" class="w-5 h-5" />
            <XCircle v-else class="w-5 h-5" />
          </div>
          <span class="text-xs font-black uppercase tracking-widest">{{ toastMessage }}</span>
        </div>
      </Transition>
    </div>
  </div>
</template>

<style scoped>
.toast-enter-active, .toast-leave-active { transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1); }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translate(-50%, 20px) scale(0.9); }

.modal-enter-active, .modal-leave-active { transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1); }
.modal-enter-from { opacity: 0; transform: translateY(40px) scale(0.95); }
.modal-leave-to { opacity: 0; transform: translateY(40px) scale(0.95); }

@keyframes modal-in {
  from { opacity: 0; transform: translateY(40px) scale(0.95); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
.animate-modal-in { animation: modal-in 0.5s cubic-bezier(0.4, 0, 0.2, 1); }
</style>
