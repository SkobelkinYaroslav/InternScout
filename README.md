
# InternScout

InternScout - это бот-агрегатор телеграм каналов со стажировками для студентов. Бот автоматически парсит каналы и отправляет уведомления о новых стажировках, соответствующих вашим ключевым словам.

## Установка

Следуйте этим шагам для установки и запуска InternScout:

### 1. Клонируйте репозиторий

```sh
git clone https://github.com/SkobelkinYaroslav/InternScout.git
```

### 2. Перейдите в директорию проекта

```sh
cd InternScout
```

### 3. Создайте файл .env

Создайте файл `.env` в корневой директории проекта со следующим содержимым:

```
API_KEY=<ваш апи ключ>
```

### 4. Создайте файл конфигураций

Создайте файл `config.json` в корневой директории проекта следующего вида:

```json
[
  {
    "id": <ваш телеграм айди>,
    "categories": [
      "ключевые слова",
      "со",
      "со строчной буквы",
      "которые вы хотите искать"
    ]
  }
]
```

### 5. Соберите проект

```sh
go build -o main ./cmd/scout/main.go
```

### 6. Запустите проект

```sh
./main
```

## Добавление каналов

Для добавления телеграм-каналов используйте команду ```/addchannels <телеграм канал>``` в сообщениях боту

## Пример использования

### Пример .env файла

```
API_KEY=123456789:ABCDEFGHIJKLMNOPQRSTUVWXYZ
```

### Пример config.json файла

```json
[
  {
    "id": 123456789,
    "categories": [
      "internship",
      "junior",
      "стажировка",
      "практика"
    ]
  }
]
```

## Contributing

Если у вас есть предложения или исправления, пожалуйста, создайте pull request или откройте issue на GitHub.

