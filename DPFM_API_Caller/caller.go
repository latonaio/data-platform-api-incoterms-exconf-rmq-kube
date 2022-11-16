package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-incoterms-exconf-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-incoterms-exconf-rmq-kube/DPFM_API_Output_Formatter"
	"data-platform-api-incoterms-exconf-rmq-kube/database"
	"sync"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
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

func (e *ExistenceConf) Conf(input *dpfm_api_input_reader.SDC) *dpfm_api_output_formatter.Incoterms {
	incoterms := *input.Incoterms.Incoterms
	notKeyExistence := make([]string, 0, 1)
	KeyExistence := make([]string, 0, 1)

	existData := &dpfm_api_output_formatter.Incoterms{
		Incoterms:     incoterms,
		ExistenceConf: false,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if !e.confIncoterms(incoterms) {
			notKeyExistence = append(notKeyExistence, incoterms)
			return
		}
		KeyExistence = append(KeyExistence, incoterms)
	}()

	wg.Wait()

	if len(KeyExistence) == 0 {
		return existData
	}
	if len(notKeyExistence) > 0 {
		return existData
	}

	existData.ExistenceConf = true
	return existData
}

func (e *ExistenceConf) confIncoterms(val string) bool {
	rows, err := e.db.Query(
		`SELECT Incoterms 
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_incoterms_incoterms_data 
		WHERE Incoterms = ?;`, val,
	)
	if err != nil {
		e.l.Error(err)
		return false
	}

	for rows.Next() {
		var incoterms string
		err := rows.Scan(&incoterms)
		if err != nil {
			e.l.Error(err)
			continue
		}
		if incoterms == val {
			return true
		}
	}
	return false
}
