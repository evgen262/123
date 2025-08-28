## Документация Сервиса Кадровых Событий (СКС) функционального блока "Ведение реестра кадров"

При регистрации в СКС выдаются следующие данные, необходимые для формирования запросов:
- уникальный идентификатор системы-абонента, зарегистрированной в СКС (SubscriberID)
- уникальный идентификатор пользователя указанной в SubscriberID системы-абонента, зарегистрированный в СКС (UserID)
- пароль системы-абонента, зарегистрированной в СКС (SubscriberPassword)

### Параметры авторизации 
#### в Get параметры во всех запросах:
- ```SubscriberID``` идентификатор абонента
- ```UserID``` идентификатор пользователя

#### в Header:
```Authorization: Basic token``` \
где ```token``` - строка ```SubscriberID:SubscriberPassword``` закодированная в формате base64
- ```SubscriberID``` идентификатор абонента
- ```SubscriberPassword``` пароль абонента


### Для получения информации о физическом лице (организация, подразделение, должность, вид занятости, дата назначения) необходимо обратиться к СКС по методу sksMobileApp()

method: ```POST``` \
url в тестовой среде: ```https://predprod-kadry2.mos.ru/{database_name}/hs/frontend_api/execute/sksMobileApp``` \
url в прод среде: ```https://kadry2.mos.ru/{database_name}/hs/frontend_api/execute/sksMobileApp``` \
где ``{database_name}`` - имя базы данных

В теле запроса в формате JSON передаются следующие параметры:
- ```PersonIDArray``` список уникальных идентификаторов физических лиц (в списке физических лиц), 
данные которых будут включены в ответ на запрос. Обязательный параметр.
- ```OrgIDArray``` список уникальных идентификаторов организаций (в списке организаций), данные которых будут включены в ответ
- ```AttributeList``` список включаемых/исключаемых свойств метода из числа доступных пользователю \
    ```Include``` массив наименований свойств метода, которые будут включены в ответ \
    ```Exclude``` массив наименований свойств метода, которые будут исключены из ответа \
  если не указан```AttributeList```, то по умолчанию в ответ будут включены все доступные свойства 

#### Пример запроса:

```
curl --location 'https://predprod-kadry2.mos.ru/hr5_rk/hs/frontend_api/execute/sksMobileApp?SubscriberID=DIT&UserID=0c0fe343-01f9-105e-9a11-01a275ae0c6b&GetReferenceInfo=true' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic ZXhhbXBsZV9pZDpleGFtcGxlX3Bhc3M' \
--data '{
    "PersonIDArray": [
        "0cd9e619-f1ba-682a-952e-00f0562c9a04"
    ]
}'
```

#### В случае успешного ответа:

```
{
    "MessageType": "Result",
    "RequestExecuted": true,
    "RequestType": "sksMobileApp",
    "ResponceBody": {
        "MobileApp": [
            {
                "PersonID": "0cd9e619-f1ba-682a-952e-00f0562c9a04",
                "FIOPerson": "Иванов Иван Иванович",
                "SNILS": "123-456-789 00",
                "OrgID": "6cf01c27-af06-90fe-11ea-005056a2c924",
                "InnOrg": "77017654321",
                "NameOrg": "ГАУ МЕДИАЦЕНТР",
                "SubdivID": "6cfc4614-af06-90fe-11ea-005056a2c924",
                "NameSubdiv": "Главный отдел",
                "PositionID": "66f6a677-f730-90fe-11ea-005056a2c924",
                "NamePosition": "Начальник отдела",
                "EmploymentType": "",
                "DateRecept": "2013-05-09T10:00:00"
            },
            {
                "PersonID": "0cd9e619-f1ba-682a-952e-00f0562c9a04",
                "FIOPerson": "Иванов Иван Иванович",
                "SNILS": "123-456-789 00",
                "OrgID": "6528c9f9-e96f-90fe-11ea-005056a2c924",
                "InnOrg": "77011234567",
                "NameOrg": "ГКУ «Инфогород»",
                "SubdivID": "73a15191-6df6-90fe-11ea-005056a2c924",
                "NameSubdiv": "Отдел начальников",
                "PositionID": "6e9f6f0e1-6df6-90fe-11ea-005056a2c924",
                "NamePosition": "ведущий начальник",
                "EmploymentType": "Совместительство",
                "DateRecept": "2021-01-01T00:10:01"
            }
        ]
    }
}
```

#### Содержание ответа: 

```PersonID```  уникальный идентификатор физического лица (передаваемый в PersonIDArray) \
```FIOPerson```  Фамилия, Имя, Отчество физического лица \
```SNILS```  СНИЛС физического лица \
```OrgID```  уникальный идентификатор организации \
```InnOrg```  ИНН организации \
```NameOrg```  наименование организации \
```SubdivID```  уникальный идентификатор подразделения \
```NameSubdiv```  наименование подразделения \
```PositionID```  уникальный идентификатор должности \
```NamePosition```  наименование должности \
```EmploymentType```  вид занятости \
```DateRecept```  дата приема 


#### Пример запроса с фильтрацией (в ответе будут только ```PersonID``` и ```InnOrg``` ):

```
curl --location 'https://predprod-kadry2.mos.ru/hr5_rk/hs/frontend_api/execute/sksMobileApp?SubscriberID=DIT&UserID=0c0fe343-01f9-105e-9a11-01a275ae0c6b&GetReferenceInfo=true' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic ZXhhbXBsZV9pZDpleGFtcGxlX3Bhc3M' \
--data '{
    "PersonIDArray": [
        "0cd9e619-f1ba-682a-952e-00f0562c9a04"
    ],
    "AttributeList": {
        "MobileApp": {
            "Include": [
                "PersonID",
                "InnOrg"
            ]
        }
    }
}'
```
