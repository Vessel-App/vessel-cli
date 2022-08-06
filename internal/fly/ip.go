package fly

type AllocateIpRequest struct {
	App  string
	Type string
}

type GetAppIpRequest struct {
	App string
}
