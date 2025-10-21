import React from 'react'
import { 
  AppBar, 
  Toolbar, 
  Typography, 
  Button, 
  Box, 
  Container,
  IconButton,
  Menu,
  MenuItem,
  useMediaQuery,
  useTheme
} from '@mui/material'
import { Link, useNavigate } from 'react-router-dom'
import { Menu as MenuIcon, AccountCircle } from '@mui/icons-material'
import { useAuth } from '../../contexts/AuthContext'

const Navbar: React.FC = () => {
  const { isAuthenticated, user, logout } = useAuth()
  const navigate = useNavigate()
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('md'))
  
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null)
  const [mobileMenuAnchor, setMobileMenuAnchor] = React.useState<null | HTMLElement>(null)

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleClose = () => {
    setAnchorEl(null)
  }

  const handleMobileMenu = (event: React.MouseEvent<HTMLElement>) => {
    setMobileMenuAnchor(event.currentTarget)
  }

  const handleMobileClose = () => {
    setMobileMenuAnchor(null)
  }

  const handleLogout = () => {
    logout()
    navigate('/')
    handleClose()
  }

  const publicPages = [
    { title: 'Home', path: '/' },
    { title: 'About', path: '/about' },
    { title: 'Coverage & Costs', path: '/coverage' },
    { title: 'How We Work', path: '/how-we-work' },
    { title: 'Retirement', path: '/retirement' },
    { title: 'Testimonials', path: '/testimonials' },
    { title: 'FAQ', path: '/faq' },
    { title: 'Contact', path: '/contact' },
  ]

  return (
    <AppBar position="sticky" sx={{ backgroundColor: 'primary.main' }}>
      <Container maxWidth="lg">
        <Toolbar disableGutters>
          <Typography
            variant="h6"
            component={Link}
            to="/"
            sx={{
              mr: 2,
              display: 'flex',
              fontWeight: 700,
              color: 'inherit',
              textDecoration: 'none',
            }}
          >
            TheraClosure
          </Typography>

          {/* Desktop Navigation */}
          {!isMobile && (
            <Box sx={{ flexGrow: 1, display: 'flex', gap: 2, ml: 4 }}>
              {publicPages.map((page) => (
                <Button
                  key={page.path}
                  component={Link}
                  to={page.path}
                  sx={{ 
                    color: 'white', 
                    textTransform: 'none',
                    '&:hover': { backgroundColor: 'rgba(255,255,255,0.1)' }
                  }}
                >
                  {page.title}
                </Button>
              ))}
            </Box>
          )}

          {/* Mobile Navigation */}
          {isMobile && (
            <Box sx={{ flexGrow: 1, display: 'flex', justifyContent: 'flex-end' }}>
              <IconButton
                size="large"
                onClick={handleMobileMenu}
                color="inherit"
              >
                <MenuIcon />
              </IconButton>
            </Box>
          )}

          {/* User Menu */}
          <Box sx={{ display: 'flex', gap: 1 }}>
            {isAuthenticated ? (
              <>
                <Button
                  component={Link}
                  to="/dashboard"
                  sx={{ 
                    color: 'white',
                    textTransform: 'none',
                  }}
                >
                  Dashboard
                </Button>
                <IconButton
                  size="large"
                  onClick={handleMenu}
                  color="inherit"
                >
                  <AccountCircle />
                </IconButton>
              </>
            ) : (
              <Button
                component={Link}
                to="/login"
                variant="contained"
                sx={{ 
                  backgroundColor: 'secondary.main',
                  color: 'primary.main',
                  '&:hover': { backgroundColor: 'secondary.dark' }
                }}
              >
                Login
              </Button>
            )}
          </Box>

          {/* User Dropdown Menu */}
          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleClose}
          >
            <MenuItem disabled>
              {user?.firstName} {user?.lastName}
            </MenuItem>
            <MenuItem onClick={handleLogout}>Logout</MenuItem>
          </Menu>

          {/* Mobile Menu */}
          <Menu
            anchorEl={mobileMenuAnchor}
            open={Boolean(mobileMenuAnchor)}
            onClose={handleMobileClose}
          >
            {publicPages.map((page) => (
              <MenuItem
                key={page.path}
                component={Link}
                to={page.path}
                onClick={handleMobileClose}
              >
                {page.title}
              </MenuItem>
            ))}
          </Menu>
        </Toolbar>
      </Container>
    </AppBar>
  )
}

export default Navbar