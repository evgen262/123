package view

const (
	ErrStatusTextStateInvalid           string = "неверный ключ авторизации"
	ErrStatusTextCodeNotProvided        string = "не предоставлен код авторизации"
	ErrStatusTextCallbackURLNotProvided string = "не предоставлен url для перенаправления"
	ErrStatusTextBadCodeProvided        string = "предоставлен неверный код авторизации"
	ErrStatusTextServiceError           string = "внутренняя ошибка сервиса авторизации"
	ErrStatusTextCloudIDNotProvided     string = "не предоставлен cloud_id"
	ErrStatusTextBadRefreshToken        string = "refresh token недействителен или устарел"
	ErrStatusTextRefreshTokenNotFound   string = "refresh token не найден"
	ErrStatusTextCredentialsNotProvided string = "не предоставлены данные для авторизации"
	ErrStatusTextBadCredentialsProvided string = "неверные реквизиты для авторизации"
	ErrStatusTextAccessTokenNotProvided string = "не предоставлен access_token"
	ErrStatusTextAccessTokenExpired     string = "access_token истёк"
	ErrStatusTextAccessDenied           string = "доступ запрещён"
	ErrStatusTextCallbackURLMismatch    string = "callback_url не соответствует url в запросе"
)
