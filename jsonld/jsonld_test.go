package jsonld

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/cayleygraph/quad"
	"github.com/cayleygraph/quad/voc/xsd"
	"github.com/piprate/json-gold/ld"
	"github.com/stretchr/testify/require"
)

var testReadCases = []struct {
	data   string
	expect []quad.Quad
}{
	{
		`{
  "@context": {
    "ex": "http://example.org/",
    "term1": {"@id": "ex:term1", "@type": "ex:datatype"},
    "term2": {"@id": "ex:term2", "@type": "@id"},
    "term3": {"@id": "ex:term3", "@language": "en"}
  },
  "@id": "ex:id1",
  "@type": ["ex:Type1", "ex:Type2"],
  "term1": "v1",
  "term2": "ex:id2",
  "term3": "v3"
}`,
		[]quad.Quad{
			{
				Subject:   quad.IRI(`http://example.org/id1`),
				Predicate: quad.IRI(`http://example.org/term1`),
				Object: quad.TypedString{
					Value: "v1", Type: "http://example.org/datatype",
				},
				Label: nil,
			},
			{
				Subject:   quad.IRI(`http://example.org/id1`),
				Predicate: quad.IRI(`http://example.org/term2`),
				Object:    quad.IRI(`http://example.org/id2`),
				Label:     nil,
			},
			{
				Subject:   quad.IRI(`http://example.org/id1`),
				Predicate: quad.IRI(`http://example.org/term3`),
				Object: quad.LangString{
					Value: "v3", Lang: "en",
				},
				Label: nil,
			},
			{
				Subject:   quad.IRI(`http://example.org/id1`),
				Predicate: quad.IRI(`http://www.w3.org/1999/02/22-rdf-syntax-ns#type`),
				Object:    quad.IRI(`http://example.org/Type1`),
				Label:     nil,
			},
			{
				Subject:   quad.IRI(`http://example.org/id1`),
				Predicate: quad.IRI(`http://www.w3.org/1999/02/22-rdf-syntax-ns#type`),
				Object:    quad.IRI(`http://example.org/Type2`),
				Label:     nil,
			},
		},
	},
}

type ByQuad []quad.Quad

func (a ByQuad) Len() int           { return len(a) }
func (a ByQuad) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByQuad) Less(i, j int) bool { return a[i].NQuad() < a[j].NQuad() }

func TestRead(t *testing.T) {
	for i, c := range testReadCases {
		r := NewReader(strings.NewReader(c.data))
		quads, err := quad.ReadAll(r)
		if err != nil {
			t.Errorf("case %d failed: %v", i, err)
		}
		sort.Sort(ByQuad(quads))
		sort.Sort(ByQuad(c.expect))
		if !reflect.DeepEqual(quads, c.expect) {
			t.Errorf("case %d failed: wrong quads returned:\n%v\n%v", i, quads, c.expect)
		}
		r.Close()
	}
}

var testWriteCases = []struct {
	data   []quad.Quad
	ctx    interface{}
	expect string
}{
	{
		[]quad.Quad{
			{
				Subject:   quad.IRI(`http://example.org/id1`),
				Predicate: quad.IRI(`http://example.org/term1`),
				Object: quad.TypedString{
					Value: "v1", Type: "http://example.org/datatype",
				},
				Label: nil,
			},
			{
				Subject:   quad.IRI(`http://example.org/id1`),
				Predicate: quad.IRI(`http://example.org/term2`),
				Object:    quad.IRI(`http://example.org/id2`),
				Label:     nil,
			},
			{
				Subject:   quad.IRI(`http://example.org/id1`),
				Predicate: quad.IRI(`http://example.org/term3`),
				Object: quad.LangString{
					Value: "v3", Lang: "en",
				},
				Label: nil,
			},
			{
				Subject:   quad.IRI(`http://example.org/id1`),
				Predicate: quad.IRI(`http://www.w3.org/1999/02/22-rdf-syntax-ns#type`),
				Object:    quad.IRI(`http://example.org/Type1`),
				Label:     nil,
			},
			{
				Subject:   quad.IRI(`http://example.org/id1`),
				Predicate: quad.IRI(`http://www.w3.org/1999/02/22-rdf-syntax-ns#type`),
				Object:    quad.IRI(`http://example.org/Type2`),
				Label:     nil,
			},
		},
		map[string]interface{}{
			"ex": "http://example.org/",
		},
		`{
  "@context": {
    "ex": "http://example.org/"
  },
  "@id": "ex:id1",
  "@type": [
    "ex:Type1",
    "ex:Type2"
  ],
  "ex:term1": {
    "@type": "ex:datatype",
    "@value": "v1"
  },
  "ex:term2": {
    "@id": "ex:id2"
  },
  "ex:term3": {
    "@language": "en",
    "@value": "v3"
  }
}
`,
	},
}

func TestWrite(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	for i, c := range testWriteCases {
		buf.Reset()
		w := NewWriter(buf)
		w.SetLdContext(c.ctx)
		_, err := quad.Copy(w, quad.NewReader(c.data))
		if err != nil {
			t.Errorf("case %d failed: %v", i, err)
		} else if err = w.Close(); err != nil {
			t.Errorf("case %d failed: %v", i, err)
		}
		data := make([]byte, buf.Len())
		copy(data, buf.Bytes())
		buf.Reset()
		json.Indent(buf, data, "", "  ")
		if buf.String() != c.expect {
			t.Errorf("case %d failed: wrong data returned:\n%v\n%v", i, buf.String(), c.expect)
		}
	}
}

var testRoundtripCases = []struct {
	data []quad.Quad
}{
	{
		[]quad.Quad{
			{
				Subject:   quad.IRI(`http://example.org/id1`),
				Predicate: quad.IRI(`http://example.org/term1`),
				Object: quad.TypedString{
					Value: "v1", Type: "http://example.org/datatype",
				},
				Label: nil,
			},
			{
				Subject:   quad.IRI(`http://example.org/id1`),
				Predicate: quad.IRI(`http://example.org/term2`),
				Object:    quad.IRI(`http://example.org/id2`),
				Label:     nil,
			},
			{
				Subject:   quad.IRI(`http://example.org/id1`),
				Predicate: quad.IRI(`http://example.org/term3`),
				Object: quad.LangString{
					Value: "v3", Lang: "en",
				},
				Label: nil,
			},
			{
				Subject:   quad.IRI(`http://example.org/id1`),
				Predicate: quad.IRI(`http://www.w3.org/1999/02/22-rdf-syntax-ns#type`),
				Object:    quad.IRI(`http://example.org/Type1`),
				Label:     nil,
			},
			{
				Subject:   quad.IRI(`http://example.org/id1`),
				Predicate: quad.IRI(`http://www.w3.org/1999/02/22-rdf-syntax-ns#type`),
				Object:    quad.IRI(`http://example.org/Type2`),
				Label:     nil,
			},
		},
	},
}

func TestRoundtrip(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	for i, c := range testRoundtripCases {
		buf.Reset()
		w := NewWriter(buf)
		_, err := quad.Copy(w, quad.NewReader(c.data))
		if err != nil {
			t.Errorf("case %d failed: %v", i, err)
		} else if err = w.Close(); err != nil {
			t.Errorf("case %d failed: %v", i, err)
		}
		arr, err := quad.ReadAll(NewReader(buf))
		sort.Sort(quad.ByQuadString(arr))
		sort.Sort(quad.ByQuadString(c.data))
		if err != nil {
			t.Errorf("case %d failed: %v", i, err)
		} else if !reflect.DeepEqual(arr, c.data) {
			t.Errorf("case %d failed: wrong data returned:\n%v\n%v", i, arr, c.data)
		}
	}
}

var fromValueTestCases = []struct {
	name   string
	value  quad.Value
	jsonLd interface{}
}{
	{
		name:   "Simple text",
		value:  quad.String("Alice"),
		jsonLd: "Alice",
	},
	{
		name:   "Localized text",
		value:  quad.LangString{Value: "Alice", Lang: "en"},
		jsonLd: map[string]interface{}{"@value": "Alice", "@language": "en"},
	},
	{
		name:   "Known typed string",
		value:  quad.TypedString{Value: quad.String("Alice"), Type: xsd.String},
		jsonLd: "Alice",
	},
	{
		name:   "Known typed integer",
		value:  quad.Int(1),
		jsonLd: int64(1),
	},
	{
		name:   "Known typed floating-point number",
		value:  quad.Float(1.0),
		jsonLd: 1.0,
	},
	{
		name:   "Known typed boolean",
		value:  quad.Bool(true),
		jsonLd: true,
	},
	{
		name:  "Datetime",
		value: quad.Time(time.Time{}),
		jsonLd: map[string]interface{}{
			"@value": "0001-01-01T00:00:00Z",
			"@type":  xsd.DateTime,
		},
	},
}

func TestFromValue(t *testing.T) {
	for _, c := range fromValueTestCases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.jsonLd, FromValue(c.value))
		})
	}
}

var toValueTestCases = []struct {
	name   string
	jsonLd ld.Node
	value  quad.Value
}{
	{
		name:   "Simple text",
		jsonLd: ld.NewLiteral("Alice", "", ""),
		value:  quad.String("Alice"),
	},
	{
		name:   "Localized text",
		jsonLd: ld.NewLiteral("Alice", "", "en"),
		value:  quad.LangString{Value: "Alice", Lang: "en"},
	},
	{
		name:   "Known typed string",
		jsonLd: ld.NewLiteral("Alice", xsd.String, ""),
		value:  quad.String("Alice"),
	},
	{
		name:   "Known typed integer",
		jsonLd: ld.NewLiteral("1", xsd.Integer, ""),
		value:  quad.Int(1),
	},
	{
		name:   "Known typed floating-point number (xsd:double)",
		jsonLd: ld.NewLiteral("1.1", xsd.Double, ""),
		value:  quad.Float(1.1),
	},
	{
		name:   "Known typed floating-point number (xsd:float)",
		jsonLd: ld.NewLiteral("1.1", xsd.Float, ""),
		value:  quad.Float(1.1),
	},
	{
		name:   "Known typed boolean",
		jsonLd: ld.NewLiteral("true", xsd.Boolean, ""),
		value:  quad.Bool(true),
	},
	{
		name:   "Datetime",
		jsonLd: ld.NewLiteral("0001-01-01T00:00:00Z", xsd.DateTime, ""),
		value:  quad.Time(time.Time{}),
	},
}

func TestToValue(t *testing.T) {
	for _, c := range toValueTestCases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.value, toValue(c.jsonLd))
		})
	}
}
