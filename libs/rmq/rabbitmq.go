package rmq

import (
	"fmt"
	"sync"

	"github.com/assembla/cony"
	"github.com/streadway/amqp"
)

//Receiver 定义接收者接口
type Receiver interface {
	Consumer(amqp.Delivery) error
}

//RMQ mq对象
type RMQ struct {
	cli        *cony.Client
	connection *amqp.Connection
	channel    *amqp.Channel
	mu         sync.RWMutex
	addr       string // 地址
}

//QueueExchange 定义队列交换机对象
type QueueExchange struct {
	QueueName    string // 队列名称
	RoutingKey   string // key值
	ExchangeName string // 交换机名称
	ExchangeType string // 交换机类型
}

//NewRMQ 创建rmq
func NewRMQ(addr string) *RMQ {
	// Construct new client with the flag url
	// and default backoff policy
	cli := cony.NewClient(
		cony.URL(addr),
		cony.Backoff(cony.DefaultBackoff),
	)

	rmq := &RMQ{
		cli:  cli,
		addr: addr,
	}
	// rmq.connect()
	return rmq
}

//BindQueue 绑定队列
func (r *RMQ) BindQueue(exchangeName, exchangeType, routingKey, queueName string) {
	// Declarations
	// The queue name will be supplied by the AMQP server
	que := &cony.Queue{
		AutoDelete: false,
		Name:       queueName,
	}
	exc := cony.Exchange{
		Name:       exchangeName,
		Kind:       exchangeType,
		AutoDelete: false,
	}
	bnd := cony.Binding{
		Queue:    que,
		Exchange: exc,
		Key:      routingKey,
	}

	r.cli.Declare([]cony.Declaration{
		cony.DeclareQueue(que),
		cony.DeclareExchange(exc),
		cony.DeclareBinding(bnd),
	})
}

//InitPaper 初始化
func (r *RMQ) InitPaper(qe QueueExchange) *Paper {
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

	r.cli.Declare([]cony.Declaration{
		cony.DeclareQueue(que),
		cony.DeclareExchange(exc),
		cony.DeclareBinding(bnd),
	})

	r.cli.Loop()

	paper := &Paper{
		mq:           r,
		queueName:    qe.QueueName,
		routingKey:   qe.RoutingKey,
		exchangeName: qe.ExchangeName,
		exchangeType: qe.ExchangeType,
		que:          que,
		exc:          exc,
		bnd:          bnd,
	}

	// 注册生产者
	pbl := cony.NewPublisher(exc.Name, qe.RoutingKey)
	r.cli.Publish(pbl)
	paper.pbl = pbl

	return paper
}

//Paper 页 类似Table
type Paper struct {
	mq           *RMQ
	queueName    string        // 队列名称
	routingKey   string        // key名称
	exchangeName string        // 交换机名称
	exchangeType string        // 交换机类型
	que          *cony.Queue   // 队列名称
	exc          cony.Exchange // key名称
	bnd          cony.Binding  // 交换机名称
	pbl          *cony.Publisher
}

//Pub 发布消息
func (p *Paper) Pub(data []byte) {
	var err error

	// 发送任务消息
	err = p.pbl.Publish(amqp.Publishing{
		ContentType: "text/plain",
		Body:        data,
	})
	if err != nil {
		fmt.Printf("MQ任务发送失败:%s \n", err)
		return
	}
}

//Receive 消费消息
func (p *Paper) Receive(c Receiver) {
	// Declare and register a consumer
	cns := cony.NewConsumer(
		p.que,
		cony.Qos(1),
	)
	p.mq.cli.Consume(cns)

	for {
		msg := <-cns.Deliveries()
		c.Consumer(msg)
	}
}
