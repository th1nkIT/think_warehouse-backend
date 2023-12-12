CREATE SEQUENCE IF NOT EXISTS user_backoffice_role_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;

CREATE TABLE IF NOT EXISTS user_backoffice_role
(
    id bigint NOT NULL DEFAULT nextval('user_backoffice_role_id_seq'::regclass),
    name character varying(20) NOT NULL,
    access text,
    is_all_access boolean,
    created_at timestamp without time zone NOT NULL,
    created_by character varying NOT NULL,
    updated_at timestamp without time zone,
    updated_by character varying,
    deleted_at timestamp without time zone,
    deleted_by character varying,
    CONSTRAINT user_backoffice_role_pkey PRIMARY KEY (id),
    CONSTRAINT uq_level_name_role UNIQUE (level, name)
);

-- default insert value
INSERT INTO user_backoffice_role(name, created_at, created_by, "access", "is_all_access")
VALUES ('Super Admin', (now() at time zone 'UTC')::TIMESTAMP, 'System', NULL, 't');