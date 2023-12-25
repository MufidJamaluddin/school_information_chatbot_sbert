import { JSX, useState, FormEvent, useEffect, useRef } from 'react';
import { dateIDFormatter } from '../constant/date';
import { getAnswer } from '../models/answer';
import { saveNewUser } from '../models/user';
import { toast } from 'react-toastify';
import { getChatUserData } from '../models/session';
import StarRatings from 'react-star-ratings';
import { saveNewUserResponse } from '../models/user_response';

import './chat.style.css';

interface ChatItem {
  isAnswer: boolean;
  data: string;
  createdAt: Date;
  userRate?: number;
}

function Rating(
  data: Readonly<{ 
    question: string, 
    answer: string, 
    userRate?: number; 
    idx: number;
    markUserRate: (idx: number, userRate: number) => void,
  }>,
): JSX.Element {
  const { question, answer, userRate, idx, markUserRate } = data;
  const [rating, setRating] = useState(userRate ?? 0);

  const changeRating = function (newRating: number) {
    if ((userRate || 0) > 0) {
      toast(
        <>
          <b>Coba Tanya & Beri Nilai Lagi!</b>
          <hr/>
          <p>Nilai yang telah diberikan tidak dapat diubah!</p>
        </>
      );
      return;
    }

    const chatUserData = getChatUserData();

    if (chatUserData?.id) {
      toast.promise(
        saveNewUserResponse({
          userId: chatUserData.id,
          question,
          answer,
          score: newRating,
        }).then(() => {
          markUserRate(idx, newRating);
          setRating(newRating);
        }),
        {
          error: {
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-ignore
            render({ data: { response } }){
              return `Simpan Feedback Rating Gagal: ${typeof response === 'object' ? response?.data : response}`
            }
          }
        }
      );
    } else {
      toast.error('Simpan Feedback Rating Gagal, Mohon Reload Halaman');
    }
  }

  return (
    <StarRatings
      rating={rating}
      starRatedColor="blue"
      changeRating={changeRating}
      numberOfStars={5}
      name='rating'
    />
  );
}

function createChatItem(
  { isAnswer, data, userRate, createdAt }: ChatItem, 
  idx: number, 
  arr: Array<ChatItem>,
  markUserRate: (idx: number, userRate: number) => void,
): JSX.Element {
  const formattedCreatedAt = createdAt ? dateIDFormatter.format(createdAt) : '';

  if (isAnswer) {
    return (
      <div className="media ml-auto mb-3 row" key={idx}>
        <div className="media-body col-sm-11">
          <div className="bg-primary rounded py-2 px-3 mb-2">
            <p className="text-small mb-0 text-white">
              <div
                dangerouslySetInnerHTML={{__html: data}}
              />
            </p>

            {
              isAnswer && (
                <Rating 
                  question={arr[idx - 1].data} 
                  idx={idx} 
                  answer={data} 
                  userRate={userRate}
                  markUserRate={markUserRate} />
              )
            }
          </div>
          <p className="small text-muted">
            {formattedCreatedAt}
          </p>
        </div>
        
        <div className="col-sm-1">
          <img 
            src="https://res.cloudinary.com/mhmd/image/upload/v1564960395/avatar_usae7z.svg" 
            alt="user" 
            width="50" 
            className="rounded-circle"/>
        </div>
      </div>
    );
  }

  return (
    <div className="media mb-3 row" key={idx}>
      <div className="col-sm-1">
        <img 
          src="https://res.cloudinary.com/mhmd/image/upload/v1564960395/avatar_usae7z.svg" 
          alt="user" 
          width="50" 
          className="rounded-circle"/>
      </div>

      <div 
        className="media-body ml-3 col-sm-11">
        <div 
          className="bg-light rounded py-2 px-3 mb-2">
          <p className="text-small mb-0 text-muted">
            {data}
          </p>
        </div>
        <p className="small text-muted">
          {formattedCreatedAt}
        </p>
      </div>
    </div>
  );
}

function ChatPage(): JSX.Element {
  const focusItemRef = useRef<HTMLDivElement>(null);
  const modalItemRef = useRef<HTMLDivElement>(null);

  const [chats, setChats] = useState<Array<ChatItem>>([]);
  const [showClass, setShowClass] = useState(true);
  const [showCouse, setShowCourse] = useState(false);
  
  const chatUserData = getChatUserData();
  const noChatUserData = !chatUserData;

  const changeRoleAct = ((chEv: React.ChangeEvent<HTMLSelectElement>) => {
    switch (chEv.target.value) {
      case 'student':
        setShowClass(true);
        setShowCourse(false);
        return;

      case 'teacher':
        setShowClass(false);
        setShowCourse(true);
        return;

      default:
        setShowClass(false);
        setShowCourse(false);
        return;
    }
  })

  useEffect(() => {
    try {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      const registrasiModal = new bootstrap.Modal(modalItemRef.current, {
        keyboard: false,
        backdrop: 'static',
      });

      if (noChatUserData) {
        registrasiModal.show();
      } else {
        registrasiModal.hide();
      }
    } catch (e) {
      console.error('Error launch registration modal', e);
    }
  }, [noChatUserData]);

  useEffect(() => {
    async function requestAnswer() {
      let answer;
      const chatLast = chats[chats.length - 1];

      if (!chatLast) {
        const lastData = JSON.parse(localStorage.getItem('chat-data') ?? '[]');

        if (Array.isArray(lastData)) {
          const lastDataTransformed = lastData.filter(
            item => 
              typeof item.data === 'string' && 
              typeof item.createdAt === 'string' && 
              typeof item.isAnswer === 'boolean'
          ).map(item => ({
            userRate: item.userRate,
            isAnswer: item.isAnswer, 
            data: item.data, 
            createdAt: new Date(item.createdAt)
          }));

          if (lastDataTransformed.length) {
            setChats(lastDataTransformed);
          }
        }

        return;
      }

      if (chatLast.isAnswer) {
        focusItemRef?.current?.focus();

        if (chatLast.userRate) {
          localStorage.setItem('chat-data', JSON.stringify(chats)); 
        }
        return;
      }

      if (!String(chatLast.data).toLowerCase().includes('belum jelas')) {
        try {
          answer = await getAnswer(chatLast.data, chatUserData.fullName);
        } catch (e) {
          answer = 'Jaringan tidak tersedia. Mohon coba lagi nanti.'
        }
      }

      if (!answer) {
        answer = null;
      }

      const newChats = [
        ...chats,
        {
          isAnswer: true,
          data: answer ?? 'Jawaban tidak ditemukan.',
          createdAt: new Date()
        }
      ];

      localStorage.setItem('chat-data', JSON.stringify(newChats)); 
      setChats(newChats);
    }

    requestAnswer();
  }, [chats, chatUserData]);

  const markUserRate = (idx: number, userRate: number) => {
    const mChats = [...chats];
    mChats[idx] = { ...mChats[idx], userRate };

    console.debug(mChats[idx]);

    setChats(mChats);
  };

  const clearAllChat = () => {
    localStorage.removeItem('chat-data');
    setChats([]);
  };

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const question = e.currentTarget.question.value;

    setChats([
      ...chats,
      {
        isAnswer: false,
        data: question,
        createdAt: new Date()
      }
    ]);

    e.currentTarget.reset();
  }

  const handleInitIdentity = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const {
      fullName,
      mClassName,
      mUserRole,
      age,
    } = e.currentTarget;

    toast.promise(
      saveNewUser({
        fullName: fullName?.value,
        userRole: mUserRole?.value,
        className: mClassName?.value,
        age: Number(age?.value),
      }).then(() => {
        try {
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-ignore
          const registrasiModal = bootstrap.Modal.getOrCreateInstance(modalItemRef.current);

          const backdrops = document.getElementsByClassName('modal-backdrop');
          if (backdrops?.length) {
            for (let i = 0; i < backdrops.length; i++) {
              // eslint-disable-next-line @typescript-eslint/ban-ts-comment
              // @ts-ignore
              backdrops.item(i).hidden = true;
            }
          }

          registrasiModal.hide();
        } catch (e) {
          console.error(e);
        }

        setChats([]);
      }),
      {
        error: {
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-ignore
          render({ data: { response } }){
            return `Mulai Chat Gagal: ${typeof response === 'object' ? response?.data : response || 'Masalah Koneksi Internet'}`
          }
        }
      }
    );

    e.currentTarget.reset();
  }

  return (
    <div className="container py-3" style={{ height: '98%' }}>
      <header className="text-center">
        <h1 className="display-4 text-white">Tanya Sabri | SMAN Situraja</h1>
      </header>

      <div 
        ref={modalItemRef}
        className="modal" 
        id="registrasiModal" 
        tabIndex={-1} 
        data-bs-backdrop="static"
        aria-labelledby="registrasiModalLabel">

        <div className="modal-dialog modal-dialog-centered">
          <div className="modal-content">
            <div className="modal-header">
              <h5 className="modal-title" id="registrasiModalLabel">
                Kenalan dengan Kak Sabri!
              </h5>
            </div>

            <form onSubmit={handleInitIdentity} className="bg-light">
              <div className="modal-body">
                <p>
                  Perkenalkan, Sabri merupakan chat robot (chatbot) baru yang tahu tentang informasi terkait sejarah pendirian SMAN Situraja, program-program kesiswaan dan humas. 
                
                  <br/>
                  <b>Tanya yuk, dan beri nilai!</b>
                </p>

                <div className="input-group row">
                  <div className='col-sm-12'>
                    <input 
                      name='fullName'
                      type="text" 
                      placeholder="Tulis Nama Lengkap Anda Disini!"
                      required={true}
                      className="form-control rounded-0 border-0 py-4 bg-light"/>
                  </div>

                  <div className='col-sm-12'>
                    <select className="form-select" name="mUserRole" defaultValue="student" aria-label="Peran" onChange={changeRoleAct}>
                      <option value="student">Siswa</option>
                      <option value="teacher">Guru</option>
                      <option value="alumni">Alumni</option>
                      <option value="public">Masyarakat Umum</option>
                    </select>
                  </div>

                  {
                    showClass && (
                      <div className='col-sm-12'>
                        <input 
                          name='mClassName'
                          type="text" 
                          placeholder="Mohon Masukkan Kelas Anda Disini!"
                          required={true}
                          className="form-control rounded-0 border-0 py-4 bg-light"/>
                      </div>
                    )
                  }

                  {
                    showCouse && (
                      <div className='col-sm-12'>
                        <input 
                          name='mClassName'
                          type="text" 
                          placeholder="Mata Pelajaran"
                          required={true}
                          className="form-control rounded-0 border-0 py-4 bg-light"/>
                      </div>
                    )
                  }

                  <div className='col-sm-12'>
                    <input 
                      name='age'
                      type="text" 
                      placeholder="Mohon Isi Umur Anda Disini!"
                      required={true}
                      className="form-control rounded-0 border-0 py-4 bg-light"/>
                  </div>
                </div>

              </div>

              <div className="modal-footer">
                <button type="submit" className="btn btn-primary" data-dismiss="modal">
                  Mulai Tanya Kak Sabri
                </button>
              </div>

            </form>
          </div>
        </div>
      </div>

      <div className="col-12" style={{ height: '96%' }}>
        <button className="btn btn-primary" onClick={clearAllChat}> 
          <i className="fa fa-trash-alt"></i> Clear All
        </button>

        <div className="px-4 py-5 chat-box bg-white" style={{ height: '85%' }}>
          {chats.map((data, idx, arr) => createChatItem(data, idx, arr, markUserRate))}

          <div ref={focusItemRef} tabIndex={-1}/>
        </div>

        <form onSubmit={handleSubmit} className="bg-light">
          <div className="input-group">
            <input 
              name='question'
              type="text" 
              placeholder="Tulis Pertanyaan Anda, Lalu Klik Kirim!" 
              aria-describedby="button-addon2" 
              className="form-control rounded-0 border-0 py-4 bg-light"/>

            <div className="input-group-append">
              <button id="button-addon2" type="submit" className="btn btn-link"> 
                <i className="fa fa-paper-plane"></i>
              </button>
            </div>
          </div>
        </form>

      </div>
    </div>
  )
}

export default ChatPage;