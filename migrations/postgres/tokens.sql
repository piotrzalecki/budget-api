CREATE TABLE
  public.tokens (
    id serial NOT NULL,
    user_id integer NULL,
    email character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    token character varying(255) NOT NULL,
    token_hash bytea NOT NULL,
    updated_at timestamp without time zone NOT NULL DEFAULT now(),
    expiry timestamp with time zone NOT NULL
  );

ALTER TABLE
  public.tokens
ADD
  CONSTRAINT tokens_pkey PRIMARY KEY (id)