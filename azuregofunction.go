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
)

func ParseFunctionHostRequest(w http.ResponseWriter, r *http.Request) (*InvokeRequest, error) {
	fmt.Println("--------------------")
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
	for k, v := range invokeReq.Data {
		fmt.Printf("%v=%v\n", k, v)
	}
	fmt.Println("----------")
	fmt.Println("The JSON metadata is:")
	for k, v := range invokeReq.Metadata {
		fmt.Printf("%v=%v\n", k, v)
	}
	fmt.Println("----------")

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
	fmt.Println(reflect.TypeOf(inMsg))
	return fmt.Sprintf("%v", inMsg)
}
