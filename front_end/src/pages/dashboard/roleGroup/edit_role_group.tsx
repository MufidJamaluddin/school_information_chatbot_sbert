import { FormEvent } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { updateRoleGroup } from '../../../models/roleGroup';
import { toast } from "react-toastify";


function EditRoleGroupPage(): JSX.Element {
  const { 
    state: { 
      id: idValue, 
      roleGroup: roleGroupOldValue
    } = {} 
  } = useLocation();

  const navigate = useNavigate();

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const { roleGroup } = e.currentTarget;

    toast.promise(
      updateRoleGroup(
        idValue,
        {
          roleGroup: roleGroup?.value,
        }
      ),
      {
        success: {
          render(){
            return 'Perubahan Role Group Berhasil Disimpan!'
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
          <h1 className="text-dark mb-1">Ubah Role Group</h1>
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
                defaultValue={roleGroupOldValue}
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

export default EditRoleGroupPage;