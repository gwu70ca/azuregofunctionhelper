package azuregofunction

type DataHttpRequest struct {
	//Identities map[string]interface{}
	//Params     map[string]interface{}
	Url     string
	Method  string
	Query   map[string]string
	Headers map[string]interface{}
}
