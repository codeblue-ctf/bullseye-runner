package master

import (
	"bytes"
	"context"
	"log"
	"reflect"
	"regexp"
	"text/template"

	pb "gitlab.com/CBCTF/bullseye-runner/proto"
)

// EscapedTemplate executes template without replacing undefined field
func EscapedTemplate(s string, params interface{}) (string, error) {
	r := regexp.MustCompile(`\{\{\s*(\.[a-zA-Z0-9_^\}\s]*)\s*\}\}`)
	s = r.ReplaceAllString(s, `{{ if $1 }}{{$1}}{{ else }}{{"{{$1}}"}}{{ end }}`)

	tpl, err := template.New("").Parse(s)
	if err != nil {
		return "", err
	}

	val := reflect.ValueOf(params)
	typ := val.Type()
	kind := val.Kind()

	dict := make(map[string]interface{})
	if kind == reflect.Struct {
		for i := 0; i < typ.NumField(); i++ {
			field := val.Field(i)
			dict[typ.Field(i).Name] = field.Interface()
		}
	} else if kind == reflect.Map {
		for _, k := range val.MapKeys() {
			v := val.MapIndex(k)
			dict[k.String()] = v
		}
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, dict)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func SendRequest(client pb.RunnerClient, req *pb.RunnerRequest) (*pb.RunnerResponse, error) {
	// TODO: replace mock
	log.Printf("%+v", req)
	return nil, nil

	ctx := context.Background()
	res, err := client.Run(ctx, req)
	if err != nil {
		log.Fatalf("%v.Run(_) = _, %v", client, err)
	}
	log.Printf("%+v", res)

	return nil, nil
}
