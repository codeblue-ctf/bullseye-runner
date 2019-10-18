package master

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"text/template"

	"google.golang.org/grpc"
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

func NewUUID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return fmt.Sprintf("%016x", buf), nil
	// return fmt.Sprintf("%04x-%02x-%02x-%02x-%06x", buf[:4], buf[4:6], buf[6:8], buf[8:10], buf[10:]), nil
}

func CreateGrpcCli(host string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	return grpc.Dial(host, opts...)
}

func Debug(p interface{}) {
	log.Printf("%+v", p)
}
