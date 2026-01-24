import React from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { ToastProvider } from '@/components/ui/toast'
import Layout from '@/components/layout/Layout'
import HomePage from '@/pages/HomePage'
import ModulesPage from '@/pages/ModulesPage'
import ModuleDetailPage from '@/pages/ModuleDetailPage'
import AcademicTermsPage from '@/pages/AcademicTermsPage'

function App() {
  return (
    <ToastProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<HomePage />} />
            <Route path="modules" element={<ModulesPage />} />
            <Route path="modules/:id" element={<ModuleDetailPage />} />
            <Route path="academic-terms" element={<AcademicTermsPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ToastProvider>
  )
}

export default App
