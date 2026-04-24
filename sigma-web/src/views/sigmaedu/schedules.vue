<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  Calendar, Clock, BookOpen, User, 
  Plus, Pencil, Trash2, ArrowLeft,
  Check, X, Building2, ChevronRight
} from 'lucide-vue-next'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()

// State
const schedules = ref<any[]>([])
const classrooms = ref<any[]>([])
const studyHours = ref<any[]>([])
const subjects = ref<any[]>([])
const teachers = ref<any[]>([])
const isLoading = ref(true)
const selectedClassId = ref<number | null>(null)
const selectedDay = ref(1) // 1=Ahad, 2=Senin, ...

const days = [
  { id: 1, name: 'Ahad' },
  { id: 2, name: 'Senin' },
  { id: 3, name: 'Selasa' },
  { id: 4, name: 'Rabu' },
  { id: 5, name: 'Kamis' },
  { id: 6, name: 'Jumat' },
  { id: 7, name: 'Sabtu' }
]

const isModalOpen = ref(false)
const isSubmitting = ref(false)
const editingId = ref<number | null>(null)

const form = ref({
  classroom_id: null as number | null,
  subject_id: null as number | null,
  teacher_id: null as number | null,
  day_of_week: 1,
  start_time: '',
  end_time: ''
})

const fetchData = async () => {
  isLoading.value = true
  try {
    const [scRes, clRes, hrRes, sbRes, tcRes] = await Promise.all([
      axios.get('/api/v1/admin/edu/schedules'),
      axios.get('/api/v1/admin/edu/classrooms'),
      axios.get('/api/v1/admin/edu/hours'),
      axios.get('/api/v1/admin/edu/subjects'),
      axios.get('/api/v1/admin/base/teachers')
    ])
    schedules.value = scRes.data
    classrooms.value = clRes.data
    studyHours.value = hrRes.data
    subjects.value = sbRes.data
    teachers.value = tcRes.data
    
    if (classrooms.value.length > 0 && !selectedClassId.value) {
      selectedClassId.value = classrooms.value[0].ID
    }
  } catch (err) {
    console.error('Gagal mengambil data jadwal:', err)
  } finally {
    isLoading.value = false
  }
}

const filteredSchedules = computed(() => {
  if (!selectedClassId.value) return []
  return schedules.value.filter(s => 
    s.classroom_id === selectedClassId.value && 
    s.day_of_week === selectedDay.value
  ).sort((a, b) => a.start_time.localeCompare(b.start_time))
})

const openModal = (schedule: any = null) => {
  if (schedule) {
    editingId.value = schedule.ID
    form.value = {
      classroom_id: schedule.classroom_id,
      subject_id: schedule.subject_id,
      teacher_id: schedule.teacher_id,
      day_of_week: schedule.day_of_week,
      start_time: schedule.start_time,
      end_time: schedule.end_time
    }
  } else {
    editingId.value = null
    form.value = {
      classroom_id: selectedClassId.value,
      subject_id: null,
      teacher_id: null,
      day_of_week: selectedDay.value,
      start_time: '',
      end_time: ''
    }
  }
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
}

const handleSubmit = async () => {
  console.log('Data yang akan disimpan:', form.value)
  
  if (!form.value.classroom_id || Number(form.value.classroom_id) === 0) {
    alert('Mohon pilih Kelas terlebih dahulu.')
    return
  }
  if (!form.value.subject_id || Number(form.value.subject_id) === 0) {
    alert('Mohon pilih Mata Pelajaran.')
    return
  }
  if (!form.value.teacher_id || Number(form.value.teacher_id) === 0) {
    alert('Mohon pilih Guru Pengajar.')
    return
  }
  if (!form.value.start_time || !form.value.end_time) {
    alert('Mohon tentukan waktu mulai dan selesai.')
    return
  }

  isSubmitting.value = true
  try {
    const payload = {
      ...form.value,
      classroom_id: Number(form.value.classroom_id),
      subject_id: Number(form.value.subject_id),
      teacher_id: Number(form.value.teacher_id),
      day_of_week: Number(form.value.day_of_week)
    }
    
    if (editingId.value) {
      await axios.put(`/api/v1/admin/edu/schedules/${editingId.value}`, payload)
    } else {
      await axios.post('/api/v1/admin/edu/schedules', payload)
    }
    await fetchData()
    closeModal()
  } catch (err) {
    console.error('Save error:', err)
    alert('Gagal menyimpan jadwal. Pastikan semua relasi data (Guru/Mapel) valid.')
  } finally {
    isSubmitting.value = false
  }
}

const deleteSchedule = async (id: number) => {
  if (!confirm('Yakin ingin menghapus jadwal ini?')) return
  try {
    await axios.delete(`/api/v1/admin/edu/schedules/${id}`)
    await fetchData()
  } catch (err) {
    alert('Gagal menghapus jadwal.')
  }
}

const setTimeFromHour = (hour: any) => {
  form.value.start_time = hour.start_time
  form.value.end_time = hour.end_time
}

onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="p-8 space-y-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <div class="h-10 w-1 bg-blue-500 rounded-full"></div>
        <div>
          <h2 class="text-xs font-black uppercase tracking-[0.3em] text-sigma-muted mb-1">Ruang Kelas</h2>
          <h1 class="text-2xl font-black italic uppercase tracking-tight">Jadwal <span class="text-blue-500">Pelajaran</span></h1>
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

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
      <!-- Sidebar Class Selection -->
      <div class="lg:col-span-3 space-y-6">
        <div class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] p-6 shadow-sm">
          <h3 class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted mb-6 px-2 flex items-center gap-2">
            <Building2 class="w-3 h-3" /> Pilih Kelas
          </h3>
          <div class="space-y-2 max-h-[60vh] overflow-y-auto custom-scrollbar pr-2">
            <button v-for="cls in classrooms" :key="cls.ID"
                    @click="selectedClassId = cls.ID"
                    class="w-full flex items-center justify-between p-4 rounded-2xl transition-all group"
                    :class="selectedClassId === cls.ID ? 'bg-blue-600 text-white shadow-lg shadow-blue-500/20' : 'hover:bg-sigma-surface-alt text-sigma-muted hover:text-sigma-text border border-transparent hover:border-sigma-border'">
              <div class="flex flex-col items-start">
                <span class="text-[8px] font-black uppercase tracking-widest opacity-60">Level {{ cls.level }}</span>
                <span class="text-sm font-black uppercase tracking-tighter italic">{{ cls.name }}</span>
              </div>
              <ChevronRight class="w-4 h-4 opacity-40 group-hover:opacity-100 transition-opacity" />
            </button>
          </div>
        </div>

        <div class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] p-8 space-y-4">
          <div class="w-10 h-10 rounded-2xl bg-blue-500/10 flex items-center justify-center text-blue-500">
            <Clock class="w-5 h-5" />
          </div>
          <h3 class="text-sm font-black uppercase tracking-widest">Jam Pelajaran</h3>
          <p class="text-[10px] text-sigma-muted font-medium leading-relaxed uppercase tracking-widest">
            Gunakan master jam pelajaran untuk mempercepat input jadwal.
          </p>
          <div class="space-y-2">
            <div v-for="hr in studyHours" :key="hr.ID" 
                 class="flex items-center justify-between p-3 rounded-xl bg-sigma-surface-alt border border-sigma-border text-[10px] font-black uppercase tracking-widest text-sigma-muted">
              <span>{{ hr.name }}</span>
              <span class="text-sigma-text">{{ hr.start_time }} - {{ hr.end_time }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Main Content -->
      <div class="lg:col-span-9 space-y-6">
        <!-- Day Selector -->
        <div class="bg-sigma-surface border border-sigma-border p-3 rounded-[2rem] flex items-center gap-2 overflow-x-auto custom-scrollbar shadow-sm">
          <button v-for="day in days" :key="day.id"
                  @click="selectedDay = day.id"
                  class="flex-1 min-w-[100px] py-3 rounded-2xl transition-all text-[10px] font-black uppercase tracking-widest border"
                  :class="selectedDay === day.id ? 'bg-blue-600 text-white border-blue-600 shadow-lg shadow-blue-500/20' : 'bg-sigma-surface text-sigma-muted border-transparent hover:border-sigma-border hover:text-sigma-text'">
            {{ day.name }}
          </button>
        </div>

        <!-- Toolbar -->
        <div class="flex justify-between items-center bg-sigma-surface p-4 rounded-[2rem] border border-sigma-border shadow-sm">
          <div class="flex items-center gap-3 px-4">
            <Calendar class="w-5 h-5 text-blue-500" />
            <span class="text-sm font-black uppercase tracking-widest italic text-sigma-text">
              {{ days.find(d => d.id === selectedDay)?.name }} - 
              <span class="text-blue-500">{{ classrooms.find(c => c.ID === selectedClassId)?.name }}</span>
            </span>
          </div>
          <button @click="openModal()" 
                  class="px-8 py-3.5 bg-blue-600 text-white rounded-2xl font-black uppercase tracking-widest text-[10px] hover:bg-blue-700 transition-all shadow-xl shadow-blue-500/20 flex items-center justify-center gap-3">
            <Plus class="w-4 h-4" /> Tambah Jadwal
          </button>
        </div>

        <!-- Schedule Table -->
        <div class="bg-sigma-surface border border-sigma-border rounded-[3rem] overflow-hidden shadow-sm">
          <table class="w-full border-collapse">
            <thead>
              <tr class="bg-sigma-surface-alt/50 border-b border-sigma-border">
                <th class="p-6 text-left text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Waktu</th>
                <th class="p-6 text-left text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Mata Pelajaran</th>
                <th class="p-6 text-left text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Pengajar</th>
                <th class="p-6 text-right text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="sch in filteredSchedules" :key="sch.ID" 
                  class="border-b border-sigma-border last:border-0 group hover:bg-sigma-surface-alt/20 transition-all">
                <td class="p-6">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-lg bg-blue-500/10 flex items-center justify-center text-blue-500">
                      <Clock class="w-3.5 h-3.5" />
                    </div>
                    <span class="text-xs font-black tracking-tight text-sigma-text">{{ sch.start_time }} - {{ sch.end_time }}</span>
                  </div>
                </td>
                <td class="p-6">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-lg bg-sigma-surface-alt border border-sigma-border flex items-center justify-center text-sigma-muted">
                      <BookOpen class="w-3.5 h-3.5" />
                    </div>
                    <div>
                      <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted block">{{ sch.subject?.code }}</span>
                      <span class="text-sm font-black uppercase tracking-tighter italic text-sigma-text">{{ sch.subject?.name }}</span>
                    </div>
                  </div>
                </td>
                <td class="p-6">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-lg bg-sigma-surface-alt border border-sigma-border flex items-center justify-center text-sigma-muted">
                      <User class="w-3.5 h-3.5" />
                    </div>
                    <span class="text-xs font-bold text-sigma-text">{{ sch.teacher?.name || 'Unknown' }}</span>
                  </div>
                </td>
                <td class="p-6">
                  <div class="flex items-center justify-end gap-2">
                    <button @click="openModal(sch)" class="p-2.5 rounded-xl bg-sigma-surface-alt border border-sigma-border text-sigma-muted hover:text-blue-500 hover:border-blue-500/30 transition-all">
                      <Pencil class="w-4 h-4" />
                    </button>
                    <button @click="deleteSchedule(sch.ID)" class="p-2.5 rounded-xl bg-sigma-surface-alt border border-sigma-border text-sigma-muted hover:text-red-500 hover:border-red-500/30 transition-all">
                      <Trash2 class="w-4 h-4" />
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="filteredSchedules.length === 0">
                <td colspan="4" class="p-20 text-center opacity-30">
                  <div class="flex flex-col items-center gap-4">
                    <Calendar class="w-12 h-12" />
                    <p class="font-black uppercase tracking-widest text-xs">Belum ada jadwal hari ini</p>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Modal Add/Edit -->
    <div v-if="isModalOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-6 backdrop-blur-xl bg-black/40">
      <div class="bg-sigma-surface w-full max-w-xl rounded-[3rem] border border-sigma-border shadow-2xl overflow-hidden animate-in zoom-in duration-300">
        <div class="p-8 bg-gradient-to-r from-blue-600 to-blue-900 flex justify-between items-center border-b border-white/10">
          <div>
            <h3 class="text-white font-black uppercase tracking-widest text-xs">{{ editingId ? 'Update' : 'Tambah' }} Jadwal</h3>
            <p class="text-blue-100/60 text-[10px] mt-1 font-medium">Atur Mata Pelajaran, Guru, dan Waktu</p>
          </div>
          <button @click="closeModal" class="p-2 rounded-xl bg-white/10 text-white hover:bg-white/20 transition-all">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="p-10 space-y-6 max-h-[70vh] overflow-y-auto custom-scrollbar">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Hari</label>
              <select v-model="form.day_of_week" 
                      class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight appearance-none">
                <option v-for="day in days" :key="day.id" :value="day.id">{{ day.name }}</option>
              </select>
            </div>
            <div class="space-y-2">
              <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Kelas</label>
              <select v-model="form.classroom_id" 
                      class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight appearance-none">
                <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
              </select>
            </div>
          </div>

          <div class="space-y-2">
            <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Mata Pelajaran</label>
            <select v-model="form.subject_id" 
                    class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight appearance-none">
              <option :value="null">-- Pilih Mata Pelajaran --</option>
              <option v-for="sub in subjects" :key="sub.ID" :value="sub.ID">{{ sub.name }} ({{ sub.code }})</option>
            </select>
          </div>

          <div class="space-y-2">
            <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Guru Pengajar</label>
            <select v-model="form.teacher_id" 
                    class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight appearance-none">
              <option :value="null">-- Pilih Guru --</option>
              <option v-for="t in teachers" :key="t.ID" :value="t.ID">{{ t.name }}</option>
            </select>
          </div>

          <div class="space-y-4">
            <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Pilih Dari Jam Master (Opsional)</label>
            <div class="flex flex-wrap gap-2">
              <button v-for="hr in studyHours" :key="hr.ID"
                      type="button"
                      @click="setTimeFromHour(hr)"
                      class="px-4 py-2 rounded-xl bg-sigma-surface-alt border border-sigma-border text-[8px] font-black uppercase tracking-widest hover:border-blue-500 hover:text-blue-500 transition-all">
                {{ hr.name }}
              </button>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Waktu Mulai</label>
              <input v-model="form.start_time" type="time" 
                     class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight" />
            </div>
            <div class="space-y-2">
              <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Waktu Selesai</label>
              <input v-model="form.end_time" type="time" 
                     class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight" />
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
              {{ editingId ? 'Simpan Perubahan' : 'Tambah Jadwal' }}
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
