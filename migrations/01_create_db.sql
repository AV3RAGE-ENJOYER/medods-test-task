-- +goose Up

CREATE TABLE public.auth (
    id integer NOT NULL,
    password_hash text NOT NULL
);


ALTER TABLE public.auth OWNER TO postgres;

CREATE SEQUENCE public.auth_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.auth_id_seq OWNER TO postgres;

ALTER SEQUENCE public.auth_id_seq OWNED BY public.auth.id;


CREATE TABLE public.refresh_tokens (
    id integer NOT NULL,
    refresh_token_hash text NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    last_ip text NOT NULL
);


ALTER TABLE public.refresh_tokens OWNER TO postgres;

CREATE SEQUENCE public.refresh_tokens_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.refresh_tokens_id_seq OWNER TO postgres;

ALTER SEQUENCE public.refresh_tokens_id_seq OWNED BY public.refresh_tokens.id;

CREATE TABLE public.users (
    id integer NOT NULL,
    email text NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;

ALTER TABLE ONLY public.auth ALTER COLUMN id SET DEFAULT nextval('public.auth_id_seq'::regclass);

ALTER TABLE ONLY public.refresh_tokens ALTER COLUMN id SET DEFAULT nextval('public.refresh_tokens_id_seq'::regclass);

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);

ALTER TABLE ONLY public.auth
    ADD CONSTRAINT auth_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.auth
    ADD CONSTRAINT id FOREIGN KEY (id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT id FOREIGN KEY (id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;

WITH add_user AS (
		INSERT INTO public.users("email")
		VALUES ('admin@gmail.com')
		RETURNING id)
		INSERT INTO auth("id", "password_hash")
		SELECT id, '$2a$12$QQR7Ljl3i21pJ0Nl76kuPO2LA5XpF6TmPa9yBqNkChduZ/WAEBisS' FROM add_user;
-- +goose Down
DROP TABLE public.auth;
DROP table public.refresh_tokens;
DROP TABLE public.users;