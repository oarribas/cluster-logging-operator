package fluentd

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/sets"
	"strings"
)

const (
	lokiLabelKubernetesHost = "kubernetes.host"
	lokiLabelTag            = "tag"
)

var (
	defaultLabelKeys = []string{
		"log_type",

		//container labels
		"kubernetes.namespace_name",
		"kubernetes.pod_name",
		"kubernetes.container_name",
	}
	requiredLokiLabelKeys = []string{
		lokiLabelKubernetesHost,
		lokiLabelTag,
	}
)

func (conf *outputLabelConf) lokiLabelKeys() []string {
	var keys sets.String
	if conf.Target.Loki != nil && len(conf.Target.Loki.LabelKeys) != 0 {
		keys = sets.NewString(conf.Target.Loki.LabelKeys...)
	} else {
		keys = sets.NewString(defaultLabelKeys...)
	}
	// Ensure required tags for serialization
	keys.Insert(requiredLokiLabelKeys...)
	return keys.List()
}

// LokiLabelFilter generates record_modifier filter lines to copy Loki label fields.
// The Loki output plugin will remove these fields after creating Loki labels.
func (conf *outputLabelConf) LokiLabelFilter() string {
	b := &strings.Builder{}
	for _, k := range conf.lokiLabelKeys() {
		tempName := strings.Replace(k, ".", "_", -1)
		recordKeys := strings.Replace(k, ".", `","`, -1)
		switch k {
		case lokiLabelTag:
			fmt.Fprintf(b, "\n	_tag ${tag}")
		case lokiLabelKubernetesHost:
			fmt.Fprintf(b, "\n    _%v \"#{ENV['NODE_NAME']}\"", tempName)
		default:
			fmt.Fprintf(b, "\n    _%v ${record.dig(\"%v\")}", tempName, recordKeys)
		}
	}
	return b.String()
}

// LokiLabel generates the <label> entries for Loki output config.
// This consumes the fields generated by LokiLabelFilter.
func (conf *outputLabelConf) LokiLabel() string {
	b := &strings.Builder{}
	for _, k := range conf.lokiLabelKeys() {
		tempName := strings.Replace(k, ".", "_", -1)
		fmt.Fprintf(b, "\n      %v _%v", tempName, tempName)
	}
	return b.String()
}

// LokiTenantKeys returns the components of the loki tenant key.
func (conf *outputLabelConf) LokiTenantKeys() []string {
	if conf.Target.Loki != nil && conf.Target.Loki.TenantKey != "" {
		return strings.Split(conf.Target.Loki.TenantKey, ".")
	}
	return nil
}
