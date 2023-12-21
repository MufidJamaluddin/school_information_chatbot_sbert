BEGIN;

CREATE TABLE "role_group" (
    id serial not null primary key,
    role_group varchar(100) not null,
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NULL,
    created_by char(10) DEFAULT NULL,
    updated_by char(10) DEFAULT NULL
);

CREATE TABLE "admin" (
     username char(10) not null primary key,
     password varchar(60) not null,
     full_name varchar(100) not null,
     phone_no varchar(15) not null,
     position varchar(50) not null,
     role_group_id bigint default null,
     created_at timestamp DEFAULT NOW(),
     updated_at timestamp DEFAULT NULL,
     created_by char(10) DEFAULT NULL,
     updated_by char(10) DEFAULT NULL,
     FOREIGN KEY (role_group_id) REFERENCES "role_group"(id)
);

CREATE TABLE "question" (
    id serial not null primary key,
    role_group_id bigint default null,
    question varchar(255) not null,
    question_vector_sbert vector DEFAULT NULL,
    answer text not null,
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NULL,
    created_by char(10) DEFAULT NULL,
    updated_by char(10) DEFAULT NULL,
    FOREIGN KEY (created_by) REFERENCES "admin"(username),
    FOREIGN KEY (updated_by) REFERENCES "admin"(username),
    FOREIGN KEY (role_group_id) REFERENCES "role_group"(id)
);

CREATE TABLE "user" (
    id serial not null primary key,
    full_name varchar(50) not null,
    is_student boolean default FALSE,
    class_name varchar(30) not null,
    age SMALLINT not null,
    created_at timestamp DEFAULT NOW()
);

CREATE TABLE "user_response" (
    id serial not null primary key,
    user_id bigint not null,
    question text not null,
    answer text not null,
    score SMALLINT not null,
    created_at timestamp DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES "user"(id)
);

CREATE TABLE "greeting" (
    id serial not null primary key,
    greeting varchar(100) not null,
    start_time time not null,
    end_time time not null,
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NULL,
    created_by char(10) DEFAULT NULL,
    updated_by char(10) DEFAULT NULL,
    FOREIGN KEY (created_by) REFERENCES "admin"(username),
    FOREIGN KEY (updated_by) REFERENCES "admin"(username)
);

CREATE TABLE "abbreviation" (
   standard_word varchar(255) not null PRIMARY KEY,
   list_abbreviation_term varchar(255) ARRAY not null,
   created_at timestamp DEFAULT NOW(),
   updated_at timestamp DEFAULT NULL,
   created_by char(10) DEFAULT NULL,
   updated_by char(10) DEFAULT NULL,
   FOREIGN KEY (created_by) REFERENCES "admin"(username),
   FOREIGN KEY (updated_by) REFERENCES "admin"(username)
);

COMMIT;