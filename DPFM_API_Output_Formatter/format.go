package dpfm_api_output_formatter

import (
	"encoding/json"

	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
	"golang.org/x/xerrors"
)

func NewOutput(rmqMsg rabbitmq.RabbitmqMessage, exconf interface{}) (*MetaData, error) {
	output := &MetaData{}
	err := json.Unmarshal(rmqMsg.Raw(), output)
	if err != nil {
		return nil, xerrors.Errorf("output Marshal error: %w", err)
	}

	switch exconf := exconf.(type) {
	case *Incoterms:
		output.Incoterms = exconf
	default:
		return nil, xerrors.Errorf("unknown type %+v", exconf)
	}

	return output, nil
}
