package main

import (
	"encoding/base64"
	"flag"
	"fmt"

	"github.com/golang/protobuf/proto"
	structpb "github.com/golang/protobuf/ptypes/struct"
)

var (
	app               string
	version           string
	owner             string
	serviceAccount    string
	workloadName      string
	workloadNamespace string
)

func init() {
	flag.StringVar(&app, "app", "unknown", "value for 'app' label")
	flag.StringVar(&version, "version", "unknown", "value for 'version' label")
	flag.StringVar(&owner, "owner", "unknown", "value for the owner of the workload")
	flag.StringVar(&serviceAccount, "service_account", "unknown", "value for the service account running the workload")
	flag.StringVar(&workloadName, "workload_name", "unknown", "value for the name of the workload")
	flag.StringVar(&workloadNamespace, "workload_namespace", "unknown", "value for the namespace of the workload")
	flag.Parse()
}

func structString(value string) *structpb.Value {
	return &structpb.Value{
		Kind: &structpb.Value_StringValue{
			StringValue: value,
		},
	}
}

func structMap(value map[string]string) *structpb.Value {

	inner := &structpb.Struct{Fields: map[string]*structpb.Value{}}
	for k, v := range value {
		inner.Fields[k] = structString(v)
	}

	return &structpb.Value{
		Kind: &structpb.Value_StructValue{
			StructValue: inner,
		},
	}
}

func main() {
	metadata := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"OWNER":           structString(owner),
			"WORKLOAD_NAME":   structString(workloadName),
			"NAMESPACE":       structString(workloadNamespace),
			"LABELS":          structMap(map[string]string{"app": app, "version": version}),
			"SERVICE_ACCOUNT": structString(serviceAccount),
		},
	}

	var out []byte
	buf := proto.NewBuffer(out)
	buf.SetDeterministic(true)
	buf.Marshal(metadata)

	fmt.Printf("%s\n", base64.StdEncoding.EncodeToString(buf.Bytes()))
}
