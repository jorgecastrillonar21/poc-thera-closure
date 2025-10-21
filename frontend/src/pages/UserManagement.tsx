import React, { useState } from 'react'
import { Container, Typography, Box, Tabs, Tab, Paper } from '@mui/material'
import UserProfileForm from '../components/UserProfileForm'
import UserProfileList from '../components/UserProfileList'
import { UserProfileDTO } from '../services/usersService'

interface TabPanelProps {
  children?: React.ReactNode
  index: number
  value: number
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`user-tabpanel-${index}`}
      aria-labelledby={`user-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ py: 3 }}>{children}</Box>}
    </div>
  )
}

const UserManagement: React.FC = () => {
  const [currentTab, setCurrentTab] = useState(0)
  const [editingUserId, setEditingUserId] = useState<string | null>(null)

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setCurrentTab(newValue)
    // Clear editing state when switching tabs
    if (newValue !== 1) {
      setEditingUserId(null)
    }
  }

  const handleCreateNew = () => {
    setEditingUserId(null)
    setCurrentTab(1) // Switch to form tab
  }

  const handleEditProfile = (userId: string) => {
    setEditingUserId(userId)
    setCurrentTab(1) // Switch to form tab
  }

  const handleProfileSaved = (profile: UserProfileDTO) => {
    console.log('Profile saved:', profile)
    // Switch back to list view after saving
    setCurrentTab(0)
    setEditingUserId(null)
  }

  const handleFormCancel = () => {
    setCurrentTab(0) // Switch back to list view
    setEditingUserId(null)
  }

  return (
    <Container maxWidth="xl" sx={{ py: 4 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        User Management
      </Typography>
      
      <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
        Manage therapist profiles and user accounts
      </Typography>

      <Paper elevation={2}>
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs value={currentTab} onChange={handleTabChange} aria-label="user management tabs">
            <Tab label="User Profiles" id="user-tab-0" aria-controls="user-tabpanel-0" />
            <Tab 
              label={editingUserId ? "Edit Profile" : "Create Profile"} 
              id="user-tab-1" 
              aria-controls="user-tabpanel-1" 
            />
          </Tabs>
        </Box>

        <TabPanel value={currentTab} index={0}>
          <UserProfileList 
            onCreateNew={handleCreateNew}
            onEdit={handleEditProfile}
          />
        </TabPanel>

        <TabPanel value={currentTab} index={1}>
          <UserProfileForm
            userId={editingUserId || undefined}
            onSave={handleProfileSaved}
            onCancel={handleFormCancel}
          />
        </TabPanel>
      </Paper>
    </Container>
  )
}

export default UserManagement