package configs

import "github.com/namsral/flag"

var (
	S3Endpoint   string
	S3AccessKey  string
	S3SecretKey  string
	S3BacketName string
)

func ConfigS3(commandSet *flag.FlagSet) {
	if command := commandSet.Lookup("s3-endpoint"); command != nil {
		S3Endpoint = command.Value.String()
	}
	if command := commandSet.Lookup("s3-access-key"); command != nil {
		S3AccessKey = command.Value.String()
	}
	if command := commandSet.Lookup("s3-secret-key"); command != nil {
		S3SecretKey = command.Value.String()
	}
	if command := commandSet.Lookup("s3-bucket-name"); command != nil {
		S3BacketName = command.Value.String()
	}
}
