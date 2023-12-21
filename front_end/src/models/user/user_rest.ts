import axios from "axios";
import { IChatUser, getHeaders, setChatUserData } from "../session";

interface IUser {
  id?: number;
  fullName: string;
  userRole: string;
  className: string;
  age: number;
}

/**
 * Save New User
 * 
 * @param newUser
 *  
 * @returns 
 */
export async function saveNewUser(newUser: IUser): Promise<{ data: IChatUser }> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  const newData = await axios.post<IUser, { data: IChatUser }>(
    `${baseApiPath}/user`,
    newUser,
    {
      headers
    }
  );

  if (newData?.data) {
    setChatUserData(newData.data);
  }

  return newData;
}
