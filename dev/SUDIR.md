## Документация СУДИР

С документацией по работе СУДИР можно ознакомиться [по ссылке](https://sudir.mos.ru/support/).

СУДИР - сервис обеспечивающий идентификацию и аутентификацию пользователей при доступе в различные информационные 
системы и приложения ОИС города Москвы. \
Подключение к СУДИР осуществляется с использованием протокола OIDC/OAuth 2.0.

При регистрации сервиса в СУДИР выдаются следующие данные, необходимые для формирования запросов:
- идентификатор приложения (client_id)
- секрет приложения (client_secret)

Приложение взаимодействует с сервисами СУДИР, используя следующие адреса:
  - URL для проведения авторизации и аутентификации: \
    ```https://sudir-test.mos.ru/blitz/oauth/ae``` (тестовая среда) \
    ```https://sudir.mos.ru/blitz/oauth/ae``` (продуктивная среда)
  - URL для получения и обновления маркера доступа: \
    ```https://sudir-test.mos.ru/blitz/oauth/te``` (тестовая среда) \
    ```https://sudir.mos.ru/blitz/oauth/te``` (продуктивная среда)
  - URL для выполнения логаута: \
    ```https://sudir-test.mos.ru/blitz/login/logout``` (тестовая среда) \
    ```https://sudir.mos.ru/blitz/login/logout``` (продуктивная среда)

### При первоначальном подключении, пользователь перенаправляется в СУДИР на URL для получения кода авторизации.

url в тестовой среде: ```https://sudir-test.mos.ru/blitz/oauth/ae``` \
url в прод среде: ```https://sudir.mos.ru/blitz/oauth/ae``` \
method: ```GET```

#### Параметры запроса:
- ```access_type``` требуется ли приложению получать refresh_token, необходимый для получения сведений о пользователе в дальнейшем,
  когда пользователь будет офлайн. \
  Используем ```access_type=offline```
- ```client_id```  идентификатор клиента. \
  Используем ```client_id=cfc-zv.ditcloud.ru```
- ```redirect_uri``` ссылка для возврата пользователя в наш сервис, ссылка должна соответствовать одному из зарегистрированных значений. \
  Используем ```redirect_uri=https://router.ditcloud.ru/router/hs/oid2op```
- ```response_type``` тип ответа может принимать значения: \
  ```code``` - Authorization Code Flow, \
  ```token``` - OAuth 2.0 Implicit Flow, \
  ```code token```, ```code id_token```, ```code id_token token``` - Hybrid Flow, \
  ```id_token token```, ```id_token``` -  OIDC Implicit Flow \
  Используем ```response_type=code```
- ```scope``` запрашиваемые разрешения, для проведения аутентификации должно быть
  передано разрешение ```openid``` и необходимые дополнительные scope для получения
  данных пользователя (см Таблица 1). \
  Используем ```scope=openid+profile+email+userinfo+employee+groups```
- ```state``` набор случайных символов, имеющий вид 128-битного идентификатора
  запроса (используется для защиты от перехвата), это же значение будет возвращено
  в ответе.

<b>Таблица 1 - доступные разрешения (scope)</b>
<table>
<thead>
  <tr>
    <th>scope</th>
    <th>описание</th>
    <th>получаемые аттрибуты</th>
  </tr>
</thead>
<tbody>
<tr>
    <td>openid</td>
    <td> разрешение, указывающее на то, что аутентификация проводится согласно спецификации OIDC 1.0</td>
    <td>При запросе этого scope СУДИР предоставляет приложению id_token</td>
</tr>
<tr>
    <td>profile</td>
    <td>Основные данные профиля пользователя</td>
    <td>
- <b>sub</b> – уникальный идентификатор (userPrincipalName из ЕСК) <br>
- <b>family_name</b> – фамилия <br>
- <b>name</b> – имя <br>
- <b>middle_name</b> – отчество <br>
    </td>
</tr>
<tr>
    <td>email</td>
    <td>Адрес электронной почты сотрудника</td>
    <td>- <b>email</b> – служебный адрес электронной почты</td>
</tr>
<tr>
    <td>employee</td>
    <td>Служебные данные сотрудника</td>
    <td>
- <b>company</b> – название организации сотрудника. <br>
- <b>department</b> – подразделение организации, где работает сотрудник <br>
- <b>workphone</b> – служебный номер телефона сотрудника <br>
- <b>position</b> – должность <br>
- <b>email</b> – служебный адрес электронной почты <br>
</td>
</tr>
<tr>
    <td>groups</td>
    <td>Группы пользователя в ЕСК</td>
    <td> - <b>groups</b> – список групп пользователя</td>
</tr>
</tbody>
</table>

#### Пример запроса:
```
curl --location 'https://sudir-test.mos.ru/blitz/oauth/ae?access_type=offline&client_id=cfc-zv.ditcloud.ru&redirect_uri=https%3A%2F%2Frouter.ditcloud.ru%2Frouter%2Fhs%2Foid2op&response_type=code&scope=openid%2Bprofile%2Bemail%2Buserinfo%2Bemployee%2Bgroups&state=342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f'
```

#### В случае успешного ответа пользователь будет переадресован по адресу со значением кода авторизации (code) и параметром state:
```
https://router.ditcloud.ru/router/hs/oid2op?code=XXijKLnoayN8xX9ap3JUMLBATcv2m8V9irtIvrKEtofDfkH0eGe2XAb9qzvKJXXuuJHtVokGL_iW5urIGNtQgWuhV1MK46C4D4jPEMl77-Q&state=342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f
```

### Получение маркеров из СУДИР

После получения кода авторизации необходимо обменять его на маркеры.

url в тестовой среде: ```https://sudir-test.mos.ru/blitz/oauth/te``` \
url в прод среде: ```https://sudir.mos.ru/blitz/oauth/te``` \
method: ```POST``` 

Заголовки ```Authorization: Basic token``` \
```token``` - строка ```client_id:client_secret``` закодированная в формате base64
 - ```client_id``` идентификатор приложения
 - ```client_secret``` секрет приложения \
полученные при регистрации в СУДИР

В запросе передаются следующие параметры:
  - ```code``` код авторизации полученный в предыдущем запросе
  - ```grant_type``` принимает значение ```authorization_code``` т.к. код авторизации меняем на маркер доступа
  - ```redirect_uri``` ссылка для возврата пользователя в наш сервис

#### Пример запроса:

```
curl --location --request POST 'https://sudir-test.mos.ru/blitz/oauth/te?grant_type=authorization_code&code=XXijKLnoayN8xX9ap3JUMLBATcv2m8V9irtIvrKEtofDfkH0eGe2XAb9qzvKJXXuuJHtVokGL_iW5urIGNtQgWuhV1MK46C4D4jPEMl77-Q&redirect_uri=https%3A%2F%2Frouter.ditcloud.ru%2Frouter%2Fhs%2Foid2op' \
--header 'Authorization: Basic ZXhhbXBsZV9pZDpleGFtcGxlX3Bhc3M'
```

#### В случае успешного ответа:

```
{
    "id_token": "eyJraWQiOiJkZWZhdWx0IiwiYWxnIjoiUlMyNTYifQ.eyJlbWFpbCI6ImRldmVsb3BlckBpdC5tb3MucnUiLCJsb2dvbm5hbWUiOiJEZXZlbG9wZXIiLCJ1YV9pZCI6IlNIQTI1Nl91TVcwN0xaWWUtcmVvbHhwbjVqMGRkNmhUcm5zT09DYmJuVlNTVVlnZ3VjIiwianRpIjoiSzVIaU1nblNCYW0tSDR6dmhzM3NROHFDRjdaU0pTek9ZTUJWS2pheDJzYyIsImV4cCI6MTY4ODQxNjEwMCwiY2xvdWRHVUlEIjoiM2M1Y2JiMTYtMDExYS0zMTBlLTk3ZTItNTY1NDAwYTI2NTA2IiwiaWF0IjoxNjg4MzgzNzAwLCJhdWQiOlsiY2ZjLXp2LmRpdGNsb3VkLnJ1Il0sImFtciI6WyJwYXNzd29yZCJdLCJpc3MiOiJodHRwczovL3N1ZGlyLXRlc3QubW9zLnJ1IiwiY29tcGFueSI6ItCT0JrQoyDQmNC90YTQvtCz0L7RgNC-0LQiLCJkZXBhcnRtZW50Ijoi0J7RgtC00LXQuyDQvNC-0LHQuNC70YzQvdC-0Lkg0YDQsNC30YDQsNCx0L7RgtC60LgiLCJwb3NpdGlvbiI6ItC_0YDQvtCz0YDQsNC80LzQuNGB0YIiLCJzdWIiOiJEZXZlbG9wZXJAaHEuY29ycC5tb3MucnUiLCJjcmlkIjoiMCIsInNpZCI6IjUxZTEwMzYxLThhZTUtNGYwMS04ZjRmLWM3YjE4ZTA4ODc2ZSJ9.XOrKFMHSJFYbo-Fq8I5Yd9znJh6prA2t86JX89FrrRGO6r-0n-2T7VTeHwd0TP7loDQQOCeBd2NG-4wkSAnDRtBBQqk0jgQjcdQwpYXBaZonnWYwRFyBI6nYlxo5Iq5DSasEKt7kYJ6PpMF7Pcp8jfauYB4wGPGmFrf_PkpXUZZqftDYqWRGeCaguoPKyUoOGEtUNfyDo7pK5T2RmUdBoxu63qs-Z9ot0ZUzJ0ZxrwbIDDGv4PQH_edj4Wtix4oWP0HKxkALzAPHkQEuUv-H6gYwxSg0qXgLLOoh_Zd8alIwLhMUfjg71rpkLkN6ZRYdrb3s8dCB_FyesI_DjquY0w",
    "refresh_token": "WNjZilQssfRetWT81slSAp-KTJFYWNRjK9yFMg3YGmySKT64TEiKGuBX_kRaTRImpzt98llHXwpCl9C1HyHALg",
    "access_token": "wjsTCFTERhLa86xYLy4mPZjvx7RSHV9oUwQ8V3zxKMs1MDMxMTZlMS04Yjk1LTRmNDEtOGYwZi1jNzEwOGI4NzY4ZWU",
    "scope": "email userinfo employee groups openid profile",
    "token_type": "Bearer",
    "expires_in": 3600
}
```

#### Содержание ответа:
- ```access_token``` маркер доступа к защищенному ресурсу, например, к данным пользователя.
- ```refresh_token``` обновление маркера доступа (выдаётся при указании параметра access_type=offline \
в запросе кода авторизации). Маркер действителен до момента использования, но не дольше 365 дней.
- ```id_token``` маркер идентификации содержит всю запрашиваемую информацию в формате JWT.
- ```scope``` запрошенные разрешения
- ```token_type``` тип маркера доступа
- ```expires_in``` время в секундах через которое токен будет не валидным

Header ```id_token``` содержит следующую информацию:
```
{
  "kid": "default",
  "alg": "RS256"
}
```
- ```kid``` идентификатор ключа которым подписан токен
- ```alg``` криптографический алгоритм шифрования


Payload```id_token``` содержит следующую информацию:
```
{
  "cloudGUID": "3c5cbb16-011a-310e-97e2-565400a26506",
  "logonname": "Developer",
  "email": "developer@it.mos.ru",
  "company": "ГКУ Инфогород",
  "department": "Отдел мобильной разработки",
  "position": "программист",
  "ua_id": "SHA256_uMW07LZYe-reolxpn5j0dd6hTrnsOOCbbnVSSUYgguc",
  "exp": 1688416100,
  "iat": 1688383700,
  "jti": "K5HiMgnSBam-H4zvhs3sQ8qCF7ZSJSzOYMBVKjax2sc",
  "iss": "https://sudir-test.mos.ru",
  "sub": "Developer@hq.corp.mos.ru",
  "aud": [
    "cfc-zv.ditcloud.ru"
  ],
  "amr": [
    "password"
  ],
}
```
- ```cloudGUID``` уникальный идентификатор пользователя
- ```logonname``` логин пользователя в домене
- ```email``` служебный адрес электронной почты
- ```company``` название организации сотрудника
- ```department``` подразделение организации, где работает сотрудник
- ```position``` занимаемая должность
- ```exp``` время в формате Unix Time до которого токен будет валидным
- ```iat``` время создания токена в формате Unix Time 
- ```jti```  уникальный идентификатор токена (JWT ID)
- ```iss```  уникальный идентификатор стороны, генерирующей токен
- ```sub```  уникальный идентификатор стороны, о которой содержится информация в данном токене
- ```aud``` список получателей данного токена
- ```amr``` методы аутентификации

### Обновление маркера доступа в СУДИР

Для обновления маркера доступа необходимо сформировать запрос методом POST на URL для получения маркера в СУДИР

url в тестовой среде: ```https://sudir-test.mos.ru/blitz/oauth/te``` \
url в прод среде: ```https://sudir.mos.ru/blitz/oauth/te``` \
method: ```POST``` 

Заголовки ```Authorization: Basic token``` \
```token``` - строка ```client_id:client_secret``` закодированная в формате base64
- ```client_id``` идентификатор приложения
- ```client_secret``` секрет приложения \
  полученные при регистрации в СУДИР

В теле запроса передаются следующие параметры:
- ```refresh_token```маркера обновления полученный в предыдущем запросе
- ```grant_type``` принимает значение ```refresh_token``` т.к. маркер обновления обменивается на маркер доступа

```
curl --location --request POST 'https://sudir-test.mos.ru/blitz/oauth/te?grant_type=refresh_token&refresh_token=OpE_Tvh73mHkh4nBIEPMNQA7vuhPt9yiSBBn_vHBD1LxjZbJ-EvsR1StWzxVRYyYCkmhYLLdhBrayxDKOC7FqA' \
--header 'Authorization: Basic ZXhhbXBsZV9pZDpleGFtcGxlX3Bhc3M'
```

#### В случае успешного ответа:
```
{
    "access_token": "KPIerz_d2rvXXNBbzgUhcY8OD3R0EciVR1ifLPLpzsY",
    "refresh_token": "tPtOSo0HQD1lfWMAKBWMZ0rHZNqNq6oUI1aqKDV0tUiJXOrOyke3c9fw8MwkYnToDPxoNOo1DswG0t0HAwuFkA",
    "scope": "email userinfo employee groups openid profile",
    "token_type": "Bearer",
    "expires_in": 3600
}
```

#### Содержание ответа:
- ```access_token``` маркер доступа к защищенному ресурсу, например, к данным пользователя.
- ```refresh_token``` обновление маркера доступа. Маркер действителен до момента использования, но не дольше 365 дней.
- ```scope``` запрошенные разрешения
- ```token_type``` тип маркера доступа
- ```expires_in``` время в секундах через которое токен будет не валидным


### Сертификат открытого ключа СУДИР можно загрузить по следующим ссылкам

Сертификат необходимый для проверки подписи находится в атрибуте x5c

url в тестовой среде: ```https://sudir-test.mos.ru/blitz/oauth/.well-known/jwks``` \
url в прод среде: ```https://sudir.mos.ru/blitz/oauth/.well-known/jwks``` \
method: ```GET```

Пример ответа:

```
{
    "keys": [
        {
            "kty": "RSA",
            "n": "kk8cC1R8rwx0FGEyH0aUpnDbeIRU5b-njho_JzwSeDg3P7lDZD63w-P8vyShvW9QMC_pjeUNGLiU8GJYEZrrEh00Capn5RB-X6hXdBT64S3fHOQVB0IUBcNA_4TKIB8pdLlcKCtTwsFaGGIBmR0ghKmE3k7tHZRcx0Vltx19Bg_L0tmI56sb2qYwzSAbhpn_TTHN3dYu08Hf8mz1y0Np8GCgPglXgMSGAxe9WzYYl_oiCheqaNm9tIhpFlUgH4IRFliaVyQwmnSY15Hh5DnJY08YPiP0D95W8W7eBmvGGOm8drhGx6mKfj_iyxfNCEb7UOqQk0bZJLcnfJnpUuSE1Q",
            "e": "AQAB",
            "use": "sig",
            "alg": "RS256",
            "kid": "default",
            "x5c": [
                "MIIDBTCCAe2gAwIBAgIEeiwqwjANBgkqhkiG9w0BAQsFADAzMREwDwYDVQQDEwhCbGl0eklkUDEeMBwGA1UEAwwVandzX3JzMjU2X3JzYV9kZWZhdWx0MB4XDTE4MDUyNTA3MTcyMFoXDTI4MDUyMjA3MTcyMFowMzERMA8GA1UEAxMIQmxpdHpJZFAxHjAcBgNVBAMMFWp3c19yczI1Nl9yc2FfZGVmYXVsdDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAJJPHAtUfK8MdBRhMh9GlKZw23iEVOW/p44aPyc8Eng4Nz+5Q2Q+t8Pj/L8kob1vUDAv6Y3lDRi4lPBiWBGa6xIdNAmqZ+UQfl+oV3QU+uEt3xzkFQdCFAXDQP+EyiAfKXS5XCgrU8LBWhhiAZkdIISphN5O7R2UXMdFZbcdfQYPy9LZiOerG9qmMM0gG4aZ/00xzd3WLtPB3/Js9ctDafBgoD4JV4DEhgMXvVs2GJf6IgoXqmjZvbSIaRZVIB+CERZYmlckMJp0mNeR4eQ5yWNPGD4j9A/eVvFu3gZrxhjpvHa4Rsepin4/4ssXzQhG+1DqkJNG2SS3J3yZ6VLkhNUCAwEAAaMhMB8wHQYDVR0OBBYEFLTHv4QDOfZ+tNqFel5EN55nMWtxMA0GCSqGSIb3DQEBCwUAA4IBAQB99QaK7c3UnKz03kpKYrd9H6vMJbJLyvA0q8Vlwrxz2qhuQ+FUCZsAhs3qUQfCXdp+htGNoJC8PDsB1JgLg/6hCXBAJf+w4u/UIbpmoyve3hPCvV+RPIOInqW+Po5xLHcio1JO8iwRDfta+IL3lkvEqOrxef1Y4j48WNvaR/p319LsLJ4peZp+BB4Y/A119QKdN+9Ze1PQmPxNG9HkZFS9tOjlXeUkpdshIoWHCbpovkZefQLNuowQ8V0IKDxOAgSjgrhrWPW209hXJ2aTrIA3yE/4jd3Cv0cLpx91qMw76M8L9H4PALYD5ijhr8d0hFDrxr9fVW7eNR2o3IbVWi/y"
            ],
            "x5t": "7HmYY6dKMEaPuv05hcAbS_oWd8E"
        }
    ]
}
```