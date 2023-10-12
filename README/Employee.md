**Справочник сотрудники**

```json
{
    "Controller" : "Employee",
    ...
}
```

*Запрос сотрудников*

```json
{
    ...
    "Action" : "GetList",
    ...
}
```


Ответ от сервера приходит в формате json

```json
{
    ...
    "Data" : [{
      "ID": 2,
      "DepartmentID": 1,
      "ManualID": 0,
      "LastName": "Иванов",
      "FirstName": "Иван",
      "PatronymicName": "Иванович",
      "Active": false,
      "SortNumber": 200,
      "Position": "",
      "PhoneNumber": "+",
      "Telegram": "",
      "Email": "a@a.ru",
      "Login": "admin",
      "Password": ""
    }]
}
```

p.s. Нет состояний при которых возможно ошибка

| Ключ           | Тип данных | Описание                                                      |
|----------------|------------|---------------------------------------------------------------|
| ID             | Number     | ID записи (встроенное). Если 0, будет создана новая запись    |
| DepartmentId   | Number     | ID департамента уникальный, задается вручную                  |
| ManualID       | Number     | ID сотрудника, задается вручную                               |
| LastName       | string     | Фамилия                                                       |
| FirstName      | string     | Имя                                                           |
| PatronymicName | string     | Отчество                                                      |
| Active         | bool       | Возможность входа в систему                                   |
| SortNumber     | Number     | Номер сортировки                                              |
| Position       | string     | Должность                                                     |
| PhoneNumber    | string     | Телефонный номер (обязательно, полный формат)                 |
| Telegram       | string     | Телеграм имя                                                  |
| Email          | string     | Адрес электронной почты                                       |
| Password       | string     | Пароль (не возвращается), если задан, то будет заменен в базе |


*Запрос обнволение записи*

```json
{
    ...
    "Action" : "Set",
    ...
    "Data" : [{
      "ID": 2,
      "DepartmentID": 1,
      "ManualID": 0,
      "LastName": "Иванов",
      "FirstName": "Иван",
      "PatronymicName": "Иванович",
      "Active": false,
      "SortNumber": 200,
      "Position": "",
      "PhoneNumber": "+",
      "Telegram": "",
      "Email": "a@a.ru",
      "Login": "admin",
      "Password": ""
    }]
}
```


Возможные ошибки

| Ошибка                        | Описание                                             |
|-------------------------------|------------------------------------------------------|
| Invalid structure             | Сервер получил структуру неправильного формата       |
| Some problem with db          | Проблемы возникшие воврем создания/обновления записи |
| Employee ManualID not unique  | ID не уникально                                      |
| Employee Names isn't correct  | Одна или все составляющие ФИО не корректны           |