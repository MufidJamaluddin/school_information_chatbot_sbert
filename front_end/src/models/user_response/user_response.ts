import axios from "axios";
import { getHeaders } from "../session";

interface IUserResponse {
  userId: number;
  question: string;
  answer: string;
  score: number;
}

interface IChangeData { 
  message: string
}

/**
 * Save New User Answer
 * 
 * @param newUserAnswer
 *  
 * @returns 
 */
export function saveNewUserResponse(newUser: IUserResponse): Promise<{ data: IChangeData }> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  return axios.post<IUserResponse, { data: IChangeData }>(
    `${baseApiPath}/user-response`,
    newUser,
    {
      headers
    }
  );
}
