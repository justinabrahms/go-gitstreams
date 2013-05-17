--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

--
-- Name: auth_group_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE auth_group_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.auth_group_id_seq OWNER TO gitstreams;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: auth_group; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE auth_group (
    id integer DEFAULT nextval('auth_group_id_seq'::regclass) NOT NULL,
    name character varying(80) NOT NULL
);


ALTER TABLE public.auth_group OWNER TO gitstreams;

--
-- Name: auth_group_permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE auth_group_permissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.auth_group_permissions_id_seq OWNER TO gitstreams;

--
-- Name: auth_group_permissions; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE auth_group_permissions (
    id integer DEFAULT nextval('auth_group_permissions_id_seq'::regclass) NOT NULL,
    group_id integer NOT NULL,
    permission_id integer NOT NULL
);


ALTER TABLE public.auth_group_permissions OWNER TO gitstreams;

--
-- Name: auth_permission_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE auth_permission_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.auth_permission_id_seq OWNER TO gitstreams;

--
-- Name: auth_permission; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE auth_permission (
    id integer DEFAULT nextval('auth_permission_id_seq'::regclass) NOT NULL,
    name character varying(50) NOT NULL,
    content_type_id integer NOT NULL,
    codename character varying(100) NOT NULL
);


ALTER TABLE public.auth_permission OWNER TO gitstreams;

--
-- Name: auth_user_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE auth_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.auth_user_id_seq OWNER TO gitstreams;

--
-- Name: auth_user; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE auth_user (
    id integer DEFAULT nextval('auth_user_id_seq'::regclass) NOT NULL,
    username character varying(30) NOT NULL,
    first_name character varying(30) NOT NULL,
    last_name character varying(30) NOT NULL,
    email character varying(75) NOT NULL,
    password character varying(128) NOT NULL,
    is_staff boolean NOT NULL,
    is_active boolean NOT NULL,
    is_superuser boolean NOT NULL,
    last_login timestamp without time zone NOT NULL,
    date_joined timestamp without time zone NOT NULL
);


ALTER TABLE public.auth_user OWNER TO gitstreams;

--
-- Name: auth_user_groups_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE auth_user_groups_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.auth_user_groups_id_seq OWNER TO gitstreams;

--
-- Name: auth_user_groups; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE auth_user_groups (
    id integer DEFAULT nextval('auth_user_groups_id_seq'::regclass) NOT NULL,
    user_id integer NOT NULL,
    group_id integer NOT NULL
);


ALTER TABLE public.auth_user_groups OWNER TO gitstreams;

--
-- Name: auth_user_user_permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE auth_user_user_permissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.auth_user_user_permissions_id_seq OWNER TO gitstreams;

--
-- Name: auth_user_user_permissions; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE auth_user_user_permissions (
    id integer DEFAULT nextval('auth_user_user_permissions_id_seq'::regclass) NOT NULL,
    user_id integer NOT NULL,
    permission_id integer NOT NULL
);


ALTER TABLE public.auth_user_user_permissions OWNER TO gitstreams;

--
-- Name: django_admin_log_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE django_admin_log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.django_admin_log_id_seq OWNER TO gitstreams;

--
-- Name: django_admin_log; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE django_admin_log (
    id integer DEFAULT nextval('django_admin_log_id_seq'::regclass) NOT NULL,
    action_time timestamp without time zone NOT NULL,
    user_id integer NOT NULL,
    content_type_id integer,
    object_id text,
    object_repr character varying(200) NOT NULL,
    action_flag integer NOT NULL,
    change_message text NOT NULL
);


ALTER TABLE public.django_admin_log OWNER TO gitstreams;

--
-- Name: django_content_type_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE django_content_type_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.django_content_type_id_seq OWNER TO gitstreams;

--
-- Name: django_content_type; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE django_content_type (
    id integer DEFAULT nextval('django_content_type_id_seq'::regclass) NOT NULL,
    name character varying(100) NOT NULL,
    app_label character varying(100) NOT NULL,
    model character varying(100) NOT NULL
);


ALTER TABLE public.django_content_type OWNER TO gitstreams;

--
-- Name: django_session; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE django_session (
    session_key character varying(40) NOT NULL,
    session_data text NOT NULL,
    expire_date timestamp without time zone NOT NULL
);


ALTER TABLE public.django_session OWNER TO gitstreams;

--
-- Name: django_site_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE django_site_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.django_site_id_seq OWNER TO gitstreams;

--
-- Name: django_site; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE django_site (
    id integer DEFAULT nextval('django_site_id_seq'::regclass) NOT NULL,
    domain character varying(100) NOT NULL,
    name character varying(50) NOT NULL
);


ALTER TABLE public.django_site OWNER TO gitstreams;

--
-- Name: nashvegas_migration_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE nashvegas_migration_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.nashvegas_migration_id_seq OWNER TO gitstreams;

--
-- Name: nashvegas_migration; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE nashvegas_migration (
    id integer DEFAULT nextval('nashvegas_migration_id_seq'::regclass) NOT NULL,
    migration_label character varying(200) NOT NULL,
    date_created timestamp without time zone NOT NULL,
    content text NOT NULL,
    scm_version character varying(50)
);


ALTER TABLE public.nashvegas_migration OWNER TO gitstreams;

--
-- Name: social_auth_association_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE social_auth_association_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.social_auth_association_id_seq OWNER TO gitstreams;

--
-- Name: social_auth_association; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE social_auth_association (
    id integer DEFAULT nextval('social_auth_association_id_seq'::regclass) NOT NULL,
    server_url character varying(255) NOT NULL,
    handle character varying(255) NOT NULL,
    secret character varying(255) NOT NULL,
    issued integer NOT NULL,
    lifetime integer NOT NULL,
    assoc_type character varying(64) NOT NULL
);


ALTER TABLE public.social_auth_association OWNER TO gitstreams;

--
-- Name: social_auth_nonce_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE social_auth_nonce_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.social_auth_nonce_id_seq OWNER TO gitstreams;

--
-- Name: social_auth_nonce; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE social_auth_nonce (
    id integer DEFAULT nextval('social_auth_nonce_id_seq'::regclass) NOT NULL,
    server_url character varying(255) NOT NULL,
    "timestamp" integer NOT NULL,
    salt character varying(40) NOT NULL
);


ALTER TABLE public.social_auth_nonce OWNER TO gitstreams;

--
-- Name: social_auth_usersocialauth_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE social_auth_usersocialauth_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.social_auth_usersocialauth_id_seq OWNER TO gitstreams;

--
-- Name: social_auth_usersocialauth; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE social_auth_usersocialauth (
    id integer DEFAULT nextval('social_auth_usersocialauth_id_seq'::regclass) NOT NULL,
    user_id integer NOT NULL,
    provider character varying(32) NOT NULL,
    uid character varying(255) NOT NULL,
    extra_data text NOT NULL
);


ALTER TABLE public.social_auth_usersocialauth OWNER TO gitstreams;

--
-- Name: streamer_activity_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE streamer_activity_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.streamer_activity_id_seq OWNER TO gitstreams;

--
-- Name: streamer_activity; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE streamer_activity (
    id integer DEFAULT nextval('streamer_activity_id_seq'::regclass) NOT NULL,
    type character varying(2) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    repo_id integer,
    user_id integer,
    event_id bigint NOT NULL,
    meta text NOT NULL
);


ALTER TABLE public.streamer_activity OWNER TO gitstreams;

--
-- Name: streamer_githubuser_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE streamer_githubuser_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.streamer_githubuser_id_seq OWNER TO gitstreams;

--
-- Name: streamer_githubuser; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE streamer_githubuser (
    id integer DEFAULT nextval('streamer_githubuser_id_seq'::regclass) NOT NULL,
    name character varying(255) NOT NULL,
    last_synced timestamp with time zone
);


ALTER TABLE public.streamer_githubuser OWNER TO gitstreams;

--
-- Name: streamer_repo_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE streamer_repo_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.streamer_repo_id_seq OWNER TO gitstreams;

--
-- Name: streamer_repo; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE streamer_repo (
    id integer DEFAULT nextval('streamer_repo_id_seq'::regclass) NOT NULL,
    username character varying(255) NOT NULL,
    project_name character varying(255) NOT NULL,
    last_synced timestamp with time zone,
    description character varying(255)
);


ALTER TABLE public.streamer_repo OWNER TO gitstreams;

--
-- Name: streamer_userprofile_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE streamer_userprofile_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.streamer_userprofile_id_seq OWNER TO gitstreams;

--
-- Name: streamer_userprofile; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE streamer_userprofile (
    id integer DEFAULT nextval('streamer_userprofile_id_seq'::regclass) NOT NULL,
    user_id integer NOT NULL,
    max_time_interval_between_emails character varying(1) NOT NULL,
    last_email_received timestamp with time zone,
    waitlisted boolean DEFAULT true,
    include_starred_repos boolean DEFAULT true
);


ALTER TABLE public.streamer_userprofile OWNER TO gitstreams;

--
-- Name: streamer_userprofile_followed_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE streamer_userprofile_followed_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.streamer_userprofile_followed_id_seq OWNER TO gitstreams;

--
-- Name: streamer_userprofile_followed; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE streamer_userprofile_followed (
    id integer DEFAULT nextval('streamer_userprofile_followed_id_seq'::regclass) NOT NULL,
    userprofile_id integer NOT NULL,
    githubuser_id integer NOT NULL,
    last_sent timestamp with time zone
);


ALTER TABLE public.streamer_userprofile_followed OWNER TO gitstreams;

--
-- Name: streamer_userprofile_repos_id_seq; Type: SEQUENCE; Schema: public; Owner: gitstreams
--

CREATE SEQUENCE streamer_userprofile_repos_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.streamer_userprofile_repos_id_seq OWNER TO gitstreams;

--
-- Name: streamer_userprofile_repos; Type: TABLE; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE TABLE streamer_userprofile_repos (
    id integer DEFAULT nextval('streamer_userprofile_repos_id_seq'::regclass) NOT NULL,
    userprofile_id integer NOT NULL,
    repo_id integer NOT NULL,
    last_sent timestamp with time zone
);


ALTER TABLE public.streamer_userprofile_repos OWNER TO gitstreams;

--
-- Name: auth_group_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY auth_group
    ADD CONSTRAINT auth_group_id_pkey PRIMARY KEY (id);


--
-- Name: auth_group_permissions_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY auth_group_permissions
    ADD CONSTRAINT auth_group_permissions_id_pkey PRIMARY KEY (id);


--
-- Name: auth_permission_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY auth_permission
    ADD CONSTRAINT auth_permission_id_pkey PRIMARY KEY (id);


--
-- Name: auth_user_groups_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY auth_user_groups
    ADD CONSTRAINT auth_user_groups_id_pkey PRIMARY KEY (id);


--
-- Name: auth_user_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY auth_user
    ADD CONSTRAINT auth_user_id_pkey PRIMARY KEY (id);


--
-- Name: auth_user_user_permissions_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY auth_user_user_permissions
    ADD CONSTRAINT auth_user_user_permissions_id_pkey PRIMARY KEY (id);


--
-- Name: django_admin_log_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY django_admin_log
    ADD CONSTRAINT django_admin_log_id_pkey PRIMARY KEY (id);


--
-- Name: django_content_type_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY django_content_type
    ADD CONSTRAINT django_content_type_id_pkey PRIMARY KEY (id);


--
-- Name: django_session_session_key_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY django_session
    ADD CONSTRAINT django_session_session_key_pkey PRIMARY KEY (session_key);


--
-- Name: django_site_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY django_site
    ADD CONSTRAINT django_site_id_pkey PRIMARY KEY (id);


--
-- Name: nashvegas_migration_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY nashvegas_migration
    ADD CONSTRAINT nashvegas_migration_id_pkey PRIMARY KEY (id);


--
-- Name: social_auth_association_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY social_auth_association
    ADD CONSTRAINT social_auth_association_id_pkey PRIMARY KEY (id);


--
-- Name: social_auth_nonce_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY social_auth_nonce
    ADD CONSTRAINT social_auth_nonce_id_pkey PRIMARY KEY (id);


--
-- Name: social_auth_usersocialauth_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY social_auth_usersocialauth
    ADD CONSTRAINT social_auth_usersocialauth_id_pkey PRIMARY KEY (id);


--
-- Name: streamer_activity_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY streamer_activity
    ADD CONSTRAINT streamer_activity_id_pkey PRIMARY KEY (id);


--
-- Name: streamer_githubuser_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY streamer_githubuser
    ADD CONSTRAINT streamer_githubuser_id_pkey PRIMARY KEY (id);


--
-- Name: streamer_repo_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY streamer_repo
    ADD CONSTRAINT streamer_repo_id_pkey PRIMARY KEY (id);


--
-- Name: streamer_userprofile_followed_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY streamer_userprofile_followed
    ADD CONSTRAINT streamer_userprofile_followed_id_pkey PRIMARY KEY (id);


--
-- Name: streamer_userprofile_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY streamer_userprofile
    ADD CONSTRAINT streamer_userprofile_id_pkey PRIMARY KEY (id);


--
-- Name: streamer_userprofile_repos_id_pkey; Type: CONSTRAINT; Schema: public; Owner: gitstreams; Tablespace: 
--

ALTER TABLE ONLY streamer_userprofile_repos
    ADD CONSTRAINT streamer_userprofile_repos_id_pkey PRIMARY KEY (id);


--
-- Name: auth_group_name; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE UNIQUE INDEX auth_group_name ON auth_group USING btree (name);


--
-- Name: auth_group_permissions_group_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX auth_group_permissions_group_id ON auth_group_permissions USING btree (group_id);


--
-- Name: auth_group_permissions_group_id_permission_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE UNIQUE INDEX auth_group_permissions_group_id_permission_id ON auth_group_permissions USING btree (group_id, permission_id);


--
-- Name: auth_group_permissions_permission_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX auth_group_permissions_permission_id ON auth_group_permissions USING btree (permission_id);


--
-- Name: auth_permission_content_type_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX auth_permission_content_type_id ON auth_permission USING btree (content_type_id);


--
-- Name: auth_permission_content_type_id_codename; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE UNIQUE INDEX auth_permission_content_type_id_codename ON auth_permission USING btree (content_type_id, codename);


--
-- Name: auth_user_groups_group_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX auth_user_groups_group_id ON auth_user_groups USING btree (group_id);


--
-- Name: auth_user_groups_user_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX auth_user_groups_user_id ON auth_user_groups USING btree (user_id);


--
-- Name: auth_user_groups_user_id_group_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE UNIQUE INDEX auth_user_groups_user_id_group_id ON auth_user_groups USING btree (user_id, group_id);


--
-- Name: auth_user_user_permissions_permission_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX auth_user_user_permissions_permission_id ON auth_user_user_permissions USING btree (permission_id);


--
-- Name: auth_user_user_permissions_user_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX auth_user_user_permissions_user_id ON auth_user_user_permissions USING btree (user_id);


--
-- Name: auth_user_user_permissions_user_id_permission_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE UNIQUE INDEX auth_user_user_permissions_user_id_permission_id ON auth_user_user_permissions USING btree (user_id, permission_id);


--
-- Name: auth_user_username; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE UNIQUE INDEX auth_user_username ON auth_user USING btree (username);


--
-- Name: django_admin_log_content_type_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX django_admin_log_content_type_id ON django_admin_log USING btree (content_type_id);


--
-- Name: django_admin_log_user_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX django_admin_log_user_id ON django_admin_log USING btree (user_id);


--
-- Name: django_content_type_app_label_model; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE UNIQUE INDEX django_content_type_app_label_model ON django_content_type USING btree (app_label, model);


--
-- Name: django_session_expire_date; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX django_session_expire_date ON django_session USING btree (expire_date);


--
-- Name: social_auth_usersocialauth_provider_uid; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE UNIQUE INDEX social_auth_usersocialauth_provider_uid ON social_auth_usersocialauth USING btree (provider, uid);


--
-- Name: social_auth_usersocialauth_user_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX social_auth_usersocialauth_user_id ON social_auth_usersocialauth USING btree (user_id);


--
-- Name: streamer_activity_created_at; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX streamer_activity_created_at ON streamer_activity USING btree (created_at);


--
-- Name: streamer_activity_repo_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX streamer_activity_repo_id ON streamer_activity USING btree (repo_id);


--
-- Name: streamer_activity_user_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX streamer_activity_user_id ON streamer_activity USING btree (user_id);


--
-- Name: streamer_userprofile_followed_githubuser_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX streamer_userprofile_followed_githubuser_id ON streamer_userprofile_followed USING btree (githubuser_id);


--
-- Name: streamer_userprofile_followed_userprofile_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX streamer_userprofile_followed_userprofile_id ON streamer_userprofile_followed USING btree (userprofile_id);


--
-- Name: streamer_userprofile_followed_userprofile_id_githubuser_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE UNIQUE INDEX streamer_userprofile_followed_userprofile_id_githubuser_id ON streamer_userprofile_followed USING btree (userprofile_id, githubuser_id);


--
-- Name: streamer_userprofile_repos_repo_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX streamer_userprofile_repos_repo_id ON streamer_userprofile_repos USING btree (repo_id);


--
-- Name: streamer_userprofile_repos_userprofile_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE INDEX streamer_userprofile_repos_userprofile_id ON streamer_userprofile_repos USING btree (userprofile_id);


--
-- Name: streamer_userprofile_repos_userprofile_id_repo_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE UNIQUE INDEX streamer_userprofile_repos_userprofile_id_repo_id ON streamer_userprofile_repos USING btree (userprofile_id, repo_id);


--
-- Name: streamer_userprofile_user_id; Type: INDEX; Schema: public; Owner: gitstreams; Tablespace: 
--

CREATE UNIQUE INDEX streamer_userprofile_user_id ON streamer_userprofile USING btree (user_id);


--
-- Name: auth_group_permissions_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY auth_group_permissions
    ADD CONSTRAINT auth_group_permissions_group_id_fkey FOREIGN KEY (group_id) REFERENCES auth_group(id);


--
-- Name: auth_group_permissions_permission_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY auth_group_permissions
    ADD CONSTRAINT auth_group_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES auth_permission(id);


--
-- Name: auth_permission_content_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY auth_permission
    ADD CONSTRAINT auth_permission_content_type_id_fkey FOREIGN KEY (content_type_id) REFERENCES django_content_type(id);


--
-- Name: auth_user_groups_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY auth_user_groups
    ADD CONSTRAINT auth_user_groups_group_id_fkey FOREIGN KEY (group_id) REFERENCES auth_group(id);


--
-- Name: auth_user_groups_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY auth_user_groups
    ADD CONSTRAINT auth_user_groups_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth_user(id);


--
-- Name: auth_user_user_permissions_permission_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY auth_user_user_permissions
    ADD CONSTRAINT auth_user_user_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES auth_permission(id);


--
-- Name: auth_user_user_permissions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY auth_user_user_permissions
    ADD CONSTRAINT auth_user_user_permissions_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth_user(id);


--
-- Name: django_admin_log_content_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY django_admin_log
    ADD CONSTRAINT django_admin_log_content_type_id_fkey FOREIGN KEY (content_type_id) REFERENCES django_content_type(id);


--
-- Name: django_admin_log_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY django_admin_log
    ADD CONSTRAINT django_admin_log_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth_user(id);


--
-- Name: social_auth_usersocialauth_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY social_auth_usersocialauth
    ADD CONSTRAINT social_auth_usersocialauth_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth_user(id);


--
-- Name: streamer_activity_repo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY streamer_activity
    ADD CONSTRAINT streamer_activity_repo_id_fkey FOREIGN KEY (repo_id) REFERENCES streamer_repo(id);


--
-- Name: streamer_activity_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY streamer_activity
    ADD CONSTRAINT streamer_activity_user_id_fkey FOREIGN KEY (user_id) REFERENCES streamer_githubuser(id);


--
-- Name: streamer_userprofile_followed_githubuser_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY streamer_userprofile_followed
    ADD CONSTRAINT streamer_userprofile_followed_githubuser_id_fkey FOREIGN KEY (githubuser_id) REFERENCES streamer_githubuser(id);


--
-- Name: streamer_userprofile_followed_userprofile_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY streamer_userprofile_followed
    ADD CONSTRAINT streamer_userprofile_followed_userprofile_id_fkey FOREIGN KEY (userprofile_id) REFERENCES streamer_userprofile(id);


--
-- Name: streamer_userprofile_repos_repo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY streamer_userprofile_repos
    ADD CONSTRAINT streamer_userprofile_repos_repo_id_fkey FOREIGN KEY (repo_id) REFERENCES streamer_repo(id);


--
-- Name: streamer_userprofile_repos_userprofile_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY streamer_userprofile_repos
    ADD CONSTRAINT streamer_userprofile_repos_userprofile_id_fkey FOREIGN KEY (userprofile_id) REFERENCES streamer_userprofile(id);


--
-- Name: streamer_userprofile_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: gitstreams
--

ALTER TABLE ONLY streamer_userprofile
    ADD CONSTRAINT streamer_userprofile_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth_user(id);


--
-- Name: public; Type: ACL; Schema: -; Owner: justinlilly
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM justinlilly;
GRANT ALL ON SCHEMA public TO justinlilly;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

