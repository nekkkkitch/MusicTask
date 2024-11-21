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

create schema if not exists public;

alter schema public owner to pg_database_owner;

create table if not exists public.groups(
    id serial not null primary key,
    title text not null unique
);

create table if not exists public.songs(
    id serial not null primary key,
    "group" text not null references public.groups(title),
    song text not null,
    releaseDate date,
    "text" text,
    link text
);

create index group_name on public.songs("group")