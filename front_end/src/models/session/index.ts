interface UserData { 
  fullName: string 
} 

export interface IChatUser {
  id: number;
  fullName: string;
  isStudent: boolean;
  className: string;
  age: number;
}

export function setData(token: string, userData: UserData) {
  sessionStorage.setItem('token', token);
  sessionStorage.setItem('user-data', JSON.stringify(userData));
}

export function setChatUserData(chatUserData: IChatUser) {
  localStorage.setItem('chat-user', JSON.stringify(chatUserData));
}

export function getChatUserData(): IChatUser {
  const data = localStorage.getItem('chat-user');

  return data ? JSON.parse(data) : null;
}

export function getHeaders() {
  return {
    Authorization: `${sessionStorage.getItem('token')}`
  }
}

export function getUserData(): UserData|null {
  const data = sessionStorage.getItem('user-data');

  return data ? JSON.parse(data) : null;
}

export function clearData() {
  sessionStorage.clear();
}
