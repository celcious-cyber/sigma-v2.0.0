import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import axios from 'axios'

// Global Axios Configuration - Auto inject JWT Token
axios.interceptors.request.use((config) => {
  const token = localStorage.getItem('sigma_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
}, (error) => {
  return Promise.reject(error)
})

// Global Response Interceptor - Handle 401 Unauthorized
axios.interceptors.response.use((response) => {
  return response
}, (error) => {
  if (error.response && error.response.status === 401) {
    localStorage.removeItem('sigma_token')
    localStorage.removeItem('sigma_user')
    window.location.href = '/login'
  }
  return Promise.reject(error)
})

const app = createApp(App)
app.use(router)
app.mount('#app')