import * as React from 'react';
import * as ReactDOM from 'react-dom/client';
import { ToastContainer } from 'react-toastify';

import {
  createBrowserRouter,
  RouterProvider,
} from 'react-router-dom';

import { 
  ChatPage, 
  LoginPage,
  dashboardRoute, 
} from './pages';

import './assets/bootstrap/css/bootstrap.min.css';
import './assets/css/styles.min.css';
import 'react-toastify/dist/ReactToastify.css';

const router = createBrowserRouter([
  {
    path: '',
    Component: ChatPage,
  },
  {
    path: '/login',
    Component: LoginPage,
  },
  {
    path: '/dashboard',
    ...dashboardRoute
  }
]);

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
    <ToastContainer />
  </React.StrictMode>
);
