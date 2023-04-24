package ftcutils

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	REGION              = "ap-southeast-2"
	DATA_BUCKET         = "fireteam-core-army-data"
	FACTIONS_FILE_NAME  = "factions.json"
	FIRETEAMS_FILE_NAME = "fireteams.json"
	ARMY_PATH_PARAM     = "army_id"
)

func GetFileContents(fileName string) (string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:           aws.String(REGION),
		S3ForcePathStyle: &[]bool{true}[0],
	}))
	downloader := s3manager.NewDownloader(sess)

	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(DATA_BUCKET),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}

func WriteFileContents(fileName string, bodyJSON string) (*s3manager.UploadOutput, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:           aws.String(REGION),
		S3ForcePathStyle: &[]bool{true}[0],
	}))
	reader := strings.NewReader(bodyJSON)
	uploader := s3manager.NewUploader(sess)
	return uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(DATA_BUCKET),
		Key:    aws.String(fileName),
		Body:   reader,
	})
}

func ValidateJsonBody(bodyStr string) (map[string]interface{}, error) {
	var bodyJSON map[string]interface{}
	err := json.Unmarshal([]byte(bodyStr), &bodyJSON)
	if err != nil {
		return nil, err
	}
	return bodyJSON, nil
}

func BuildResponse(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       message,
		StatusCode: statusCode,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"content-type":                "application/json",
		},
	}, nil
}

func GetPathParameter(request events.APIGatewayProxyRequest, parameterName string) (string, error) {
	if request.PathParameters == nil {
		return "", errors.New("no path parameters")
	}
	pathParameter, ok := request.PathParameters[parameterName]
	if !ok {
		return "", errors.New("No path parameter with name " + parameterName)
	}
	return pathParameter, nil
}
