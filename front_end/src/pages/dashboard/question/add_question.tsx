import { FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import { saveNewQuestion } from '../../../models/question';
import { toast } from "react-toastify";


function AddQuestionAnswerPage(): JSX.Element {
  const navigate = useNavigate();

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const {
      question,
      answer
    } = e.currentTarget;

    toast.promise(
      saveNewQuestion({
        question: question?.value,
        answer: answer?.value
      }),
      {
        success: {
          render(){
            return 'Pertanyaan Berhasil Disimpan!'
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
          <h1 className="text-dark mb-1">Tambah Pertanyaan &amp; Jawaban</h1>
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
                required
              />
            </div>

            <div className="mb-3">
              <label htmlFor="answer" className="form-label">Jawaban</label>
              <input
                type="text"
                className="form-control"
                id="answer"
                name="answer"
                required
              />
            </div>
            
            <button type="submit" className="btn btn-primary">Tambah</button>
          </form>
        </div>
      </div>
    </div>
  )
}

export default AddQuestionAnswerPage;