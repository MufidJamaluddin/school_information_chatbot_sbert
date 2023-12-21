import { FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import { saveNewGreeting } from '../../../models/greeting';
import { toast } from "react-toastify";


function AddGreetingPage(): JSX.Element {
  const navigate = useNavigate();

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const {
      greeting,
      startTime,
      endTime
    } = e.currentTarget;

    toast.promise(
      saveNewGreeting({
        greeting: greeting?.value,
        startTime: startTime?.value,
        endTime: endTime?.value
      }),
      {
        success: {
          render(){
            return 'Salam Berhasil Disimpan!'
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
    ).then(() => navigate('/dashboard/greeting'));
  }
  
  return (
    <div className="container">
      <div className="row">
        <div className="col-md-12">
          <h1 className="text-dark mb-1">Tambah Salam</h1>
        </div>

        <hr/>

        <div className="col-md-12">
          <form onSubmit={handleSubmit}>
            <div className="mb-3">
              <label htmlFor="greeting" className="form-label">Salam</label>
              <input
                type="text"
                className="form-control"
                id="greeting"
                name="greeting"
                required
              />
            </div>

            <div className="mb-3">
              <label htmlFor="startTime" className="form-label">Waktu Mulai</label>
              <input
                type="text"
                className="form-control"
                id="startTime"
                name="startTime"
                required
              />
            </div>

            <div className="mb-3">
              <label htmlFor="endTime" className="form-label">Waktu Selesai</label>
              <input
                type="text"
                className="form-control"
                id="endTime"
                name="endTime"
                required
              />
            </div>
            
            <button type="submit" className="btn btn-primary">Tambah</button>
          </form>
        </div>
      </div>
    </div>
  )
}

export default AddGreetingPage;