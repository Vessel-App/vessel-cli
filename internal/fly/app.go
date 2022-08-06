package fly

type CreateAppRequest struct {
	AppName string
	OrgName string // todo: Not needed?
}

type GetAppRequest struct {
	AppName string
}

type DeleteAppRequest struct {
	AppName string
}
