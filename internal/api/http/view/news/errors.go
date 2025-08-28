package news

import "git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"

const (
	ErrAuthorNotFound        diterrors.StringError = "Автор не найден"
	ErrCategoryAlreadyExists diterrors.StringError = "Категория уже существует"
	ErrCategoryNotFound      diterrors.StringError = "Категория не найдена"
	ErrTitleRequired         diterrors.StringError = "Заголовок обязателен"
	ErrSlugRequired          diterrors.StringError = "Адрес ссылки обязателен"
	ErrCategoryRequired      diterrors.StringError = "Категория обязательна"
	ErrPublishTimeBeforeNow  diterrors.StringError = "Дата публикации не может быть в прошлом"
	ErrStatus                diterrors.StringError = "Некорректный статус новости"
	ErrNewsNotFound          diterrors.StringError = "Новость не найдена"
)
