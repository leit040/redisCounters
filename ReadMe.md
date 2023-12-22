import 
rc "github.com/leit040/redisCounters"

countersMap := rc.NewCountersMap() 
возвращает структуру данных для хранения обьектов 

type CountersGroup struct {
prefix  string
keys    []string
connect *RedisPull.Connect
}

prefix - идентификатор для данной группы счетчиков, может быть домен, id feed, и так далее.
keys    []string - это имена непосредственно счетчиков, в базе они будут в формате prefix_key.

все счетчики одной группы будут использовать одно соединение с редисом 
описанное в connect *RedisPull.Connect


для добавления непосредственно группы счетчиков вызываем метод
countersMap.AddCountersGroup(prefix string, keys []string, connect *RedisPull.Connect)
где указываем префикс, массив ключей ну и передаем ссылку на connect *RedisPull.Connect
который будет использовантся для этой группы

метод создает счетчики по ключам и добавляет эту группу в countersMap

после этого можно получить группу счетчиков вызвав
countersGroup := countersMap.GetCountersGroup(prefix string)

доступны методы:
countersGroup.IncreaseCounter(key string) - увеличииь значение счетчика на 1

countersGroup.GetCounterValue(key string, interval string) 
тут опять таки, key - это ключ счетчика, (без префикса, он сам его подставит) а interval - это "day" or "hour"
что вернет соответсвенно int значение счетчика за текущий день или текущий час (начиная с 00:00 или с 14-00 допустим если сейчас 14-45)







