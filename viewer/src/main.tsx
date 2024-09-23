import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import TrafficGraph from './Graph.tsx'
import Meta from './Meta.tsx'
import Menus from './Menus.tsx'
import { Router, Route } from 'react-router'
import './index.css'
import SourceView from "./Source.tsx";

createRoot(document.getElementById('root')!).render(
  <StrictMode>
      <Menus/>
      <TrafficGraph/>
      <Meta/>
      <SourceView/>
  </StrictMode>,
)
