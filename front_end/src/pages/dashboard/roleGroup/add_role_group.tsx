import { FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import { saveNewRoleGroup } from '../../../models/roleGroup';
import { toast } from "react-toastify";


function AddRoleGroupPage(): JSX.Element {
  const navigate = useNavigate();

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const { roleGroup } = e.currentTarget;

    toast.promise(
      saveNewRoleGroup({
        roleGroup: roleGroup?.value
      }),
      {
        success: {
          render(){
            return `Role Group '${roleGroup?.value}' Berhasil Disimpan!`
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
    ).then(() => navigate('/dashboard/role_group'));
  }
  
  return (
    <div className="container">
      <div className="row">
        <div className="col-md-12">
          <h1 className="text-dark mb-1">Tambah Role Group</h1>
        </div>

        <hr/>

        <div className="col-md-12">
          <form onSubmit={handleSubmit}>
            <div className="mb-3">
              <label htmlFor="roleGroup" className="form-label">Role Group</label>
              <input
                type="text"
                className="form-control"
                id="roleGroup"
                name="roleGroup"
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

export default AddRoleGroupPage;