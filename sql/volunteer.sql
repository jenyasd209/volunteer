drop table if exists freelancers cascade;
drop table if exists orders cascade;
drop table if exists requests cascade;
drop table if exists customers cascade;
drop table if exists companies cascade;
drop table if exists complete_orders cascade;
drop table if exists freelancer_session cascade;
drop table if exists customers_session cascade;

create table freelancers (
  id serial primary key,
  first_name varchar(255) not null,
  last_name varchar(255) not null,
  email varchar(255) not null unique,
  password varchar(255) not null,
  phone varchar(255) DEFAULT '',
  facebook varchar(255) DEFAULT '',
  skype varchar(255) DEFAULT '',
  about text DEFAULT '',
  rait float DEFAULT 0,
  created_at timestamp not null
);

create table companies (
  id serial primary key,
  title varchar(255) not null
);

create table customers (
  id serial primary key,
  email varchar(255) not null unique,
  password varchar(255) not null,
  phone varchar(255) DEFAULT '',
  facebook varchar(255) DEFAULT '',
  skype varchar(255) DEFAULT '',
  about text DEFAULT '',
  rait float DEFAULT 0,
  company_id integer references companies(id),
  created_at timestamp not null
);

create table orders (
  id serial primary key,
  title varchar(255) not null,
  content text not null,
  id_customer integer references customers(id),
  status varchar(255) DEFAULT 'not done',
  created_at timestamp not null
);

create table requests (
  id serial primary key,
  id_freelancer integer references freelancers(id),
  id_order integer references orders(id),
  created_at timestamp not null
);

create table complete_orders(
  id serial primary key,
  order_id integer references orders(id),
  freelancer_id integer references freelancers(id),
  data_complete timestamp not null,
  rait float DEFAULT 0,
  comment text
);

create table freelancer_session(
  id serial primary key,
  uuid varchar(64) not null unique,
  email varchar(255) not null,
  freelancer_id integer references freelancers(id),
  created_at timestamp not null
);

create table customers_session(
  id serial primary key,
  uuid varchar(64) not null unique,
  email varchar(255) not null,
  customer_id integer references customers(id),
  created_at timestamp not null
);

INSERT INTO companies (title) values ('No company');
