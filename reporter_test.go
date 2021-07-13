package datadog

import (
	"github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitNameAndTagsNoTags(t *testing.T) {
	reporter := &Reporter{
		Registry: metrics.DefaultRegistry,
	}

	var metricName = "test.metric_name"
	name, tags := reporter.splitNameAndTags(metricName)
	assert.NotNil(t, name, "name")
	assert.Nil(t, tags, "tags")

	assert.Equal(t, "test.metric_name", name)
}

func TestSplitNameAndTagsSingleTag(t *testing.T) {
	reporter := &Reporter{
		Registry: metrics.DefaultRegistry,
	}

	var metricName = "test.httpcall[method:GET]"

	name, tags := reporter.splitNameAndTags(metricName)
	assert.NotNil(t, name, "name")
	assert.NotNil(t, tags, "tags")
	assert.Equal(t, 1, len(tags))

	assert.Equal(t, "test.httpcall", name)
	assert.Equal(t, "method:GET", tags[0])
}

func TestSplitNameAndTagsMultipleTag(t *testing.T) {
	globalTags := []string{"globaltag:true"}
	reporter := &Reporter{
		Registry: metrics.DefaultRegistry,
		tags:     globalTags,
	}

	var metricName = "test.httpcall[method:GET]"

	name, tags := reporter.splitNameAndTags(metricName)
	assert.NotNil(t, name, "name")
	assert.NotNil(t, tags, "tags")
	//expect two tags: the global tag and the metric level tag
	assert.Equal(t, 2, len(tags))

	assert.Equal(t, "test.httpcall", name)
	assert.Equal(t, "method:GET", tags[0])
	assert.Equal(t, "globaltag:true", tags[1])
}
