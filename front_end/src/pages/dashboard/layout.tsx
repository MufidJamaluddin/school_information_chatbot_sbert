import { JSX  } from 'react';
import { Link, Outlet } from 'react-router-dom';
import { getUserData } from '../../models/session';

interface DashboardLayoutProp {
  links: Array<{
    path: string,
    icon: JSX.Element,
    title: string
  }>

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  children?: any 
}

function DashboardLayout(props: DashboardLayoutProp): JSX.Element {
  const userData = getUserData();

  return (
    <div id="wrapper">
      <nav className="navbar align-items-start sidebar sidebar-dark accordion bg-gradient-primary p-0 navbar-dark">
        <div className="container-fluid d-flex flex-column p-0">
          <div className="navbar-brand d-flex justify-content-center align-items-center sidebar-brand m-0">
              <div className="sidebar-brand-text mx-3">
                <span>SMAN Situraja</span>
              </div>
          </div>
          
          <hr className="sidebar-divider my-0"/>

          <ul 
            className="navbar-nav text-light" 
            id="accordionSidebar">
            <li className="nav-item">
              {
                props.links.map(item => (
                  <Link to={item.path} key={item.path} className='nav-link'>
                    {item.icon}
                    &nbsp;
                    {item.title}
                  </Link>
                ))
              }
            </li>
          </ul>
        </div>
      </nav>

    <div 
      className="d-flex flex-column" 
      id="content-wrapper">

      <div id="content">
        <nav className="navbar navbar-expand bg-white shadow mb-4 topbar static-top navbar-light">
          <div className="container-fluid">
            <button className="btn btn-link d-md-none rounded-circle me-3" id="sidebarToggleTop" type="button">
              <i className="fas fa-bars"></i>
            </button>
            <ul className="navbar-nav flex-nowrap ms-auto">
              <li className="nav-item dropdown no-arrow">
                <div className="nav-item dropdown no-arrow">
                  <a className="dropdown-toggle nav-link" aria-expanded="false" data-bs-toggle="dropdown" href="#">
                    <span className="d-none d-lg-inline me-2 text-gray-600 small">
                      {userData?.fullName}
                    </span>
                  </a>
                  <div className="dropdown-menu shadow dropdown-menu-end animated--grow-in">
                    <a className="dropdown-item" href="#">
                      <i className="fas fa-sign-out-alt fa-sm fa-fw me-2 text-gray-400"></i>
                      &nbsp;Logout
                    </a>
                  </div>
                </div>
              </li>
            </ul>
          </div>
        </nav>
          
        {props.children}

        <Outlet />
      </div>

      <footer className="bg-white sticky-footer">
        <div className="container my-auto">
          <div className="text-center my-auto copyright">
            <span>Copyright © Tanya SMAN Situraja 2023</span>
          </div>
        </div>
      </footer>
    
  
    </div>
      <a 
        className="border rounded d-inline scroll-to-top" 
        href="#page-top">
        <i className="fas fa-angle-up"></i>
      </a>
    </div>
  )
}

export default DashboardLayout;