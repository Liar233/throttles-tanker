package protocol

const AddWhiteRuleUrl = "/add/white"
const AddBlackRuleUrl = "/add/black"
const DeleteWhiteRuleUrl = "/delete/white"
const DeleteBlackRuleUrl = "/delete/white"
const CheckUrl = "/check"
const FlushUrl = "/flush"

type AddRuleRequest struct {
	Subnet string `json:"subnet"`
}

type CheckRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
}

type DeleteRuleRequest struct {
	Subnet string `json:"subnet"`
}

type FlushRequest struct {
	Login string `json:"login"`
	Ip    string `json:"ip"`
}
