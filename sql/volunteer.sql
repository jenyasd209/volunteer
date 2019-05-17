drop table if exists users cascade;
drop table if exists roles cascade;
drop table if exists freelancers cascade;
drop table if exists specialization cascade;
drop table if exists customers cascade;
drop table if exists order_status cascade;
drop table if exists orders cascade;
drop table if exists requests cascade;
drop table if exists performed_orders cascade;
drop table if exists complete_orders cascade;
drop table if exists comments cascade;
drop table if exists messages cascade;
drop table if exists session cascade;

create table roles (
  id serial primary key,
  name varchar(255) not null unique
);
INSERT INTO roles (name) values ('Freelancer');
INSERT INTO roles (name) values ('Customer');
INSERT INTO roles (name) values ('Moderator');

create table users(
  id serial primary key,
  email varchar(255) not null unique,
  password varchar(255) not null,
  first_name varchar(255) not null,
  last_name varchar(255) not null,
  about text,
  rait float DEFAULT 0,
  phone varchar(255) DEFAULT '',
  facebook varchar(255) DEFAULT '',
  skype varchar(255) DEFAULT '',
  role_id integer references roles(id) on delete cascade on update cascade,
  photo_url varchar(255) DEFAULT '/static/image/profile.jpg',
  created_at timestamp not null
);

create table freelancers (
  id serial primary key,
  user_id integer unique references users(id) on delete cascade on update cascade,
  specialization integer[] not null
);

create table specialization(
  id serial primary key,
  name varchar(255) not null unique
);
INSERT INTO specialization (name) values ('Web');
INSERT INTO specialization (name) values ('Mobile');
INSERT INTO specialization (name) values ('Desktop');

create table customers (
  id serial primary key,
  user_id integer unique references users(id) on delete cascade on update cascade,
  organization varchar(255)
);

create table order_status(
  id serial primary key,
  name varchar(255) not null
);
INSERT INTO order_status (name) values ('Available');
INSERT INTO order_status (name) values ('Performed');
INSERT INTO order_status (name) values ('Done');

create table orders (
  id serial primary key,
  title varchar(255) not null,
  content text not null,
  customer_id integer references customers(user_id) on delete cascade on update cascade,
  status_id integer references order_status(id) on delete cascade on update cascade DEFAULT 1,
  created_at timestamp not null
);

create table requests (
  id serial primary key,
  freelancer_id integer references freelancers(user_id) on delete cascade on update cascade,
  order_id integer references orders(id) on delete cascade on update cascade,
  comment text,
  created_at timestamp not null
);

create table performed_orders(
  id serial primary key,
  order_id integer references orders(id) on delete cascade on update cascade,
  freelancer_id integer references freelancers(user_id) on delete cascade on update cascade
);

create table comments(
  id serial primary key,
  comment_text text,
  rait float not null,
  user_id integer references users(id) on delete cascade on update cascade,
  created_at timestamp not null
);

create table complete_orders(
  id serial primary key,
  order_id integer references orders(id) on delete cascade on update cascade,
  freelancer_id integer references freelancers(user_id) on delete cascade on update cascade,
  freelancer_comment_id integer unique references comments(id) on delete cascade on update cascade,
  customer_comment_id integer unique references comments(id) on delete cascade on update cascade,
  date_complete timestamp not null
);

create table messages(
  id serial primary key,
  sender_id integer references users(id) on delete cascade on update cascade,
  receiver_id integer references users(id) on delete cascade on update cascade,
  text_message text not null,
  read boolean DEFAULT FALSE,
  date_send timestamp not null
);

create table session(
  id serial primary key,
  uuid varchar(64) not null unique,
  email varchar(255) not null,
  user_id integer references users(id) on delete cascade on update cascade,
  created_at timestamp not null
);
