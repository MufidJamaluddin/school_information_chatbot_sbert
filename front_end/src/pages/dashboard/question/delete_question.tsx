import { FormEvent } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { deleteQuestion } from '../../../models/question';
import { toast } from "react-toastify";


function DeleteQuestionAnswerPage(): JSX.Element {
  const { 
    state: { 
      id: idValue, 
      question: questionOldValue,
      answer: answerOldValue
    } = {} 
  } = useLocation();

  const navigate = useNavigate();

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    toast.promise(
      deleteQuestion(
        idValue
      ),
      {
        success: {
          render(){
            return `Pertanyaan '${questionOldValue}' Berhasil Dihapus Permanen!`
          }
        },
        error: {
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-ignore
          render({ data: { response: { data = 'Koneksi anda terputus!' } = {} } = {} }){
            return data;
          }
        }
      }
    ).then(() => navigate('/dashboard/question_answer'));
  }
  
  return (
    <div className="container">
      <div className="row">
        <div className="col-md-12">
          <h1 className="text-dark mb-1">Apakah Anda Yakin Ingin Menghapus Pertanyaan & Jawaban Ini?</h1>
        </div>

        <hr/>

        <div className="col-md-12">
          <form onSubmit={handleSubmit}>
            <div className="mb-3">
              <label htmlFor="question" className="form-label">Pertanyaan</label>
              <input
                type="text"
                className="form-control"
                id="question"
                name="question"
                value={questionOldValue}
                disabled
              />
            </div>

            <div className="mb-3">
              <label htmlFor="answer" className="form-label">Jawaban</label>
              <input
                type="text"
                className="form-control"
                id="answer"
                name="answer"
                value={answerOldValue}
                disabled
              />
            </div>
            
            <button type="submit" className="btn btn-danger">Hapus Permanen</button>
          </form>
        </div>
      </div>
    </div>
  )
}

export default DeleteQuestionAnswerPage;