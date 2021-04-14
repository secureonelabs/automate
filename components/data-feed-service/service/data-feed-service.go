package service

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	cfgmgmtRequest "github.com/chef/automate/api/interservice/cfgmgmt/request"
	cfgmgmt "github.com/chef/automate/api/interservice/cfgmgmt/service"
	rrule "github.com/teambition/rrule-go"

	"github.com/chef/automate/components/data-feed-service/config"
	"github.com/chef/automate/components/data-feed-service/dao"
	"github.com/chef/automate/lib/cereal"
	grpccereal "github.com/chef/automate/lib/cereal/grpc"

	"github.com/chef/automate/lib/grpc/secureconn"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const version = "1"

var creds *credentials.Credentials

const (
	// add your refion eg: us-east-2
	AWS_S3_REGION = ""
	AWS_S3_BUCKET = ""
	AccessKey     = ""
	SecretKey     = ""
)

var sess = connectAWS()

type datafeedNotification struct {
	credentials Credentials
	url         string
	data        bytes.Buffer
	contentType string
}

type DataClient struct {
	client              http.Client
	acceptedStatusCodes []int32
}

func NewDataClient(statusCodes []int32) DataClient {
	return DataClient{
		client:              http.Client{},
		acceptedStatusCodes: statusCodes,
	}
}

type StructData struct {
	Data []string `json:"data"`
}

type NotificationSender interface {
	sendNotification(notification datafeedNotification) error
}

func Start(dataFeedConfig *config.DataFeedConfig, connFactory *secureconn.Factory, db *dao.DB) error {
	fmt.Println("s3:::::", dataFeedConfig.S3.Accept)
	log.Info("Starting data-feed-service")
	conn, err := connFactory.Dial("cereal-service", dataFeedConfig.CerealConfig.Target)
	if err != nil {
		return err
	}

	backend := grpccereal.NewGrpcBackendFromConn("data-feed-service", conn)
	manager, err := cereal.NewManager(backend)
	if err != nil {
		return err
	}

	cfgMgmtConn, err := connFactory.Dial("config-mgmt-service", dataFeedConfig.CfgmgmtConfig.Target)
	if err != nil {
		return errors.Wrap(err, "could not connect to config-mgmt-service")
	}

	complianceConn, err := connFactory.Dial("compliance-service", dataFeedConfig.ComplianceConfig.Target)
	if err != nil {
		return errors.Wrap(err, "could not connect to compliance-service")
	}

	secretsConn, err := connFactory.Dial("secrets-service", dataFeedConfig.SecretsConfig.Target)
	if err != nil {
		return errors.Wrap(err, "could not connect to secrets-service")
	}

	dataFeedPollTask, err := NewDataFeedPollTask(dataFeedConfig, cfgMgmtConn, complianceConn, db, manager)
	if err != nil {
		return errors.Wrap(err, "could not create data feed poll task")
	}

	dataFeedAggregateTask := NewDataFeedAggregateTask(dataFeedConfig, cfgMgmtConn, complianceConn, secretsConn, db)

	err = manager.RegisterWorkflowExecutor(dataFeedWorkflowName, &DataFeedWorkflowExecutor{workflowName: dataFeedWorkflowName, dataFeedConfig: dataFeedConfig, manager: manager})
	if err != nil {
		return err
	}
	err = manager.RegisterTaskExecutor(dataFeedPollTaskName, dataFeedPollTask, cereal.TaskExecutorOpts{
		Workers: 1,
	})
	if err != nil {
		return err
	}
	err = manager.RegisterTaskExecutor(dataFeedAggregateTaskName, dataFeedAggregateTask, cereal.TaskExecutorOpts{
		Workers: 1,
	})
	if err != nil {
		return err
	}

	r, err := rrule.NewRRule(rrule.ROption{
		Freq:     rrule.SECONDLY,
		Interval: int(dataFeedConfig.ServiceConfig.FeedInterval.Seconds()),
		Dtstart:  time.Now().Round(dataFeedConfig.ServiceConfig.FeedInterval).Add(30 * time.Second),
	})
	if err != nil {
		return err
	}

	dataFeedWorkflowParams := DataFeedWorkflowParams{}

	err = manager.CreateWorkflowSchedule(context.Background(),
		dataFeedScheduleName, dataFeedWorkflowName,
		dataFeedWorkflowParams, true, r)
	if err == cereal.ErrWorkflowScheduleExists {
		err = manager.UpdateWorkflowScheduleByName(context.Background(), dataFeedScheduleName, dataFeedWorkflowName,
			cereal.UpdateParameters(dataFeedWorkflowParams),
			cereal.UpdateEnabled(true),
			cereal.UpdateRecurrence(r))
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return manager.Start(context.Background())
}

func handleSendErr(notification datafeedNotification, startTime time.Time, endTime time.Time, err error) {
	log.Errorf("Failed to send notification to %v. Start: %v, End: %v. %v", notification.url, startTime, endTime, err)
}

func send(sender NotificationSender, notification datafeedNotification) error {
	return sender.sendNotification(notification)
}

func (client DataClient) sendNotification(notification datafeedNotification) error {

	log.Debugf("sendNotification bytes length %v", notification.data.Len())
	var contentBuffer bytes.Buffer
	zip := gzip.NewWriter(&contentBuffer)
	_, err := zip.Write(notification.data.Bytes())
	if err != nil {
		return err
	}
	err = zip.Close()
	if err != nil {
		return err
	}

	t := time.Now().UTC()
	year := t.Year()
	month := int(t.Month())
	day := t.Day()
	hr := t.Hour()
	min := t.Minute()
	sec := t.Second()

	filename :=
		strconv.Itoa(year) + "/" +
			strconv.Itoa(month) + "/" +
			strconv.Itoa(day) + "/" +
			strconv.Itoa(hr) + ":" + strconv.Itoa(min) + ":" + strconv.Itoa(sec) + ".json"
	// p := make([]byte, len(s))

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(AWS_S3_BUCKET),                  // Bucket to be used
		Key:    aws.String(filename),                       // Name of the file to be saved
		Body:   bytes.NewReader(notification.data.Bytes()), // File
	})
	if err != nil {
		// Do your error handling here
		return err
	}

	// request, err := http.NewRequest("POST", notification.url, &contentBuffer)
	// if err != nil {
	// 	log.Error("Error creating request")
	// 	return err
	// }
	// request.Header.Add("Authorization", notification.credentials.GetAuthorizationHeaderValue())
	// request.Header.Add("Content-Type", notification.contentType)
	// request.Header.Add("Content-Encoding", "gzip")
	// request.Header.Add("Accept", notification.contentType)
	// request.Header.Add("Chef-Data-Feed-Message-Version", version)

	// response, err := client.client.Do(request)

	// if err != nil {
	// 	log.Errorf("Error sending message %v", err)
	// 	return err
	// } else {
	// 	log.Infof("Asset data posted to %v, Status %v", notification.url, response.Status)
	// }
	// if !config.IsAcceptedStatusCode(int32(response.StatusCode), client.acceptedStatusCodes) {
	// 	return errors.New(response.Status)
	// }
	// err = response.Body.Close()
	// if err != nil {
	// 	log.Warnf("Error closing response body %v", err)
	// }
	return nil
}

func addDataContent(nodeDataContent map[string]interface{}, attributes map[string]interface{}) {
	os, _ := attributes["os"].(string)
	if strings.ToLower(os) == "windows" {
		kernel, ok := attributes["kernel"].(map[string]interface{})
		if !ok {
			nodeDataContent["serial_number"] = ""
			nodeDataContent["os_service_pack"] = ""
			return
		}

		osInfo, ok := kernel["os_info"].(map[string]interface{})
		if !ok {
			nodeDataContent["serial_number"] = ""
			nodeDataContent["os_service_pack"] = ""
			return
		}
		nodeDataContent["serial_number"] = osInfo["serial_number"]
		nodeDataContent["os_service_pack"] = ""
		majorVersion, ok := osInfo["service_pack_major_version"].(float64)
		if !ok {
			return
		}
		minorVersion, ok := osInfo["service_pack_minor_version"].(float64)
		if !ok {
			return
		}
		servicePackMajorVersion := fmt.Sprintf("%g", majorVersion)
		servicePackMinorVersion := fmt.Sprintf("%g", minorVersion)
		servicePack := strings.Join([]string{servicePackMajorVersion, servicePackMinorVersion}, ".")
		nodeDataContent["os_service_pack"] = servicePack
	} else {
		// assume linux
		dmi, _ := attributes["dmi"].(map[string]interface{})
		system, _ := dmi["system"].(map[string]interface{})
		serialNumber := system["serial_number"]
		if serialNumber == nil {
			serialNumber = ""
		}
		nodeDataContent["serial_number"] = serialNumber
	}
}

func getNodeFields(ctx context.Context, client cfgmgmt.CfgMgmtServiceClient, filters []string) (string, string, error) {

	nodeFilters := &cfgmgmtRequest.Nodes{Filter: filters}
	nodes, err := client.GetNodes(ctx, nodeFilters)
	if err != nil {
		log.Errorf("Error getting cfgmgmt/nodes %v", err)
		return "", "", err
	}

	if len(nodes.Values) == 0 {
		log.Debug("no node data exists for this node")
		return "", "", nil
	}
	node := nodes.Values[0].GetStructValue()
	id := node.Fields["id"].GetStringValue()
	lastRunId := node.Fields["latest_run_id"].GetStringValue()

	return id, lastRunId, nil

}

func getNodeAttributes(ctx context.Context, client cfgmgmt.CfgMgmtServiceClient, nodeId string) (map[string]interface{}, error) {

	attributesJson := make(map[string]interface{})

	nodeAttributes, err := client.GetAttributes(ctx, &cfgmgmtRequest.Node{NodeId: nodeId})
	if err != nil {
		log.Warnf("Error getting attributes %v", err)
		return attributesJson, err
	}

	attributesJson["automatic"] = getAttributesAsJson(nodeAttributes.Automatic, "automatic")
	attributesJson["default"] = getAttributesAsJson(nodeAttributes.Default, "default")
	attributesJson["normal"] = getAttributesAsJson(nodeAttributes.Normal, "normal")
	attributesJson["override"] = getAttributesAsJson(nodeAttributes.Override, "override")
	attributesJson["all_value_count"] = nodeAttributes.AllValueCount
	attributesJson["automatic_value_count"] = nodeAttributes.AutomaticValueCount
	attributesJson["default_value_count"] = nodeAttributes.DefaultValueCount
	attributesJson["normal_value_count"] = nodeAttributes.NormalValueCount
	attributesJson["override_value_count"] = nodeAttributes.OverrideValueCount
	attributesJson["node_id"] = nodeAttributes.NodeId
	attributesJson["name"] = nodeAttributes.Name
	attributesJson["run_list"] = nodeAttributes.RunList
	attributesJson["chef_environment"] = nodeAttributes.ChefEnvironment

	return attributesJson, nil
}

func getAttributesAsJson(attributes string, attributeType string) map[string]interface{} {
	attributesJson := make(map[string]interface{})
	err := json.Unmarshal([]byte(attributes), &attributesJson)
	if err != nil {
		log.Errorf("Could not parse %v attributes from json: %v", attributeType, err)
	}
	return attributesJson
}

func getNodeHostFields(ctx context.Context, client cfgmgmt.CfgMgmtServiceClient, filters []string) (string, string, string, error) {
	nodeId, _, err := getNodeFields(ctx, client, filters)
	if err != nil {
		return "", "", "", err
	}
	attributesJson, err := getNodeAttributes(ctx, client, nodeId)
	if err != nil {
		return "", "", "", err
	}
	ipaddress, macAddress, hostname := getHostAttributes(attributesJson["automatic"].(map[string]interface{}))
	return ipaddress, macAddress, hostname, nil
}

func getHostAttributes(attributesJson map[string]interface{}) (string, string, string) {

	ipAddress, _ := attributesJson["ipaddress"].(string)
	macAddress, _ := attributesJson["macaddress"].(string)
	hostname, _ := attributesJson["hostname"].(string)

	return ipAddress, macAddress, hostname
}

func connectAWS() *session.Session {
	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(AWS_S3_REGION),
			Credentials: credentials.NewStaticCredentials(string(AccessKey), string(SecretKey), ""),
		},
	)
	if err != nil {
		panic(err)
	}
	return sess
}
