CREATE SEQUENCE IF NOT EXISTS user_handheld_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;


CREATE TABLE IF NOT EXISTS user_handheld
(
    id bigint NOT NULL DEFAULT nextval('user_handheld_id_seq'::regclass),
    guid character varying NOT NULL,
    name character varying(100) NOT NULL,
    profile_picture_image_url text,
    phone character varying(20),
    email character varying(50) NOT NULL,
    gender character varying(20) NOT NULL,
    address text,
    salt character varying(50) NOT NULL,
    password character varying(200) NOT NULL,
    is_active boolean,
    fcm_token text,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    last_login timestamp without time zone,
    CONSTRAINT user_handheld_pkey PRIMARY KEY (guid),
    CONSTRAINT uq_email_user_handheld UNIQUE (email)
);
