import axios, { AxiosHeaders } from "axios";
import { getHeaders } from "../session";
import { IPagination } from "../shared";
import { parse } from "content-range";

interface IUpdateGreetingData {
  greeting: string
  startTime: string
  endTime: string
}

interface IGreetingData {
  id?: number,
  greeting: string
  startTime: string
  endTime: string
  updatedAt?: string
  updatedBy?: string
  createdAt?: string
  createdBy?: string
}

interface IChangeData { 
  questionId: number 
  message: string
}

/**
 * Greeting Data
 * 
 * @returns 
 */
export async function getGreetingData(start: number, size: number, keyword: string): Promise<IPagination<IGreetingData>> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  const data = await axios.get<{ 
    params: {
      start: number, 
      size: number, 
      keyword: string
    }
  }, { 
    data: Array<IGreetingData>
    headers: AxiosHeaders
  }>(
    `${baseApiPath}/greeting`,
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
 * Save New Greeting
 * 
 * @param newGreeting
 *  
 * @returns 
 */
export function saveNewGreeting(newGreeting: IGreetingData): Promise<{ data: IChangeData }> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  return axios.post<IGreetingData, { data: IChangeData }>(
    `${baseApiPath}/greeting`,
    newGreeting,
    {
      headers
    }
  );
}

/**
 * Update Greeting
 * 
 * @param greetingId 
 * @param updatedData
 *  
 * @returns 
 */
export function updateGreeting(greetingId: number, updatedData: IUpdateGreetingData): Promise<{ data: IChangeData }> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  return axios.put<IUpdateGreetingData, { data: IChangeData }>(
    `${baseApiPath}/greeting/${greetingId}`,
    updatedData,
    {
      headers
    }
  );
}

/**
 * Delete Greeting
 * 
 * @param greetingId 
 * 
 * @returns 
 */
export function deleteGreeting(greetingId: number): Promise<{ data: IChangeData }> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  return axios.delete<null, { data: IChangeData }>(
    `${baseApiPath}/greeting/${greetingId}`,
    {
      headers
    }
  );
}
