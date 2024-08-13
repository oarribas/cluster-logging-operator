package azuremonitor

import (
	_ "embed"
	"testing"

	. "github.com/openshift/cluster-logging-operator/internal/constants"
	vectorhelpers "github.com/openshift/cluster-logging-operator/internal/generator/vector/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	loggingv1 "github.com/openshift/cluster-logging-operator/api/logging/v1"
	. "github.com/openshift/cluster-logging-operator/test/matchers"
	corev1 "k8s.io/api/core/v1"
)

const (
	sharedKey           = "z9ndQSFH1RLDnS6WR35m84u326p3"
	sharedKeyWithDollar = "mykeyWith$Sign"
	azureId             = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/otherResourceGroup/providers/Microsoft.Storage/storageAccounts/examplestorage"
	hostCN              = "ods.opinsights.azure.cn"
	customerId          = "6vzw6sHc-0bba-6sHc-4b6c-8bz7sr5eggRt"
	secretName          = "azure-monitor-secret"
	secretTlsName       = "azure-monitor-secret-tls"
	outputName          = "azure_monitor_logs"
	secretWithDollar    = "azure-monitor-secret-dollar"
	logType             = "myLogType"
)

//go:embed azm_common.toml
var ExpectedAzureCommonToml string

//go:embed azm_common_dollar.toml
var ExpectedAzureCommonDollarToml string

//go:embed azm_advance.toml
var ExpectedAzureAdvanceToml string

//go:embed azm_tls.toml
var ExpectedAzureTlsToml string

//go:embed azm_skip_tls.toml
var ExpectedAzureSkipTls string

var _ = Describe("Generating vector config for Azure Monitor Logs output:", func() {

	var (
		outputCommon = loggingv1.OutputSpec{
			Type: loggingv1.OutputTypeAzureMonitor,
			Name: outputName,
			OutputTypeSpec: loggingv1.OutputTypeSpec{
				AzureMonitor: &loggingv1.AzureMonitor{
					CustomerId: customerId,
					LogType:    logType,
				},
			},
			Secret: &loggingv1.OutputSecretSpec{
				Name: secretName,
			},
		}

		outputCommonWithDollar = loggingv1.OutputSpec{
			Type: loggingv1.OutputTypeAzureMonitor,
			Name: outputName,
			OutputTypeSpec: loggingv1.OutputTypeSpec{
				AzureMonitor: &loggingv1.AzureMonitor{
					CustomerId: customerId,
					LogType:    logType,
				},
			},
			Secret: &loggingv1.OutputSecretSpec{
				Name: secretWithDollar,
			},
		}

		outputCommonTls = loggingv1.OutputSpec{
			Type: loggingv1.OutputTypeAzureMonitor,
			Name: outputName,
			OutputTypeSpec: loggingv1.OutputTypeSpec{
				AzureMonitor: &loggingv1.AzureMonitor{
					CustomerId: customerId,
					LogType:    logType,
				},
			},
			Secret: &loggingv1.OutputSecretSpec{
				Name: secretTlsName,
			},
		}

		outputAdvance = loggingv1.OutputSpec{
			Type: loggingv1.OutputTypeAzureMonitor,
			Name: outputName,
			OutputTypeSpec: loggingv1.OutputTypeSpec{
				AzureMonitor: &loggingv1.AzureMonitor{
					CustomerId:      customerId,
					LogType:         logType,
					AzureResourceId: azureId,
					Host:            hostCN,
				},
			},
			Secret: &loggingv1.OutputSecretSpec{
				Name: secretName,
			},
		}

		outputWithTlsSkipVerify = loggingv1.OutputSpec{
			Type: loggingv1.OutputTypeAzureMonitor,
			Name: outputName,
			OutputTypeSpec: loggingv1.OutputTypeSpec{
				AzureMonitor: &loggingv1.AzureMonitor{
					CustomerId: customerId,
					LogType:    logType,
				},
			},
			Secret: &loggingv1.OutputSecretSpec{
				Name: secretName,
			},
			TLS: &loggingv1.OutputTLSSpec{
				InsecureSkipVerify: true,
			},
		}

		secrets = map[string]*corev1.Secret{
			secretName: {
				Data: map[string][]byte{
					SharedKey: []byte(sharedKey),
				},
			},
			secretTlsName: {
				Data: map[string][]byte{
					SharedKey:          []byte(sharedKey),
					ClientPrivateKey:   []byte("junk"),
					ClientCertKey:      []byte("junk"),
					TrustedCABundleKey: []byte("junk"),
					Passphrase:         []byte("foo"),
				},
			},
			secretWithDollar: {
				Data: map[string][]byte{
					SharedKey: []byte(sharedKeyWithDollar),
				},
			},
		}
	)

	DescribeTable("should generate valid config", func(outputSpec loggingv1.OutputSpec, secret, expValue string) {
		conf := New(vectorhelpers.MakeOutputID(outputSpec.Name), outputSpec, []string{"pipelineName"}, secrets[secret], nil, nil)
		Expect(expValue).To(EqualConfigFrom(conf))
	},
		Entry("for common case", outputCommon, secretName, ExpectedAzureCommonToml),
		Entry("for advance case", outputAdvance, secretName, ExpectedAzureAdvanceToml),
		Entry("for common with tls case", outputCommonTls, secretTlsName, ExpectedAzureTlsToml),
		Entry("for common with skip tls verify true case", outputWithTlsSkipVerify, secretName, ExpectedAzureSkipTls),
		Entry("for common case with sharedKey containing `$`", outputCommonWithDollar, secretWithDollar, ExpectedAzureCommonDollarToml),
	)
})

func TestVectorConfGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vector for Azure Log Monitor New Generation")
}
