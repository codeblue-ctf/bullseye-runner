package master

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"regexp"
	"sync"
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

func NewUUID() string {
	buf := make([]byte, 16)
	rand.Read(buf)
	return fmt.Sprintf("%016x", buf)
	// return fmt.Sprintf("%04x-%02x-%02x-%02x-%06x", buf[:4], buf[4:6], buf[6:8], buf[8:10], buf[10:]), nil
}

func CreateGrpcCli(host string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	return grpc.Dial(host, opts...)
}

type CancelManager struct {
	mut sync.Mutex
	c   map[string]context.CancelFunc
}

func NewCancelManager() *CancelManager {
	return &CancelManager{
		c: make(map[string]context.CancelFunc),
	}
}

func (cm *CancelManager) Has(key string) bool {
	cm.mut.Lock()
	defer cm.mut.Unlock()
	_, ok := cm.c[key]
	return ok
}

func (cm *CancelManager) Keys() []string {
	cm.mut.Lock()
	defer cm.mut.Unlock()
	res := []string{}
	for k := range cm.c {
		res = append(res, k)
	}
	return res
}

func (cm *CancelManager) Add(key string, _ctx context.Context) (context.Context, error) {
	cm.mut.Lock()
	defer cm.mut.Unlock()
	if _, ok := cm.c[key]; ok {
		return nil, fmt.Errorf("key %s already exists", key)
	}

	ctx, cancel := context.WithCancel(_ctx)
	cm.c[key] = cancel
	return ctx, nil
}

func (cm *CancelManager) Cancel(key string) error {
	cm.mut.Lock()
	defer cm.mut.Unlock()
	cancel, ok := cm.c[key]
	if !ok {
		return fmt.Errorf("key %s does not exist", key)
	}
	cancel()
	delete(cm.c, key)
	return nil
}

func Debug(p interface{}) {
	log.Printf("%+v", p)
}
