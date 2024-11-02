
CREATE TABLE user_credentials (
    id       SERIAL NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE matask_user (
	id                  SERIAL NOT NULL,
	name                TEXT   NOT NULL,
	email               TEXT   NOT NULL,
    birthday            DATE,
	user_credentials_fk SERIAL NOT NULL
);

CREATE TABLE task (
    id      SERIAL    NOT NULL,
    name    TEXT      NOT NULL,
    type    TEXT      NOT NULL,
    started DATE,
    ended   DATE,
    created TIMESTAMP NOT NULL,
    user_fk SERIAL    NOT NULL
);

CREATE TABLE project (
    id             SERIAL  NOT NULL,
    description    TEXT,
    progress       INTEGER NOT NULL,
    dynamic_fields JSONB,
    task_fk        SERIAL  NOT NULL
);

CREATE TABLE book (
    id         SERIAL  NOT NULL,
    progress   INTEGER NOT NULL,
    author     TEXT,
    synopsis   TEXT,
    comments   TEXT,
    year       TEXT,
    genre      TEXT,
    rate       INTEGER,
    cover_path TEXT,
    task_fk    SERIAL  NOT NULL
);

CREATE TABLE movie (
    id          SERIAL   NOT NULL,
    synopsis    TEXT,
    comments    TEXT,
    year        TEXT,
    rate        INTEGER,
    actors      JSONB,
    director    TEXT,
    genre       TEXT,
    poster_path TEXT,
    task_fk     SERIAL   NOT NULL
);

ALTER TABLE task    ADD CONSTRAINT task_pkey    PRIMARY KEY (id);
ALTER TABLE project ADD CONSTRAINT project_pkey PRIMARY KEY (id);
ALTER TABLE book    ADD CONSTRAINT book_pkey    PRIMARY KEY (id);
ALTER TABLE movie   ADD CONSTRAINT movie_pkey   PRIMARY KEY (id);

ALTER TABLE project ADD CONSTRAINT project_task_fkey FOREIGN KEY (task_fk) REFERENCES task (id);
ALTER TABLE book    ADD CONSTRAINT book_task_fkey    FOREIGN KEY (task_fk) REFERENCES task (id);
ALTER TABLE movie   ADD CONSTRAINT movie_task_fkey   FOREIGN KEY (task_fk) REFERENCES task (id);

ALTER TABLE user_credentials ADD CONSTRAINT user_credentials_username_unq UNIQUE (username);
ALTER TABLE matask_user ADD CONSTRAINT user_email_unq UNIQUE (email);
