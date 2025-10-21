import { createTheme } from '@mui/material/styles'

// TheraClosure brand colors
const colors = {
  primary: {
    main: '#2C5F5D', // Teal
    light: '#4A8B8A',
    dark: '#1A3B3A',
    contrastText: '#FFFFFF',
  },
  secondary: {
    main: '#F5F1E8', // Cream
    light: '#F8F5ED',
    dark: '#E8E2D6',
    contrastText: '#2C5F5D',
  },
  accent: {
    main: '#D4A574', // Gold accent
    light: '#E0B68A',
    dark: '#B8935F',
  },
  background: {
    default: '#FAFAFA',
    paper: '#FFFFFF',
  },
  text: {
    primary: '#2C3E50',
    secondary: '#5D6D7E',
  },
}

export const theme = createTheme({
  palette: {
    primary: colors.primary,
    secondary: colors.secondary,
    background: colors.background,
    text: colors.text,
  },
  typography: {
    fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
    h1: {
      fontSize: '2.5rem',
      fontWeight: 700,
      color: colors.primary.main,
    },
    h2: {
      fontSize: '2rem',
      fontWeight: 600,
      color: colors.primary.main,
    },
    h3: {
      fontSize: '1.5rem',
      fontWeight: 600,
      color: colors.primary.main,
    },
    h4: {
      fontSize: '1.25rem',
      fontWeight: 500,
      color: colors.primary.main,
    },
    h5: {
      fontSize: '1.125rem',
      fontWeight: 500,
      color: colors.primary.main,
    },
    h6: {
      fontSize: '1rem',
      fontWeight: 500,
      color: colors.primary.main,
    },
    body1: {
      fontSize: '1rem',
      lineHeight: 1.6,
    },
    body2: {
      fontSize: '0.875rem',
      lineHeight: 1.5,
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          borderRadius: 8,
          fontWeight: 500,
          padding: '12px 24px',
        },
        contained: {
          boxShadow: '0 2px 8px rgba(44, 95, 93, 0.15)',
          '&:hover': {
            boxShadow: '0 4px 12px rgba(44, 95, 93, 0.25)',
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          boxShadow: '0 2px 12px rgba(0, 0, 0, 0.08)',
          '&:hover': {
            boxShadow: '0 4px 20px rgba(0, 0, 0, 0.12)',
          },
        },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            borderRadius: 8,
          },
        },
      },
    },
  },
})