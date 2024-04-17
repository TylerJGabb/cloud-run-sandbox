package logging

import (
	"encoding/json"
	"log"
)

type Entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`
	Extras map[string]interface{} `json:"extra,omitempty"`
}

type Logger struct {
	// TODO: pass this in when you build a logger through a constructor
	// Extract an interface!
	Trace string
}

// ToJson renders an entry structure to the JSON format expected by Cloud Logging.
func (e Entry) ToJson() string {
	if e.Severity == "" {
		e.Severity = "INFO"
	}
	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}

func keyValuesToExtras(keyValue ...any) map[string]any {
	var key string
	result := map[string]any{}
	for idx, item := range keyValue {
		if idx%2 == 0 {
			key = item.(string)
		} else {
			result[key] = item
		}
	}
	return result
}

func (l Logger) Error(m string, keyValues ...interface{}) {
	e := Entry{
		Message:     m,
		Severity:    "ERROR",
		Trace:       l.Trace,
		Extras: keyValuesToExtras(keyValues...),
	}
	log.Println(e.ToJson())
}

func (l Logger) Info(m string, keyValues ...interface{}) {
	e := Entry{
		Message:     m,
		Trace:       l.Trace,
		Extras: keyValuesToExtras(keyValues...),
	}
	log.Println(e.ToJson())
}

func (l Logger) Warn(m string, keyValues ...interface{}) {
	e := Entry{
		Message:     m,
		Severity:    "WARNING",
		Trace:       l.Trace,
		Extras: keyValuesToExtras(keyValues...),
	}
	log.Println(e.ToJson())
}

func (l Logger) Debug(m string, keyValues ...interface{}) {
	e := Entry{
		Message:     m,
		Severity:    "DEBUG",
		Trace:       l.Trace,
		Extras: keyValuesToExtras(keyValues...),
	}
	log.Println(e.ToJson())
}