import { JSX  } from 'react';

import { RouteObject } from 'react-router-dom';

import DashboardLayout from './layout';
import DashboardMainPage from './dashboard';

import { 
  RoleGroupPage, 
  AddRoleGroupPage, 
  EditRoleGroupPage 
} from './roleGroup';

import { 
  GreetingPage, 
  AddGreetingPage, 
  EditGreetingPage, 
  DeleteGreetingPage 
} from './greeting';

import { 
  QuestionAnswerPage, 
  AddQuestionAnswerPage, 
  UpdateQuestionAnswerPage, 
  DeleteQuestionAnswerPage 
} from './question';

import { 
  AbbreviationPage, 
  AddAbbreviationPage, 
  DeleteAbbreviationPage, 
  EditAbbreviationPage,
} from './abbreviation';

const links: Array<{
  path: string,
  icon: JSX.Element,
  title: string
}> = [
  {
    path: '',
    icon: <i className="far fa-newspaper"></i>,
    title: 'Dashboard'
  },
  {
    path: 'role_group',
    icon: <i className="fa-solid fa-building-user"></i>,
    title: 'Role Group'
  },
  // {
  //   path: 'abbreviation',
  //   icon: <i className="fa-solid fa-book"></i>,
  //   title: 'Abbreviation'
  // },
  {
    path: 'question_answer',
    icon: <i className="far fa-question-circle"></i>,
    title: 'Pertanyaan & Jawaban'
  },
  {
    path: 'greeting',
    icon: <i className="far fa-grin-beam"></i>,
    title: 'Salam'
  }
]

const dashboardRoute: Partial<RouteObject> = {
  Component: (props) => <DashboardLayout links={links} {...props}/>,
  children: [
    {
      path: '',
      Component: DashboardMainPage
    },
    {
      path: 'role_group',
      children: [
        {
          path: '',
          Component: RoleGroupPage
        },
        {
          path: 'new',
          Component: AddRoleGroupPage
        },
        {
          path: 'edit',
          Component: EditRoleGroupPage
        }
      ]
    },
    {
      path: 'abbreviation',
      children: [
        {
          path: '',
          Component: AbbreviationPage,
        },
        {
          path: 'new',
          Component: AddAbbreviationPage,
        },
        {
          path: 'edit',
          Component: EditAbbreviationPage,
        },
        {
          path: 'delete',
          Component: DeleteAbbreviationPage,
        },
      ]
    },
    {
      path: 'question_answer',
      children: [
        {
          path: '',
          Component: QuestionAnswerPage
        },
        {
          path: 'new',
          Component: AddQuestionAnswerPage
        },
        {
          path: 'edit',
          Component: UpdateQuestionAnswerPage
        },
        {
          path: 'delete',
          Component: DeleteQuestionAnswerPage
        }
      ]
    },
    {
      path: 'greeting',
      children: [
        {
          path: '',
          Component: GreetingPage
        },
        {
          path: 'new',
          Component: AddGreetingPage
        },
        {
          path: 'edit',
          Component: EditGreetingPage
        },
        {
          path: 'delete',
          Component: DeleteGreetingPage
        },
      ]
    },

  ]
};

export default dashboardRoute;