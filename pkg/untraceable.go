package untraceable

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

//GetTraced retrieves a list of sms sended in last month
func (h *Handler) GetTraced() error {
	db := h.Storer
	sel := []Untraceable{}
	date := time.Now().Add(time.Duration(-720) * time.Hour) // -30 days

	err := db.Instance().Debug().Where("date(sms_date) >= ?", date.Format("2006-01-02")).Find(&sel).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	h.Leads = sel
	return nil
}

//GetUntraceables retrives a list of untraceable leads
func (h *Handler) GetUntraceables() error {
	leontel := []Leontel{}

	duration := -24
	if time.Now().Weekday() == time.Monday {
		duration = -72
	}
	date := time.Now().Add(time.Duration(duration) * time.Hour)

	traced := []int64{1}
	for _, l := range h.Leads {
		traced = append(traced, l.LeaID)
	}

	err := h.Storer.Instance().Debug().
		Table("crmti.lea_leads").
		Select("lea_id, TELEFONO, lea_source").
		Joins("INNER JOIN crmti.act_activity ON crmti.act_activity.act_id = crmti.lea_leads.lea_id").
		Joins("INNER JOIN crmti.sub_subcategories ON crmti.sub_subcategories.sub_id = crmti.act_activity.act_last_cat").
		Where("crmti.lea_leads.lea_source IN (?) ", []int{73, 74, 75}).
		Where("crmti.sub_subcategories.sub_id in (?)", []int{575, 562}).
		Where("crmti.lea_leads.TELEFONO like ? or crmti.lea_leads.TELEFONO like ?", "6%", "7%").
		Where("crmti.lea_leads.lea_id not in (?)", traced).
		Where("date(crmti.act_activity.act_ts) >= ?", date.Format("2006-01-02")).
		Order("lea_id desc").
		Find(&leontel).
		Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}

	candR := Candidates{
		Desc:  "R",
		DDI:   "881550607",
		Leads: []Untraceable{},
	}

	candK := Candidates{
		Desc:  "Euskatel",
		DDI:   "945551061",
		Leads: []Untraceable{},
	}

	candT := Candidates{
		Desc:  "Telecable",
		DDI:   "984851473",
		Leads: []Untraceable{},
	}

	for _, l := range leontel {
		un := l.MapToUntraceable()

		switch un.SouID {
		case 73:
			un.DDI = candR.DDI
			candR.Leads = append(candR.Leads, un)
		case 75:
			un.DDI = candK.DDI
			candK.Leads = append(candK.Leads, un)
		case 74:
			un.DDI = candT.DDI
			candT.Leads = append(candT.Leads, un)
		default:
			un.DDI = ""
		}
	}

	log.Printf("candidate %s: %d", candR.Desc, len(candR.Leads))
	h.Candidates = append(h.Candidates, candR)

	log.Printf("candidate %s: %d", candK.Desc, len(candK.Leads))
	h.Candidates = append(h.Candidates, candK)

	log.Printf("candidate %s: %d", candT.Desc, len(candT.Leads))
	h.Candidates = append(h.Candidates, candT)

	return nil
}

// MapToUntraceable blablba
func (l *Leontel) MapToUntraceable() Untraceable {
	un := Untraceable{}
	un.LeaID = l.LeaID
	un.Phone = l.Phone
	un.SouID = l.SouID
	return un
}

// MapToLLeida maps phones in leads array to LLeidanet structure
func MapToLLeida(candidates Candidates) Destination {
	log.Printf("MapToLLeida => %s", candidates.Desc)

	dest := Destination{}
	numbers := []string{}
	for _, l := range candidates.Leads {
		numbers = append(numbers, *l.Phone)
	}
	log.Println(numbers)
	// dest.Number = numbers
	// TODO DEV => REMOVE!
	// dest.Number = []string{"665932355", "685243280", "606677113"}
	dest.Number = []string{"665932355"}

	log.Printf("numbers: %s", dest.Number)
	return dest
}

// Fire starts sms sending process for each campaign
func (h *Handler) Fire() {
	for _, cand := range h.Candidates {
		if len(cand.Leads) > 0 {
			lleida := h.LLeidanet
			lleida.Sms.Destination = MapToLLeida(cand)
			lleida.Sms.Source = cand.Desc
			lleida.Sms.Text = lleida.Sms.Text + " " + cand.DDI

			// resp, err := lleida.Send()
			// if err != nil {
			// 	h.pushError(err)
			// }
			// TODO DEV => REMOVE!
			resp := LLeidaResp{
				Status: "Success",
			}

			if resp.Status == "Success" {
				// store candidates
				log.Printf("SMS send [fake] %s success! => storing data", cand.Desc)

				for _, lead := range cand.Leads {
					lead.SmsDate = time.Now()
					h.Storer.Insert(&lead)
				}
			}
		}
	}
}

// Send launch POST request to LLeidanet API
func (ll *LLeidanet) Send() (LLeidaResp, error) {
	log.Println("sending sms")

	// TODO set endpoint as env var
	endpoint := "https://api.lleida.net/sms/v2/"
	// endpoint, ok := os.LookupEnv("LEAD_LEONTEL_ENDPOINT")
	// if !ok {
	// 	err := errors.New("unable to load LleidaNet URL endpoint")
	// 	return err
	// }
	ko := LLeidaResp{
		Request: "sms",
		Code:    500,
		Status:  "Error",
	}

	bytevalues, err := json.Marshal(ll)
	if err != nil {
		return ko, err
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(bytevalues))
	if err != nil {
		return ko, err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	lleidaR := LLeidaResp{}

	if err := xml.Unmarshal(data, &lleidaR); err != nil {
		return ko, err
	}
	return lleidaR, nil
}

func (h *Handler) pushError(err error) {
	h.Errors = append(h.Errors, err)
}
