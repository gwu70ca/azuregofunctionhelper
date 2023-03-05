package azuregofunction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
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

type ReturnValue struct {
	Data string
}
type InvokeResponse struct {
	Outputs     map[string]interface{}
	Logs        []string
	ReturnValue interface{}
}

type InvokeResponseStringReturnValue struct {
	Outputs     map[string]interface{} //shows as Http response
	Logs        []string               //shows in log
	ReturnValue string                 //saved to output binding
}

type InvokeRequest struct {
	Data     map[string]interface{}
	Metadata map[string]interface{}
}

const (
	BlobNameKey = "name"
	BlobUriKey  = "Uri"
	HttpReqKey  = "req"
)

func ParseFunctionHostRequest(w http.ResponseWriter, r *http.Request) (*InvokeRequest, error) {
	fmt.Println("+--------------------+")
	fmt.Println("Parsing request from function host")

	var invokeReq InvokeRequest
	d := json.NewDecoder(r.Body)
	decodeErr := d.Decode(&invokeReq)
	if decodeErr != nil {
		// bad JSON or unrecognized json field
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return nil, decodeErr
	}
	fmt.Println("The JSON data is:")
	fmt.Println("----------")
	fmt.Println(fmt.Sprintf("Type of invokeReq.Data: %v", reflect.TypeOf(invokeReq.Data)))
	for k, v := range invokeReq.Data {
		fmt.Printf("%v=%v\n", k, v)
	}

	fmt.Println("The JSON metadata is:")
	fmt.Println("----------")
	fmt.Println(fmt.Sprintf("Type of invokeReq.Metadata: %v", reflect.TypeOf(invokeReq.Metadata)))
	for k, v := range invokeReq.Metadata {
		fmt.Printf("%v=%v\n", k, v)
	}
	fmt.Println("+--------------------+")

	return &invokeReq, nil
}

// Return blob data
func BlobData(ir *InvokeRequest, bindingName string) interface{} {
	return ir.Data[bindingName]
}

func BlobName(ir *InvokeRequest) string {
	return fmt.Sprintf("%v", ir.Metadata[BlobNameKey])
}

func BlobUri(ir *InvokeRequest) string {
	return fmt.Sprintf("%v", ir.Metadata[BlobUriKey])
}

func QueueMessage(ir *InvokeRequest, bindingName string) string {
	return fmt.Sprintf("%v", ir.Data[bindingName])
}

func EventHubMessage(ir *InvokeRequest, bindingName string) string {
	inMsg := ir.Data[bindingName]
	fmt.Println(fmt.Sprintf("Type of event hub message: %v", reflect.TypeOf(inMsg)))
	return fmt.Sprintf("%v", inMsg)
}

func HttpRequestDataWithBinding(ir *InvokeRequest, bindingName string) *DataHttpRequest {
	return parseDataHttpRequest(ir.Data[bindingName])
}

func HttpRequestData(ir *InvokeRequest) *DataHttpRequest {
	return parseDataHttpRequest(ir.Data[HttpReqKey])
}

/*
func HttpRequestMetaData(ir *InvokeRequest, bindingName string) *DataHttpRequest {
	name := HttpReq
	if bindingName != "" {
		name = bindingName
	}

	return parseDataHttpRequest(ir.Data[name])
}
*/

func parseDataHttpRequest(req interface{}) *DataHttpRequest {
	httpRequest := DataHttpRequest{}

	v := req.(map[string]interface{})
	for k, v := range v {
		fmt.Printf("%v=%v\n", k, v)
		if k == "Url" {
			httpRequest.Url = v.(string)
		} else if k == "Method" {
			httpRequest.Method = v.(string)
		} else if k == "Query" {
			m := v.(map[string]interface{})
			pm := make(map[string]string)
			for mk, mv := range m {
				pm[mk] = mv.(string)
			}
			httpRequest.Query = pm
		}
	}

	return &httpRequest
}
