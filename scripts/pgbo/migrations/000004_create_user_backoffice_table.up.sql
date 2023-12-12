CREATE SEQUENCE IF NOT EXISTS user_backoffice_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;

CREATE TABLE IF NOT EXISTS user_backoffice
(
    id bigint NOT NULL DEFAULT nextval('user_backoffice_id_seq'::regclass),
    guid character varying NOT NULL,
    name character varying(100),
    profile_picture_image_url text,
    phone character varying(20) NOT NULL,
    email character varying(50) NOT NULL,
    role_id integer NOT NULL,
    password character varying(200) NOT NULL,
    salt character varying(50) NOT NULL,
    is_active boolean,
    created_at timestamp without time zone NOT NULL,
    created_by character varying NOT NULL,
    updated_at timestamp without time zone,
    updated_by character varying,
    deleted_at timestamp without time zone,
    deleted_by character varying,
    last_login timestamp without time zone,
    CONSTRAINT user_backoffice_pkey PRIMARY KEY (guid),
    CONSTRAINT uq_email_user_backoffice UNIQUE (email)
);

INSERT INTO
     user_backoffice("guid", "name", "profile_picture_image_url", "phone", "email", "role_id", "password", "salt", "is_active", "created_at", "created_by")
     VALUES ('31eb7f5b-1ddd-4746-9930-efa89951e8c3', 'Super Admin WIT', NULL, '+6281122334455', 'admin@wit.id', 1, 'ba0ac3dc230b79c22a7159ab0bfafd74dd658080', 'Ae5Xn', 't', (now() at time zone 'UTC')::TIMESTAMP, 'System')
ON CONFLICT DO NOTHING;