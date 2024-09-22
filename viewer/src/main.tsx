import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import TrafficGraph from './Graph.tsx'
import './index.css'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <TrafficGraph />
  </StrictMode>,
)
