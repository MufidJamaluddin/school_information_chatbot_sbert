import { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";
import { saveNewAbbreviation } from '../../../models/abbreviation';
import { toast } from "react-toastify";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import { ChipsInput } from 'chips-input-lib';


function AddAbbreviationPage(): JSX.Element {
  const navigate = useNavigate();
  const [abbreviations, setAbbreviations] = useState<Array<string>>([]);

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const { standardWord } = e.currentTarget;

    toast.promise(
      saveNewAbbreviation({
        standardWord: standardWord?.value,
        listAbbreviationTerm: abbreviations,
      }),
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
          <h1 className="text-dark mb-1">Tambah Sinonim</h1>
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
                required
              />
            </div>

            <div className="mb-3">
              <ChipsInput 
                onAddChips={setAbbreviations} 
                placeholder="List Kata-Kata Sinonim"/>
            </div>
            
            <button type="submit" className="btn btn-primary">Tambah</button>
          </form>
        </div>
      </div>
    </div>
  )
}

export default AddAbbreviationPage;