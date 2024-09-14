
-- TAGS

CREATE TABLE public.tags (
    id integer NOT NULL,
    name character varying(25) NOT NULL,
    description character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

ALTER TABLE public.tags OWNER TO postgres;

CREATE SEQUENCE public.tags_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.tags_id_seq OWNED BY public.tags.id;

ALTER TABLE ONLY public.tags ALTER COLUMN id SET DEFAULT nextval('public.tags_id_seq'::regclass);

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);

-- BUDGETS

CREATE TABLE public.budgets (
    id integer NOT NULL,
    name character varying(25) NOT NULL,
    description character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

ALTER TABLE public.budgets OWNER TO postgres;

CREATE SEQUENCE public.budgets_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.budgets_id_seq OWNED BY public.budgets.id;

ALTER TABLE ONLY public.budgets ALTER COLUMN id SET DEFAULT nextval('public.budgets_id_seq'::regclass);

ALTER TABLE ONLY public.budgets
    ADD CONSTRAINT budgets_pkey PRIMARY KEY (id);

-- TRANSACTIONS RECURENCES

CREATE TABLE public.transactions_recurrences (
    id integer NOT NULL,
    name character varying(25) NOT NULL,
    description character varying(255) NOT NULL,
    add_time character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

ALTER TABLE public.transactions_recurrences OWNER TO postgres;

CREATE SEQUENCE public.transactions_recurrences_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.transactions_recurrences_id_seq OWNED BY public.transactions_recurrences.id;

ALTER TABLE ONLY public.transactions_recurrences ALTER COLUMN id SET DEFAULT nextval('public.transactions_recurrences_id_seq'::regclass);

ALTER TABLE ONLY public.transactions_recurrences
    ADD CONSTRAINT transactions_recurrences_pkey PRIMARY KEY (id);


-- TRANSACTIONS

CREATE TABLE public.transactions (
    id integer NOT NULL,
    name character varying(25) NOT NULL,
    description character varying(255) NOT NULL,
    budget integer NOT NULL,
    quote real NOT NULL,
    transaction_recurrence integer NOT NULL,
    active bool NOT NULL,
    starts timestamp without time zone NOT NULL,
    ends timestamp without time zone NOT NULL,
    tag integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

ALTER TABLE public.transactions OWNER TO postgres;

CREATE SEQUENCE public.transactions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.transactions_id_seq OWNED BY public.transactions.id;

ALTER TABLE ONLY public.transactions ALTER COLUMN id SET DEFAULT nextval('public.transactions_id_seq'::regclass);

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_budget_budget_id_fk FOREIGN KEY (budget) REFERENCES public.budgets(id);

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_transaction_recurrence_transaction_recurrence_id_fk FOREIGN KEY (transaction_recurrence) REFERENCES public.transactions_recurrences(id);

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_tag_tags_id_fk FOREIGN KEY (tag) REFERENCES public.tags(id);

-- LOG

CREATE TABLE public.logs (
    id integer NOT NULL,
    log character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

ALTER TABLE public.logs OWNER TO postgres;

CREATE SEQUENCE public.logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.logs_id_seq OWNED BY public.logs.id;

ALTER TABLE ONLY public.logs ALTER COLUMN id SET DEFAULT nextval('public.logs_id_seq'::regclass);

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT logs_pkey PRIMARY KEY (id);


