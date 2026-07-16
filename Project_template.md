## Изучите [README.md](README.md) файл и структуру проекта.

## Задание 1

Cinema Abyss - Диаграмма контейнеров - TO-BE (через 2 месяца)

![Cinema Abyss - Диаграмма контейнеров - TO-BE (через 2 месяца)](./arch/to-be-container.png)

[PUML: Cinema Abyss - Диаграмма контейнеров - TO-BE (через 2 месяца)](./arch/to-be-container.puml)


## Задание 2

Proxy и events сервисы написаны на языке golang.

В docker-compose.yml добавлены healthcheck для kafka и зависимости для kafka и БД (service_healthy) для сервисов monolith и movies-service, иначе они запускались слишком рано и падали.

**Результаты выполнения тестов**

![Результаты выполнения тестов](./images/task-2-local-postman-tests.png)

**Результаты запросов curl**

```
MAXIM@LAPTOP-123 MINGW64 /
$ curl http://localhost:8000/api/movies
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  1566  100  1566    0     0   194k      0 --:--:-- --:--:-- --:--:--  218k[{"id":1,"title":"The Shawshank Redemption","description":"Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.","genres":["Drama"],"rating":9.3},{"id":2,"title":"The Godfather","description":"The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.","genres":["Crime","Drama"],"rating":9.2},{"id":3,"title":"The Dark Knight","description":"When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.","genres":["Action","Crime","Drama"],"rating":9},{"id":4,"title":"Pulp Fiction","description":"The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.","genres":["Crime","Drama"],"rating":8.9},{"id":5,"title":"Forrest Gump","description":"The presidencies of Kennedy and Johnson, the Vietnam War, the Watergate scandal and other historical events unfold from the perspective of an Alabama man with an IQ of 75, whose only desire is to be reunited with his childhood sweetheart.","genres":["Drama","Romance"],"rating":8.8},{"id":6,"title":"Test Movie 496","description":"A test movie created by automated tests","genres":["Action","Drama"],"rating":4.5},{"id":7,"title":"Microservice Test Movie 769","description":"A test movie created by automated tests for the microservice","genres":["Sci-Fi","Thriller"],"rating":4.8}]


MAXIM@LAPTOP-123 MINGW64 /
$ curl http://localhost:8000/api/users
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   238  100   238    0     0  12936      0 --:--:-- --:--:-- --:--:-- 13222[{"id":1,"username":"user1","email":"user1@example.com"},{"id":2,"username":"user2","email":"user2@example.com"},{"id":3,"username":"user3","email":"user3@example.com"},{"id":4,"username":"testuser703","email":"testuser379@example.com"}]


MAXIM@LAPTOP-123 MINGW64 /
$ curl http://localhost:8000/health
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    30  100    30    0     0   8145      0 --:--:-- --:--:-- --:--:-- 10000Strangler Fig Proxy is healthy

```

**Логи docker для случая MOVIES_MIGRATION_PERCENT: "90"**

```
cinemaabyss-kafka-ui        | 2026-07-16 19:52:02,298 DEBUG [parallel-5] c.p.k.u.s.ClustersStatisticsScheduler: Metrics updated for cluster: cinemaabyss
cinemaabyss-proxy-service   | 2026/07/16 19:52:10 GET /api/movies from 172.21.0.1:46384
cinemaabyss-movies-service  | get movies from movies
cinemaabyss-proxy-service   | 2026/07/16 19:52:11 GET /api/movies from 172.21.0.1:46384
cinemaabyss-movies-service  | get movies from movies
cinemaabyss-proxy-service   | 2026/07/16 19:52:12 GET /api/movies from 172.21.0.1:46384
cinemaabyss-movies-service  | get movies from movies
cinemaabyss-proxy-service   | 2026/07/16 19:52:12 GET /api/movies from 172.21.0.1:46384
cinemaabyss-movies-service  | get movies from movies
cinemaabyss-proxy-service   | 2026/07/16 19:52:13 GET /api/movies from 172.21.0.1:46384
cinemaabyss-movies-service  | get movies from movies
cinemaabyss-proxy-service   | 2026/07/16 19:52:14 GET /api/movies from 172.21.0.1:46384
cinemaabyss-movies-service  | get movies from movies
cinemaabyss-proxy-service   | 2026/07/16 19:52:15 GET /api/movies from 172.21.0.1:46384
cinemaabyss-monolith        | get movies from monolith
cinemaabyss-proxy-service   | 2026/07/16 19:52:15 GET /api/movies from 172.21.0.1:46384
cinemaabyss-movies-service  | get movies from movies
cinemaabyss-proxy-service   | 2026/07/16 19:52:16 GET /api/movies from 172.21.0.1:46384
cinemaabyss-movies-service  | get movies from movies
cinemaabyss-proxy-service   | 2026/07/16 19:52:17 GET /api/movies from 172.21.0.1:46384
cinemaabyss-monolith        | get movies from monolith
cinemaabyss-proxy-service   | 2026/07/16 19:52:18 GET /api/movies from 172.21.0.1:46384
cinemaabyss-movies-service  | get movies from movies
cinemaabyss-proxy-service   | 2026/07/16 19:52:18 GET /api/movies from 172.21.0.1:46384
cinemaabyss-movies-service  | get movies from movies
cinemaabyss-proxy-service   | 2026/07/16 19:52:19 GET /api/movies from 172.21.0.1:46384
cinemaabyss-movies-service  | get movies from movies
cinemaabyss-proxy-service   | 2026/07/16 19:52:20 GET /api/movies from 172.21.0.1:46384
cinemaabyss-movies-service  | get movies from movies
```

**Скрины из Kafka UI**

Топики:

![kafka-topics](images/task-2-kafka-topics.png)

Топик movie-events:

![movie-events](images/task-2-kafka-topic-movie-events.png)

Топик user-events:

![user-events](images/task-2-kafka-topic-user-events.png)

Топик user-events - сообщения:

![user-events-messages](images/task-2-kafka-topic-user-events-messages.png)

Топик user-events - консьюмер:

![user-events-consumers](images/task-2-kafka-topic-user-events-consumer-info.png)

Топик payment-events:

![payment-events](images/task-2-kafka-topic-payment-events.png)

Консьюмеры:

![kafka-consumers](images/task-2-kafka-consumers.png)



## Задание 3

Команда начала переезд в Kubernetes для лучшего масштабирования и повышения надежности. 
Вам, как архитектору осталось самое сложное:
 - реализовать CI/CD для сборки прокси сервиса
 - реализовать необходимые конфигурационные файлы для переключения трафика.


### CI/CD

**Скриншот сборки**

![github-build](images/task-3-github-build.png)

**Скриншот тестов**

![github-tests](images/task-3-github-tests.png)

### Proxy в Kubernetes

#### Шаг 1
Для деплоя в kubernetes необходимо залогиниться в docker registry Github'а.
1. Создайте Personal Access Token (PAT) https://github.com/settings/tokens . Создавайте class с правом read:packages
2. В src/kubernetes/*.yaml (event-service, monolith, movies-service и proxy-service)  отредактируйте путь до ваших образов 
```bash
 spec:
      containers:
      - name: events-service
        image: ghcr.io/ваш логин/имя репозитория/events-service:latest
```
3. Добавьте в секрет src/kubernetes/dockerconfigsecret.yaml в поле
```bash
 .dockerconfigjson: значение в base64 файла ~/.docker/config.json
```

4. Если в ~/.docker/config.json нет значения для аутентификации
```json
{
        "auths": {
                "ghcr.io": {
                       тут пусто
                }
        }
}
```
то выполните 

и добавьте

```json 
 "auth": "имя пользователя:токен в base64"
```

Чтобы получить значение в base64 можно выполнить команду
```bash
 echo -n ваш_логин:ваш_токен | base64
```

После заполнения config.json, также прогоните содержимое через base64

```bash
cat .docker/config.json | base64
```

и полученное значение добавляем в

```bash
 .dockerconfigjson: значение в base64 файла ~/.docker/config.json
```

#### Шаг 2

  Доработайте src/kubernetes/event-service.yaml и src/kubernetes/proxy-service.yaml

  - Необходимо создать Deployment и Service 
  - Доработайте ingress.yaml, чтобы можно было с помощью тестов проверить создание событий
  - Выполните дальшейшие шаги для поднятия кластера:

  1. Создайте namespace:
  ```bash
  kubectl apply -f src/kubernetes/namespace.yaml
  ```
  2. Создайте секреты и переменные
  ```bash
  kubectl apply -f src/kubernetes/configmap.yaml
  kubectl apply -f src/kubernetes/secret.yaml
  kubectl apply -f src/kubernetes/dockerconfigsecret.yaml
  kubectl apply -f src/kubernetes/postgres-init-configmap.yaml
  ```

  3. Разверните базу данных:
  ```bash
  kubectl apply -f src/kubernetes/postgres.yaml
  ```

  На этом этапе если вызвать команду
  ```bash
  kubectl -n cinemaabyss get pod
  ```
  Вы увидите

  NAME         READY   STATUS    
  postgres-0   1/1     Running   

  4. Разверните Kafka:
  ```bash
  kubectl apply -f src/kubernetes/kafka/kafka.yaml
  ```

  Проверьте, теперь должно быть запущено 3 пода, если что-то не так, то посмотрите логи
  ```bash
  kubectl -n cinemaabyss logs имя_пода (например - kafka-0)
  ```

  5. Разверните монолит:
  ```bash
  kubectl apply -f src/kubernetes/monolith.yaml
  ```
  6. Разверните микросервисы:
  ```bash
  kubectl apply -f src/kubernetes/movies-service.yaml
  kubectl apply -f src/kubernetes/events-service.yaml
  ```
  7. Разверните прокси-сервис:
  ```bash
  kubectl apply -f src/kubernetes/proxy-service.yaml
  ```

  После запуска и поднятия подов вывод команды 
  ```bash
  kubectl -n cinemaabyss get pod
  ```

  Будет наподобие такого

  NAME                              READY   STATUS    

  events-service-7587c6dfd5-6whzx   1/1     Running  

  kafka-0                           1/1     Running   

  monolith-8476598495-wmtmw         1/1     Running  

  movies-service-6d5697c584-4qfqs   1/1     Running  

  postgres-0                        1/1     Running  

  proxy-service-577d6c549b-6qfcv    1/1     Running  

  zookeeper-0                       1/1     Running 

  8. Добавим ingress

  - добавьте аддон
  ```bash
  minikube addons enable ingress
  ```
  ```bash
  kubectl apply -f src/kubernetes/ingress.yaml
  ```
  9. Добавьте в /etc/hosts
  127.0.0.1 cinemaabyss.example.com

  10. Вызовите
  ```bash
  minikube tunnel
  ```
  11. Вызовите https://cinemaabyss.example.com/api/movies
  Вы должны увидеть вывод списка фильмов
  Можно поэкспериментировать со значением   MOVIES_MIGRATION_PERCENT в src/kubernetes/configmap.yaml и убедится, что вызовы movies уходят полностью в новый сервис

  12. Запустите тесты из папки tests/postman
  ```bash
   npm run test:kubernetes
  ```
  Часть тестов с health-чек упадет, но создание событий отработает.
  Откройте логи event-service и сделайте скриншот обработки событий

#### Шаг 3
Добавьте сюда скриншота вывода при вызове https://cinemaabyss.example.com/api/movies и  скриншот вывода event-service после вызова тестов.


## Задание 4
Для простоты дальнейшего обновления и развертывания вам как архитектуру необходимо так же реализовать helm-чарты для прокси-сервиса и проверить работу 

Для этого:
1. Перейдите в директорию helm и отредактируйте файл values.yaml

```yaml
# Proxy service configuration
proxyService:
  enabled: true
  image:
    repository: ghcr.io/db-exp/cinemaabysstest/proxy-service
    tag: latest
    pullPolicy: Always
  replicas: 1
  resources:
    limits:
      cpu: 300m
      memory: 256Mi
    requests:
      cpu: 100m
      memory: 128Mi
  service:
    port: 80
    targetPort: 8000
    type: ClusterIP
```

- Вместо ghcr.io/db-exp/cinemaabysstest/proxy-service напишите свой путь до образа для всех сервисов
- для imagePullSecret проставьте свое значение (скопируйте из конфигурации kubernetes)
  ```yaml
  imagePullSecrets:
      dockerconfigjson: ewoJImF1dGhzIjogewoJCSJnaGNyLmlvIjogewoJCQkiYXV0aCI6ICJaR0l0Wlhod09tZG9jRjl2UTJocVZIa3dhMWhKVDIxWmFVZHJOV2hRUW10aFVXbFZSbTVaTjJRMFNYUjRZMWM9IgoJCX0KCX0sCgkiY3JlZHNTdG9yZSI6ICJkZXNrdG9wIiwKCSJjdXJyZW50Q29udGV4dCI6ICJkZXNrdG9wLWxpbnV4IiwKCSJwbHVnaW5zIjogewoJCSIteC1jbGktaGludHMiOiB7CgkJCSJlbmFibGVkIjogInRydWUiCgkJfQoJfSwKCSJmZWF0dXJlcyI6IHsKCQkiaG9va3MiOiAidHJ1ZSIKCX0KfQ==
  ```

2. В папке ./templates/services заполните шаблоны для proxy-service.yaml и events-service.yaml (опирайтесь на свою kubernetes конфигурацию - смысл helm'а сделать шаблоны для быстрого обновления и установки)

```yaml
template:
    metadata:
      labels:
        app: proxy-service
    spec:
      containers:
       Тут ваша конфигурация
```

3. Проверьте установку
Сначала удалим установку руками

```bash
kubectl delete all --all -n cinemaabyss
kubectl delete  namespace cinemaabyss
```
Запустите 
```bash
helm install cinemaabyss .\src\kubernetes\helm --namespace cinemaabyss --create-namespace
```
Если в процессе будет ошибка
```code
[2025-04-08 21:43:38,780] ERROR Fatal error during KafkaServer startup. Prepare to shutdown (kafka.server.KafkaServer)
kafka.common.InconsistentClusterIdException: The Cluster ID OkOjGPrdRimp8nkFohYkCw doesn't match stored clusterId Some(sbkcoiSiQV2h_mQpwy05zQ) in meta.properties. The broker is trying to join the wrong cluster. Configured zookeeper.connect may be wrong.
```

Проверьте развертывание:
```bash
kubectl get pods -n cinemaabyss
minikube tunnel
```

Потом вызовите 
https://cinemaabyss.example.com/api/movies
и приложите скриншот развертывания helm и вывода https://cinemaabyss.example.com/api/movies


# Задание 5
Компания планирует активно развиваться и для повышения надежности, безопасности, реализации сетевых паттернов типа Circuit Breaker и канареечного деплоя вам как архитектору необходимо развернуть istio и настроить circuit breaker для monolith и movies сервисов.

```bash

helm repo add istio https://istio-release.storage.googleapis.com/charts
helm repo update

helm install istio-base istio/base -n istio-system --set defaultRevision=default --create-namespace
helm install istio-ingressgateway istio/gateway -n istio-system
helm install istiod istio/istiod -n istio-system --wait

helm install cinemaabyss .\src\kubernetes\helm --namespace cinemaabyss --create-namespace

kubectl label namespace cinemaabyss istio-injection=enabled --overwrite

kubectl get namespace -L istio-injection

kubectl apply -f .\src\kubernetes\circuit-breaker-config.yaml -n cinemaabyss

```

Тестирование

# fortio
```bash
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.25/samples/httpbin/sample-client/fortio-deploy.yaml -n cinemaabyss
```

# Get the fortio pod name
```bash
FORTIO_POD=$(kubectl get pod -n cinemaabyss | grep fortio | awk '{print $1}')

kubectl exec -n cinemaabyss $FORTIO_POD -c fortio -- fortio load -c 50 -qps 0 -n 500 -loglevel Warning http://movies-service:8081/api/movies
```
Например,

```bash
kubectl exec -n cinemaabyss fortio-deploy-b6757cbbb-7c9qg  -c fortio -- fortio load -c 50 -qps 0 -n 500 -loglevel Warning http://movies-service:8081/api/movies
```

Вывод будет типа такого

```bash
IP addresses distribution:
10.106.113.46:8081: 421
Code 200 : 79 (15.8 %)
Code 500 : 22 (4.4 %)
Code 503 : 399 (79.8 %)
```
Можно еще проверить статистику

```bash
kubectl exec -n cinemaabyss fortio-deploy-b6757cbbb-7c9qg -c istio-proxy -- pilot-agent request GET stats | grep movies-service | grep pending
```

И там смотрим 

```bash
cluster.outbound|8081||movies-service.cinemaabyss.svc.cluster.local;.upstream_rq_pending_total: 311 - столько раз срабатывал circuit breaker
You can see 21 for the upstream_rq_pending_overflow value which means 21 calls so far have been flagged for circuit breaking.
```

Приложите скриншот работы circuit breaker'а

Удаляем все
```bash
istioctl uninstall --purge
kubectl delete namespace istio-system
kubectl delete all --all -n cinemaabyss
kubectl delete namespace cinemaabyss
```
