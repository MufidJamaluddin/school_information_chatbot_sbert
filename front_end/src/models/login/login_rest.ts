import axios from "axios";
import { setData, clearData } from "../session";

interface ResponseData { 
  data: { 
    token: string, 
    userData: { 
      fullName: string 
    } 
  } 
}

interface IChangeData { 
  message: string
}

/**
 * Login Action
 * 
 * @param username string
 * @param password string
 * @returns 
 */
export async function loginAction(username: string, password: string): Promise<ResponseData> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';

  const data = await axios.post<{ 
    username: string, 
    password: string
  }, ResponseData>(
    `${baseApiPath}/login`,
      {
        username,
        password
      }
  );

  setData(data?.data?.token, data?.data?.userData);
  
  return data;
}

/**
 * Logout Action
 * 
 * @param username string
 * @param password string
 * @returns 
 */
export async function logoutAction(): Promise<IChangeData> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';

  const data = await axios.delete<void, { data: IChangeData }>(
    `${baseApiPath}/login`,
  );

  clearData();

  return data.data;
}