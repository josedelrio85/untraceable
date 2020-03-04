package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	untraceable "github.com/bysidecar/untraceable/pkg"
)

func main() {
	log.Printf("Untraceable process launched at %s", time.Now().Format("2006-01-02 15-04-05"))

	port := getSetting("DB_PORT")
	portInt, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing to string Database's port %s, Err: %s", port, err)
	}

	database := &untraceable.Database{
		Host:      getSetting("DB_HOST"),
		Port:      portInt,
		User:      getSetting("DB_USER"),
		Pass:      getSetting("DB_PASS"),
		Dbname:    getSetting("DB_NAME"),
		Charset:   "utf8",
		ParseTime: "True",
		Loc:       "Local",
	}

	if err := database.Open(); err != nil {
		msg := "error opening database connection"
		e := untraceable.ErrorLogger{
			Msg:    msg,
			Status: http.StatusInternalServerError,
			Err:    err,
			Log:    untraceable.LogError(err),
		}
		e.SendAlarm()
		log.Fatalf("error opening database connection. err: %s", err)
	}
	defer database.Close()

	if err := database.AutoMigrate(); err != nil {
		msg := "error creating the table"
		e := untraceable.ErrorLogger{
			Msg:    msg,
			Status: http.StatusInternalServerError,
			Err:    err,
			Log:    untraceable.LogError(err),
		}
		e.SendAlarm()
		log.Fatalf("error creating the table. err: %s", err)
	}

	// TODO set user/password as env var
	handler := untraceable.Handler{
		Storer: database,
		LLeidanet: untraceable.LLeidanet{
			Sms: untraceable.Parameters{
				User:     "bysidecar",
				Password: "wwvkys",
				Text:     "Te hemos llamado pero no hemos podido localizarte. Queremos ofrecerte una promo exclusiva para ti. ¿Nos devuelves la llamada? Este es nuestro número ",
			},
		},
	}

	if err := handler.GetTraced(); err != nil {
		msg := "error retriving traced leads"
		e := untraceable.ErrorLogger{
			Msg:    msg,
			Status: http.StatusInternalServerError,
			Err:    err,
			Log:    untraceable.LogError(err),
		}
		e.SendAlarm()
		log.Fatalf("error retriving traced leads. err: %s", err)
	}

	if err := handler.GetUntraceables(); err != nil {
		msg := "error retrieving untraceable leads"
		e := untraceable.ErrorLogger{
			Msg:    msg,
			Status: http.StatusInternalServerError,
			Err:    err,
			Log:    untraceable.LogError(err),
		}
		e.SendAlarm()
		log.Fatalf("error retrieving untraceable leads. err: %s", err)
	}

	handler.Fire()

	//test send
	// handler.LLeidanet.Sms.Source = "Test"
	// handler.LLeidanet.Sms.Destination.Number = []string{"665932355"}
	// resp, err := handler.LLeidanet.Send()
	// log.Println(resp)

	if len(handler.Errors) > 0 {
		for _, err := range handler.Errors {
			msg := "Error sending sms"
			e := untraceable.ErrorLogger{
				Msg:    msg,
				Status: http.StatusInternalServerError,
				Err:    err,
				Log:    untraceable.LogError(err),
			}
			e.SendAlarm()
		}
	}

	log.Printf("Untraceable process finished OK at %s", time.Now().Format("2006-01-02 15-04-05"))
}

func getSetting(setting string) string {
	value, ok := os.LookupEnv(setting)
	if !ok {
		log.Fatalf("Init error, %s ENV var not found", setting)
	}

	return value
}
