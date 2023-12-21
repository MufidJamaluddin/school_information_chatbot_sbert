import { FormEvent } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { deleteGreeting } from '../../../models/greeting';
import { toast } from "react-toastify";


function DeleteGreetingPage(): JSX.Element {
  const { 
    state: { 
      id: idValue, 
      greeting: greetingOldValue, 
      startTime: startTimeOldValue, 
      endTime: endTimeOldValue
    } = {} 
  } = useLocation();

  const navigate = useNavigate();

  const handleDelete = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    toast.promise(
      deleteGreeting(
        idValue
      ),
      {
        success: {
          render(){
            return `Salam '${greetingOldValue}' Berhasil Dihapus Permanen!`
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
          <h1 className="text-dark mb-1">Apakah Anda Yakin Ingin Menghapus Permanen Salam Ini?</h1>
        </div>

        <hr/>

        <div className="col-md-12">
          <form onSubmit={handleDelete}>
            <div className="mb-3">
              <label htmlFor="greeting" className="form-label">Salam</label>
              <input
                type="text"
                className="form-control"
                id="greeting"
                name="greeting"
                value={greetingOldValue}
                disabled
              />
            </div>

            <div className="mb-3">
              <label htmlFor="startTime" className="form-label">Waktu Mulai</label>
              <input
                type="text"
                className="form-control"
                id="startTime"
                name="startTime"
                value={startTimeOldValue}
                disabled
              />
            </div>

            <div className="mb-3">
              <label htmlFor="endTime" className="form-label">Waktu Selesai</label>
              <input
                type="text"
                className="form-control"
                id="endTime"
                name="endTime"
                value={endTimeOldValue}
                disabled
              />
            </div>
            
            <button type="submit" className="btn btn-danger">Hapus Permanen</button>
          </form>
        </div>
      </div>
    </div>
  )
}

export default DeleteGreetingPage;