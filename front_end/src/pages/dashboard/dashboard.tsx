import { JSX, useEffect, useState, useRef } from 'react';
import { useNavigate } from 'react-router';
import { getDashboardResume } from '../../models/dashboard';
import { toast } from 'react-toastify';
interface IDashboardResumeData {
  totalAdmin: number;
  totalGreeting: number;
  totalQuestion: number;
  totalUser: number;
}

function DashboardMainPage(): JSX.Element {
  const navigate = useNavigate();

  const initialized = useRef(false)
  const [is401, setIs401] = useState(false);
  const [dashboardData, setDashboardData] = useState<IDashboardResumeData|null>(null);

  useEffect(() => {
    if (!initialized.current) {
      initialized.current = true;

      toast.promise(
        getDashboardResume().then(data => {
          setDashboardData(data.data);
        }),
        {
          error: {
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-ignore
            render({ data: { response: { data = 'Koneksi anda terputus!' } = {} } }){
              return data;
            }
          }
        },
      ).catch((err) => {
        console.log(err);
        if (err?.response?.status === 401) {
          setIs401(true);
        } else {
          initialized.current = false;
        }
      });
    }
  }, []);

  useEffect(() => {
    if (is401) {
      navigate('/login');
    }
  }, [is401, navigate]);

  return (
    <div className="container">
      <h3 className="text-dark mb-1">Dashboard</h3>

      <div className="row">
        <div className="col-md-3">
          <div className="card">
            <div className="card-body">
              <h4 className="card-title">Pertanyaan</h4>
              <h1 className="text-muted card-subtitle mb-2">
                {dashboardData?.totalQuestion ?? '...'}
              </h1>
            </div>
          </div>
        </div>

        <div className="col-md-3">
          <div className="card">
            <div className="card-body">
              <h4 className="card-title">Salam</h4>
              <h1 className="text-muted card-subtitle mb-2">
                {dashboardData?.totalGreeting ?? '...'}
              </h1>
            </div>
          </div>
        </div>

        <div className="col-md-3">
          <div className="card">
            <div className="card-body">
              <h4 className="card-title">Pengguna</h4>
              <h1 className="text-muted card-subtitle mb-2">
                {dashboardData?.totalUser ?? '...'}
              </h1>
            </div>
          </div>
        </div>
        
        <div className="col-md-3">
          <div className="card">
            <div className="card-body">
              <h4 className="card-title">Admin</h4>
              <h1 className="text-muted card-subtitle mb-2">
                {dashboardData?.totalAdmin ?? '...'}
              </h1>
            </div>
          </div>
        </div>
      </div>
    
      <div className="card mt-5">
        <div className="card-body">
          <h4 className="card-title">Pertanyaan Belum Dijawab (Manual)</h4>
          <div className="table-responsive">
            <table className="table">
              <thead>
                <tr>
                  <th>Pertanyaan</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  )
}

export default DashboardMainPage;