package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func printSuccess(value string) {
	if jsonOutput == false {
		fmt.Println(value)
	} else {
		buf := new(bytes.Buffer)
		response := struct {
			Result string
		}{Result: value}

		responseString, _ := json.Marshal(response)
		fmt.Fprintf(buf, "%s", responseString)
		buf.WriteTo(os.Stdout)
	}
}

func printFailure(value string) {
	if jsonOutput == false {
		fmt.Println(value)
	} else {
		buf := new(bytes.Buffer)
		response := struct {
			Error string
		}{Error: value}

		responseString, _ := json.Marshal(response)
		fmt.Fprintf(buf, "%s", responseString)
		buf.WriteTo(os.Stderr)
	}
}
