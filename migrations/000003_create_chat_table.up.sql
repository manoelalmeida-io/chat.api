CREATE TABLE IF NOT EXISTS chat (
  id varchar(50) primary key,
  user_ref varchar(255) not null,
  user_id int not null,
  foreign key (user_id) references user (id)
);

CREATE TABLE IF NOT EXISTS chat_message (
  id varchar(50) primary key,
  content TEXT not null,
  user_ref varchar(255) not null,
  delivery_type varchar(50) not null,
  chat_id varchar(50) not null,
  foreign key (chat_id) references chat (id)
);
