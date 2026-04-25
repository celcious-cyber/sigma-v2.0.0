<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { 
  User, Mail, Phone, ShieldCheck, Save, Camera, 
  Lock, KeyRound, CheckCircle2, AlertCircle 
} from 'lucide-vue-next'
import axios from 'axios'

const user = ref<any>(null)
const isLoading = ref(true)
const isSavingProfile = ref(false)
const isSavingPassword = ref(false)
const isUploadingAvatar = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

const profileForm = ref({
  name: '',
  email: '',
  phone: ''
})

const passwordForm = ref({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const toast = ref({
  show: false,
  message: '',
  type: 'success'
})

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => toast.value.show = false, 3000)
}

const fetchProfile = async () => {
  try {
    const res = await axios.get('/api/v1/auth/me')
    user.value = res.data.data
    profileForm.value = {
      name: user.value.name,
      email: user.value.email || '',
      phone: user.value.phone || ''
    }
  } catch (err) {
    showToast('Gagal memuat profil', 'error')
  } finally {
    isLoading.value = false
  }
}

const updateProfile = async () => {
  isSavingProfile.value = true
  try {
    await axios.put('/api/v1/auth/profile', profileForm.value)
    showToast('Profil berhasil diperbarui')
    fetchProfile()
  } catch (err: any) {
    showToast(err.response?.data?.error || 'Gagal memperbarui profil', 'error')
  } finally {
    isSavingProfile.value = false
  }
}

const updatePassword = async () => {
  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
    showToast('Konfirmasi password tidak cocok', 'error')
    return
  }

  isSavingPassword.value = true
  try {
    await axios.put('/api/v1/auth/password', {
      old_password: passwordForm.value.old_password,
      new_password: passwordForm.value.new_password
    })
    showToast('Password berhasil diperbarui')
    passwordForm.value = { old_password: '', new_password: '', confirm_password: '' }
  } catch (err: any) {
    showToast(err.response?.data?.error || 'Gagal memperbarui password', 'error')
  } finally {
    isSavingPassword.value = false
  }
}

const triggerFileUpload = () => {
  fileInput.value?.click()
}

const handleFileUpload = async (event: any) => {
  const file = event.target.files[0]
  if (!file) return

  const formData = new FormData()
  formData.append('avatar', file)

  isUploadingAvatar.value = true
  try {
    const res = await axios.post('/api/v1/auth/avatar', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    showToast('Foto profil berhasil diunggah')
    // Update local state immediately or refetch
    if (user.value) user.value.avatar_url = res.data.url
  } catch (err: any) {
    showToast(err.response?.data?.error || 'Gagal mengunggah foto', 'error')
  } finally {
    isUploadingAvatar.value = false
  }
}

onMounted(fetchProfile)
</script>

<template>
  <div class="p-8 max-w-6xl mx-auto space-y-8 animate-in fade-in duration-700">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-4">
      <div class="flex items-center gap-6">
        <div class="relative group">
          <input type="file" ref="fileInput" class="hidden" accept="image/*" @change="handleFileUpload" />
          <div class="w-24 h-24 rounded-3xl bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white text-3xl font-black shadow-2xl shadow-blue-500/20 transform -rotate-3 group-hover:rotate-0 transition-transform duration-500 overflow-hidden relative">
            <div v-if="isUploadingAvatar" class="absolute inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center z-10">
              <div class="w-8 h-8 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
            </div>
            <img v-if="user?.avatar_url" :src="user.avatar_url" class="w-full h-full object-cover" />
            <span v-else>{{ user?.name?.charAt(0) }}</span>
          </div>
          <button @click="triggerFileUpload" class="absolute -bottom-2 -right-2 w-10 h-10 rounded-2xl bg-white shadow-xl flex items-center justify-center text-blue-600 hover:bg-blue-600 hover:text-white transition-all border border-blue-50 z-20">
            <Camera class="w-5 h-5" />
          </button>
        </div>
        <div>
          <h2 class="text-[10px] font-black uppercase tracking-[0.4em] text-slate-400 mb-1 opacity-60">Sistem Keamanan • Profil</h2>
          <h1 class="text-4xl font-black italic uppercase tracking-tighter leading-none text-slate-800">
            Profile <span class="text-blue-600">Settings</span>
          </h1>
          <p class="text-xs font-bold text-slate-400 mt-2 uppercase tracking-widest">Kelola informasi pribadi dan keamanan akun Anda</p>
        </div>
      </div>
    </div>

    <div v-if="isLoading" class="flex flex-col items-center justify-center py-40 gap-4">
      <div class="w-12 h-12 border-4 border-blue-500/20 border-t-blue-600 rounded-full animate-spin"></div>
      <p class="text-[10px] font-black uppercase tracking-widest text-slate-400">Menyiapkan Ruang Pribadi...</p>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-8 items-start">
      <!-- Left: Personal Information -->
      <div class="lg:col-span-2 space-y-8">
        <div class="bg-white/60 backdrop-blur-xl border border-white rounded-[3rem] p-10 shadow-2xl shadow-blue-500/5 relative overflow-hidden group">
          <div class="absolute -right-20 -top-20 w-60 h-60 bg-blue-500/5 rounded-full blur-3xl group-hover:scale-110 transition-transform duration-1000"></div>
          
          <div class="flex items-center gap-4 mb-10 relative z-10">
            <div class="w-12 h-12 rounded-2xl bg-blue-500/10 flex items-center justify-center text-blue-600">
              <User class="w-6 h-6" />
            </div>
            <div>
              <h3 class="text-lg font-black uppercase italic text-slate-800">Informasi Pribadi</h3>
              <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Detail identitas utama Anda</p>
            </div>
          </div>

          <form @submit.prevent="updateProfile" class="space-y-8 relative z-10">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
              <div class="space-y-3">
                <label class="text-[10px] font-black uppercase tracking-widest text-slate-400 ml-1">Nama Lengkap</label>
                <div class="relative">
                  <User class="absolute left-5 top-1/2 -translate-y-1/2 w-4 h-4 text-blue-500" />
                  <input v-model="profileForm.name" type="text" placeholder="Masukkan nama..."
                         class="w-full pl-12 pr-6 py-4 bg-white/50 border border-slate-100 rounded-2xl text-xs font-black outline-none focus:border-blue-500/50 focus:ring-4 focus:ring-blue-500/5 transition-all" />
                </div>
              </div>
              <div class="space-y-3">
                <label class="text-[10px] font-black uppercase tracking-widest text-slate-400 ml-1">Alamat Email</label>
                <div class="relative">
                  <Mail class="absolute left-5 top-1/2 -translate-y-1/2 w-4 h-4 text-blue-500" />
                  <input v-model="profileForm.email" type="email" placeholder="nama@sekolah.id"
                         class="w-full pl-12 pr-6 py-4 bg-white/50 border border-slate-100 rounded-2xl text-xs font-black outline-none focus:border-blue-500/50 focus:ring-4 focus:ring-blue-500/5 transition-all" />
                </div>
              </div>
              <div v-if="user?.role" class="space-y-3">
                <label class="text-[10px] font-black uppercase tracking-widest text-slate-400 ml-1">Role / Jabatan</label>
                <div class="px-6 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-xs font-black text-slate-400 uppercase tracking-widest flex items-center gap-3">
                  <ShieldCheck class="w-4 h-4" />
                  {{ user.role }}
                </div>
              </div>
              <div v-if="user?.phone !== undefined" class="space-y-3">
                <label class="text-[10px] font-black uppercase tracking-widest text-slate-400 ml-1">Nomor WhatsApp</label>
                <div class="relative">
                  <Phone class="absolute left-5 top-1/2 -translate-y-1/2 w-4 h-4 text-blue-500" />
                  <input v-model="profileForm.phone" type="text" placeholder="08xxxx"
                         class="w-full pl-12 pr-6 py-4 bg-white/50 border border-slate-100 rounded-2xl text-xs font-black outline-none focus:border-blue-500/50 focus:ring-4 focus:ring-blue-500/5 transition-all" />
                </div>
              </div>
            </div>

            <div class="flex justify-end pt-4">
              <button type="submit" :disabled="isSavingProfile"
                      class="flex items-center gap-3 px-10 py-4 bg-blue-600 text-white rounded-2xl text-[10px] font-black uppercase tracking-widest hover:bg-blue-700 transition-all shadow-xl shadow-blue-500/20 active:scale-95 disabled:opacity-50">
                <div v-if="isSavingProfile" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                <Save v-else class="w-4 h-4" />
                Simpan Perubahan
              </button>
            </div>
          </form>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
           <div class="bg-blue-600 rounded-[2.5rem] p-8 text-white shadow-xl shadow-blue-600/10 relative overflow-hidden group">
              <Activity class="absolute -right-4 -bottom-4 w-24 h-24 opacity-10 group-hover:scale-125 transition-transform duration-500" />
              <h4 class="text-xs font-black uppercase tracking-widest mb-1 opacity-60">Status Keaktifan</h4>
              <p class="text-2xl font-black italic">Active <span class="text-blue-200">User</span></p>
              <div class="flex items-center gap-2 mt-4">
                 <div class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></div>
                 <span class="text-[10px] font-black uppercase tracking-widest">Sistem Online</span>
              </div>
           </div>
           <div class="bg-white border border-slate-100 rounded-[2.5rem] p-8 shadow-xl shadow-slate-200/50">
              <h4 class="text-xs font-black uppercase tracking-widest text-slate-400 mb-1">Terakhir Login</h4>
              <p class="text-xl font-black italic text-slate-800">Hari ini, <span class="text-blue-600">22:28</span></p>
              <p class="text-[10px] font-bold text-slate-400 mt-2 uppercase tracking-widest">Dari Perangkat: Windows (Chrome)</p>
           </div>
        </div>
      </div>

      <!-- Right: Security / Password Change -->
      <div class="space-y-8">
        <div class="bg-slate-900 rounded-[3rem] p-10 shadow-2xl relative overflow-hidden group">
          <div class="absolute inset-0 bg-gradient-to-br from-blue-600/20 to-transparent opacity-50"></div>
          
          <div class="flex items-center gap-4 mb-10 relative z-10">
            <div class="w-12 h-12 rounded-2xl bg-white/10 flex items-center justify-center text-blue-400">
              <Lock class="w-6 h-6" />
            </div>
            <div>
              <h3 class="text-lg font-black uppercase italic text-white">Ganti Password</h3>
              <p class="text-[9px] font-bold text-slate-500 uppercase tracking-widest">Perbarui keamanan akun Anda</p>
            </div>
          </div>

          <form @submit.prevent="updatePassword" class="space-y-6 relative z-10">
            <div class="space-y-3">
              <label class="text-[10px] font-black uppercase tracking-widest text-slate-500 ml-1">Password Sekarang</label>
              <div class="relative">
                <KeyRound class="absolute left-5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-600" />
                <input v-model="passwordForm.old_password" type="password" placeholder="••••••••"
                       class="w-full pl-12 pr-6 py-4 bg-white/5 border border-white/10 rounded-2xl text-xs font-black text-white outline-none focus:border-blue-500/50 transition-all placeholder:text-slate-700" />
              </div>
            </div>

            <div class="space-y-3">
              <label class="text-[10px] font-black uppercase tracking-widest text-slate-500 ml-1">Password Baru</label>
              <div class="relative">
                <ShieldCheck class="absolute left-5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-600" />
                <input v-model="passwordForm.new_password" type="password" placeholder="••••••••"
                       class="w-full pl-12 pr-6 py-4 bg-white/5 border border-white/10 rounded-2xl text-xs font-black text-white outline-none focus:border-blue-500/50 transition-all placeholder:text-slate-700" />
              </div>
            </div>

            <div class="space-y-3">
              <label class="text-[10px] font-black uppercase tracking-widest text-slate-500 ml-1">Konfirmasi Password Baru</label>
              <div class="relative">
                <ShieldCheck class="absolute left-5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-600" />
                <input v-model="passwordForm.confirm_password" type="password" placeholder="••••••••"
                       class="w-full pl-12 pr-6 py-4 bg-white/5 border border-white/10 rounded-2xl text-xs font-black text-white outline-none focus:border-blue-500/50 transition-all placeholder:text-slate-700" />
              </div>
            </div>

            <button type="submit" :disabled="isSavingPassword"
                    class="w-full mt-4 flex items-center justify-center gap-3 py-4 bg-white text-slate-900 rounded-2xl text-[10px] font-black uppercase tracking-widest hover:bg-blue-500 hover:text-white transition-all shadow-xl active:scale-95 disabled:opacity-50">
              <div v-if="isSavingPassword" class="w-4 h-4 border-2 border-slate-900/30 border-t-slate-900 rounded-full animate-spin"></div>
              Perbarui Password
            </button>
          </form>
        </div>

        <div class="bg-white/40 backdrop-blur-md border border-white rounded-[2.5rem] p-8 shadow-sm">
           <div class="flex items-center gap-3 mb-4">
              <AlertCircle class="w-4 h-4 text-amber-500" />
              <p class="text-[10px] font-black uppercase tracking-widest text-slate-800">Tips Keamanan</p>
           </div>
           <p class="text-[10px] font-bold text-slate-500 leading-relaxed uppercase">Gunakan kombinasi huruf besar, kecil, angka dan simbol untuk keamanan maksimal. Jangan bagikan password Anda kepada siapapun.</p>
        </div>
      </div>
    </div>

    <!-- Toast Notification -->
    <Transition name="toast">
      <div v-if="toast.show" 
           class="fixed bottom-10 left-1/2 -translate-x-1/2 z-[200] flex items-center gap-4 px-8 py-4 rounded-[2rem] border shadow-2xl backdrop-blur-2xl transition-all"
           :class="toast.type === 'success' ? 'bg-emerald-500/90 border-emerald-400/50 text-white shadow-emerald-500/20' : 'bg-rose-500/90 border-rose-400/50 text-white shadow-rose-500/20'">
        <div class="w-8 h-8 rounded-full bg-white/20 flex items-center justify-center">
          <CheckCircle2 v-if="toast.type === 'success'" class="w-5 h-5" />
          <AlertCircle v-else class="w-5 h-5" />
        </div>
        <span class="text-xs font-black uppercase tracking-widest">{{ toast.message }}</span>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.toast-enter-active, .toast-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}
.toast-enter-from, .toast-leave-to {
  opacity: 0;
  transform: translate(-50%, 20px) scale(0.9);
}

input::placeholder {
  transition: opacity 0.3s;
}
input:focus::placeholder {
  opacity: 0.3;
}
</style>
