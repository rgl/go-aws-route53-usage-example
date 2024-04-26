package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

type Zone struct {
	Name        string
	NameServers []string
}

func getZones(ctx context.Context, client *route53.Client) ([]Zone, error) {
	response, err := client.ListHostedZones(ctx, nil)
	if err != nil {
		return nil, err
	}

	zones := []Zone{}

	for _, hostedZone := range response.HostedZones {
		zone, err := client.GetHostedZone(ctx, &route53.GetHostedZoneInput{Id: hostedZone.Id})
		if err != nil {
			return nil, err
		}
		zones = append(zones, Zone{
			Name:        *zone.HostedZone.Name,
			NameServers: zone.DelegationSet.NameServers})
	}

	return zones, nil
}

func main() {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load the default configuration: %v", err)
	}

	client := route53.NewFromConfig(cfg)

	zones, err := getZones(ctx, client)
	if err != nil {
		log.Fatalf("failed to retrieve the aws route53 hosted dns zones: %v", err)
	}

	for _, zone := range zones {
		fmt.Printf("zone=%s nameservers=%s\n", zone.Name, strings.Join(zone.NameServers, ","))
	}

	if len(zones) == 0 {
		log.Printf("NB The AWS account does not have any DNS Zones.")
	}
}
