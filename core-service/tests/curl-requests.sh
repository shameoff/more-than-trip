#!/bin/bash

API_ENDPOINT="https://mtt.shameoff.ru"

# Фото
# Загрузка фото
curl -X POST "${API_ENDPOINT}/api/photo" \
-F "file=@/path/to/photo.jpg" \
-F "metadata={\"Coords\":\"51.5074,-0.1278\",\"Description\":\"A beautiful view of London\",\"Place\":\"London\",\"RegionId\":\"550e8400-e29b-41d4-a716-446655440000\",\"TripId\":\"550e8400-e29b-41d4-a716-446655440001\",\"UserId\":\"550e8400-e29b-41d4-a716-446655440002\"}"

# Получение фото по UUID
curl -X GET "${API_ENDPOINT}/api/photo/550e8400-e29b-41d4-a716-446655440000"

# Получение всех фото с фильтрами (по региону или тегу)
curl -X GET "${API_ENDPOINT}/api/photos?region=550e8400-e29b-41d4-a716-446655440000&tag=nature"

# Удаление фото по UUID
curl -X DELETE "${API_ENDPOINT}/api/photo/550e8400-e29b-41d4-a716-446655440000"

# Обновление фото по UUID
curl -X PUT "${API_ENDPOINT}/api/photo/550e8400-e29b-41d4-a716-446655440000" \
-H "Content-Type: application/json" \
-d '{"Coords": "51.5074,-0.1278", "Description": "Updated description", "Place": "London", "RegionId": "550e8400-e29b-41d4-a716-446655440000", "TripId": "550e8400-e29b-41d4-a716-446655440001", "UserId": "550e8400-e29b-41d4-a716-446655440002"}'

# Лайк фото
curl -X POST "${API_ENDPOINT}/api/photo/550e8400-e29b-41d4-a716-446655440000/like" \
-H "User-Id: 550e8400-e29b-41d4-a716-446655440002"

# Дизлайк фото
curl -X POST "${API_ENDPOINT}/api/photo/550e8400-e29b-41d4-a716-446655440000/dislike" \
-H "User-Id: 550e8400-e29b-41d4-a716-446655440002"

# Регионы
# Создание региона
curl -X POST "${API_ENDPOINT}/api/region" \
-H "Content-Type: application/json" \
-d '{"Id":"550e8400-e29b-41d4-a716-446655440000", "Name":"London", "Country":"UK", "ObjectKey":"london"}'

# Удаление региона по UUID
curl -X DELETE "${API_ENDPOINT}/api/region/550e8400-e29b-41d4-a716-446655440000"

# Получение региона по UUID
curl -X GET "${API_ENDPOINT}/api/region/550e8400-e29b-41d4-a716-446655440000"

# Получение всех регионов
curl -X GET "${API_ENDPOINT}/api/regions"

# Обновление региона
curl -X PUT "${API_ENDPOINT}/api/region/550e8400-e29b-41d4-a716-446655440000" \
-H "Content-Type: application/json" \
-d '{"Name":"London Updated", "Country":"UK", "ObjectKey":"london-updated"}'

# Теги
# Создание тега
curl -X POST "${API_ENDPOINT}/api/tag" \
-H "Content-Type: application/json" \
-d '{"Id":"550e8400-e29b-41d4-a716-446655440000", "Name":"nature"}'

# Удаление тега по UUID
curl -X DELETE "${API_ENDPOINT}/api/tag/550e8400-e29b-41d4-a716-446655440000"

# Получение всех тегов
curl -X GET "${API_ENDPOINT}/api/tags"

# Вызовы (Challenges)
# Создание вызова
curl -X POST "${API_ENDPOINT}/api/challenge"

# Удаление вызова
curl -X DELETE "${API_ENDPOINT}/api/challenge/550e8400-e29b-41d4-a716-446655440000"

# Обновление вызова
curl -X PUT "${API_ENDPOINT}/api/challenge/550e8400-e29b-41d4-a716-446655440000"

# Получение вызова
curl -X GET "${API_ENDPOINT}/api/challenge/550e8400-e29b-41d4-a716-446655440000"

# Получение всех вызовов
curl -X GET "${API_ENDPOINT}/api/challenges"

# Получение моих вызовов
curl -X GET "${API_ENDPOINT}/api/mychallenges"

# Пользователи
# Создание пользователя
curl -X POST "${API_ENDPOINT}/api/user"

# Удаление пользователя
curl -X DELETE "${API_ENDPOINT}/api/user/550e8400-e29b-41d4-a716-446655440000"

# Обновление пользователя
curl -X PUT "${API_ENDPOINT}/api/user/550e8400-e29b-41d4-a716-446655440000"

# Получение пользователя
curl -X GET "${API_ENDPOINT}/api/user/550e8400-e29b-41d4-a716-446655440000"

# Получение всех пользователей
curl -X GET "${API_ENDPOINT}/api/users"

# Поездки
# Создание поездки
curl -X POST "${API_ENDPOINT}/api/trip"

# Удаление поездки
curl -X DELETE "${API_ENDPOINT}/api/trip/550e8400-e29b-41d4-a716-446655440000"

# Обновление поездки
curl -X PUT "${API_ENDPOINT}/api/trip/550e8400-e29b-41d4-a716-446655440000"

# Получение поездки
curl -X GET "${API_ENDPOINT}/api/trip/550e8400-e29b-41d4-a716-446655440000"

# Получение всех поездок
curl -X GET "${API_ENDPOINT}/api/trips"