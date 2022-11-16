package dpfm_api_input_reader

import (
	"data-platform-api-incoterms-exconf-rmq-kube/DPFM_API_Caller/requests"
)

func (sdc *SDC) ConvertToIncoterms() *requests.Incoterms {
	data := sdc.Incoterms
	return &requests.Incoterms{
		Incoterms: data.Incoterms,
	}
}
