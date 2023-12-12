CREATE SEQUENCE IF NOT EXISTS app_key_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;

CREATE TABLE IF NOT EXISTS app_key
(
    id bigint NOT NULL DEFAULT nextval('app_key_id_seq'::regclass),
    name character varying(100) NOT NULL,
    key character varying(200) NOT NULL,
    CONSTRAINT app_key_pkey PRIMARY KEY (id),
    CONSTRAINT uq_app_key_name UNIQUE (name),
    CONSTRAINT uq_app_key_key UNIQUE (key)
);

-- default insert value
INSERT INTO app_key(name, key)
VALUES ('wit-dev', 'w1t-d3V') ON CONFLICT DO NOTHING;