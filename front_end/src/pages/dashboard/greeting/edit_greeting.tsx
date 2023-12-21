import { FormEvent } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { updateGreeting } from '../../../models/greeting';
import { toast } from "react-toastify";


function EditGreetingPage(): JSX.Element {
  const { 
    state: { 
      id: idValue, 
      greeting: greetingOldValue, 
      startTime: startTimeOldValue, 
      endTime: endTimeOldValue
    } = {} 
  } = useLocation();

  const navigate = useNavigate();

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const {
      greeting,
      startTime,
      endTime
    } = e.currentTarget;

    toast.promise(
      updateGreeting(
        idValue,
        {
          greeting: greeting?.value,
          startTime: startTime?.value,
          endTime: endTime?.value
        }
      ),
      {
        success: {
          render(){
            return 'Perubahan Salam Berhasil Disimpan!'
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
          <h1 className="text-dark mb-1">Ubah Salam</h1>
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
                defaultValue={greetingOldValue}
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
                defaultValue={startTimeOldValue}
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
                defaultValue={endTimeOldValue}
                required
              />
            </div>
            
            <button type="submit" className="btn btn-info">Ubah</button>
          </form>
        </div>
      </div>
    </div>
  )
}

export default EditGreetingPage;