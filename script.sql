

CREATE TABLE user_credentials (
    id       SERIAL NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE matask_user (
	id                  SERIAL NOT NULL,
	name                TEXT   NOT NULL,
	email               TEXT   NOT NULL,
	user_credentials_fk SERIAL NOT NULL
);

CREATE TABLE task (
    id      SERIAL NOT NULL,
    name    TEXT   NOT NULL,
    started DATE,
    ended   DATE
    --user_fk SERIAL NOT NULL
    --created_at timestamp
);

CREATE TABLE project (
    id       SERIAL  NOT NULL,
    progress INTEGER NOT NULL,
    task_fk  SERIAL  NOT NULL
);

CREATE TABLE book (
    id       SERIAL NOT NULL,
    progress INTEGER NOT NULL,
    author   TEXT,
    task_fk  SERIAL NOT NULL
);

CREATE TABLE movie (
    id       SERIAL NOT NULL,
    year     TEXT,
    director TEXT,
    task_fk  SERIAL NOT NULL
);

ALTER TABLE task    ADD CONSTRAINT task_pkey    PRIMARY KEY (id);
ALTER TABLE project ADD CONSTRAINT project_pkey PRIMARY KEY (id);
ALTER TABLE book    ADD CONSTRAINT book_pkey    PRIMARY KEY (id);
ALTER TABLE movie   ADD CONSTRAINT movie_pkey   PRIMARY KEY (id);

ALTER TABLE project ADD CONSTRAINT project_task_fkey FOREIGN KEY (task_fk) REFERENCES task (id);
ALTER TABLE book    ADD CONSTRAINT book_task_fkey    FOREIGN KEY (task_fk) REFERENCES task (id);
ALTER TABLE movie   ADD CONSTRAINT movie_task_fkey   FOREIGN KEY (task_fk) REFERENCES task (id);