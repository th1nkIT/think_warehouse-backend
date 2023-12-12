CREATE SEQUENCE IF NOT EXISTS configs_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;

CREATE TABLE IF NOT EXISTS configs
(
    id bigint NOT NULL DEFAULT nextval('configs_id_seq'::regclass),
    key character varying(50) NOT NULL,
    description character varying(200),
    value text NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    updated_by character varying(50),
    CONSTRAINT configs_pkey PRIMARY KEY (id),
    CONSTRAINT key_config UNIQUE (key)
);

-- default insert config
INSERT INTO configs(key, description, value, created_at)
VALUES
    ('config_routes', 'config for routes access', '', (now() at time zone 'UTC')::TIMESTAMP)
ON CONFLICT DO NOTHING;