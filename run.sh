#!/usr/bin/env sh

# Установите наименование микро-сервиса
APP_NAME=auth
PROTO_PATHS="internal/api/grpc internal/client/grpc"
UNIT_COVERAGE_MIN=39

####
# Далее менять по согласованию
####

GOARCH=
GOOS=

# Запуск прототул для генерации протобаф-файлов
run_prototool(){
  if [ $OSTYPE == "msys" ]; then
    MSYS_NO_PATHCONV=1 docker run --rm --platform=linux/x86_64 -v "$(pwd):/work" dockerhub.mos.ru/citilink/prototool:1.11.1 $@
  else
    docker run --rm --platform=linux/x86_64 -v "$(pwd):/work" dockerhub.mos.ru/citilink/prototool:1.11.1 $@
  fi
}

# Обрабатывает proto-файлы prototool
process_proto_files(){
  local COMMAND="$1"
  local PROTO_DIR="$2"

  if [ ! -d "$PROTO_DIR" ]; then
    return 0
  fi

  run_prototool prototool "$COMMAND" "$PROTO_DIR"
}

# Генерация proto-файлов
gen_proto(){
  for CURPATH in ${PROTO_PATHS}; do
    echo "start process $CURPATH..."
    rm -Rf "$CURPATH/gen/*"
    process_proto_files all "$CURPATH"
    if [ -d "$CURPATH/gen" ]; then
      run_prototool chown -R "$(id -u)":"$(id -g)" "/work/$CURPATH/gen"
    fi
    echo "finish process $CURPATH..."
  done
}

# Запуск линтера proto-файлов
lint_proto(){
  echo "run proto linter"
  for CURPATH in ${PROTO_PATHS}; do
    process_proto_files lint "$CURPATH"
  done
}

build_win32() {
  echo "Build for Win32"
  GOARCH=386
  GOOS=windows
  return
}

build_win64() {
  echo "Build for Win64"
  GOARCH=amd64
  GOOS=windows
  return
}

build_linux() {
  echo "Build for Linux"
  GOARCH=amd64
  GOOS=linux
  go tool dist install -v pkg/runtime
  go install -v -a std
  return
}

build_darwin() {
    echo "Build for MacOS"
    GOARCH=amd64
    GOOS=darwin
    go install -v -a std
    return
}

build_darwin_arm() {
    echo "Build for MacOS arm"
    GOARCH=arm64
    GOOS=darwin
    go install -v -a std
    return
}

# Сборка приложения
build() {
  if [ -z "$2" ]; then
    echo "Не выбрана система для компиляции"
    return
  fi

  unit
  local APP_PATH="./cmd/"$APP_NAME

  if [ -z "$RELEASE" ]; then RELEASE="undefined"; fi
  echo "RELEASE=$RELEASE"
  if [ -z "$COMMIT" ]; then COMMIT="undefined"; fi
  echo "COMMIT=$COMMIT"

  local BUILD_TIME=$(date '+%Y-%m-%d_%H:%M:%S')

  local LDFFLAGS="-s -w -X 'main.AppName=$APP_NAME' \
    -X 'main.AppRelease=$RELEASE' \
    -X 'main.AppCommit=$COMMIT' \
    -X 'main.AppBuildTime=$BUILD_TIME'"

  case $2 in
    "win32")
        build_win32
        GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="$LDFFLAGS" -o ./bin/$APP_NAME.exe $APP_PATH

        return
        ;;
    "win64")
        build_win64
        local APP_NAME64=$APP_NAME"64"
        GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="$LDFFLAGS" -o ./bin/$APP_NAME64.exe $APP_PATH

        return
        ;;
    "linux")
        build_linux
        GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="$LDFFLAGS" -o ./bin/$APP_NAME $APP_PATH

        return
        ;;
    "darwin")
        build_darwin
        GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="$LDFFLAGS" -o ./bin/$APP_NAME $APP_PATH
        return
        ;;
    "darwin_arm")
        build_darwin_arm
        GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="$LDFFLAGS" -o ./bin/$APP_NAME $APP_PATH
        return
        ;;
  esac

  echo "Неизвестная система для компиляции"
}

# Запуск unit-тестов
unit() {
  go test $(go list ./... | grep -e "git.mos.ru/buch-cloud/moscow-team-2.0/pud/$APP_NAME.git/")
}

# Запуск race-тестов
unit_race() {
  go test -race $(go list ./... | grep -e "git.mos.ru/buch-cloud/moscow-team-2.0/pud/$APP_NAME.git/")
}

# тест на покрытие
unit_coverage() {
  echo "run test coverage"
  go test $(go list ./... | grep -e "git.mos.ru/buch-cloud/moscow-team-2.0/pud/$APP_NAME.git/") -coverprofile=cover_profile.out.tmp $(go list ./internal/...)
  # удаляем protobuf, validate и моки из тестов покрытия
  grep < cover_profile.out.tmp -v -e "mock" -e "\.gen\.go" -e "\.pb\.go" -e "\.pb\.validate\.go" | grep -e "$APP_NAME.git\/" -e "mode:" > cover_profile.out
  rm cover_profile.out.tmp
  CUR_COVERAGE=$( go tool cover -func=cover_profile.out | tail -n 1 | awk '{ print $3 }' | sed -e 's/^\([0-9]*\).*$/\1/g' )
  rm cover_profile.out
  if [ "$CUR_COVERAGE" -lt $UNIT_COVERAGE_MIN ]
  then
    echo "coverage is not enough $CUR_COVERAGE < $UNIT_COVERAGE_MIN"
    return 1
  else
    echo "coverage is enough $CUR_COVERAGE >= $UNIT_COVERAGE_MIN"
  fi
}

# генерация теста на покрытие в виде html
html_unit_coverage() {
  echo "run test coverage to html"
  go test $(go list ./... | grep -e "git.mos.ru/buch-cloud/moscow-team-2.0/pud/$APP_NAME.git/") -coverprofile=cover_profile.out.tmp $(go list ./internal/...)
  # удаляем protobuf, validate и моки из тестов покрытия
  grep < cover_profile.out.tmp -v -e "mock" -e "\.gen\.go" -e "\.pb\.go" -e "\.pb\.validate\.go" | grep -e "$APP_NAME.git\/" -e "mode:" > cover_profile.out
  rm cover_profile.out.tmp
  CUR_COVERAGE=$( go tool cover -func=cover_profile.out | tail -n 1 | awk '{ print $3 }' | sed -e 's/^\([0-9]*\).*$/\1/g' )
  go tool cover -html=cover_profile.out -o cover.html
  rm cover_profile.out
  if [ "$CUR_COVERAGE" -lt $UNIT_COVERAGE_MIN ]
  then
    echo "coverage is not enough $CUR_COVERAGE < $UNIT_COVERAGE_MIN"
    return 1
  else
    echo "coverage is enough $CUR_COVERAGE >= $UNIT_COVERAGE_MIN"
  fi
}

# Настроить прокси
set_private_repo() {
    git config --global credential.helper store
    echo "https://$GOPROXY_LOGIN:$GOPROXY_TOKEN@git.mos.ru" >~/.git-credentials
    go env -w GOPROXY="https://repo-mirror.mos.ru/repository/go-public"
    go env -w GOPRIVATE="git.mos.ru/buch-cloud/moscow-team-2.0/*"
}

# Подтянуть зависимости
deps() {
  go get -t ./...
}

deps_check() {
  echo "run security-check of dependencies"
  go list -json -deps ./cmd/"$APP_NAME" | docker run --rm -i -v $(pwd)/.nancy-ignore:/.nancy-ignore sonatypecommunity/nancy:latest sleuth
}

deps_check_pipe() {
  echo "run security-check of dependencies pipe"
  go list -json -deps ./cmd/"$APP_NAME" | ./.bin/nancy sleuth --exclude-vulnerability-file=$(pwd)/.nancy-ignore
}

newmigrate() {
    local MIGRATENAME=$2
    migrate create -ext sql -dir ./internal/app/migrations -seq $MIGRATENAME
}

# Добавьте сюда список команд
using() {
  echo "Укажите команду при запуске: ./run.sh [command]"
  echo "Список команд:"
  echo "-build <СИСТЕМА> (win32/win64/linux/darwin/darwin_arm) - сборка приложения с необходимой архитектурой"
  echo "-gen_proto - генерация protobuf-файлов"
  echo "-unit - запуск unit-тестов"
  echo "-unit_race - запуск тестов гонки"
  echo "-unit_coverage - запуск тестов покрытия"
  echo "-html_unit_coverage - запуск тестов покрытия с генерацией html файла"
  echo "-deps - загрузка зависимостей"
  echo "-deps_check - проверка зависимостей на уязвимости"
  echo "-newmigrate - создание миграции"
}

############### НЕ МЕНЯЙТЕ КОД НИЖЕ ЭТОЙ СТРОКИ #################

command="$1"
if [ -z "$command" ]; then
  using
  exit 0
else
  $command $@
fi
