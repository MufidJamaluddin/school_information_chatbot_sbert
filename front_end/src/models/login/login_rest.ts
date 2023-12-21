import axios from "axios";
import { setData } from "../session";

interface ResponseData { 
  data: { 
    token: string, 
    userData: { 
      fullName: string 
    } 
  } 
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
