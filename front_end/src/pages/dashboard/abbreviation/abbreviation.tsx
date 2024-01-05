import { JSX, useEffect, useState, useRef } from 'react';
import { useNavigate } from 'react-router';
import { getAbbreviationData } from '../../../models/abbreviation';
import { toast } from 'react-toastify';
import { IPagination, SearchParam } from '../../../models/shared';

interface IAbbreviationData {
  standardWord: string
  listAbbreviationTerm: Array<string>
  updatedAt?: string
  updatedBy?: string
  createdAt?: string
  createdBy?: string
}

function AbbreviationPage(): JSX.Element {
  const navigate = useNavigate();

  const typingTimer = useRef<NodeJS.Timeout>();

  const [is401, setIs401] = useState(false);
  
  const [searchParam, setSearchParam] = useState<SearchParam>({
    page: 1,
    size: 5,
    keyword: ''
  });

  const [data, setData] = useState<IPagination<IAbbreviationData>>({
    start: 0,
    end: 0,
    length: 0,
    data: []
  });

  useEffect(() => {
    if (searchParam) {
      toast.promise(
        getAbbreviationData(
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

  const goToEdit = (state: IAbbreviationData) => {
    navigate('/dashboard/abbreviation/edit', {
      state
    });
  }

  const goToDelete = (state: IAbbreviationData) => {
    navigate('/dashboard/abbreviation/delete', {
      state
    });
  }

  const pages = [];
  const totalPage = Math.ceil((data?.length || 0) / (searchParam.size || 1)); 

  if (totalPage > 1) {
    if (searchParam.page > 5) {
      pages.push(1);
    }

    if (searchParam.page > 2) {
      pages.push(searchParam.page - 2);
    }

    if (searchParam.page > 1) {
      pages.push(searchParam.page - 1);
    }
  
    pages.push(searchParam.page);
  
    if (searchParam.page < totalPage) {
      pages.push(searchParam.page + 1);
    }

    if (searchParam.page + 1 < totalPage) {
      pages.push(searchParam.page + 2);
    }
  
    if (searchParam.page + 3 < totalPage) {
      pages.push(totalPage);
    }
  }

  return (
    <div className="container">
      <div className="row">
        <div className="col">
          <h1 className="text-dark mb-1">Singkatan</h1>
          <p>Tidak Dipakai Langsung, Terpisah dari Aplikasi, perlu Fine-Tuned Ulang</p>
        </div>
        <div className="col-md-6 col-xl-2">
          <button className="btn btn-primary" type="button" onClick={() => navigate('/dashboard/abbreviation/new')}>
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
                  <th>Kata Baku Standar</th>
                  <th>Kata Sinonim</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                {
                  data.data.map(item => (
                    <tr key={item.standardWord}>
                      <td>{item.standardWord}</td>
                      <td>
                        {
                          item.listAbbreviationTerm.map(abbreviationItem => (
                            <div className="alert alert-secondary" role="alert" key={`${item.standardWord}_${abbreviationItem}`}>
                              {abbreviationItem}
                            </div>
                          ))
                        }
                      </td>
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

export default AbbreviationPage;