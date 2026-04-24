<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { 
  Building2, Plus, Edit, Trash2, X, Check
} from 'lucide-vue-next'
import axios from 'axios'

// Config API
const API_URL = '/api/v1/admin/base/units'

// State
const units = ref<any[]>([])
const isLoading = ref(false)
const isModalOpen = ref(false)
const isEditing = ref(false)
const currentId = ref<number | null>(null)

const form = ref({
  name: '',
  code: '',
  description: '',
  status: 'Active'
})

const fetchUnits = async () => {
  try {
    isLoading.value = true
    const response = await axios.get(API_URL)
    units.value = response.data
    
    // If no units found, offer to seed
    if (units.value.length === 0) {
      await seedUnits()
    }
  } catch (err) {
    console.error('Gagal mengambil data unit:', err)
  } finally {
    isLoading.value = false
  }
}

const seedUnits = async () => {
  try {
    await axios.post(`${API_URL}/seed`)
    const response = await axios.get(API_URL)
    units.value = response.data
  } catch (err) {
    console.error('Gagal seeding unit:', err)
  }
}

const openNewModal = () => {
  isEditing.value = false
  currentId.value = null
  form.value = {
    name: '',
    code: '',
    description: '',
    status: 'Active'
  }
  isModalOpen.value = true
}

const handleEdit = (unit: any) => {
  isEditing.value = true
  currentId.value = unit.ID
  form.value = {
    name: unit.name,
    code: unit.code,
    description: unit.description || '',
    status: unit.status || 'Active'
  }
  isModalOpen.value = true
}

const handleSubmit = async () => {
  try {
    isLoading.value = true
    if (isEditing.value && currentId.value) {
      await axios.put(`${API_URL}/${currentId.value}`, form.value)
    } else {
      await axios.post(API_URL, form.value)
    }
    await fetchUnits()
    isModalOpen.value = false
  } catch (err) {
    alert('Gagal menyimpan data unit.')
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

const handleDelete = async (id: number) => {
  if (!confirm('Apakah Anda yakin ingin menghapus unit ini?')) return
  try {
    isLoading.value = true
    await axios.delete(`${API_URL}/${id}`)
    await fetchUnits()
  } catch (err) {
    alert('Gagal menghapus data.')
  } finally {
    isLoading.value = false
  }
}

onMounted(fetchUnits)
</script>

<template>
  <div class="flex flex-col">
    <!-- Sticky Header -->
      <div class="sticky top-0 z-40 bg-sigma-app/80 backdrop-blur-xl border-b border-sigma-border px-8 py-6 space-y-6 shadow-sm">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-purple-500/10 rounded-2xl flex items-center justify-center text-purple-400 border border-purple-500/20">
              <Building2 class="w-6 h-6" />
            </div>
            <div>
              <h2 class="text-3xl font-black tracking-tight italic">UNIT <span class="text-sigma-emerald not-italic uppercase text-2xl">LEMBAGA</span></h2>
              <p class="text-[10px] text-sigma-muted font-bold uppercase tracking-widest mt-1">Struktur Organisasi Al-Hikmah</p>
            </div>
          </div>

          <div class="flex items-center gap-4">
            <button @click="openNewModal" class="flex items-center gap-2 px-8 py-4 bg-sigma-emerald text-white rounded-2xl font-bold uppercase tracking-widest text-xs hover:bg-emerald-600 transition-all shadow-xl shadow-sigma-emerald/20">
              <Plus class="w-4 h-4" /> Tambah Unit
            </button>
          </div>
        </div>
      </div>

      <!-- Main Content Area -->
      <div class="p-8 space-y-10">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <div v-for="unit in units" :key="unit.ID" class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] p-8 space-y-6 hover:bg-sigma-surface-alt transition-all group shadow-sm relative overflow-hidden">
             <!-- Background Pattern -->
             <div class="absolute -right-4 -bottom-4 opacity-[0.03] group-hover:opacity-[0.08] transition-all duration-700">
                <Building2 class="w-40 h-40 rotate-12" />
             </div>

             <div class="flex justify-between items-start relative z-10">
                <div class="w-14 h-14 bg-sigma-emerald/10 rounded-2xl flex items-center justify-center text-sigma-emerald border border-sigma-emerald/20 group-hover:bg-sigma-emerald group-hover:text-white transition-all duration-500">
                  <span class="text-xl font-black italic">{{ unit.name.charAt(0) }}</span>
                </div>
                <div class="flex gap-2">
                  <button @click="handleEdit(unit)" class="p-2.5 bg-sigma-surface border border-sigma-border text-sigma-muted hover:text-sigma-emerald rounded-xl transition-all"><Edit class="w-4 h-4" /></button>
                  <button @click="handleDelete(unit.ID)" class="p-2.5 bg-sigma-surface border border-sigma-border text-sigma-muted hover:text-red-500 rounded-xl transition-all"><Trash2 class="w-4 h-4" /></button>
                </div>
             </div>

             <div class="space-y-2 relative z-10">
                <div class="flex items-center gap-3">
                  <h3 class="text-2xl font-black text-sigma-text group-hover:text-sigma-emerald transition-colors">{{ unit.name }}</h3>
                  <span class="px-2 py-0.5 rounded-md bg-sigma-surface-alt border border-sigma-border text-[9px] font-black text-sigma-muted uppercase tracking-widest">{{ unit.code }}</span>
                </div>
                <p class="text-sm text-sigma-muted font-medium line-clamp-2 leading-relaxed">
                  {{ unit.description || 'Tidak ada deskripsi untuk unit ini.' }}
                </p>
             </div>

             <div class="pt-4 border-t border-sigma-border/50 flex items-center justify-between relative z-10">
                <div class="flex items-center gap-2">
                   <div class="w-2 h-2 rounded-full" :class="unit.status === 'Active' ? 'bg-emerald-500' : 'bg-red-500'"></div>
                   <span class="text-[10px] font-black uppercase tracking-widest text-sigma-muted">{{ unit.status }}</span>
                </div>
                <div class="flex items-center gap-1.5 px-3 py-1 rounded-full bg-sigma-surface-alt border border-sigma-border text-[9px] font-black text-sigma-emerald uppercase tracking-tighter">
                   <Check class="w-3 h-3" /> Managed
                </div>
             </div>
          </div>

          <!-- Empty State -->
          <div v-if="units.length === 0" class="md:col-span-3 py-20 text-center space-y-4">
             <div class="w-20 h-20 bg-sigma-surface-alt rounded-full flex items-center justify-center mx-auto border border-sigma-border">
                <Building2 class="w-10 h-10 text-sigma-muted/20" />
             </div>
             <p class="text-sigma-muted font-bold uppercase tracking-widest text-sm">Belum ada data unit lembaga.</p>
             <button @click="seedUnits" class="text-sigma-emerald font-black text-[10px] uppercase tracking-widest hover:underline">Jalankan Seeder Otomatis</button>
          </div>
        </div>
      </div>

      <!-- Modal -->
      <Transition name="fade">
        <div v-if="isModalOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
          <div class="absolute inset-0 bg-sigma-app/80 backdrop-blur-md" @click="isModalOpen = false"></div>
          
          <div class="bg-sigma-surface w-full max-w-2xl rounded-[3rem] border border-sigma-border relative z-10 overflow-hidden shadow-2xl animate-in zoom-in-95 duration-500">
            <div class="p-10 pb-6 border-b border-sigma-border flex justify-between items-start bg-gradient-to-br from-sigma-surface-alt to-transparent">
              <div>
                <h3 class="text-3xl font-black text-sigma-text uppercase italic tracking-tighter">
                  {{ isEditing ? 'Edit' : 'Tambah' }} <span class="text-sigma-emerald not-italic">Unit</span>
                </h3>
                <p class="text-sigma-muted text-sm mt-1 uppercase font-bold tracking-widest">Sigmabase Module / Institutions</p>
              </div>
              <button @click="isModalOpen = false" class="p-3 bg-sigma-surface-alt hover:bg-sigma-emerald/10 rounded-full transition-all text-sigma-muted hover:text-sigma-emerald">
                <X class="w-6 h-6" />
              </button>
            </div>

            <div class="p-10 space-y-6">
               <div class="grid grid-cols-2 gap-6">
                 <div class="space-y-2">
                    <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Nama Unit</label>
                    <input v-model="form.name" type="text" placeholder="Ex: SMP, MA, KMI" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all text-sm font-bold text-sigma-text" />
                 </div>
                 <div class="space-y-2">
                    <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Kode Unit</label>
                    <input v-model="form.code" type="text" placeholder="Ex: SMP-AH" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all text-sm font-mono font-bold text-sigma-text" />
                 </div>
               </div>

               <div class="space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Deskripsi</label>
                  <textarea v-model="form.description" rows="3" placeholder="Keterangan mengenai unit lembaga ini..." class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all resize-none text-sm font-bold text-sigma-text"></textarea>
               </div>

               <div class="space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Status Operasional</label>
                  <div class="flex gap-4">
                    <button v-for="s in ['Active', 'Inactive']" :key="s" @click="form.status = s" 
                            class="flex-1 py-4 rounded-xl border font-bold text-xs uppercase tracking-widest transition-all"
                            :class="form.status === s ? 'bg-sigma-emerald text-white border-sigma-emerald shadow-lg shadow-sigma-emerald/20' : 'bg-sigma-surface-alt text-sigma-muted border-sigma-border hover:border-sigma-emerald/30'">
                      {{ s }}
                    </button>
                  </div>
               </div>
            </div>

            <div class="p-10 pt-6 border-t border-sigma-border bg-gradient-to-br from-transparent to-sigma-surface-alt flex gap-4">
              <button @click="isModalOpen = false" :disabled="isLoading" class="flex-1 py-5 rounded-2xl font-black uppercase tracking-[0.2em] text-xs bg-sigma-surface-alt hover:bg-sigma-emerald/10 transition-all text-sigma-muted hover:text-sigma-emerald disabled:opacity-50">
                Cancel
              </button>
              <button @click="handleSubmit" :disabled="isLoading" class="flex-1 [flex-grow:2] py-5 rounded-2xl font-black uppercase tracking-[0.2em] text-xs bg-sigma-emerald hover:bg-emerald-600 transition-all shadow-xl shadow-sigma-emerald/20 disabled:opacity-50 text-white">
                {{ isLoading ? 'Syncing...' : 'Simpan Data Unit' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(255, 255, 255, 0.05); border-radius: 10px; }
.fade-enter-active, .fade-leave-active { transition: opacity 0.5s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
