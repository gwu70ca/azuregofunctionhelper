package azuregofunction

import (
	"bytes"
	"fmt"
)

type DataHttpRequest struct {
	//Identities map[string]interface{}
	//Params     map[string]interface{}
	Url     string
	Method  string
	Query   map[string]string
	Headers map[string]interface{}
}

func (r DataHttpRequest) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("URL:%v\n" + r.Url))
	buffer.WriteString(fmt.Sprintf("Method:%v\n" + r.Method))
	buffer.WriteString(fmt.Sprintf("Headers:\n"))
	for k, v := range r.Headers {
		buffer.WriteString(fmt.Sprintf("\t%v=%v\n", k, v))
	}
	buffer.WriteString(fmt.Sprintf("Query:\n"))
	for k, v := range r.Query {
		buffer.WriteString(fmt.Sprintf("\t%v=%v\n", k, v))
	}
	return buffer.String()
}
