create database task_checker owner= postgres encoding= 'UTF8' TEMPLATE='template0' LC_COLLATE = 'C' LC_CTYPE='C';
\c task_checker;
create table users(
id serial not null
, user_name varchar(50) not null
, password varchar(255) not null
, mail_address varchar(255) not null
, created_id varchar(50)
, created_at timestamp
, updated_id varchar(50)
, updated_at timestamp
, primary key (id)
, unique (mail_address)
);
create table task(
id serial not null
, user_id integer not null
, title varchar(255) not null
, done boolean default false not null
, del_flag boolean default false not null
, sort integer not null
, created_at timestamp not null
, updated_at timestamp not null
, primary key (id)
);