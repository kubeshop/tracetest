package junit_test

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/junit"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConversion(t *testing.T) {
	test := model.Test{
		Name: "Example test",
	}

	results := model.RunResults{
		Results: (maps.Ordered[model.SpanQuery, []model.AssertionResult]{}).MustAdd(
			model.SpanQuery(`span[tracetest.span.type = "database"]`), []model.AssertionResult{
				{
					Assertion: `attr:db.statement contains "INSERT"`,
					Results: []model.SpanAssertionResult{
						{
							ObservedValue: "INSERT into whatever",
							CompareErr:    nil,
						},
					},
				},
				{
					Assertion: `attr:tracetest.span.duration < 500`,
					Results: []model.SpanAssertionResult{
						{
							ObservedValue: "notANumber",
							CompareErr:    errors.New(`cannot parse "notANumber" as integer`),
						},
					},
				},
				{
					Assertion: `attr:tracetest.span.type = "http"`,
					Results: []model.SpanAssertionResult{
						{
							ObservedValue: "database",
							CompareErr:    comparator.ErrNoMatch,
						},
					},
				},
			}),
	}

	run := model.Run{
		CreatedAt:   time.Date(2022, 05, 23, 14, 55, 07, 0, time.UTC),
		CompletedAt: time.Date(2022, 05, 23, 14, 55, 18, 0, time.UTC),
		Results:     &results,
	}

	expected, err := os.ReadFile("testdata/junit_result.xml")
	require.NoError(t, err)

	actual, err := junit.FromRunResult(test, run)
	require.NoError(t, err)

	t.Log("File:", string(actual))

	assertXMLEqual(t, bytes.NewReader(expected), bytes.NewReader(actual))

}

func assertXMLEqual(t *testing.T, expected io.Reader, obtained io.Reader) {
	var expec node
	var obt node
	var err error

	err = xml.NewDecoder(expected).Decode(&expec)
	assert.Nil(t, err)
	err = xml.NewDecoder(obtained).Decode(&obt)
	assert.Nil(t, err)

	assert.Equal(t, expec, obt)
}

type node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Content string     `xml:",innerxml"`
	Nodes   []node     `xml:",any"`
}

func (n *node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type inNode node

	err := d.DecodeElement((*inNode)(n), &start)
	if err != nil {
		return err
	}

	//Discard content if there are child nodes
	if len(n.Nodes) > 0 {
		n.Content = ""
	}
	return nil
}
