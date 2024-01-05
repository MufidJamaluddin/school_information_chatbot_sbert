import { JSX, useEffect, useState, useRef } from 'react';
import { useNavigate } from 'react-router';
import { getGreetingData } from '../../../models/greeting';
import { toast } from 'react-toastify';
import { IPagination, SearchParam } from '../../../models/shared';

interface IGreetingData {
  id?: number,
  greeting: string
  startTime: string
  endTime: string
  updatedAt?: string
  updatedBy?: string
  createdAt?: string
  createdBy?: string
}

function GreetingPage(): JSX.Element {
  const navigate = useNavigate();

  const typingTimer = useRef<NodeJS.Timeout>();

  const [is401, setIs401] = useState(false);
  
  const [searchParam, setSearchParam] = useState<SearchParam>({
    page: 1,
    size: 5,
    keyword: ''
  });

  const [data, setData] = useState<IPagination<IGreetingData>>({
    start: 0,
    end: 0,
    length: 0,
    data: []
  });

  useEffect(() => {
    if (searchParam) {
      toast.promise(
        getGreetingData(
          (searchParam.page - 1) * searchParam.size,
          searchParam.size,
          searchParam.keyword
        ).then(data => {
          setData(data);
        }),
        {
          error: {
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-ignore
            render({ data: { response: { data = 'Koneksi anda terputus!' } = {} } = {} }){
              return data;
            }
          }
        }
      ).catch((err) => {
        console.log(err);
        if (err?.response?.status === 401) {
          setIs401(true);
        }
      });
    }
  }, [searchParam]);

  useEffect(() => {
    if (is401) {
      navigate('/login');
    }
  }, [is401, navigate]);

  const changeParams = (params: Partial<SearchParam>) => {
    setSearchParam({ ...searchParam, ...params });
  }

  const changeKeyword = (newKeyword: string) => {
    const cleanKeyword = String(newKeyword).trim();

    clearTimeout(typingTimer.current);

    typingTimer.current = setTimeout(() => changeParams({ keyword: cleanKeyword }), 600);
  }

  const goToEdit = (state: IGreetingData) => {
    navigate('/dashboard/greeting/edit', {
      state
    });
  }

  const goToDelete = (state: IGreetingData) => {
    navigate('/dashboard/greeting/delete', {
      state
    });
  }

  const pages = [];
  const totalPage = Math.ceil((data?.length || 0) / (searchParam.size || 1)); 

  if (totalPage > 1) {
    pages.push(1);
  
    if (searchParam.page > 1) {
      pages.push(searchParam.page - 1);
    }
  
    pages.push(searchParam.page);
  
    if (searchParam.page < totalPage) {
      pages.push(searchParam.page + 1);
    }
  
    pages.push(totalPage);
  }

  return (
    <div className="container">
      <div className="row">
        <div className="col">
          <h1 className="text-dark mb-1">Salam</h1>
        </div>
        <div className="col-md-6 col-xl-2">
          <button className="btn btn-primary" type="button" onClick={() => navigate('/dashboard/greeting/new')}>
            Buat Baru
          </button>
        </div>
      </div>

      <hr/>

      <div className="row">
        <div className='col-md-2'>
          Page Size
        </div>

        <div className='col-md-3'>
          <select className="form-select" onChange={(e) => changeParams({ size: Number(e.target.value) })}>
            <option value="5">5</option>
            <option value="10">10</option>
            <option value="20">20</option>
            <option value="50">50</option>
            <option value="100">100</option>
          </select>
        </div>

        <div className='col-md-2 text-end'>
          Cari
        </div>

        <div className='col-md-5'>
          <input 
            type="text" 
            className="form-control" 
            onKeyUp={(e) => changeKeyword(e.currentTarget.value)}
          />
        </div>

      </div>
      
      <hr/>

      <div className="row">

        <div className="col-xl-12">
          <div className="table-responsive">
            <table className="table">
              <thead>
                <tr>
                  <th>Pesan Salam</th>
                  <th>Mulai</th>
                  <th>Selesai</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                {
                  data.data.map(item => (
                    <tr key={item.id}>
                      <td>{item.greeting}</td>
                      <td>{item.startTime}</td>
                      <td>{item.endTime}</td>
                      <td>
                          <div className="col-xl-5">
                            <button className="btn btn-success" type="button" onClick={() => goToEdit(item)}>
                              <i className="far fa-edit"></i>
                            </button>
                            &nbsp;
                            <button className="btn btn-danger" type="button" onClick={() => goToDelete(item)}>
                              <i className="far fa-trash-alt"></i>
                            </button>
                          </div>
                      </td>
                  </tr>
                  ))
                }
              </tbody>
            </table>

            <nav aria-label="Page navigation example">
              <ul className="pagination">
                {
                  pages.map(pageNum => (
                    <li className={pageNum === searchParam?.page ? "page-item active" : "page-item"} key={pageNum}>
                      <a className="page-link" onClick={() => changeParams({ page: pageNum })}>
                        {pageNum}
                      </a>
                    </li>
                  ))
                }
              </ul>
            </nav>
          </div>
        </div>
      </div>
    </div>
  )
}

export default GreetingPage;