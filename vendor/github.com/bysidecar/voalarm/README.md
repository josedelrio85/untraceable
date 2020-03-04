# VOALARM

Component to send alarms to VictorOps plattform. This plattform will advise through several methods if any problem occurs.

```
  // Create a custom error to test component
	anerror := fmt.Errorf("Error. Artificial error: %v", errors.New("emit macho dwarf: elf header corrupted"))

  // Creates an instance of voalarm client, setting the platform params automatically.
	alarm := voalarm.NewClient("")
  
  // Send alarm and handle the error
	resp, err := alarm.SendAlarm(voalarm.Acknowledgement, anerror)
	if err != nil {
		log.Fatalf("Error creating alarm. Err: %s", err)
	}
  
```