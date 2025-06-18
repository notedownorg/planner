import React from 'react'
import ReactDOM from 'react-dom/client'
import App from '@/App.tsx'
import { ThemeProvider } from '@/components/providers/ThemeProvider.tsx'
import '@/index.css'

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <ThemeProvider defaultTheme="system">
      <App />
    </ThemeProvider>
  </React.StrictMode>,
)
