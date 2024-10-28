package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type RDSInfo struct {
	ID       string
	Endpoint string
	Port     int32
}

func GetRDSInstancesByProfile(profile string) ([]RDSInfo, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, err
	}

	rdsClient := rds.NewFromConfig(cfg)

	result, err := rdsClient.DescribeDBInstances(context.TODO(), &rds.DescribeDBInstancesInput{})
	if err != nil {
		return nil, err
	}

	var instances []RDSInfo
	for _, dbInstance := range result.DBInstances {
		instances = append(instances, RDSInfo{
			ID:       *dbInstance.DBInstanceIdentifier,
			Endpoint: *dbInstance.Endpoint.Address,
			Port:     *dbInstance.Endpoint.Port,
		})
	}
	return instances, nil
}

func GetRDSInstanceByDatabaseID(profile, databaseID string) (*RDSInfo, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, err
	}

	rdsClient := rds.NewFromConfig(cfg)

	result, err := rdsClient.DescribeDBInstances(context.TODO(), &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: &databaseID,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching RDS instance %s: %v", databaseID, err)
	}

	if len(result.DBInstances) == 0 {
		return nil, fmt.Errorf("RDS instance %s not found", databaseID)
	}

	dbInstance := result.DBInstances[0]
	return &RDSInfo{
		ID:       *dbInstance.DBInstanceIdentifier,
		Endpoint: *dbInstance.Endpoint.Address,
		Port:     *dbInstance.Endpoint.Port,
	}, nil
}
