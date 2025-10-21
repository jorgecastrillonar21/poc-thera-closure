import React, { useState, useEffect } from 'react'
import {
  Box,
  Paper,
  Typography,
  TextField,
  Button,
  Grid,
  MenuItem,
  FormControl,
  InputLabel,
  Select,
  Divider,
  Alert,
  CircularProgress,
} from '@mui/material'
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider'
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns'
import { UserProfileDTO, CreateUserProfileDTO, UpdateUserProfileDTO, usersService } from '../services/usersService'

interface UserProfileFormProps {
  userId?: string // For editing existing profile
  onSave?: (profile: UserProfileDTO) => void
  onCancel?: () => void
}

const US_STATES = [
  'AL', 'AK', 'AZ', 'AR', 'CA', 'CO', 'CT', 'DE', 'FL', 'GA',
  'HI', 'ID', 'IL', 'IN', 'IA', 'KS', 'KY', 'LA', 'ME', 'MD',
  'MA', 'MI', 'MN', 'MS', 'MO', 'MT', 'NE', 'NV', 'NH', 'NJ',
  'NM', 'NY', 'NC', 'ND', 'OH', 'OK', 'OR', 'PA', 'RI', 'SC',
  'SD', 'TN', 'TX', 'UT', 'VT', 'VA', 'WA', 'WV', 'WI', 'WY'
]

const UserProfileForm: React.FC<UserProfileFormProps> = ({ userId, onSave, onCancel }) => {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)
  
  const [formData, setFormData] = useState<CreateUserProfileDTO & { id?: string; user_id?: string }>({
    first_name: '',
    last_name: '',
    email: '',
    phone: '',
    date_of_birth: '',
    license_number: '',
    license_state: '',
    license_expiration: '',
    years_of_experience: 0,
    practice_name: '',
    practice_address: '',
    practice_phone: '',
  })

  const isEditing = Boolean(userId)

  useEffect(() => {
    if (isEditing && userId) {
      loadProfile(userId)
    }
  }, [userId, isEditing])

  const loadProfile = async (profileUserId: string) => {
    try {
      setLoading(true)
      setError(null)
      const profile = await usersService.getProfile(profileUserId)
      setFormData({
        first_name: profile.first_name || '',
        last_name: profile.last_name || '',
        email: profile.email || '',
        phone: profile.phone || '',
        date_of_birth: profile.date_of_birth ? profile.date_of_birth.split('T')[0] : '',
        license_number: profile.license_number || '',
        license_state: profile.license_state || '',
        license_expiration: profile.license_expiration ? profile.license_expiration.split('T')[0] : '',
        years_of_experience: profile.years_of_experience || 0,
        practice_name: profile.practice_name || '',
        practice_address: profile.practice_address || '',
        practice_phone: profile.practice_phone || '',
        id: profile.id,
        user_id: profile.user_id,
      })
    } catch (err: any) {
      setError(`Failed to load profile: ${err.response?.data?.error || err.message}`)
    } finally {
      setLoading(false)
    }
  }

  const handleInputChange = (field: keyof typeof formData, value: string | number) => {
    setFormData(prev => ({
      ...prev,
      [field]: value
    }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      setLoading(true)
      setError(null)
      setSuccess(null)

      const submitData = {
        first_name: formData.first_name,
        last_name: formData.last_name,
        email: formData.email,
        phone: formData.phone || undefined,
        date_of_birth: formData.date_of_birth ? `${formData.date_of_birth}T00:00:00Z` : undefined,
        license_number: formData.license_number || undefined,
        license_state: formData.license_state || undefined,
        license_expiration: formData.license_expiration ? `${formData.license_expiration}T00:00:00Z` : undefined,
        years_of_experience: formData.years_of_experience || 0,
        practice_name: formData.practice_name || undefined,
        practice_address: formData.practice_address || undefined,
        practice_phone: formData.practice_phone || undefined,
      }

      let savedProfile: UserProfileDTO

      if (isEditing && userId) {
        savedProfile = await usersService.updateProfile(userId, submitData as UpdateUserProfileDTO)
        setSuccess('Profile updated successfully!')
      } else {
        savedProfile = await usersService.createProfile(submitData as CreateUserProfileDTO)
        setSuccess('Profile created successfully!')
      }

      if (onSave) {
        onSave(savedProfile)
      }
    } catch (err: any) {
      setError(`Failed to ${isEditing ? 'update' : 'create'} profile: ${err.response?.data?.error || err.message}`)
    } finally {
      setLoading(false)
    }
  }

  if (loading && isEditing) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
        <CircularProgress />
      </Box>
    )
  }

  return (
    <LocalizationProvider dateAdapter={AdapterDateFns}>
      <Paper elevation={3} sx={{ p: 4, maxWidth: 800, mx: 'auto' }}>
        <Typography variant="h5" gutterBottom>
          {isEditing ? 'Edit Profile' : 'Create New Profile'}
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mb: 3 }}>
            {error}
          </Alert>
        )}

        {success && (
          <Alert severity="success" sx={{ mb: 3 }}>
            {success}
          </Alert>
        )}

        <Box component="form" onSubmit={handleSubmit}>
          {/* Personal Information */}
          <Typography variant="h6" gutterBottom sx={{ mt: 2 }}>
            Personal Information
          </Typography>
          <Divider sx={{ mb: 2 }} />

          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <TextField
                required
                fullWidth
                label="First Name"
                value={formData.first_name}
                onChange={(e) => handleInputChange('first_name', e.target.value)}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                required
                fullWidth
                label="Last Name"
                value={formData.last_name}
                onChange={(e) => handleInputChange('last_name', e.target.value)}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                required
                fullWidth
                type="email"
                label="Email"
                value={formData.email}
                onChange={(e) => handleInputChange('email', e.target.value)}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Phone"
                value={formData.phone}
                onChange={(e) => handleInputChange('phone', e.target.value)}
                placeholder="+1-555-0123"
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                type="date"
                label="Date of Birth"
                value={formData.date_of_birth}
                onChange={(e) => handleInputChange('date_of_birth', e.target.value)}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>
          </Grid>

          {/* Professional Information */}
          <Typography variant="h6" gutterBottom sx={{ mt: 4 }}>
            Professional Information
          </Typography>
          <Divider sx={{ mb: 2 }} />

          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="License Number"
                value={formData.license_number}
                onChange={(e) => handleInputChange('license_number', e.target.value)}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <FormControl fullWidth>
                <InputLabel>License State</InputLabel>
                <Select
                  value={formData.license_state}
                  onChange={(e) => handleInputChange('license_state', e.target.value)}
                  label="License State"
                >
                  {US_STATES.map(state => (
                    <MenuItem key={state} value={state}>{state}</MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                type="date"
                label="License Expiration"
                value={formData.license_expiration}
                onChange={(e) => handleInputChange('license_expiration', e.target.value)}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                type="number"
                label="Years of Experience"
                value={formData.years_of_experience}
                onChange={(e) => handleInputChange('years_of_experience', parseInt(e.target.value) || 0)}
                inputProps={{ min: 0, max: 50 }}
              />
            </Grid>
          </Grid>

          {/* Practice Information */}
          <Typography variant="h6" gutterBottom sx={{ mt: 4 }}>
            Practice Information
          </Typography>
          <Divider sx={{ mb: 2 }} />

          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Practice Name"
                value={formData.practice_name}
                onChange={(e) => handleInputChange('practice_name', e.target.value)}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Practice Phone"
                value={formData.practice_phone}
                onChange={(e) => handleInputChange('practice_phone', e.target.value)}
                placeholder="+1-555-0456"
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Practice Address"
                value={formData.practice_address}
                onChange={(e) => handleInputChange('practice_address', e.target.value)}
                multiline
                rows={2}
              />
            </Grid>
          </Grid>

          {/* Actions */}
          <Box sx={{ display: 'flex', justifyContent: 'flex-end', gap: 2, mt: 4 }}>
            {onCancel && (
              <Button onClick={onCancel} disabled={loading}>
                Cancel
              </Button>
            )}
            <Button
              type="submit"
              variant="contained"
              disabled={loading}
              startIcon={loading ? <CircularProgress size={20} /> : undefined}
            >
              {loading ? (isEditing ? 'Updating...' : 'Creating...') : (isEditing ? 'Update Profile' : 'Create Profile')}
            </Button>
          </Box>
        </Box>
      </Paper>
    </LocalizationProvider>
  )
}

export default UserProfileForm