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
insert into users (user_name, password, mail_address, created_id, created_at, updated_id, updated_at) values ('test_user', '$2a$10$9zrx9F8L0tXicCULPt6Qh.2chvH3RXX9wXF5nmJ7RIlCJt8uzQ59C', 'test@example.com', 'admin', now(), 'admin', now());
insert into task  (user_id, title, sort, created_at, updated_at) SELECT id, 'task1', 1, now(), now() FROM users WHERE user_name='test_user';
insert into task  (user_id, title, sort, created_at, updated_at) SELECT id, 'task2', 2, now(), now() FROM users WHERE user_name='test_user';
insert into task  (user_id, title, done, sort, created_at, updated_at) SELECT id, 'task3', true, 3, now(), now() FROM users WHERE user_name='test_user';