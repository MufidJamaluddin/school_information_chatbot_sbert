import { FormEvent } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { deleteAbbreviation } from '../../../models/abbreviation';
import { toast } from "react-toastify";


function DeleteAbbreviationPage(): JSX.Element {
  const { 
    state: { 
      standardWord: standardWordValue,
    } = {} 
  } = useLocation();

  const navigate = useNavigate();

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const { standardWord } = e.currentTarget;

    toast.promise(
      deleteAbbreviation(standardWordValue),
      {
        success: {
          render(){
            return `Abbreviation '${standardWord?.value}' Berhasil Disimpan!`
          }
        },
        error: {
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-ignore
          render({ data: { response: { data = 'Koneksi anda terputus!' } = {} } = {} }){
            return data;
          }
        }
      }
    ).then(() => navigate('/dashboard/abbreviation'));
  }
  
  return (
    <div className="container">
      <div className="row">
        <div className="col-md-12">
          <h1 className="text-dark mb-1">Ubah Sinonim</h1>
        </div>

        <hr/>

        <div className="col-md-12">
          <form onSubmit={handleSubmit}>
            <div className="mb-3">
              <label htmlFor="standardWord" className="form-label">Kata Baku Standar</label>
              <input
                type="text"
                className="form-control"
                id="standardWord"
                name="standardWord"
                value={standardWordValue}
                disabled
              />
            </div>
            
            <button type="submit" className="btn btn-danger">Delete</button>
          </form>
        </div>
      </div>
    </div>
  )
}

export default DeleteAbbreviationPage;