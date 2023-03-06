package azuregofunction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type DataHttpRequest struct {
	Request http.Request
	//Identities map[string]interface{}
	//Params     map[string]interface{}
	//Url     string
	//Method  string
	//Query   map[string]string
	//Headers map[string][]string
}

func (r DataHttpRequest) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("URL:%v\n" + r.Request.URL.String()))
	buffer.WriteString(fmt.Sprintf("Method:%v\n" + r.Request.Method))
	buffer.WriteString(fmt.Sprintf("Headers:\n"))
	for k, v := range r.Request.Header {
		buffer.WriteString(fmt.Sprintf("\t%v=%v\n", k, v))
	}
	buffer.WriteString(fmt.Sprintf("Query:\n"))
	for k, v := range r.Request.URL.Query() {
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
	fmt.Println("+--------------------+")
	fmt.Println("Generating data http request")
	dataHttpRequest := DataHttpRequest{Request: http.Request{}}

	var queryValues string

	var err error
	v := req.(map[string]interface{})
	for k, v := range v {
		fmt.Printf("%v=%v\n", k, v)
		if k == "Url" {
			dataHttpRequest.Request.URL, err = url.Parse(v.(string))
			if err != nil {
				fmt.Println(err)
			}
		} else if k == "Method" {
			dataHttpRequest.Request.Method = v.(string)
		} else if k == "Query" {
			var sb strings.Builder
			m := v.(map[string]interface{})
			//pm := make(map[string]string)
			for mk, mv := range m {
				//pm[mk] = mv.(string)
				sb.WriteString(fmt.Sprintf("%v=%v", mk, mv))
			}
			queryValues = sb.String()
		}
	}

	if dataHttpRequest.Request.URL != nil {
		fmt.Println("Set raw query :" + queryValues)
		dataHttpRequest.Request.URL.RawQuery = queryValues
	}
	fmt.Println("+--------------------+")
	return &dataHttpRequest
}
