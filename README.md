## Описание
Очень легкий пакет по работе с Zabbix протоколом, который появился в процессе исследования сетевого взаимодействия с Zabbix Агентом.<br>
Работает с не сжатыми пакетами, имеющими флаг `0x01`<br>
Подробнее о структуре пакета можно почитать в [официальной документации](https://www.zabbix.com/documentation/current/en/manual/appendix/protocols/header_datalen).

## Важно
Zabbix Агент разрывает TCP соединение сразу после окончания передачи ответа.

## Как пользоваться
Установка
```bash
go get -u github.com/rzolotuhin/go-zabbix-agent-protocol
```

Объявляем транспорт


TCP
```go
conn, err := net.DialTimeout("tcp", "127.0.0.1:10050", time.Second * 3)
```

TLS
```go
conn, err := tls.Dial("tcp", "127.0.0.1:10050", &tls.Config{
    // Cert or PSK config
})
```

Указываем транспорт в качестве параметра
```go
agent := zabbix.Agent{
    Transport: conn,
}
```

Запрашиваем данные по ключу
```go
key := "net.if.discovery"
answer, err := agent.Get(key)
```

## Транспорт
Может быть использован любой транспорт, соответствующий следующему интерфейсу
```go
type NetIO interface {
	Read(b []byte) (n int, err error)
	Write(p []byte) (n int, err error)
}

type Agent struct {
	Transport NetIO
}
```
