import axios, { AxiosHeaders } from "axios";
import { getHeaders } from "../session";
import { IPagination } from "../shared";
import { parse } from "content-range";

interface IUpdateRoleGroupData {
  roleGroup: string
}

interface IRoleGroupData {
  id?: number
  roleGroup: string
  createdBy?: string
  roleGroupId?: number
  updatedBy?: string
}

interface IChangeData { 
  id: number 
  message: string
}

/**
 * Role Group Data
 * 
 * @returns 
 */
export async function getRoleGroupData(start: number, size: number, keyword: string): Promise<IPagination<IRoleGroupData>> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  const data = await axios.get<{ 
    params: {
      start: number, 
      size: number, 
      keyword: string
    }
  }, { 
    data: Array<IRoleGroupData>,
    headers: AxiosHeaders
  }>(
    `${baseApiPath}/role-group`,
    {
      headers,
      params: {
        start,
        size,
        keyword
      }
    }
  );

  const meta = parse(data.headers['content-range']);

  return {
    start: meta?.start || 0,
    end: meta?.end || 0,
    length: meta?.size || 0,
    data: data.data
  };
}

/**
 * Save New Role Group
 * 
 * @param newRoleGroup
 *  
 * @returns 
 */
export function saveNewRoleGroup(newRoleGroup: IRoleGroupData): Promise<{ data: IChangeData }> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  return axios.post<IRoleGroupData, { data: IChangeData }>(
    `${baseApiPath}/role-group`,
    newRoleGroup,
    {
      headers
    }
  );
}

/**
 * Update Role Group
 * 
 * @param roleGroupId 
 * @param updatedData
 *  
 * @returns 
 */
export function updateRoleGroup(roleGroupId: number, updatedData: IUpdateRoleGroupData): Promise<{ data: IChangeData }> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  return axios.put<IUpdateRoleGroupData, { data: IChangeData }>(
    `${baseApiPath}/role-group/${roleGroupId}`,
    updatedData,
    {
      headers
    }
  );
}
