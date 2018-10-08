package common

import (
	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
)

func MachineryInstance() (*machinery.Server, error) {
	var cnf = config.Config{
		Broker:        "amqp://guest:guest@localhost:5672",
		DefaultQueue:  "machinery_tasks",
		ResultBackend: "amqp://guest:guest@localhost:5672",
		AMQP: &config.AMQPConfig{
			Exchange:     "machinery_exchange",
			ExchangeType: "direct",
			BindingKey:   "machinery_task",
		},
	}
	server, error := machinery.NewServer(&cnf)
	if error != nil {
		panic("Could not create server")
	}

	return server, nil
}
