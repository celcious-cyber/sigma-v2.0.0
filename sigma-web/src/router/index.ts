import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/portal'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/auth/LoginView.vue')
  },
  {
    path: '/portal',
    name: 'Portal',
    component: () => import('../views/portal/portal.vue')
  },
  // Sigmabase Module
  {
    path: '/sigmabase',
    component: () => import('../layouts/SigmabaseLayout.vue'),
    children: [
      {
        path: '',
        name: 'Sigmabase',
        component: () => import('../views/sigmabase/dashboard.vue')
      },
      {
        path: 'students',
        name: 'Students',
        component: () => import('../views/sigmabase/students.vue')
      },
      {
        path: 'teachers',
        name: 'Teachers',
        component: () => import('../views/sigmabase/teachers.vue')
      },
      {
        path: 'alumni',
        name: 'Alumni',
        component: () => import('../views/sigmabase/alumni.vue')
      },
      {
        path: 'units',
        name: 'Units',
        component: () => import('../views/sigmabase/units.vue')
      }
    ]
  },

  {
    path: '/sigmaedu',
    component: () => import('../layouts/SigmaeduLayout.vue'),
    children: [
      {
        path: '',
        name: 'Sigmaedu',
        component: () => import('../views/sigmaedu/dashboard.vue')
      },
      {
        path: 'subjects',
        name: 'Subjects',
        component: () => import('../views/sigmaedu/subjects.vue')
      },
      {
        path: 'hours',
        name: 'StudyHours',
        component: () => import('../views/sigmaedu/hours.vue')
      },
      {
        path: 'calendar',
        name: 'AcademicCalendar',
        component: () => import('../views/sigmaedu/calendar.vue')
      },
      {
        path: 'classes',
        name: 'Classrooms',
        component: () => import('../views/sigmaedu/classes.vue')
      },
      {
        path: 'students',
        name: 'EduStudents',
        component: () => import('../views/sigmaedu/students.vue')
      },
      {
        path: 'schedules',
        name: 'EduSchedules',
        component: () => import('../views/sigmaedu/schedules.vue'),
        meta: { title: 'Jadwal Pelajaran' }
      },
      {
        path: 'attendance',
        name: 'EduAttendance',
        component: () => import('../views/sigmaedu/attendance.vue'),
        meta: { title: 'Presensi Santri' }
      },
      {
        path: 'grades',
        name: 'EduGrades',
        component: () => import('../views/sigmaedu/assessments.vue'),
        meta: { title: 'Nilai Pelajaran' }
      },
      {
        path: 'attitude',
        name: 'EduAttitude',
        component: () => import('../views/sigmaedu/attitude.vue'),
        meta: { title: 'Nilai Karakter' }
      },
      {
        path: 'tahfidz',
        name: 'EduTahfidz',
        component: () => import('../views/sigmaedu/tahfidz.vue'),
        meta: { title: 'Hafalan Quran' }
      },
      {
        path: 'memorization',
        name: 'EduMemorization',
        component: () => import('../views/sigmaedu/memorization.vue'),
        meta: { title: 'Hafalan Pelajaran' }
      },
      {
        path: 'teacher-attendance',
        name: 'EduTeacherAttendance',
        component: () => import('../views/sigmaedu/teacher_attendance.vue'),
        meta: { title: 'Presensi Guru' }
      },
      {
        path: 'journal',
        name: 'EduTeachingJournal',
        component: () => import('../views/sigmaedu/journal.vue'),
        meta: { title: 'Jurnal Mengajar' }
      },
      {
        path: 'reports/kmi',
        name: 'EduReportsKMI',
        component: () => import('../views/sigmaedu/reports.vue'),
        meta: { title: 'Rapot KMI', type: 'kmi' }
      },
      {
        path: 'reports/tahfidz',
        name: 'EduReportsTahfidz',
        component: () => import('../views/sigmaedu/reports_tahfidz.vue'),
        meta: { title: 'Rapot Tahfidz', type: 'tahfidz' }
      },
      {
        path: 'reports/ekskul',
        name: 'EduReportsEkskul',
        component: () => import('../views/sigmaedu/reports_ekskul.vue'),
        meta: { title: 'Rapot Ekskul', type: 'ekskul' }
      }
    ]
  },
  {
    path: '/dashboard',
    component: () => import('../layouts/AppLayout.vue'),
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('../views/DashboardView.vue')
      },
      // Sigmaflow
      {
        path: '/flow/invoices',
        name: 'Invoices',
        component: () => import('../views/sigmaflow/dashboard.vue')
      },
      // Sigmaguard
      {
        path: '/guard/violations',
        name: 'Violations',
        component: () => import('../views/sigmaguard/dashboard.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation Guard
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('sigma_token')
  if (to.name !== 'Login' && !token) next({ name: 'Login' })
  else next()
})

export default router