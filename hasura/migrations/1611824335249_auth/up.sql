CREATE SCHEMA auth;
CREATE TABLE auth.role (
    role text NOT NULL
);
COMMENT ON TABLE auth.role IS 'Contain Roles for Authorization';
CREATE TABLE auth."user" (
    id uuid DEFAULT public.gen_random_uuid() NOT NULL,
    username text NOT NULL,
    password text NOT NULL
);
CREATE TABLE auth.user_role (
    id uuid DEFAULT public.gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    role text NOT NULL,
    is_default boolean DEFAULT false NOT NULL
);
ALTER TABLE ONLY auth.role
    ADD CONSTRAINT roles_pkey PRIMARY KEY (role);
ALTER TABLE ONLY auth.user_role
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);
ALTER TABLE ONLY auth.user_role
    ADD CONSTRAINT "user_roles_user_id_isDefault_key" UNIQUE (user_id, is_default);
ALTER TABLE ONLY auth.user_role
    ADD CONSTRAINT user_roles_user_id_role_key UNIQUE (user_id, role);
ALTER TABLE ONLY auth."user"
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
ALTER TABLE ONLY auth."user"
    ADD CONSTRAINT users_username_key UNIQUE (username);
ALTER TABLE ONLY auth.user_role
    ADD CONSTRAINT user_roles_role_fkey FOREIGN KEY (role) REFERENCES auth.role(role) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY auth.user_role
    ADD CONSTRAINT user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth."user"(id) ON UPDATE CASCADE ON DELETE CASCADE;
