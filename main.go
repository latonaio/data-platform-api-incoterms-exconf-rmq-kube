package main

import (
	"context"
	dpfm_api_caller "data-platform-api-incoterms-exconf-rmq-kube/DPFM_API_Caller"
	dpfm_api_input_reader "data-platform-api-incoterms-exconf-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-incoterms-exconf-rmq-kube/DPFM_API_Output_Formatter"
	"data-platform-api-incoterms-exconf-rmq-kube/config"
	"data-platform-api-incoterms-exconf-rmq-kube/database"
	"encoding/json"
	"fmt"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
)

func main() {
	ctx := context.Background()
	l := logger.NewLogger()
	c := config.NewConf()
	db, err := database.NewMySQL(c.DB)
	if err != nil {
		l.Error(err)
		return
	}

	rmq, err := rabbitmq.NewRabbitmqClient(c.RMQ.URL(), c.RMQ.QueueFrom(), "", nil, -1)
	if err != nil {
		l.Fatal(err.Error())
	}
	iter, err := rmq.Iterator()
	if err != nil {
		l.Fatal(err.Error())
	}
	defer rmq.Stop()
	for msg := range iter {
		go dataCallProcess(ctx, c, db, msg)
	}
}

func dataCallProcess(
	ctx context.Context,
	c *config.Conf,
	db *database.Mysql,
	rmqMsg rabbitmq.RabbitmqMessage,
) {
	defer rmqMsg.Success()
	l := logger.NewLogger()
	sessionId := getBodyHeader(rmqMsg.Data())
	l.AddHeaderInfo(map[string]interface{}{"runtime_session_id": sessionId})
	input := &dpfm_api_input_reader.SDC{}
	err := json.Unmarshal(rmqMsg.Raw(), input)
	if err != nil {
		l.Error(rmqMsg.Data())
		return
	}

	conf := dpfm_api_caller.NewExistenceConf(ctx, db, l)
	exist := conf.Conf(input)
	rmqMsg.Respond(exist)

	out := dpfm_api_output_formatter.MetaData{}
	err = json.Unmarshal(rmqMsg.Raw(), &out)
	if err != nil {
		l.Error(rmqMsg.Data())
		return
	}
	out.Incoterms = *exist
	l.Info(out)
}

func getBodyHeader(data map[string]interface{}) string {
	id := fmt.Sprintf("%v", data["runtime_session_id"])
	return id
}
