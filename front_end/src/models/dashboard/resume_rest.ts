import axios from "axios";
import { getHeaders } from "../session";

interface IDashboardResumeData {
  totalAdmin: number;
  totalGreeting: number;
  totalQuestion: number;
  totalUser: number;
}

/**
 * Resume Rest
 * 
 * @returns 
 */
export async function getDashboardResume() {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  return axios.get<{ 
    username: string, 
    password: string
  }, { 
    data: IDashboardResumeData
  }>(
    `${baseApiPath}/dashboard-resume`,
    {
      headers
    }
  );
}
