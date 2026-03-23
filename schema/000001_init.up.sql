CREATE TABLE users
(
    id            serial       not null unique,
    name          varchar(255) not null,
    username      varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE todo_lists
(
    id serial not null unique,
    title varchar(255) not null,
    description varchar(255)
);

CREATE TABLE users_lists
(
    id serial not null unique,
    user_id int references users (id) on delete cascade not null,
    list_id int references todo_lists(id) on delete cascade not null
);

CREATE TABLE todo_items
(
    id serial not null unique,
    title varchar(255) not null,
    description varchar(255),
    done boolean not null default false
);

CREATE TABLE lists_items
(
    id serial not null unique,
    item_id int references todo_items (id) on delete cascade not null,
    list_id int references todo_lists(id) on delete cascade not null
);

/*
 База данных для Todo-приложения: пользователи, списки задач и сами задачи, связанные между собой через промежуточные таблицы.

 Описываем структуру нашей базы данных, обьявили 5 таблиц.

Публичные:

Регистрация → users
Вход → users

Приватные:

Создание/просмотр/удаление списков → todo_lists, users_lists
Создание/просмотр/удаление задач → todo_items, lists_items
 */