package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

type CloudWatchClient struct {
	Region    string
	AccessKey string
	SecretKey string
	client    *cloudwatch.CloudWatch
}

func NewCloudWatchClient(region string, accesskey string, secretkey string) CloudWatchClient {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(accesskey, secretkey, ""),
		},
	}))
	return CloudWatchClient{
		client:    cloudwatch.New(sess),
		Region:    region,
		AccessKey: accesskey,
		SecretKey: secretkey,
	}
}
func (cw *CloudWatchClient) SendMetric(namespace string, key string, value float64) (*cloudwatch.PutMetricDataOutput, error) {
	result, err := cw.client.PutMetricData(&cloudwatch.PutMetricDataInput{
		MetricData: []*cloudwatch.MetricDatum{
			&cloudwatch.MetricDatum{
				MetricName: aws.String(key),
				Unit:       aws.String(cloudwatch.StandardUnitNone),
				Value:      aws.Float64(value),
				Dimensions: []*cloudwatch.Dimension{},
			},
		},
		Namespace: aws.String(namespace),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
