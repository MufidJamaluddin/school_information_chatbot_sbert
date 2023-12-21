import axios from "axios";


/**
 * Get Answer
 * 
 * @param question string
 * @returns 
 */
export async function getAnswer(question: string) {
  const baseApiPath = import.meta.env.VITE_BASE_API_PATH || '';

  const { data: { answer } = {} } = await axios.get<{ question: string }, { data: { answer: string } }>(
    `${baseApiPath}/answer`,
    {
      method: 'GET',
      params: {
        question,
      }
    }
  );

  return answer;
}

