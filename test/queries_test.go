package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/checkmarxDev/ice/pkg/engine"
	"github.com/checkmarxDev/ice/pkg/engine/mock"
	"github.com/checkmarxDev/ice/pkg/model"
	"github.com/checkmarxDev/ice/pkg/parser"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type expectedResult struct {
	line     int
	severity model.Severity
	name     string
}

type testCase struct {
	query           string
	file            string
	expectedResults []expectedResult
}

var testCases = []testCase{
	{
		query: "ALB_Protocol_is_HTTP.q",
		file:  "ALB_protocol_is_HTTP.tf",
		expectedResults: []expectedResult{
			{
				line:     25,
				severity: model.SeverityHigh,
				name:     "ALB protocol is HTTP",
			},
			{
				line:     19,
				severity: model.SeverityHigh,
				name:     "ALB protocol is HTTP",
			},
		},
	},
	{
		query: "ALB_Protocol_is_HTTP.q",
		file:  "ALB_protocol_is_HTTP_success.tf",
	},
	{
		query: "Cloudfront_Configuration_Allow_HTTP.q",
		file:  "Cloudfront_configuration_allow_HTTP.tf",
		expectedResults: []expectedResult{
			{
				line:     49,
				severity: model.SeverityHigh,
				name:     "Cloudfront configuration allow HTTP",
			},
			{
				line:     76,
				severity: model.SeverityHigh,
				name:     "Cloudfront configuration allow HTTP",
			},
		},
	},
	{
		query: "Cloudfront_Configuration_Allow_HTTP.q",
		file:  "Cloudfront_configuration_allow_HTTP_success.tf",
	},
	{
		query: "Cloudwatch_without_Retention_Days.q",
		file:  "Cloudwatch_without_retention_days.tf",
		expectedResults: []expectedResult{
			{
				line:     15,
				severity: model.SeverityLow,
				name:     "Cloudwatch without retention days",
			},
		},
	},
	{
		query: "Cloudwatch_without_Retention_Days.q",
		file:  "Cloudwatch_without_retention_days_success.tf",
	},
	{
		query: "Cloudfront_without_WAF.q",
		file:  "Cloudfront_without_WAF.tf",
		expectedResults: []expectedResult{
			{
				line:     15,
				severity: model.SeverityLow,
				name:     "Cloudfront without WAF",
			},
		},
	},
	{
		query: "Cloudfront_without_WAF.q",
		file:  "Cloudfront_without_WAF_success.tf",
	},
	{
		query: "EKS_Cluster_Public_Access_cidrs.q",
		file:  "Eks_Cluster_Public_Access_cidrs.tf",
		expectedResults: []expectedResult{
			{
				line:     8,
				severity: model.SeverityHigh,
				name:     "EKS cluster public access cidrs",
			},
		},
	},
	{
		query: "EKS_Cluster_Public_Access_cidrs.q",
		file:  "Eks_Cluster_Public_Access_cidrs_success.tf",
	},
	{
		query: "EKS_Cluster_Public_Access.q",
		file:  "Eks_Cluster_Public_Access.tf",
		expectedResults: []expectedResult{
			{
				line:     7,
				severity: model.SeverityMedium,
				name:     "EKS cluster public access",
			},
		},
	},
	{
		query: "EKS_Cluster_Public_Access.q",
		file:  "Eks_Cluster_Public_Access_success.tf",
	},
	{
		query: "Fully_Open_Ingress.q",
		file:  "Fully_Open_Ingress.tf",
		expectedResults: []expectedResult{
			{
				line:     91,
				severity: model.SeverityHigh,
				name:     "Fully open Ingress",
			},
		},
	},
	{
		query: "Fully_Open_Ingress.q",
		file:  "Fully_Open_Ingress_success.tf",
	},
	{
		query: "Hard_Coded_AWS_Access_Key.q",
		file:  "Hard_Coded_AWS_Access_Key.tf",
		expectedResults: []expectedResult{
			{
				line:     14,
				severity: model.SeverityLow,
				name:     "Hardcoded AWS access key",
			},
		},
	},
	{
		query: "Hard_Coded_AWS_Access_Key.q",
		file:  "Hard_Coded_AWS_Access_Key_success.tf",
	},
	{
		query: "IAM_Policies_Allow_All.q",
		file:  "IAM_policies_allow_all.tf",
		expectedResults: []expectedResult{
			{
				line:     32,
				severity: model.SeverityMedium,
				name:     "Allow all IAM policies",
			},
		},
	},
	{
		query: "IAM_Policies_Allow_All.q",
		file:  "IAM_policies_allow_all_success.tf",
	},
	{
		query: "IAM_Policies_Attached_to_User.q",
		file:  "IAM_policies_attached_to_User.tf",
		expectedResults: []expectedResult{
			{
				line:     16,
				severity: model.SeverityLow,
				name:     "IAM policies attached to user",
			},
		},
	},
	{
		query: "IAM_Policies_Attached_to_User.q",
		file:  "IAM_policies_attached_to_User_success.tf",
	},
	{
		query: "IAM_Policies_with_Full_Pivileges.q",
		file:  "IAM_Policies_with_Full_Pivileges.tf",
		expectedResults: []expectedResult{
			{
				line:     1,
				severity: model.SeverityMedium,
				name:     "IAM policies with full privileges",
			},
		},
	},
	{
		query: "IAM_Policies_with_Full_Pivileges.q",
		file:  "IAM_Policies_with_Full_Pivileges_success.tf",
	},
	{
		query: "IAM_Role_Allows_Public_Assume.q",
		file:  "IAM_Role_Allows_Public_Assume.tf",
		expectedResults: []expectedResult{
			{
				line:     15,
				severity: model.SeverityLow,
				name:     "IAM role allows public assume",
			},
		},
	},
	{
		query: "IAM_Role_Allows_Public_Assume.q",
		file:  "IAM_Role_Allows_Public_Assume_success.tf",
	},
	{
		query: "IAM_Role_Assumed_by_All.q",
		file:  "IAM_Role_Assumed_by_All.tf",
		expectedResults: []expectedResult{
			{
				line:     26,
				severity: model.SeverityLow,
				name:     "IAM role allows all principals to assume",
			},
		},
	},
	{
		query: "IAM_Role_Assumed_by_All.q",
		file:  "IAM_Role_Assumed_by_All_success.tf",
	},
	{
		query: "Incorrect_Password_Policy_Experation.q",
		file:  "Incorrect_Password_Policy_Experation.tf",
		expectedResults: []expectedResult{
			{
				line:     1,
				severity: model.SeverityMedium,
				name:     "Incorrect password policy expiration",
			},
		},
	},
	{
		query: "Incorrect_Password_Policy_Experation.q",
		file:  "Incorrect_Password_Policy_Experation_success.tf",
	},
	{
		query: "Insufficient_Password_Length.q",
		file:  "Insufficient_Password_Length.tf",
		expectedResults: []expectedResult{
			{
				line:     1,
				severity: model.SeverityHigh,
				name:     "Insufficient password length",
			},
		},
	},
	{
		query: "Insufficient_Password_Length.q",
		file:  "Insufficient_Password_Length_success.tf",
	},
	{
		query: "Lamda_Hardcoded_AWS_Access_Key.q",
		file:  "Lamda_Hardcoded_AWS_Access_Key.tf",
		expectedResults: []expectedResult{
			{
				line:     21,
				severity: model.SeverityLow,
				name:     "Lambda hardcoded AWS access key",
			},
		},
	},
	{
		query: "Lamda_Hardcoded_AWS_Access_Key.q",
		file:  "Lamda_Hardcoded_AWS_Access_Key_success.tf",
	},
	{
		query: "Missing_Cluster_Log_Types.q",
		file:  "Missing_Cluster_Log_Types.tf",
		expectedResults: []expectedResult{
			{
				line:     6,
				severity: model.SeverityLow,
				name:     "Missing cluster log types",
			},
		},
	},
	{
		query: "Missing_Cluster_Log_Types.q",
		file:  "Missing_Cluster_Log_Types_success.tf",
	},
	{
		query: "No_Password_Reuse_Prevention.q",
		file:  "No_Password_Reuse_Prevention.tf",
		expectedResults: []expectedResult{
			{
				line:     3,
				severity: model.SeverityMedium,
				name:     "No password reuse prevention",
			},
		},
	},
	{
		query: "No_Password_Reuse_Prevention.q",
		file:  "No_Password_Reuse_Prevention_success.tf",
	},
	{
		query: "Not_Encypted_Data_in_Launch_Configuration.q",
		file:  "Not_Encypted_Data_in_Launch_Configuration.tf",
		expectedResults: []expectedResult{
			{
				line:     2,
				severity: model.SeverityMedium,
				name:     "Not encrypted data in launch configuration",
			},
		},
	},
	{
		query: "Not_Encypted_Data_in_Launch_Configuration.q",
		file:  "Not_Encypted_Data_in_Launch_Configuration_success.tf",
	},
	{
		query: "Open_Access_to_Resources_through_API.q",
		file:  "Open_Access_to_Resources_through_API.tf",
		expectedResults: []expectedResult{
			{
				line:     12,
				severity: model.SeverityLow,
				name:     "Open access to resources through API",
			},
		},
	},
	{
		query: "Open_Access_to_Resources_through_API.q",
		file:  "Open_Access_to_Resources_through_API_success.tf",
	},
	{
		query: "Public_ECR_Policy.q",
		file:  "Public_ECR_Policy.tf",
		expectedResults: []expectedResult{
			{
				line:     5,
				severity: model.SeverityMedium,
				name:     "Public ECR policy",
			},
		},
	},
	{
		query: "Public_ECR_Policy.q",
		file:  "Public_ECR_Policy_success.tf",
	},
	{
		query: "S3_Bucket_with_Ignore_Public_ACL.q",
		file:  "S3_Bucket_with_Ignore_Public_ACL.tf",
		expectedResults: []expectedResult{
			{
				line:     5,
				severity: model.SeverityLow,
				name:     "S3 bucket with ignore public ACL",
			},
		},
	},
	{
		query: "S3_Bucket_with_Ignore_Public_ACL.q",
		file:  "S3_Bucket_with_Ignore_Public_ACL_success.tf",
	},
	{
		query: "S3_Bucket_with_Public_ACL.q",
		file:  "S3_Bucket_with_Public_ACL.tf",
		expectedResults: []expectedResult{
			{
				line:     5,
				severity: model.SeverityMedium,
				name:     "S3 bucket allows public ACL",
			},
		},
	},
	{
		query: "S3_Bucket_with_Public_ACL.q",
		file:  "S3_Bucket_with_Public_ACL_success.tf",
	},
	{
		query: "S3_Bucket_with_Public_Policy.q",
		file:  "S3_Bucket_with_Public_Policy.tf",
		expectedResults: []expectedResult{
			{
				line:     5,
				severity: model.SeverityHigh,
				name:     "S3 bucket allows public policy",
			},
		},
	},
	{
		query: "S3_Bucket_with_Public_Policy.q",
		file:  "S3_Bucket_with_Public_Policy_success.tf",
	},
	{
		query: "S3_Bucket_with_any_Principal.q",
		file:  "S3_Bucket_with_any_Principal.tf",
		expectedResults: []expectedResult{
			{
				line:     5,
				severity: model.SeverityHigh,
				name:     "S3 bucket with any principal",
			},
		},
	},
	{
		query: "S3_Bucket_with_any_Principal.q",
		file:  "S3_Bucket_with_any_Principal_success.tf",
	},
	{
		query: "S3_Bucket_without_Enabled_MFA_Delete.q",
		file:  "S3_Bucket_without_Enabled_MFA_Delete.tf",
		expectedResults: []expectedResult{
			{
				line:     2,
				severity: model.SeverityHigh,
				name:     "S3 bucket without enabled MFA Delete",
			},
		},
	},
	{
		query: "S3_Bucket_without_Enabled_MFA_Delete.q",
		file:  "S3_Bucket_without_Enabled_MFA_Delete_success.tf",
	},
	{
		query: "S3_Bucket_without_Encryption_at_REST.q",
		file:  "S3_Bucket_without_Encryption_at_REST.tf",
		expectedResults: []expectedResult{
			{
				line:     2,
				severity: model.SeverityHigh,
				name:     "S3 bucket without encryption at REST",
			},
		},
	},
	{
		query: "S3_Bucket_without_Encryption_at_REST.q",
		file:  "S3_Bucket_without_Encryption_at_REST_success.tf",
	},
	{
		query: "S3_Bucket_without_Logging.q",
		file:  "S3_Bucket_without_Logging.tf",
		expectedResults: []expectedResult{
			{
				line:     2,
				severity: model.SeverityLow,
				name:     "S3 no logging",
			},
		},
	},
	{
		query: "S3_Bucket_without_Logging.q",
		file:  "S3_Bucket_without_Logging_success.tf",
	},
	{
		query: "S3_Bucket_without_Restriction_of_Public_Bucket.q",
		file:  "S3_Bucket_without_Restriction_of_Public_Bucket.tf",
		expectedResults: []expectedResult{
			{
				line:     5,
				severity: model.SeverityHigh,
				name:     "S3 bucket without restriction of public buckety",
			},
		},
	},
	{
		query: "S3_Bucket_without_Restriction_of_Public_Bucket.q",
		file:  "S3_Bucket_without_Restriction_of_Public_Bucket_success.tf",
	},
	{
		query: "S3_Bucket_without_Versioning.q",
		file:  "S3_Bucket_without_Versioning.tf",
		expectedResults: []expectedResult{
			{
				line:     2,
				severity: model.SeverityHigh,
				name:     "S3 bucket without versioning",
			},
		},
	},
	{
		query: "S3_Bucket_without_Versioning.q",
		file:  "S3_Bucket_without_Versioning_success.tf",
	},
	{
		query: "S3_Bucket_wth_Public_RW.q",
		file:  "S3_Bucket_wth_Public_RW.tf",
		expectedResults: []expectedResult{
			{
				line:     2,
				severity: model.SeverityInfo,
				name:     "S3 bucket with public RW access",
			},
		},
	},
	{
		query: "S3_Bucket_wth_Public_RW.q",
		file:  "S3_Bucket_wth_Public_RW_success.tf",
	},
	{
		query: "SQS_Policy_with_ALL_Actions.q",
		file:  "SQS_Policy_with_ALL_Actions.tf",
		expectedResults: []expectedResult{
			{
				line:     5,
				severity: model.SeverityMedium,
				name:     "SQS policy allows ALL (*) actions",
			},
		},
	},
	{
		query: "SQS_Policy_with_ALL_Actions.q",
		file:  "SQS_Policy_with_ALL_Actions_success.tf",
	},
	{
		query: "Unchangeable_Password.q",
		file:  "Unchangeable_Password.tf",
		expectedResults: []expectedResult{
			{
				line:     6,
				severity: model.SeverityMedium,
				name:     "Unchangeable password",
			},
		},
	},
	{
		query: "Unchangeable_Password.q",
		file:  "Unchangeable_Password_success.tf",
	},
}

func TestQueries(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: ioutil.Discard})

	for _, testCase := range testCases {
		t.Run(testCase.query, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			scanID := "scanID"
			fileID := 111

			storage := mock.NewMockFilesStorage(ctrl)
			storage.EXPECT().GetFiles(gomock.Eq(ctx), gomock.Eq(scanID), gomock.Any()).
				DoAndReturn(func(ctx context.Context, scanID, filter string) (model.FileMetadatas, error) {
					filePath := path.Join("./test-data", testCase.file)
					f, err := os.Open(filePath)
					if err != nil {
						return nil, err
					}

					content, err := ioutil.ReadAll(f)
					if err != nil {
						return nil, err
					}

					jsonContent, err := parser.NewDefault().Parse(filePath, content)
					if err != nil {
						return nil, err
					}

					return []model.FileMetadata{
						{
							ID:           fileID,
							ScanID:       scanID,
							JSONData:     jsonContent,
							OriginalData: string(content),
							Kind:         model.KindTerraform,
							FileName:     filePath,
							JSONHash:     0,
						},
					}, nil
				})

			storage.EXPECT().SaveVulnerabilities(gomock.Eq(ctx), gomock.Any()).
				DoAndReturn(func(_ context.Context, results []model.Vulnerability) error {
					require.Len(t, results, len(testCase.expectedResults), "Found issues and expected doesn't match")

					for i, item := range testCase.expectedResults {
						if i > len(results)-1 {
							t.Fatalf("Not enough results detected, expected %d, found %d", len(testCase.expectedResults), len(results))
						}

						result := results[i]
						require.Equal(t, item.line, result.Line, "Not corrected detected line")
						require.Equal(t, item.severity, result.Severity, "Invalid severity")
						require.Equal(t, item.name, result.QueryName, "Invalid query name")
						require.Equal(t, fileID, result.FileID)
					}

					return nil
				})

			queriesSource := mock.NewMockQueriesSource(ctrl)
			queriesSource.EXPECT().GetQueries().
				DoAndReturn(func() ([]model.QueryMetadata, error) {
					qCode, err := ioutil.ReadFile(path.Join("../assets/queries", testCase.query))
					if err != nil {
						return nil, fmt.Errorf("query source: %w", err)
					}

					return []model.QueryMetadata{
						{
							FileName: testCase.query,
							Content:  string(qCode),
						},
					}, nil
				})

			inspector, err := engine.NewInspector(ctx, queriesSource, storage)
			require.Nil(t, err)
			require.NotNil(t, inspector)

			err = inspector.Inspect(ctx, scanID)
			assert.Nil(t, err)
		})
	}
}