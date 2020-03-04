package untraceable

// LLeidanet struct
type LLeidanet struct {
	Sms Parameters `json:"sms"`
}

// Parameters represents params structure to send data to Lleidanet platform
type Parameters struct {
	User        string      `json:"user"`
	Password    string      `json:"password"`
	Destination Destination `json:"dst"`
	Text        string      `json:"txt"`
	Source      string      `json:"src"`
	Schedule    string      `json:"schedule,omitempty"`
}

// Destination is an array of phone numbers
type Destination struct {
	Number []string `json:"num"`
}

// {
//   "sms": {
//     "user":"{:name}",
//     "password":"{:pass}",
//     "dst":{
//       "num":["+34600000000", "+34666666666"]
//     },
//     "txt":"Message text",
//		 "src":"Sender"
//   }
// }

// LLeidaResp represents Lleidanet API XML response
type LLeidaResp struct {
	Request   string `xml:"request"`
	Code      int    `xml:"code"`
	Status    string `xml:"status"`
	Newcredit string `xml:"newcredit"`
}
