CREATE TABLE IF NOT EXISTS user_contact (
  id int primary key auto_increment,
  first_name varchar(50) not null,
  last_name varchar(50) null,
  email varchar(255) not null,
  user_id int not null,
  foreign key (user_id) references user (id)
);
