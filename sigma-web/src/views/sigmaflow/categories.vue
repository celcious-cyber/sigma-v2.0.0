<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { 
  Building2, Plus, Search, MoreVertical, 
  Trash2, Edit3, AlertCircle, CheckCircle2, X,
  Layers, CreditCard, Clock
} from 'lucide-vue-next'
import axios from 'axios'

const categories = ref<any[]>([])
const isLoading = ref(true)
const isSubmitting = ref(false)
const showModal = ref(false)
const modalMode = ref<'create' | 'edit'>('create')
const editingId = ref<number | null>(null)

const form = ref({
  name: '',
  amount: 0,
  is_recurring: false
})

const formatRupiah = (val: number) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(val)
}

const fetchCategories = async () => {
  isLoading.value = true
  try {
    const res = await axios.get('/api/v1/admin/flow/categories')
    categories.value = Array.isArray(res.data) ? res.data : []
  } catch (err) {
    console.error('Failed to fetch categories', err)
  } finally {
    isLoading.value = false
  }
}

const openCreateModal = () => {
  modalMode.value = 'create'
  form.value = { name: '', amount: 0, is_recurring: false }
  showModal.value = true
}

const openEditModal = (cat: any) => {
  modalMode.value = 'edit'
  editingId.value = cat.ID || cat.id
  form.value = { 
    name: cat.Name || cat.name, 
    amount: cat.Amount || cat.amount, 
    is_recurring: cat.IsRecurring || cat.is_recurring 
  }
  showModal.value = true
}

const handleSubmit = async () => {
  if (!form.value.name || form.value.amount <= 0) {
    alert('Nama dan Nominal harus diisi dengan benar!')
    return
  }

  isSubmitting.value = true
  try {
    if (modalMode.value === 'create') {
      await axios.post('/api/v1/admin/flow/categories', form.value)
    } else {
      await axios.put(`/api/v1/admin/flow/categories/${editingId.value}`, form.value)
    }
    showModal.value = false
    await fetchCategories()
  } catch (err: any) {
    alert('Gagal menyimpan kategori: ' + (err.response?.data?.error || err.message))
  } finally {
    isSubmitting.value = false
  }
}

const deleteCategory = async (id: number) => {
  if (!confirm('Apakah Anda yakin ingin menghapus kategori ini?')) return
  try {
    await axios.delete(`/api/v1/admin/flow/categories/${id}`)
    await fetchCategories()
  } catch (err: any) {
    alert('Gagal menghapus kategori: ' + (err.response?.data?.error || err.message))
  }
}

onMounted(fetchCategories)
</script>

<template>
  <div class="min-h-screen bg-[#F8FAFC] font-sans pb-20">
    <main class="max-w-[1400px] mx-auto p-8 lg:p-10 flex flex-col gap-8 animate-reveal">
      
      <!-- Header -->
      <section class="flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="flex items-center gap-5">
          <div class="p-4 bg-indigo-600 rounded-[2rem] text-white shadow-2xl shadow-indigo-500/20">
            <Building2 class="w-8 h-8" />
          </div>
          <div>
            <div class="flex items-center gap-3 mb-1">
              <h1 class="text-3xl font-black italic tracking-tighter uppercase text-slate-900">Kategori <span class="text-indigo-600">Pembayaran</span></h1>
              <span class="px-3 py-1 bg-indigo-500/10 border border-indigo-500/20 rounded-full text-[9px] font-black uppercase tracking-widest text-indigo-600">Konfigurasi Penagihan</span>
            </div>
            <p class="text-xs font-bold text-slate-400 uppercase tracking-[0.2em]">Atur jenis tagihan &amp; nominal standar</p>
          </div>
        </div>
        
        <button @click="openCreateModal" class="flex items-center gap-3 px-6 py-3 bg-[#0F2942] hover:bg-indigo-600 text-white rounded-2xl shadow-xl shadow-indigo-600/20 transition-all active:scale-95">
          <Plus class="w-4 h-4" />
          <span class="text-[10px] font-black uppercase tracking-widest">Tambah Kategori</span>
        </button>
      </section>

      <!-- Stats Cards (Mini) -->
      <section class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="bg-white/40 backdrop-blur-2xl p-6 rounded-[2.5rem] border border-white/60 shadow-xl">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-2xl bg-indigo-50 flex items-center justify-center text-indigo-600">
              <Layers class="w-6 h-6" />
            </div>
            <div>
              <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Total Kategori</p>
              <h3 class="text-xl font-black text-slate-900 italic">{{ categories.length }} Jenis Tagihan</h3>
            </div>
          </div>
        </div>
        <div class="bg-white/40 backdrop-blur-2xl p-6 rounded-[2.5rem] border border-white/60 shadow-xl">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-2xl bg-emerald-50 flex items-center justify-center text-emerald-600">
              <Clock class="w-6 h-6" />
            </div>
            <div>
              <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Kategori Rutin</p>
              <h3 class="text-xl font-black text-slate-900 italic">{{ categories.filter(c => c.IsRecurring || c.is_recurring).length }} Tagihan Bulanan</h3>
            </div>
          </div>
        </div>
        <div class="bg-indigo-600 p-6 rounded-[2.5rem] shadow-2xl shadow-indigo-500/20 text-white flex items-center justify-between">
          <div>
            <p class="text-[9px] font-black text-indigo-200 uppercase tracking-widest">Info</p>
            <h3 class="text-lg font-bold italic tracking-tight leading-tight">Gunakan kategori untuk mempercepat input faktur manual.</h3>
          </div>
        </div>
      </section>

      <!-- Category Grid -->
      <section class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-if="isLoading" v-for="i in 3" :key="i" class="h-48 bg-white/40 rounded-[2.5rem] border border-white/60 animate-pulse"></div>
        
        <template v-else>
          <div v-for="cat in categories" :key="cat.ID || cat.id" 
               class="bg-white group hover:border-indigo-500/30 transition-all duration-500 rounded-[2.5rem] border border-slate-100 p-8 shadow-sm hover:shadow-2xl hover:shadow-indigo-500/10 flex flex-col justify-between">
            <div>
              <div class="flex items-start justify-between mb-6">
                <div class="w-14 h-14 rounded-2xl bg-slate-50 group-hover:bg-indigo-50 flex items-center justify-center text-slate-400 group-hover:text-indigo-600 transition-colors">
                  <CreditCard class="w-7 h-7" />
                </div>
                <div class="flex items-center gap-1">
                  <button @click="openEditModal(cat)" class="p-2 text-slate-300 hover:text-indigo-600 hover:bg-indigo-50 rounded-xl transition-all">
                    <Edit3 class="w-4 h-4" />
                  </button>
                  <button @click="deleteCategory(cat.ID || cat.id)" class="p-2 text-slate-300 hover:text-rose-600 hover:bg-rose-50 rounded-xl transition-all">
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </div>
              
              <h4 class="text-lg font-black text-slate-900 uppercase tracking-tight mb-1">{{ cat.Name || cat.name }}</h4>
              <p class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-4">Nominal Standar</p>
              
              <div class="p-4 bg-slate-50 rounded-2xl border border-slate-100 group-hover:border-indigo-100 group-hover:bg-indigo-50/30 transition-all">
                <p class="text-2xl font-black text-slate-900 italic tracking-tighter">{{ formatRupiah(cat.Amount || cat.amount) }}</p>
              </div>
            </div>

            <div class="mt-6 pt-6 border-t border-slate-50 flex items-center justify-between">
              <span :class="(cat.IsRecurring || cat.is_recurring) ? 'bg-indigo-500/10 text-indigo-600 border-indigo-500/20' : 'bg-slate-100 text-slate-400 border-slate-200'" 
                    class="px-4 py-1.5 rounded-full text-[9px] font-black uppercase tracking-widest border">
                {{ (cat.IsRecurring || cat.is_recurring) ? 'Rutin Bulanan' : 'Sekali Bayar' }}
              </span>
              <div class="flex -space-x-2">
                <div v-for="i in 3" :key="i" class="w-6 h-6 rounded-full border-2 border-white bg-slate-200"></div>
              </div>
            </div>
          </div>
        </template>

        <!-- Empty State -->
        <div v-if="!isLoading && categories.length === 0" class="col-span-full py-20 flex flex-col items-center justify-center text-slate-300">
          <AlertCircle class="w-16 h-16 mb-4 opacity-20" />
          <p class="text-xs font-black uppercase tracking-widest">Belum ada kategori pembayaran</p>
        </div>
      </section>
    </main>

    <!-- Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
        <div class="modal-container">
          <div class="p-10 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
            <div>
              <h3 class="text-2xl font-black text-slate-900 italic uppercase">
                {{ modalMode === 'create' ? 'Tambah' : 'Edit' }} <span class="text-indigo-600">Kategori</span>
              </h3>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1">Konfigurasi tagihan standar</p>
            </div>
            <button @click="showModal = false" class="p-3 bg-white hover:bg-rose-50 text-slate-400 hover:text-rose-500 rounded-2xl transition-all border border-slate-100">
              <X class="w-5 h-5" />
            </button>
          </div>

          <div class="p-10 space-y-8">
            <div class="space-y-2">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Nama Kategori</label>
              <input v-model="form.name" type="text" placeholder="Contoh: SPP Bulanan" class="w-full px-6 py-4 bg-slate-50 border border-slate-100 rounded-[1.5rem] text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all font-bold" />
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="space-y-2">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Nominal Standar (Rp)</label>
                <input v-model.number="form.amount" type="number" placeholder="500000" class="w-full px-6 py-4 bg-slate-50 border border-slate-100 rounded-[1.5rem] text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all font-black italic text-indigo-600 text-lg" />
              </div>
              <div class="space-y-2">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Tipe Pembayaran</label>
                <div class="flex bg-slate-100 p-1 rounded-2xl border border-slate-200">
                  <button @click="form.is_recurring = false" :class="!form.is_recurring ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-400'" class="flex-1 py-3 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all">Sekali</button>
                  <button @click="form.is_recurring = true" :class="form.is_recurring ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-400'" class="flex-1 py-3 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all">Rutin</button>
                </div>
              </div>
            </div>
          </div>

          <div class="p-10 bg-slate-50/50 border-t border-slate-100 flex items-center justify-end gap-4">
            <button @click="showModal = false" class="px-8 py-4 rounded-2xl text-[10px] font-black text-slate-400 uppercase tracking-widest hover:bg-slate-100 transition-all">Batal</button>
            <button @click="handleSubmit" :disabled="isSubmitting" class="px-10 py-4 rounded-2xl bg-indigo-600 text-white text-[10px] font-black uppercase tracking-widest shadow-xl shadow-indigo-600/20 hover:bg-indigo-700 transition-all active:scale-95 disabled:opacity-50 flex items-center gap-3">
              <CheckCircle2 v-if="!isSubmitting" class="w-4 h-4" />
              <div v-else class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
              {{ isSubmitting ? 'Menyimpan...' : (modalMode === 'create' ? 'Tambah Kategori' : 'Simpan Perubahan') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.animate-reveal {
  animation: reveal 0.8s cubic-bezier(0, 0, 0.2, 1) forwards;
}
@keyframes reveal {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(12px);
}
.modal-container {
  background: white;
  border-radius: 3rem;
  width: 100%;
  max-width: 36rem;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  overflow: hidden;
  animation: modal-in 0.3s ease-out;
}
@keyframes modal-in {
  from { opacity: 0; transform: scale(0.95) translateY(10px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}
</style>
