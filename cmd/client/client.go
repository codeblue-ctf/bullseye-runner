package main

import (
	"context"
	"flag"
	"log"

	"google.golang.org/grpc/credentials"

	pb "gitlab.com/CBCTF/bullseye-runner/proto"
	"google.golang.org/grpc"
)

var (
	caFile = flag.String("ca", "", "CA root cert file")
	host   = flag.String("host", "192.168.121.51:10080", "Server address")
)

func sendRequest(client pb.RunnerClient, req *pb.RunnerRequest) (*pb.RunnerResponse, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	ctx := context.Background()
	res, err := client.Run(ctx, req)
	if err != nil {
		log.Fatalf("%v.Run(_) = _, %v", client, err)
	}
	log.Printf("%+v", res)

	return nil, nil
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption

	if *caFile != "" {
		creds, err := credentials.NewClientTLSFromFile(*caFile, *host)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(*host, opts...)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewRunnerClient(conn)

	req := pb.RunnerRequest{
		Uuid:    "hoge",
		Timeout: 1000,
		Yml: `
version: '2'

services:
  exploit:
    image: localhost:5000/team01/test-exploit{{if .exploitHash}}@{{.exploitHash}}{{else}}{{end}}
    depends_on:
      - challenge
      - flag-submit
  challenge:
    image: localhost:5000/test-challenge
    volumes:
      - "{{.flagPath}}:/flag"
    expose:
      - "8080"
  flag-submit:
    image: localhost:5000/flag-submit
    volumes:
      - "{{.submitPath}}:/flag"
    expose:
      - "1337"
`[1:],
		RegistryToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IkdSWVM6WTVKSDpXQlJLOktMNkc6TFJBRDpOQ0pCOk5ZQ0I6SU5ZTjpaU0E1OllDRlc6SkxQRTo3UkpYIn0.eyJpc3MiOiJBdXRoIFNlcnZpY2UiLCJzdWIiOiJhZG1pbiIsImF1ZCI6IkRvY2tlciByZWdpc3RyeSIsImV4cCI6MTU2OTE4MDg4MCwibmJmIjoxNTY5MTc5OTcwLCJpYXQiOjE1NjkxNzk5ODAsImp0aSI6Ijc2OTc2NjM1NDkzMjk3NDY0MjciLCJhY2Nlc3MiOlt7InR5cGUiOiJyZXBvc2l0b3J5IiwibmFtZSI6IjEvc2FuaXR5LWNoZWNrIiwiYWN0aW9ucyI6WyJwdWxsIl19XX0.mcdiut68HrnXxkRG8GW5IT17hk8JMAy5L9MMRRld59YS6PXshFiIcw3AWDsjYy8ejQ3rJoNWSR90pYhmD29CSo2pYVL0qJQOLMxWORFYUG1fD_acq8UExPLnjgJ2e96jhg9ZRS2WJ8C-qN8CkwLnKdx5-2mEE9jZMephdHSGSmmsGvT3ficxjiiyRzi1xpK5tdGoE3V0gv5NQcDVhJM8iKrNM0PBL8uDgBj7AFJrQ2y_IPTimcRHNtAylR-vaWFLwga6ASdtGO88BLJW4L_OG5vxwtSxE_lkf1jqhYgr_XJCN5IBYe548uuiESYzkuG2qycjwJi6I_yE01EM1klLikKnk4UYThOiuJ1kOOs3v69JshoZhSqNWQvYA2-VQXN99yq8BmBMEMfmqdZAiEBsyAM6m6jTg-AmjSWU8u4Mb4t2RPDLHqxxyZWZ3bcym0X4DPwOR1HIvnwIVjEtbXR8FPTLYyCmsZha302dyltS71D0q6EFMWg_fIMt1yYYzhLR2wgFG0xf8UN6Dk-BL5MvgCKw45Hd01llSh9JtNRfv_dYkwKu_A7E9kqPXmUe4FRdRhIipdnEsV_d4432uzYCizaO0V7VEGImbKX07MXulbqJPp3OZVpYcLT08BQxkBliIEll_yQMN6vlNdsNhL0u6afQRtRKZohua3I3D-GRMWM",
		FlagTemplate:  "CBCTF{[a-f0-9]{16}}",
	}

	sendRequest(client, &req)

	// for i := 0; i < 10; i++ {
	// 	go sendRequest(client, &req)
	// 	time.Sleep(1 * time.Second)
	// }

	// time.Sleep(60 * time.Second)

}
