import { JSX, useEffect, useState, useRef } from 'react';
import { toast } from 'react-toastify';
import { getQuestionData } from '../../../models/question';
import { useNavigate } from 'react-router-dom';
import { IPagination, SearchParam } from '../../../models/shared';

interface IQuestionData {
  id?: number,
  question: string
  answer: string
  createdBy?: string
  roleGroupId?: number
  updatedBy?: string
}


function QuestionAnswerPage(): JSX.Element {
  const navigate = useNavigate();

  const typingTimer = useRef<NodeJS.Timeout>();

  const [is401, setIs401] = useState(false);
  
  const [searchParam, setSearchParam] = useState<SearchParam>({
    page: 1,
    size: 5,
    keyword: ''
  });

  const [data, setData] = useState<IPagination<IQuestionData>>({
    start: 0,
    end: 0,
    length: 0,
    data: []
  });

  useEffect(() => {
    if (searchParam) {
      toast.promise(
        getQuestionData(
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

  const goToEdit = (state: IQuestionData) => {
    navigate('/dashboard/question_answer/edit', {
      state
    });
  }

  const goToDelete = (state: IQuestionData) => {
    navigate('/dashboard/question_answer/delete', {
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
          <h1 className="text-dark mb-1">Pertanyaan &amp; Jawaban</h1>
        </div>
        <div className="col-md-6 col-xl-2">
          <button className="btn btn-primary" type="button" onClick={() => navigate('/dashboard/question_answer/new')}>
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
    
      <div className="row mt-5">
        <div className="col-xl-12">

          {
            data.data.map(item => (
              <div className="card border-0 shadow todo-card" style={{"width": "100%"}}>
                <div className="card-header d-flex align-items-center bg-primary row p-1">
                  <a className="mw-85 d-flex align-items-center text-light col-sm-10" 
                    data-bs-toggle="collapse" 
                    href="#cardMyCollapse" 
                    aria-expanded="false" 
                    aria-controls="cardMyCollapse" 
                    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                    // @ts-ignore
                    onclick="rotate(this)">
                      <svg 
                        xmlns="http://www.w3.org/2000/svg" 
                        width="1em" 
                        height="1em" 
                        fill="currentColor" 
                        viewBox="0 0 16 16" 
                        className="bi bi-chevron-down fs-6 rotate me-1">
                        <path fillRule="evenodd" d="M1.646 4.646a.5.5 0 0 1 .708 0L8 10.293l5.646-5.647a.5.5 0 0 1 .708.708l-6 6a.5.5 0 0 1-.708 0l-6-6a.5.5 0 0 1 0-.708z"></path>
                      </svg>
                      <h4 className="m-0">
                        {item.question}
                      </h4>
                  </a>
                  <div className="col-xl-2 text-end">
                    <button className="btn btn-success m-1" type="button" onClick={() => goToEdit(item)}>
                      <i className="far fa-edit"></i>
                    </button>
                    &nbsp;
                    <button className="btn btn-danger m-1" type="button" onClick={() => goToDelete(item)}>
                      <i className="far fa-trash-alt"></i>
                    </button>
                  </div>
                </div>
  
                <div className="card-body collapse show row p-1" id="cardMyCollapse">
                  <ul className="list-group list-group-flush">
                    <li className="list-group-item d-flex align-items-center">
                      <p>
                        {item.answer}
                      </p>
                    </li>
                  </ul>
                </div>
              </div>
            ))
          }

        </div>
      </div>

      <div className='row mt-5'>
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
  )
}

export default QuestionAnswerPage;