import React from 'react'
import { Routes, Route } from 'react-router-dom'
import { Box } from '@mui/material'
import Layout from './components/Layout/Layout'
import ProtectedRoute from './components/ProtectedRoute'
import Home from './pages/Home'
import About from './pages/About'
import Coverage from './pages/Coverage'
import HowWeWork from './pages/HowWeWork'
import { Retirement, Testimonials, FAQ, Contact, Templates, Billing, Support } from './pages/AllPages'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import Enrollment from './pages/Enrollment'
import UserManagement from './pages/UserManagement'

const App: React.FC = () => {
  return (
    <Box>
      <Routes>
        <Route path="/" element={<Layout />}>
          {/* Public routes */}
          <Route index element={<Home />} />
          <Route path="about" element={<About />} />
          <Route path="coverage" element={<Coverage />} />
          <Route path="how-we-work" element={<HowWeWork />} />
          <Route path="retirement" element={<Retirement />} />
          <Route path="testimonials" element={<Testimonials />} />
          <Route path="faq" element={<FAQ />} />
          <Route path="contact" element={<Contact />} />
          <Route path="login" element={<Login />} />
          
          {/* Protected routes */}
          <Route path="dashboard" element={<ProtectedRoute />}>
            <Route index element={<Dashboard />} />
          </Route>
          <Route path="enrollment" element={<Enrollment />} />
          <Route path="user-management" element={<UserManagement />} />
          <Route path="templates" element={<Templates />} />
          <Route path="billing" element={<Billing />} />
          <Route path="support" element={<Support />} />
        </Route>
      </Routes>
    </Box>
  )
}

export default App