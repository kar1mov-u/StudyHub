import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { ToastProvider } from '@/components/ui/toast'
import { AuthProvider } from '@/context/AuthContext'
import ProtectedRoute from '@/components/auth/ProtectedRoute'
import Layout from '@/components/layout/Layout'
import HomePage from '@/pages/HomePage'
import ModulesPage from '@/pages/ModulesPage'
import ModuleDetailPage from '@/pages/ModuleDetailPage'
import WeekDetailPage from '@/pages/WeekDetailPage'
import AcademicTermsPage from '@/pages/AcademicTermsPage'
import UserProfilePage from '@/pages/UserProfilePage'
import LoginPage from '@/pages/LoginPage'
import RegisterPage from '@/pages/RegisterPage'

function App() {
  return (
    <ToastProvider>
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            {/* Public routes */}
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            
            {/* Protected routes */}
            <Route path="/" element={
              <ProtectedRoute>
                <Layout />
              </ProtectedRoute>
            }>
              <Route index element={<Navigate to="/modules" replace />} />
              <Route path="home" element={<HomePage />} />
              <Route path="modules" element={<ModulesPage />} />
              <Route path="modules/:id" element={<ModuleDetailPage />} />
              <Route path="modules/:moduleId/weeks/:weekId" element={<WeekDetailPage />} />
              <Route path="academic-terms" element={<AcademicTermsPage />} />
              <Route path="users/:userId" element={<UserProfilePage />} />
            </Route>

            {/* Catch all - redirect to modules */}
            <Route path="*" element={<Navigate to="/modules" replace />} />
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </ToastProvider>
  )
}

export default App
