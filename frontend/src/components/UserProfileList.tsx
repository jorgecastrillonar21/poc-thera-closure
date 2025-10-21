import React, { useState, useEffect } from 'react'
import {
  Box,
  Paper,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  Button,
  TextField,
  InputAdornment,
  IconButton,
  Chip,
  Alert,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from '@mui/material'
import {
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  PersonAddOutlined,
  RefreshOutlined,
} from '@mui/icons-material'
import { UserProfileDTO, usersService } from '../services/usersService'

interface UserProfileListProps {
  onCreateNew?: () => void
  onEdit?: (userId: string) => void
}

const UserProfileList: React.FC<UserProfileListProps> = ({ onCreateNew, onEdit }) => {
  const [profiles, setProfiles] = useState<UserProfileDTO[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [searchQuery, setSearchQuery] = useState('')
  const [page, setPage] = useState(0)
  const [rowsPerPage, setRowsPerPage] = useState(10)
  const [totalCount, setTotalCount] = useState(0)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [profileToDelete, setProfileToDelete] = useState<UserProfileDTO | null>(null)
  const [deleting, setDeleting] = useState(false)

  useEffect(() => {
    loadProfiles()
  }, [page, rowsPerPage])

  useEffect(() => {
    // Reset page when search query changes
    if (page > 0) {
      setPage(0)
    } else {
      loadProfiles()
    }
  }, [searchQuery])

  const loadProfiles = async () => {
    try {
      setLoading(true)
      setError(null)
      
      const offset = page * rowsPerPage
      
      if (searchQuery.trim()) {
        const response = await usersService.searchProfiles(searchQuery, rowsPerPage, offset)
        setProfiles(response.profiles || [])
        setTotalCount(response.count || 0)
      } else {
        const response = await usersService.listProfiles(rowsPerPage, offset)
        setProfiles(response.profiles || [])
        setTotalCount(response.total || 0)
      }
    } catch (err: any) {
      setError(`Failed to load profiles: ${err.response?.data?.error || err.message}`)
      setProfiles([])
      setTotalCount(0)
    } finally {
      setLoading(false)
    }
  }

  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(event.target.value)
  }

  const handlePageChange = (_event: unknown, newPage: number) => {
    setPage(newPage)
  }

  const handleRowsPerPageChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10))
    setPage(0)
  }

  const handleDeleteClick = (profile: UserProfileDTO) => {
    setProfileToDelete(profile)
    setDeleteDialogOpen(true)
  }

  const handleDeleteConfirm = async () => {
    if (!profileToDelete) return

    try {
      setDeleting(true)
      await usersService.deleteProfile(profileToDelete.user_id)
      setDeleteDialogOpen(false)
      setProfileToDelete(null)
      await loadProfiles() // Refresh the list
    } catch (err: any) {
      setError(`Failed to delete profile: ${err.response?.data?.error || err.message}`)
    } finally {
      setDeleting(false)
    }
  }

  const handleDeleteCancel = () => {
    setDeleteDialogOpen(false)
    setProfileToDelete(null)
  }

  const formatDate = (dateString?: string) => {
    if (!dateString) return 'N/A'
    return new Date(dateString).toLocaleDateString()
  }

  const getStatusChip = (profile: UserProfileDTO) => {
    const status = profile.status || 'unknown'
    const color = status === 'active' ? 'success' : status === 'inactive' ? 'error' : 'default'
    return <Chip label={status.toUpperCase()} color={color} size="small" />
  }

  return (
    <Box>
      <Paper elevation={3} sx={{ overflow: 'hidden' }}>
        {/* Header */}
        <Box sx={{ p: 3, borderBottom: '1px solid #e0e0e0' }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
            <Typography variant="h5">
              User Profiles
            </Typography>
            <Box sx={{ display: 'flex', gap: 1 }}>
              <IconButton onClick={loadProfiles} disabled={loading}>
                <RefreshOutlined />
              </IconButton>
              {onCreateNew && (
                <Button
                  variant="contained"
                  startIcon={<PersonAddOutlined />}
                  onClick={onCreateNew}
                >
                  New Profile
                </Button>
              )}
            </Box>
          </Box>

          {/* Search */}
          <TextField
            fullWidth
            placeholder="Search profiles by name, email, or practice..."
            value={searchQuery}
            onChange={handleSearchChange}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <SearchOutlined />
                </InputAdornment>
              ),
            }}
          />
        </Box>

        {/* Error Alert */}
        {error && (
          <Alert severity="error" sx={{ mx: 3, mt: 2 }}>
            {error}
          </Alert>
        )}

        {/* Loading */}
        {loading && (
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
            <CircularProgress />
          </Box>
        )}

        {/* Table */}
        {!loading && (
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Name</TableCell>
                  <TableCell>Email</TableCell>
                  <TableCell>Practice</TableCell>
                  <TableCell>License</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Created</TableCell>
                  <TableCell align="right">Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {profiles.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} align="center" sx={{ py: 4 }}>
                      <Typography variant="body2" color="text.secondary">
                        {searchQuery ? 'No profiles found matching your search' : 'No profiles found'}
                      </Typography>
                    </TableCell>
                  </TableRow>
                ) : (
                  profiles.map((profile) => (
                    <TableRow key={profile.id} hover>
                      <TableCell>
                        <Box>
                          <Typography variant="body2" fontWeight="medium">
                            {profile.first_name} {profile.last_name}
                          </Typography>
                          {profile.professional_title && (
                            <Typography variant="caption" color="text.secondary">
                              {profile.professional_title}
                            </Typography>
                          )}
                        </Box>
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {profile.email}
                        </Typography>
                        {profile.phone && (
                          <Typography variant="caption" color="text.secondary" display="block">
                            {profile.phone}
                          </Typography>
                        )}
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {profile.practice_name || 'N/A'}
                        </Typography>
                        {profile.practice_address && (
                          <Typography variant="caption" color="text.secondary" display="block">
                            {profile.practice_address}
                          </Typography>
                        )}
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {profile.license_number || 'N/A'}
                        </Typography>
                        {profile.license_state && (
                          <Typography variant="caption" color="text.secondary" display="block">
                            {profile.license_state}
                          </Typography>
                        )}
                      </TableCell>
                      <TableCell>
                        {getStatusChip(profile)}
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {formatDate(profile.created_at)}
                        </Typography>
                      </TableCell>
                      <TableCell align="right">
                        <Box sx={{ display: 'flex', gap: 1 }}>
                          {onEdit && (
                            <IconButton
                              size="small"
                              onClick={() => onEdit(profile.user_id)}
                              title="Edit Profile"
                            >
                              <EditOutlined />
                            </IconButton>
                          )}
                          <IconButton
                            size="small"
                            onClick={() => handleDeleteClick(profile)}
                            title="Delete Profile"
                            color="error"
                          >
                            <DeleteOutlined />
                          </IconButton>
                        </Box>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </TableContainer>
        )}

        {/* Pagination */}
        {!loading && profiles.length > 0 && (
          <TablePagination
            rowsPerPageOptions={[5, 10, 25, 50]}
            component="div"
            count={totalCount}
            rowsPerPage={rowsPerPage}
            page={page}
            onPageChange={handlePageChange}
            onRowsPerPageChange={handleRowsPerPageChange}
          />
        )}
      </Paper>

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteDialogOpen} onClose={handleDeleteCancel}>
        <DialogTitle>Confirm Deletion</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to delete the profile for{' '}
            <strong>
              {profileToDelete?.first_name} {profileToDelete?.last_name}
            </strong>
            ? This action cannot be undone.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleDeleteCancel} disabled={deleting}>
            Cancel
          </Button>
          <Button
            onClick={handleDeleteConfirm}
            color="error"
            disabled={deleting}
            startIcon={deleting ? <CircularProgress size={16} /> : undefined}
          >
            {deleting ? 'Deleting...' : 'Delete'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}

export default UserProfileList