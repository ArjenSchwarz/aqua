package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
)

func printSuccess(value string) {
	if !aws.BoolValue(settings.JSONOutput) {
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

func printSliceMaps(values []map[string]string) {
	buf := new(bytes.Buffer)
	if !aws.BoolValue(settings.JSONOutput) {
		for _, item := range values {
			for key, value := range item {
				fmt.Fprintf(buf, "%s: %s\n", key, value)
			}
		}
	} else {
		responseString, _ := json.Marshal(values)
		fmt.Fprintf(buf, "%s", responseString)
	}
	buf.WriteTo(os.Stdout)
}

func printFailure(value string) {
	if !aws.BoolValue(settings.JSONOutput) {
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
