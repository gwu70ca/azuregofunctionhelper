package azuregofunction

type NormalHttpRequest struct {
	//Headers    map[string]interface{}
	//Identities map[string]interface{}
	//Params     map[string]interface{}
	Method string
	Query  map[string]string
	Url    string
}
