package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-incoterms-exconf-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-incoterms-exconf-rmq-kube/DPFM_API_Output_Formatter"
	"encoding/json"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	database "github.com/latonaio/golang-mysql-network-connector"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
	"golang.org/x/xerrors"
)

type ExistenceConf struct {
	ctx context.Context
	db  *database.Mysql
	l   *logger.Logger
}

func NewExistenceConf(ctx context.Context, db *database.Mysql, l *logger.Logger) *ExistenceConf {
	return &ExistenceConf{
		ctx: ctx,
		db:  db,
		l:   l,
	}
}

func (e *ExistenceConf) Conf(msg rabbitmq.RabbitmqMessage) interface{} {
	var ret interface{}
	ret = map[string]interface{}{
		"ExistenceConf": false,
	}
	input := make(map[string]interface{})
	err := json.Unmarshal(msg.Raw(), &input)
	if err != nil {
		return ret
	}

	_, ok := input["Incoterms"]
	if ok {
		input := &dpfm_api_input_reader.SDC{}
		err = json.Unmarshal(msg.Raw(), input)
		ret = e.confIncoterms(input)
		goto endProcess
	}

	err = xerrors.Errorf("can not get exconf check target")
endProcess:
	if err != nil {
		e.l.Error(err)
	}
	return ret
}

func (e *ExistenceConf) confIncoterms(input *dpfm_api_input_reader.SDC) *dpfm_api_output_formatter.Incoterms {
	exconf := dpfm_api_output_formatter.Incoterms{
		ExistenceConf: false,
	}
	if input.Incoterms.Incoterms == nil {
		return &exconf
	}
	exconf = dpfm_api_output_formatter.Incoterms{
		Incoterms:     *input.Incoterms.Incoterms,
		ExistenceConf: false,
	}

	rows, err := e.db.Query(
		`SELECT Incoterms 
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_incoterms_incoterms_data 
		WHERE Incoterms = ?;`, exconf.Incoterms,
	)
	if err != nil {
		e.l.Error(err)
		return &exconf
	}
	defer rows.Close()

	exconf.ExistenceConf = rows.Next()
	return &exconf
}
