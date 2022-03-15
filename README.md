# ad-service-demo
REST API в формате JSON для сайта объявлений

## Установка
Скопируйте исходный код из репозитория 
```
git clone https://github.com/pshvedko/adservice.git
cd adservice
```
## Запуск
Соберите сервис в докере и запустите
```
docker-compose up --build
```
## Использование
Сервис можно использовать, например, с помощью `curl`
### Метод получения списка объявлений
Параметры запроса: 
* `sort` - разделенный через запятую список полей `price` `date`, 
`-` для обратной сортировки
* `limit` - ограничение размера списка
* `offset` - пропуск элементов списка
```
curl -v "http://localhost:8080/api/v1/?sort=price,-date&limit=3&offset=3"
```
### Метод получения конкретного объявления
Параметры запроса:
* `{ID}` - идентификатор объявления
* `field` - разделенный через запятую список _дополнительных_ полей 
в ответе `photo` `date` `description`
```
curl -v "http://localhost:8080/api/v1/{ID}?field=photo,date,description"
```
### Метод создания объявления
Параметры запроса:
* `subject` - название объявления
* `description` - содержание 
* `price` - цена
* `photo` - добавляет ссылку на фото в массив ссылок
```
curl -v -X POST http://localhost:8080/api/v1/ \
        -F photo=http://localhost/1.png \
        -F photo=http://localhost/2.png \
        -F photo=http://localhost/3.png \
        -F subject="..." \
        -F description="..." \
        -F price=1.23 
```
Посылка JSON запроса:
* `ad.json` - файл с объявлением в JSON формате
```
{
  "price": 1.23,
  "subject": "...",
  "description": "...",
  "photo": [
    "https://localhost/1.png",
    "https://localhost/2.png",
    "https://localhost/3.png"
  ]
}
```
```
curl -v -X POST http://localhost:8080/api/v1/ \
        -H "Content-Type: application/json; charset=utf-8" \
        --data "@ad.json" 
```