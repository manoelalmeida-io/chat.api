CREATE TABLE IF NOT EXISTS user (
  id int primary key auto_increment,
  first_name varchar(50) not null,
  last_name varchar(50) not null,
  email varchar(255) not null unique,
  google_sub varchar(255) not null unique
);
