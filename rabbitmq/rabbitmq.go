package rabbitmq

import (
	"fmt"

	"github.com/haobird/golibs/cony"
	"github.com/streadway/amqp"
)

// 封装基于cony 的易用包

//MsgHandle 注册函数
type MsgHandle func(amqp.Delivery)

//Consumer 消费者
type Consumer struct {
	cli *cony.Client
	qe  QueueExchange
	cns *cony.Consumer
}

//Publisher 生产者
type Publisher struct {
	cli *cony.Client
	qe  QueueExchange
	pbl *cony.Publisher
}

//Declarer 定义
type Declarer struct {
	que *cony.Queue   // 队列名称
	exc cony.Exchange // 交换机名称
	bnd cony.Binding  //
}

//QueueExchange 定义队列交换机对象
type QueueExchange struct {
	QueueName    string // 队列名称
	RoutingKey   string // key值
	ExchangeName string // 交换机名称
	ExchangeType string // 交换机类型
}

//RabbitMQ mq对象
type RabbitMQ struct {
	cli  *cony.Client
	addr string // 地址
}

//New 创建rmq
func New(addr string) *RabbitMQ {
	// Construct new client with the flag url
	// and default backoff policy
	cli := cony.NewClient(
		cony.URL(addr),
		cony.Backoff(cony.DefaultBackoff),
	)

	rmq := &RabbitMQ{
		cli:  cli,
		addr: addr,
	}
	return rmq
}

//BindQueue 绑定队列（创建队列）
func (r *RabbitMQ) BindQueue(exchangeName, exchangeType, routingKey, queueName string) {
	qe := QueueExchange{
		QueueName:    queueName,
		RoutingKey:   routingKey,
		ExchangeName: exchangeName,
		ExchangeType: exchangeType,
	}
	de := r.Declare(qe)
	if qe.QueueName == "" {
		r.cli.Declare([]cony.Declaration{
			cony.DeclareExchange(de.exc),
		})
	} else {
		r.cli.Declare([]cony.Declaration{
			cony.DeclareQueue(de.que),
			cony.DeclareExchange(de.exc),
			cony.DeclareBinding(de.bnd),
		})
	}
}

//Declare 定义队列或交换机
func (r *RabbitMQ) Declare(qe QueueExchange) Declarer {
	// Declarations
	// The queue name will be supplied by the AMQP server
	que := &cony.Queue{
		AutoDelete: false,
		Durable:    true,
		Name:       qe.QueueName,
	}
	exc := cony.Exchange{
		Name:       qe.ExchangeName,
		Kind:       qe.ExchangeType,
		Durable:    true,
		AutoDelete: false,
	}
	bnd := cony.Binding{
		Queue:    que,
		Exchange: exc,
		Key:      qe.RoutingKey,
	}

	return Declarer{
		que: que,
		exc: exc,
		bnd: bnd,
	}

}

//Bind 绑定
func (r *RabbitMQ) Bind(de Declarer) {
	r.cli.Declare([]cony.Declaration{
		cony.DeclareQueue(de.que),
		cony.DeclareExchange(de.exc),
		cony.DeclareBinding(de.bnd),
	})
}

//NewConsumer 创建
func (r *RabbitMQ) NewConsumer(qe QueueExchange) *Consumer {
	// Construct new client with the flag url
	// and default backoff policy
	cli := cony.NewClient(
		cony.URL(r.addr),
		cony.Backoff(cony.DefaultBackoff),
	)

	// Declarations
	de := r.Declare(qe)

	cli.Declare([]cony.Declaration{
		cony.DeclareQueue(de.que),
		cony.DeclareExchange(de.exc),
		cony.DeclareBinding(de.bnd),
	})
	// Declare and register a consumer
	cns := cony.NewConsumer(
		de.que,
		cony.Qos(1),
	)

	cli.Consume(cns)

	go func() {
		for cli.Loop() {
			select {
			case err := <-cns.Errors():
				fmt.Printf("Consumer error: %v\n", err)
			case err := <-cli.Errors():
				fmt.Printf("Consumer Client error: %v\n", err)
			}
		}
	}()

	return &Consumer{
		cli: cli,
		qe:  qe,
		cns: cns,
	}
}

//Receive 消费消息
func (c *Consumer) Receive(h MsgHandle) {
	for {
		msg := <-c.cns.Deliveries()
		h(msg)
	}
}

//NewPublisher 创建生产者
func (r *RabbitMQ) NewPublisher(qe QueueExchange) *Publisher {
	// Construct new client with the flag url
	// and default backoff policy
	cli := cony.NewClient(
		cony.URL(r.addr),
		cony.Backoff(cony.DefaultBackoff),
	)

	// Declarations
	de := r.Declare(qe)

	cli.Declare([]cony.Declaration{
		cony.DeclareExchange(de.exc),
	})

	pbl := cony.NewPublisher(de.exc.Name, qe.RoutingKey)
	cli.Publish(pbl)

	// Client loop sends out declarations(exchanges, queues, bindings
	// etc) to the AMQP server. It also handles reconnecting.
	go func() {
		for cli.Loop() {
			select {
			case err := <-cli.Errors():
				fmt.Printf("Publisher Client error: %v\n", err)
			case blocked := <-cli.Blocking():
				fmt.Printf("Publisher Client is blocked %v\n", blocked)
			}
		}
	}()

	return &Publisher{
		cli: cli,
		qe:  qe,
		pbl: pbl,
	}

}

//Pub 发布消息
func (p *Publisher) Pub(data []byte) {
	var err error

	// 发送任务消息
	err = p.pbl.Publish(amqp.Publishing{
		ContentType: "text/plain",
		Body:        data,
	})
	if err != nil {
		fmt.Printf("Client publish error: %v\n", err)
		return
	}
}
