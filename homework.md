# ДОМАШКА

Все очень просто, нужно просто настроить group replication для поднимаемых нод =)


## Небольшой чек-лист перед тем как начать выполнять именно дз
- Сколнить репу с проектом(эту)
- Выполнить `make uninstall` и `make install`
- Выполнить `make deploy` и убедится что создался `namespace` и в нем запустился наш под с оператором
- Меняем контекст чтобы `kubectl` "смотрел" в наш namespace.
 `kubectl config set-context --current --namespace=NAMESPACE_NAME`
- Делаем apply нашего custom resource. `kubectl apply -f config/samples/workshop_v1alpha1_custommysql.yaml `
- После этого, должно поднятся 2 пода с mysql (может занять время, так что не торопимся тут)
## ЕСЛИ ЧТО ТО ПОШЛО НЕ ТАК, ПИШЕМ В ЧАТ)

## После, смотрим файл `controllers/mysql.go`