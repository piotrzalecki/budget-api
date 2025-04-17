CREATE TABLE
  public.users (
    id serial NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    email character varying(255) NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    password character varying(60) NOT NULL,
    updated_at timestamp without time zone NOT NULL DEFAULT now(),
    user_active integer NOT NULL DEFAULT 0
  );

ALTER TABLE
  public.users
ADD
  CONSTRAINT users_pkey PRIMARY KEY (id)

-- Adds defult admind@admin.com user with password: password
insert into "public"."users" ("created_at", "email", "first_name", "id", "last_name", "password", "updated_at", "user_active") values ('2024-08-01 18:18:36.377487', 'admin@admin.com', 'Admin', 1, 'User', '$2a$12$/izY/i03vRWj.aEnBdwt0ue4JnTRqQWuXdFufcS4StNr7fbfsu8Fi', '2024-08-01 18:18:36.377487', 1)