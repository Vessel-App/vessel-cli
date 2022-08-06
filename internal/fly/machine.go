package fly

type RunMachineRequest struct {
	App    string
	Region string
	Image  string
	Env    map[string]string
}
type ListMachinesRequest struct {
	App string
}

type GetMachineRequest struct {
	App     string
	Machine string
}

type StartMachineRequest struct {
	App     string
	Machine string
}

type StopMachineRequest struct {
	App     string
	Machine string
}

type DeleteMachineRequest struct {
	App     string
	Machine string
}
