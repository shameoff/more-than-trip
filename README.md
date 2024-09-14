# More than trip 

Приложение для шеринга фотографий от пользователей программы больше чем путешествие


Привет!



### Пересоздание базы
```bash
psql -U more_than_trip -h localhost -p 45432 <<EOF
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
EOF
```