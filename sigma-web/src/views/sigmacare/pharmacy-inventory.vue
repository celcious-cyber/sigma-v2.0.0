<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  Pill, Package, AlertTriangle, 
  Plus, Search, Filter, History,
  ShieldCheck, Clock, MoreVertical, X
} from 'lucide-vue-next'
import axios from 'axios'

const medicines = ref<any[]>([])
const loading = ref(true)
const showAddModal = ref(false)
const showMutationModal = ref(false)

const newMedicine = ref({
  code: '',
  name: '',
  category: 'Obat Bebas',
  unit: 'Tablet',
  stock: 0,
  description: ''
})

const newMutation = ref({
  medicine_id: 0,
  type: 'IN',
  quantity: 0,
  notes: ''
})

const fetchMedicines = async () => {
  try {
    const response = await axios.get('/admin/care/medicines')
    medicines.value = response.data
  } catch (error) {
    console.error('Failed to fetch medicines:', error)
  } finally {
    loading.value = false
  }
}

const addMedicine = async () => {
  try {
    await axios.post('/admin/care/medicines', newMedicine.value)
    showAddModal.value = false
    alert('Obat Berhasil Ditambahkan')
    fetchMedicines()
  } catch (error) {
    alert('Gagal menambah obat')
  }
}

const submitMutation = async () => {
  try {
    await axios.post('/admin/care/medicines/mutation', newMutation.value)
    showMutationModal.value = false
    alert('Stok Berhasil Diperbarui')
    fetchMedicines()
  } catch (error) {
    alert('Gagal memperbarui stok')
  }
}

onMounted(() => {
  fetchMedicines()
})

const getStockStatus = (stock: number) => {
  if (stock <= 0) return { label: 'Habis', color: 'bg-rose-50 text-rose-600 border-rose-100' }
  if (stock < 10) return { label: 'Kritis', color: 'bg-amber-50 text-amber-600 border-amber-100' }
  return { label: 'Tersedia', color: 'bg-emerald-50 text-emerald-600 border-emerald-100' }
}

const lowStockCount = computed(() => medicines.value.filter(m => m.stock < 10).length)
const goodConditionPercent = computed(() => {
  if (medicines.value.length === 0) return 0
  return Math.round(((medicines.value.length - lowStockCount.value) / medicines.value.length) * 100)
})
</script>

<template>
  <div class="p-8 space-y-8 min-h-screen bg-slate-50/50">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
      <div>
        <h1 class="text-3xl font-black tracking-tight text-slate-800 italic uppercase">Inventaris <span class="text-rose-600">Farmakologi</span></h1>
        <p class="text-slate-500 font-bold mt-1 tracking-wide uppercase text-xs flex items-center gap-2">
          <Pill class="w-4 h-4" /> Manajemen Perbekalan Medis & Obat-obatan
        </p>
      </div>
      
      <div class="flex items-center gap-3">
        <button 
          @click="showMutationModal = true"
          class="px-6 py-3.5 bg-white border border-slate-200 text-slate-600 rounded-2xl font-black text-xs uppercase tracking-widest hover:bg-slate-50 transition-all flex items-center gap-3 shadow-sm"
        >
          <History class="w-4 h-4" /> Mutasi Stok
        </button>
        <button 
          @click="showAddModal = true"
          class="px-8 py-4 bg-rose-600 text-white rounded-2xl font-black shadow-lg shadow-rose-500/20 hover:bg-rose-700 active:scale-95 transition-all flex items-center gap-3"
        >
          <Plus class="w-5 h-5" /> <span class="text-xs uppercase tracking-widest">Tambah Obat Baru</span>
        </button>
      </div>
    </div>

    <!-- Quick Stats -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <div class="bg-white p-6 rounded-3xl border border-slate-100 shadow-sm flex items-center gap-5">
        <div class="w-12 h-12 bg-rose-50 text-rose-600 rounded-2xl flex items-center justify-center">
          <Package class="w-6 h-6" />
        </div>
        <div>
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Total SKU</p>
          <p class="text-xl font-black text-slate-800 tracking-tighter">{{ medicines.length }} Item</p>
        </div>
      </div>
      <div class="bg-white p-6 rounded-3xl border border-slate-100 shadow-sm flex items-center gap-5">
        <div class="w-12 h-12 bg-amber-50 text-amber-600 rounded-2xl flex items-center justify-center">
          <AlertTriangle class="w-6 h-6" />
        </div>
        <div>
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Stok Kritis</p>
          <p class="text-xl font-black text-slate-800 tracking-tighter">{{ lowStockCount }} Item</p>
        </div>
      </div>
      <div class="bg-white p-6 rounded-3xl border border-slate-100 shadow-sm flex items-center gap-5">
        <div class="w-12 h-12 bg-emerald-50 text-emerald-600 rounded-2xl flex items-center justify-center">
          <ShieldCheck class="w-6 h-6" />
        </div>
        <div>
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Kondisi Baik</p>
          <p class="text-xl font-black text-slate-800 tracking-tighter">{{ goodConditionPercent }}% SKU</p>
        </div>
      </div>
      <div class="bg-white p-6 rounded-3xl border border-slate-100 shadow-sm flex items-center gap-5">
        <div class="w-12 h-12 bg-slate-50 text-slate-400 rounded-2xl flex items-center justify-center">
          <Clock class="w-6 h-6" />
        </div>
        <div>
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Expired Dekat</p>
          <p class="text-xl font-black text-slate-800 tracking-tighter">0 Item</p>
        </div>
      </div>
    </div>

    <!-- Medicine Table -->
    <div class="bg-white rounded-[2.5rem] border border-slate-100 shadow-sm overflow-hidden">
      <div class="p-8 border-b border-slate-50 flex items-center justify-between bg-slate-50/30">
        <div class="flex items-center gap-6">
          <h2 class="text-lg font-black text-slate-800 uppercase italic">Daftar Persediaan Obat</h2>
          <div class="relative">
            <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
            <input type="text" placeholder="Cari Kode atau Nama Obat..." class="pl-11 pr-6 py-2.5 bg-white border border-slate-200 rounded-xl text-xs font-bold w-64 focus:border-rose-500 outline-none" />
          </div>
        </div>
        <div class="flex items-center gap-2">
          <button class="p-2.5 bg-white border border-slate-200 text-slate-400 rounded-xl hover:text-rose-600 transition-colors">
            <Filter class="w-4 h-4" />
          </button>
        </div>
      </div>
      
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead>
            <tr class="bg-slate-50/50">
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest">Kode / Nama Obat</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest">Kategori</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest text-center">Stok Saat Ini</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest">Satuan</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest">Status Stok</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-50">
            <tr v-for="m in medicines" :key="m.ID" class="hover:bg-rose-50/20 transition-colors group">
              <td class="px-8 py-6">
                <div class="flex items-center gap-4">
                  <div class="w-10 h-10 bg-slate-50 text-slate-400 rounded-xl flex items-center justify-center group-hover:bg-rose-600 group-hover:text-white transition-all">
                    <Pill class="w-5 h-5" />
                  </div>
                  <div>
                    <p class="font-black text-slate-800 uppercase tracking-tight">{{ m.name }}</p>
                    <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">{{ m.code }}</p>
                  </div>
                </div>
              </td>
              <td class="px-8 py-6">
                <span class="px-3 py-1 bg-slate-100 text-slate-500 rounded-lg text-[10px] font-black uppercase">{{ m.category }}</span>
              </td>
              <td class="px-8 py-6 text-center">
                <span class="text-lg font-black text-slate-800 tracking-tighter">{{ m.stock }}</span>
              </td>
              <td class="px-8 py-6">
                <span class="text-xs font-bold text-slate-500 uppercase tracking-widest italic">{{ m.unit }}</span>
              </td>
              <td class="px-8 py-6">
                <span 
                  class="px-4 py-1.5 rounded-full text-[10px] font-black uppercase tracking-widest border"
                  :class="getStockStatus(m.stock).color"
                >
                  {{ getStockStatus(m.stock).label }}
                </span>
              </td>
              <td class="px-8 py-6 text-right">
                <button class="p-2 text-slate-300 hover:text-rose-600 transition-colors">
                  <MoreVertical class="w-5 h-5" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Add Medicine -->
    <div v-if="showAddModal" class="fixed inset-0 z-[100] flex items-center justify-center p-6 bg-slate-900/60 backdrop-blur-sm">
      <div class="bg-white w-full max-w-lg rounded-[2.5rem] shadow-2xl overflow-hidden animate-in zoom-in-95 duration-300">
        <div class="p-8 border-b border-slate-100 bg-rose-600 text-white font-black italic uppercase tracking-tight">
          Tambah Item Obat
        </div>
        <div class="p-8 space-y-6">
          <div class="space-y-2">
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Kode Obat (SKU)</label>
            <input v-model="newMedicine.code" type="text" placeholder="Misal: PCT-500" class="w-full px-5 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold" />
          </div>
          <div class="space-y-2">
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Nama Obat</label>
            <input v-model="newMedicine.name" type="text" placeholder="Nama lengkap obat..." class="w-full px-5 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold" />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Kategori</label>
              <select v-model="newMedicine.category" class="w-full px-5 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold">
                <option>Obat Bebas</option>
                <option>Obat Keras</option>
                <option>Alat Kesehatan</option>
              </select>
            </div>
            <div class="space-y-2">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Satuan</label>
              <input v-model="newMedicine.unit" type="text" placeholder="Tablet/Botol" class="w-full px-5 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold" />
            </div>
          </div>
        </div>
        <div class="p-8 bg-slate-50 border-t border-slate-100 flex gap-4">
          <button @click="showAddModal = false" class="flex-1 py-4 bg-white border border-slate-200 rounded-2xl font-black text-xs uppercase">Batal</button>
          <button @click="addMedicine" class="flex-1 py-4 bg-rose-600 text-white rounded-2xl font-black text-xs uppercase shadow-lg shadow-rose-500/20">Simpan Item</button>
        </div>
      </div>
    </div>

    <!-- Modal Mutation -->
    <div v-if="showMutationModal" class="fixed inset-0 z-[100] flex items-center justify-center p-6 bg-slate-900/60 backdrop-blur-sm">
      <div class="bg-white w-full max-w-lg rounded-[2.5rem] shadow-2xl overflow-hidden animate-in zoom-in-95 duration-300">
        <div class="p-8 border-b border-slate-100 bg-slate-800 text-white font-black italic uppercase tracking-tight flex items-center justify-between">
          Mutasi Stok Obat
          <button @click="showMutationModal = false" class="p-2 hover:bg-white/10 rounded-xl transition-colors">
            <X class="w-6 h-6" />
          </button>
        </div>
        <div class="p-8 space-y-6">
          <div class="space-y-2">
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Pilih Obat</label>
            <select v-model="newMutation.medicine_id" class="w-full px-5 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold">
              <option :value="0">Pilih...</option>
              <option v-for="m in medicines" :key="m.ID" :value="m.ID">{{ m.name }} ({{ m.code }})</option>
            </select>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Jenis Mutasi</label>
              <select v-model="newMutation.type" class="w-full px-5 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold">
                <option value="IN">Stok Masuk (+)</option>
                <option value="OUT">Stok Keluar (-)</option>
              </select>
            </div>
            <div class="space-y-2">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Jumlah</label>
              <input v-model.number="newMutation.quantity" type="number" class="w-full px-5 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold" />
            </div>
          </div>
          <div class="space-y-2">
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Catatan</label>
            <textarea v-model="newMutation.notes" rows="2" placeholder="Alasan mutasi..." class="w-full px-5 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold resize-none"></textarea>
          </div>
        </div>
        <div class="p-8 bg-slate-50 border-t border-slate-100 flex gap-4">
          <button @click="showMutationModal = false" class="flex-1 py-4 bg-white border border-slate-200 rounded-2xl font-black text-xs uppercase">Batal</button>
          <button @click="submitMutation" class="flex-1 py-4 bg-slate-800 text-white rounded-2xl font-black text-xs uppercase shadow-lg shadow-slate-800/20">Simpan Mutasi</button>
        </div>
      </div>
    </div>
  </div>
</template>
