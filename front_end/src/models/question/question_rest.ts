import axios, { AxiosHeaders } from "axios";
import { parse } from "content-range";
import { getHeaders } from "../session";
import { IPagination } from "../shared";

interface IUpdateQuestionData {
  question: string
  answer: string
}

interface IQuestionData {
  id?: number,
  question: string
  answer: string
  createdBy?: string
  roleGroupId?: number
  updatedBy?: string
}

interface IChangeData { 
  questionId: number 
  message: string
}

/**
 * Question Data
 * 
 * @returns 
 */
export async function getQuestionData(start: number, size: number, keyword: string): Promise<IPagination<IQuestionData>> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  const data = await axios.get<{ 
    params: {
      start: number, 
      size: number, 
      keyword: string
    }
  }, { 
    data: Array<IQuestionData>,
    headers: AxiosHeaders
  }>(
    `${baseApiPath}/question`,
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
 * Save New Question
 * 
 * @param newQuestion
 *  
 * @returns 
 */
export function saveNewQuestion(newQuestion: IQuestionData): Promise<{ data: IChangeData }> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  return axios.post<IQuestionData, { data: IChangeData }>(
    `${baseApiPath}/question`,
    newQuestion,
    {
      headers
    }
  );
}

/**
 * Update Question
 * 
 * @param questionId 
 * @param updatedData
 *  
 * @returns 
 */
export function updateQuestion(questionId: number, updatedData: IUpdateQuestionData): Promise<{ data: IChangeData }> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  return axios.put<IUpdateQuestionData, { data: IChangeData }>(
    `${baseApiPath}/question/${questionId}`,
    updatedData,
    {
      headers
    }
  );
}

export function deleteQuestion(questionId: number): Promise<{ data: IChangeData }> {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';
  const headers = getHeaders();

  return axios.delete<null, { data: IChangeData }>(
    `${baseApiPath}/question/${questionId}`,
    {
      headers
    }
  );
}
