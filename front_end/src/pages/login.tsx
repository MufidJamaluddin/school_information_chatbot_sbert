import { JSX, useState, useEffect, FormEvent } from 'react';
import { useNavigate } from 'react-router';
import { toast } from 'react-toastify';
import { loginAction } from '../models/login';
import { getUserData } from '../models/session';

function LoginPage(): JSX.Element {
  const [loginData, setLoginData] = useState<{ username: string, password: string }|null>(null)
  const navigate = useNavigate();
  const userData = getUserData();

  useEffect(() => {
    if (!loginData || !loginData.username || !loginData.password) {
      return;
    }

    toast.promise(
      loginAction(loginData.username, loginData.password).then(() => {
        navigate('/dashboard');
      }),
      {
        error: {
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-ignore
          render({ data: { response } }){
            return `Login Gagal: ${typeof response === 'object' ? response?.data : response}`
          }
        }
      }
    );
  }, [loginData, navigate])

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const username = e.currentTarget.username.value;
    const password = e.currentTarget.password.value;

    setLoginData({ username, password });
  }

  if (userData?.fullName) {
    navigate('/dashboard');
    return <>Loading...</>
  }

  return (
    <div className="container" style={{'height': '100%'}}>
      <div className="row g-0 d-xl-flex justify-content-xl-center align-items-xl-center" style={{'height': '100%'}}>
        <div className="col-auto col-md-9 col-lg-12 col-xl-10 text-center d-xl-flex justify-content-xl-center align-items-xl-center" style={{'height': '100%'}}>
          <div className="card shadow-lg o-hidden border-0 my-5 text-center" style={{'height': '50%', 'width': '50%'}}>
            <div className="card-body p-0">
              <div className="row d-xl-flex align-items-xl-center">
                <div className="col-auto col-lg-8 col-xl-12">
                  <div className="p-5">
                    <div className="text-center">
                      <h4 className="text-dark mb-4">Selamat Datang</h4>
                    </div>
                    <form className="user" onSubmit={handleSubmit}>
                      <div className="mb-3">
                        <input 
                          className="form-control form-control-user" 
                          type="text" 
                          id="username" 
                          aria-describedby="username" 
                          placeholder="Masukkan Username" 
                          required={true}
                          name="username"/>
                      </div>
                      <div className="mb-3">
                        <input 
                          className="form-control form-control-user" 
                          type="password" 
                          id="password" 
                          placeholder="Masukkan Password" 
                          required={true}
                          name="password"/>
                      </div>
                      <div className="mb-3">
                        <div className="custom-control custom-checkbox small"></div>
                      </div>
                      <button 
                        className="btn btn-primary d-block btn-user w-100" 
                        type="submit">
                        Login
                      </button>
                    </form>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default LoginPage;