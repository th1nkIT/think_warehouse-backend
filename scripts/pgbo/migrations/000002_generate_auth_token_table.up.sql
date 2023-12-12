CREATE SEQUENCE IF NOT EXISTS auth_token_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;

CREATE TABLE IF NOT EXISTS auth_token
(
    id bigint NOT NULL DEFAULT nextval('auth_token_id_seq'::regclass),
    name character varying NOT NULL,
    device_id character varying NOT NULL,
    device_type character varying NOT NULL,
    token character varying NOT NULL,
    token_expired timestamp without time zone NOT NULL,
    refresh_token character varying NOT NULL,
    refresh_token_expired timestamp without time zone NOT NULL,
    is_login boolean NOT NULL DEFAULT false,
    user_login character varying ,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone,
    CONSTRAINT token_pkey PRIMARY KEY (name, device_id, device_type)
);

