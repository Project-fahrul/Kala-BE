--
-- PostgreSQL database dump
--

-- Dumped from database version 10.22 (Ubuntu 10.22-0ubuntu0.18.04.1)
-- Dumped by pg_dump version 10.22 (Ubuntu 10.22-0ubuntu0.18.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: DATABASE postgres; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON DATABASE postgres IS 'default administrative connection database';


--
-- Name: kala; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA kala;


ALTER SCHEMA kala OWNER TO postgres;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: role; Type: TYPE; Schema: kala; Owner: postgres
--

CREATE TYPE kala.role AS ENUM (
    'admin',
    'sales'
);


ALTER TYPE kala.role OWNER TO postgres;

--
-- Name: customer_seq; Type: SEQUENCE; Schema: kala; Owner: postgres
--

CREATE SEQUENCE kala.customer_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE kala.customer_seq OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: customers; Type: TABLE; Schema: kala; Owner: postgres
--

CREATE TABLE kala.customers (
    id bigint DEFAULT nextval('kala.customer_seq'::regclass) NOT NULL,
    name character varying(100),
    address character varying(100),
    no_hp character varying(20),
    sales_id bigint,
    tgl_dec date NOT NULL,
    tgl_lahir date NOT NULL,
    tgl_stnk date NOT NULL,
    tgl_angsuran date NOT NULL,
    type_angsuran character varying(10),
    no_rangka character varying(100),
    type_kendaraan character varying(50) NOT NULL,
    leasing character varying(50),
    new_customer boolean DEFAULT true,
    total_angsuran integer DEFAULT 0
);


ALTER TABLE kala.customers OWNER TO postgres;

--
-- Name: evidances; Type: TABLE; Schema: kala; Owner: postgres
--

CREATE TABLE kala.evidances (
    sales_id bigint,
    customer_id bigint,
    submit_date date,
    due_date date,
    content character varying(250),
    comment character varying(200),
    type_evidance character varying(15)
);


ALTER TABLE kala.evidances OWNER TO postgres;

--
-- Name: notifications; Type: TABLE; Schema: kala; Owner: postgres
--

CREATE TABLE kala.notifications (
    sales_id bigint,
    customer_id bigint,
    message character varying(200),
    type_notification character varying(15),
    due_date date
);


ALTER TABLE kala.notifications OWNER TO postgres;

--
-- Name: user_seq; Type: SEQUENCE; Schema: kala; Owner: postgres
--

CREATE SEQUENCE kala.user_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE kala.user_seq OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: kala; Owner: postgres
--

CREATE TABLE kala.users (
    id bigint DEFAULT nextval('kala.user_seq'::regclass) NOT NULL,
    name character varying(30),
    email character varying(30),
    password character varying(70),
    role character varying(15),
    token character varying(100),
    token_expired timestamp without time zone,
    login_delay timestamp without time zone,
    phone_number character varying(20),
    verified boolean DEFAULT false
);


ALTER TABLE kala.users OWNER TO postgres;

--
-- Data for Name: customers; Type: TABLE DATA; Schema: kala; Owner: postgres
--

COPY kala.customers (id, name, address, no_hp, sales_id, tgl_dec, tgl_lahir, tgl_stnk, tgl_angsuran, type_angsuran, no_rangka, type_kendaraan, leasing, new_customer, total_angsuran) FROM stdin;
1	Junettt	Jakarta	081273645743	5	2022-12-12	2022-01-12	2022-12-12	2022-12-12	Kredit	0986765462	mobil	kala	t	0
6	hh	Sleman Yogyakarta	085225077747	9	2022-09-18	2020-12-12	2022-09-19	2022-09-18	Tunai	as	as	as	t	4
\.


--
-- Data for Name: evidances; Type: TABLE DATA; Schema: kala; Owner: postgres
--

COPY kala.evidances (sales_id, customer_id, submit_date, due_date, content, comment, type_evidance) FROM stdin;
5	1	2022-09-17	2022-09-17	\N	saya sudah mengirim	birthday
11	1	\N	2022-09-19	\N	\N	birthday
\.


--
-- Data for Name: notifications; Type: TABLE DATA; Schema: kala; Owner: postgres
--

COPY kala.notifications (sales_id, customer_id, message, type_notification, due_date) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: kala; Owner: postgres
--

COPY kala.users (id, name, email, password, role, token, token_expired, login_delay, phone_number, verified) FROM stdin;
5	fahrul	fahrulputraa40@gmail.com	$2a$12$D3IfoTjEk8wwN2u8nJSqy.c49ExU67D6fg/n3BJmMBy1.fTSoxdY6	sales		0001-01-01 00:00:00	0001-01-01 00:00:00	123456	f
9	fahruel	fahrulputra40@gmail.com	$2a$12$osf/g3o8tfXjXs3E/EZ7BeniUel2p6am37OZWn7eYwO58foUldA/.	admin	\N	\N	\N	12121212	f
13	fahrul putra 2	fahrulputrae40@gmail.com		sales		0001-01-01 00:00:00	0001-01-01 00:00:00	085225077747	f
11	fahrul3	fahrul.putra40@gmail.com		sales		0001-01-01 00:00:00	0001-01-01 00:00:00	09928282	f
\.


--
-- Name: customer_seq; Type: SEQUENCE SET; Schema: kala; Owner: postgres
--

SELECT pg_catalog.setval('kala.customer_seq', 6, true);


--
-- Name: user_seq; Type: SEQUENCE SET; Schema: kala; Owner: postgres
--

SELECT pg_catalog.setval('kala.user_seq', 13, true);


--
-- Name: customers customers_no_rangka_key; Type: CONSTRAINT; Schema: kala; Owner: postgres
--

ALTER TABLE ONLY kala.customers
    ADD CONSTRAINT customers_no_rangka_key UNIQUE (no_rangka);


--
-- Name: customers customers_pkey; Type: CONSTRAINT; Schema: kala; Owner: postgres
--

ALTER TABLE ONLY kala.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: kala; Owner: postgres
--

ALTER TABLE ONLY kala.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: kala; Owner: postgres
--

ALTER TABLE ONLY kala.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: evidances evidances_customer_id_fkey; Type: FK CONSTRAINT; Schema: kala; Owner: postgres
--

ALTER TABLE ONLY kala.evidances
    ADD CONSTRAINT evidances_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES kala.customers(id);


--
-- Name: evidances evidances_sales_id_fkey; Type: FK CONSTRAINT; Schema: kala; Owner: postgres
--

ALTER TABLE ONLY kala.evidances
    ADD CONSTRAINT evidances_sales_id_fkey FOREIGN KEY (sales_id) REFERENCES kala.users(id);


--
-- Name: customers fk_sales; Type: FK CONSTRAINT; Schema: kala; Owner: postgres
--

ALTER TABLE ONLY kala.customers
    ADD CONSTRAINT fk_sales FOREIGN KEY (sales_id) REFERENCES kala.users(id);


--
-- Name: notifications notifications_customer_id_fkey; Type: FK CONSTRAINT; Schema: kala; Owner: postgres
--

ALTER TABLE ONLY kala.notifications
    ADD CONSTRAINT notifications_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES kala.customers(id);


--
-- Name: notifications notifications_sales_id_fkey; Type: FK CONSTRAINT; Schema: kala; Owner: postgres
--

ALTER TABLE ONLY kala.notifications
    ADD CONSTRAINT notifications_sales_id_fkey FOREIGN KEY (sales_id) REFERENCES kala.users(id);


--
-- PostgreSQL database dump complete
--

