create table if not exists message(id INTEGER PRIMARY KEY AUTOINCREMENT, message text not null, user_id integer not null, message_date integer, FOREIGN KEY(user_id) REFERENCES user(id));