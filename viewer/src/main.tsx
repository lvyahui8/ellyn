import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import TrafficGraph from './Graph.tsx'
import Meta from './Meta.tsx'
import Menus from './Menus.tsx'
import './index.css'
import SourceView from "./Source.tsx";
import TrafficList from "./TrafficList.tsx"
import {BrowserRouter as Router,Routes,Route} from 'react-router-dom';
import SiteLayout from "./Layout.tsx";

createRoot(document.getElementById('root')!).render(
  // <StrictMode>
         <Router>

             <SiteLayout/>
         </Router>
  // </StrictMode>,
)
