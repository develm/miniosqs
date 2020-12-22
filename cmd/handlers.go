package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// reference: https://docs.aws.amazon.com/AmazonS3/latest/dev/notification-content-structure.html
type Event struct {
	EventVersion string `json:"eventVersion"`
	EventSource  string `json:"eventSource"`
	AwsRegion    string `json:"awsRegion"`
	EventTime    string `json:"eventTime"`
	EventName    string `json: "eventName"`
	UserIdentity struct {
		PrincipalId string `json:"principalId"`
	}
	RequestParameters struct {
		SourceIPAddress string `json:"sourceIPAddress"`
	}
	ResponseElements struct {
		XAmzRequestId string `json:"x-amz-request-id"`
		XAmzId2       string `json:"x-amz-id-2"`
	}
	S3 struct {
		S3SchemaVersion string `json:"s3SchemaVersion"`
		ConfigurationId string `json:"configurationId"`
		Bucket          struct {
			Name          string `json:"name"`
			OwnerIdentity struct {
				PrincipalId string `json:"principalId"`
			}
			Arn string `json:"arn"`
		}
		Object struct {
			Key       string `json:"key"`
			Size      int    `json:"size"`
			ETag      string `json:"eTag"`
			VersionId string `json:"versionId"`
			Sequencer string `json:"sequencer"`
		}
	}
}

type Records struct {
	Records []Event
}

type Application struct {
	ErrorLog  *log.Logger
	InfoLog   *log.Logger
	SqsClient *sqs.SQS
	QueueUrl  string
}

func (app *Application) notifySqs(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	if r.ContentLength == 0 {
		// minio appears to first make an empty post with each event
		app.InfoLog.Printf("Empty request body")
		fmt.Fprintf(w, "Empty body")
		return
	}

	dec := json.NewDecoder(r.Body)

	var e Records
	err := dec.Decode(&e)
	if err != nil {
		app.ErrorLog.Printf("Failed to decode JSON: %s", err)
		app.serverError(w, err)
		return
	}
	app.InfoLog.Printf("Event: %+v", e)

	j, err := json.Marshal(e)
	if err != nil {
		app.ErrorLog.Printf("JSON marshal error: %s", err)
		app.serverError(w, err)
		return
	}

	result, err := app.SqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(j)),
		QueueUrl:    &app.QueueUrl,
	})
	if err != nil {
		app.ErrorLog.Printf("Failed to send message: %s", err)
		app.serverError(w, err)
		return
	}
	fmt.Fprintf(w, "Message ID: %s", *result.MessageId)
	app.InfoLog.Printf("SQS message created with ID: %s", *result.MessageId)
}

//TODO setup bucket (create bucket & configure notification)
func (app *Application) setupMinioBucket(w http.ResponseWriter, r *http.Request) {

}
