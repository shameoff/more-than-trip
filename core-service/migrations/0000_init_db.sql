-- Включает автогенерацию id средствами СУБД
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Таблица пользователя
CREATE TABLE user_account (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),    -- Идентификатор пользователя
    username VARCHAR(250) UNIQUE NOT NULL, -- Имя пользователя (логин)
    full_name VARCHAR(250) NULL, -- Полное имя пользователя
    avatar_url TEXT NULL,          -- Ссылка на изображение профиля пользователя
    birth_date TIMESTAMP NULL,    -- Дата рождения пользователя
    travels_count INTEGER DEFAULT 0, -- Количество поездок
    education VARCHAR(250) NULL, -- Образование
    city VARCHAR(250) NULL,      -- Город проживания
    likes INTEGER DEFAULT 0      -- Количество лайков на фото пользователя
);

-- Таблица региона
CREATE TABLE region (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),    -- Идентификатор региона
    country VARCHAR(250) NOT NULL, -- Страна
    object_key VARCHAR(250) UNIQUE NOT NULL, -- Ключ региона (короткое уникальное имя)
    name VARCHAR(250) NOT NULL,   -- Название региона
    img_url TEXT NULL,            -- Ссылка на изображение региона
    tag VARCHAR(100) NULL          -- Тег региона
);

-- Таблица тега
CREATE TABLE tag (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),    -- Идентификатор тега
    object_key VARCHAR(250) UNIQUE, -- Ключ тега (короткое уникальное имя)
    name VARCHAR(50) NOT NULL   -- Название тега
);

-- Таблица места
CREATE TABLE place (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),        -- Идентификатор места
    name VARCHAR(100) NOT NULL, -- Название места
    region_id UUID NOT NULL    -- Ссылка на регион
    -- CONSTRAINT fk_region FOREIGN KEY (region_id) REFERENCES region(id)
);

-- Таблица поездки
CREATE TABLE trip (
    id UUID PRIMARY KEY  DEFAULT gen_random_uuid(),         -- Идентификатор поездки
    name VARCHAR(250) NOT NULL,  -- Название поездки
    description TEXT NULL,       -- Описание поездки
    region_id UUID NOT NULL,      -- Идентификатор региона
    place VARCHAR(250) NULL      -- Место поездки
    -- CONSTRAINT fk_place FOREIGN KEY (region_id) REFERENCES region(id)
);

-- Таблица фотографии
CREATE TABLE photo (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),          -- Идентификатор фото
    img_url TEXT NOT NULL,       -- Ссылка на фото
    description TEXT NULL,        -- Описание фото
    user_id UUID NOT NULL,        -- Ссылка на пользователя
    region_id UUID NOT NULL,      -- Ссылка на регион
    place VARCHAR(250) NULL,      -- Место съемки
    trip_id UUID NULL,            -- Ссылка на поездку
    watched_amount INTEGER DEFAULT 0, -- Количество просмотров фото
    likes INTEGER DEFAULT 0,      -- Количество лайков фото
    coords VARCHAR(100) NULL      -- Координаты места съемки
    -- CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES user(id),
    -- CONSTRAINT fk_region FOREIGN KEY (region_id) REFERENCES region(id),
    -- CONSTRAINT fk_trip FOREIGN KEY (trip_id) REFERENCES trip(id)
);

-- Связь между фото и тегами
CREATE TABLE photo_tag (
    photo_id UUID NOT NULL,       -- Ссылка на фото
    tag_id VARCHAR(50) NOT NULL,  -- Ссылка на тег
    PRIMARY KEY (photo_id, tag_id)
    -- CONSTRAINT fk_photo FOREIGN KEY (photo_id) REFERENCES photo(id),
    -- CONSTRAINT fk_tag FOREIGN KEY (tag_id) REFERENCES tag(id)
);

-- Таблица лайков
CREATE TABLE photo_likes (
    
    photo_id UUID NOT NULL,       -- Ссылка на фото
    user_id UUID NOT NULL,        -- Ссылка на пользователя, который поставил лайк
    PRIMARY KEY (photo_id, user_id)
    -- CONSTRAINT fk_photo FOREIGN KEY (photo_id) REFERENCES photo(id),
    -- CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES user(id)
);

-- Таблица просмотров
CREATE TABLE view (
    photo_id UUID NOT NULL,       -- Ссылка на фото
    user_id UUID NOT NULL,        -- Ссылка на пользователя, который просмотрел фото
    PRIMARY KEY (photo_id, user_id)
    -- CONSTRAINT fk_photo FOREIGN KEY (photo_id) REFERENCES photo(id),
    -- CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES user(id)
);

-- Таблица вызова
CREATE TABLE challenge (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),         -- Идентификатор вызова
    user_id UUID NOT NULL,       -- Ссылка на пользователя
    trip_id UUID NOT NULL,       -- Ссылка на поездку
    description TEXT NULL,       -- Описание вызова
    status VARCHAR(50) NOT NULL -- Статус вызова (например, "active", "completed")
    -- CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES user(id),
    -- CONSTRAINT fk_trip FOREIGN KEY (trip_id) REFERENCES trip(id)
);

-- Таблица для подсчета поездок пользователя
CREATE TABLE user_trip (
    user_id UUID NOT NULL,       -- Ссылка на пользователя
    trip_id UUID NOT NULL,       -- Ссылка на поездку
    PRIMARY KEY (user_id, trip_id)
    -- CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES user(id),
    -- CONSTRAINT fk_trip FOREIGN KEY (trip_id) REFERENCES trip(id)
);