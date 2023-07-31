import Downloader from '../pages/downloader'
import LoraInfoComponent from '../pages/loraInfo'
import NotFound from '../pages/notFound'
import { useEffect, useState } from 'react'

import { Routes, Route, Link, useLocation } from 'react-router-dom'
import Upscale from '../pages/upscale'

const NavBar = () => {
  interface Path {
    path: string
    name: string
  }

  const paths: Path[] = [
    { path: '/upscale', name: 'Upscale' },
    { path: '/lora', name: 'LoraInfo' },
    { path: '/download', name: 'Download' },
  ]

  const [pageName, setPageName] = useState(paths[0].name)

  const pathToName = new Map<string, string>([])
  paths.forEach((p) => pathToName.set(p.path, p.name))

  const usePageViews = () => {
    const location = useLocation()

    useEffect(() => {
      const name = pathToName.get(location.pathname)
      if (name) {
        setPageName(name)
        document.title = name
        return
      }
      document.title = 'Not Found'
    }, [location])
  }
  usePageViews()

  return (
    <div className="drawer">
      <input id="my-drawer-3" type="checkbox" className="drawer-toggle" />
      <div className="drawer-content flex flex-col">
        {/* Navbar */}
        <div className="w-full navbar bg-base-300">
          <div className="flex-none lg:hidden">
            <label htmlFor="my-drawer-3" className="btn btn-square btn-ghost">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                className="inline-block w-6 h-6 stroke-current"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M4 6h16M4 12h16M4 18h16"
                ></path>
              </svg>
            </label>
          </div>
          <div className="flex-1 px-2 mx-2">{pageName}</div>
          <div className="flex-none hidden lg:block">
            <ul className="menu menu-horizontal">
              {paths.map((p) => (
                <li key={p.path}>
                  <Link to={p.path}>{p.name}</Link>
                </li>
              ))}
            </ul>
          </div>
        </div>
        {/* Content */}
        <Routes>
          <Route path="download" element={<Downloader />} />
          <Route path="lora" element={<LoraInfoComponent />} />
          <Route path="upscale" element={<Upscale />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </div>
      {/* Sidebar */}
      <div className="drawer-side">
        <label htmlFor="my-drawer-3" className="drawer-overlay"></label>
        <ul className="menu p-4 w-80 h-full bg-base-200">
          {paths.map((p) => (
            <li key={p.path}>
              <Link to={p.path}>{p.name}</Link>
            </li>
          ))}
        </ul>
      </div>
    </div>
  )
}

export default NavBar
