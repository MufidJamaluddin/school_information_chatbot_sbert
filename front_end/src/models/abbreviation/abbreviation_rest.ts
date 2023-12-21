import axios, { AxiosHeaders } from "axios";
import { getHeaders } from "../session";
import { parse } from "content-range";
import { IPagination } from "../shared";

interface IUpdateQuestionData {
  standardWord: string
  listAbbreviationTerm: Array<string>
}

interface IAbbreviationData {
  standardWord: string
  listAbbreviationTerm: Array<string>
  updatedAt?: string
  updatedBy?: string
  createdAt?: string
  createdBy?: string
}

interface IChangeData { 
  message: string
}

/**
 * Abbreviation Data
 * 
 * @returns 
 */
export async function getAbbreviationData(start: number, size: number, keyword: string): Promise<IPagination<IAbbreviationData>> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  const data = await axios.get<{ 
    params: {
      start: number, 
      size: number, 
      keyword: string
    }
  }, { 
    data: Array<IAbbreviationData>,
    headers: AxiosHeaders
  }>(
    `${baseApiPath}/abbreviation`,
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
    start: meta?.start ?? 0,
    end: meta?.end ?? 0,
    length: meta?.size ?? 0,
    data: data.data
  };
}

/**
 * Save New Abbreviation
 * 
 * @param newAbbreviation
 *  
 * @returns 
 */
export function saveNewAbbreviation(newAbbreviation: IAbbreviationData): Promise<{ data: IChangeData }> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  return axios.post<IAbbreviationData, { data: IChangeData }>(
    `${baseApiPath}/abbreviation`,
    newAbbreviation,
    {
      headers
    }
  );
}

/**
 * Update Abbreviation
 * 
 * @param updatedData
 *  
 * @returns 
 */
export function updateAbbreviation(updatedData: IUpdateQuestionData): Promise<{ data: IChangeData }> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  return axios.put<IUpdateQuestionData, { data: IChangeData }>(
    `${baseApiPath}/abbreviation`,
    updatedData,
    {
      headers
    }
  );
}

/**
 * Delete Abbreviation
 * 
 * @param standardWord 
 * 
 * @returns 
 */
export function deleteAbbreviation(standardWord: string): Promise<{ data: IChangeData }> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  return axios.delete<null, { data: IChangeData }>(
    `${baseApiPath}/abbreviation/${standardWord}`,
    {
      headers
    }
  );
}
