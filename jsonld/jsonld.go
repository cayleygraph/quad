// Package jsonld provides an encoder/decoder for JSON-LD quad format
package jsonld

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/cayleygraph/quad"
	"github.com/cayleygraph/quad/voc"
	"github.com/cayleygraph/quad/voc/xsd"
	"github.com/piprate/json-gold/ld"
)

// AutoConvertTypedString allows to convert TypedString values to native
// equivalents directly while parsing. It will call ToNative on all TypedString values.
//
// If conversion error occurs, it will preserve original TypedString value.
var AutoConvertTypedString = true

func init() {
	quad.RegisterFormat(quad.Format{
		Name:   "jsonld",
		Ext:    []string{".jsonld"},
		Mime:   []string{"application/ld+json"},
		Writer: func(w io.Writer) quad.WriteCloser { return NewWriter(w) },
		Reader: func(r io.Reader) quad.ReadCloser { return NewReader(r) },
	})
}

// NewReader returns quad reader for JSON-LD stream.
func NewReader(r io.Reader) *Reader {
	var o interface{}
	if err := json.NewDecoder(r).Decode(&o); err != nil {
		return &Reader{err: err}
	}
	return NewReaderFromMap(o)
}

// NewReaderFromMap returns quad reader for JSON-LD map object.
func NewReaderFromMap(o interface{}) *Reader {
	opts := ld.NewJsonLdOptions("")
	processor := ld.NewJsonLdProcessor()
	data, err := processor.ToRDF(o, opts)
	if err != nil {
		return &Reader{err: err}
	}
	return &Reader{
		graphs: data.(*ld.RDFDataset).Graphs,
	}
}

var _ quad.Reader = &Reader{}

// Reader implements the quad.Reader interface
type Reader struct {
	err    error
	name   string
	n      int
	graphs map[string][]*ld.Quad
}

// ReadQuad implements the quad.Reader interface
func (r *Reader) ReadQuad() (quad.Quad, error) {
	if r.err != nil {
		return quad.Quad{}, r.err
	}
next:
	if len(r.graphs) == 0 {
		return quad.Quad{}, io.EOF
	}
	if r.name == "" {
		for gname := range r.graphs {
			r.name = gname
			break
		}
	}
	if r.n >= len(r.graphs[r.name]) {
		r.n = 0
		delete(r.graphs, r.name)
		r.name = ""
		goto next
	}
	cur := r.graphs[r.name][r.n]
	r.n++
	var graph quad.Value
	if r.name != "@default" {
		graph = quad.IRI(r.name)
	}
	return quad.Quad{
		Subject:   toValue(cur.Subject),
		Predicate: toValue(cur.Predicate),
		Object:    toValue(cur.Object),
		Label:     graph,
	}, nil
}

// Close implements quad.Reader
func (r *Reader) Close() error {
	r.graphs = nil
	return r.err
}

var _ quad.Writer = &Writer{}

// Writer implements quad.Writer
type Writer struct {
	w   io.Writer
	ds  *ld.RDFDataset
	ctx interface{}
}

// NewWriter constructs a new Writer
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w, ds: ld.NewRDFDataset()}
}

// SetLdContext defines a context for the emitted JSON-LD data
// See: https://json-ld.org/spec/latest/json-ld/#the-context
func (w *Writer) SetLdContext(ctx interface{}) {
	w.ctx = ctx
}

// WriteQuad implements quad.Writer
func (w *Writer) WriteQuad(q quad.Quad) error {
	if !q.IsValid() {
		return quad.ErrInvalid
	}
	var graph string
	if q.Label == nil {
		graph = "@default"
	} else if iri, ok := q.Label.(quad.IRI); ok {
		graph = string(iri)
	} else {
		graph = q.Label.String()
	}
	g := w.ds.Graphs[graph]
	g = append(g, ld.NewQuad(
		toTerm(q.Subject),
		toTerm(q.Predicate),
		toTerm(q.Object),
		graph,
	))
	w.ds.Graphs[graph] = g
	return nil
}

// WriteQuads implements quad.Writer
func (w *Writer) WriteQuads(buf []quad.Quad) (int, error) {
	for i, q := range buf {
		if err := w.WriteQuad(q); err != nil {
			return i, err
		}
	}
	return len(buf), nil
}

// Close implements quad.Writer
func (w *Writer) Close() error {
	opts := ld.NewJsonLdOptions("")
	api := ld.NewJsonLdApi()
	processor := ld.NewJsonLdProcessor()
	var data interface{}
	data, err := api.FromRDF(w.ds, opts)
	if err != nil {
		return err
	}
	if w.ctx != nil {
		out, err := processor.Compact(data, w.ctx, opts)
		if err != nil {
			return err
		}
		data = out
	}
	return json.NewEncoder(w.w).Encode(data)
}

func toTerm(v quad.Value) ld.Node {
	switch v := v.(type) {
	case quad.IRI:
		return ld.NewIRI(string(v))
	case quad.BNode:
		return ld.NewBlankNode(string(v))
	case quad.String:
		return ld.NewLiteral(string(v), "", "")
	case quad.TypedString:
		return ld.NewLiteral(string(v.Value), string(v.Type), "")
	case quad.LangString:
		return ld.NewLiteral(string(v.Value), "", string(v.Lang))
	case quad.TypedStringer:
		return toTerm(v.TypedString())
	default:
		return ld.NewLiteral(v.String(), "", "")
	}
}

// FromValue converts quad value to a JSON-LD compatible object.
func FromValue(v quad.Value) interface{} {
	switch v := v.(type) {
	case quad.IRI:
		return map[string]interface{}{
			"@id": string(v),
		}
	case quad.BNode:
		return map[string]interface{}{
			"@id": v.String(),
		}
	case quad.String:
		return string(v)
	case quad.LangString:
		return map[string]interface{}{
			"@value":    string(v.Value),
			"@language": string(v.Lang),
		}
	case quad.TypedString:
		return typedStringToJSON(v)
	case quad.TypedStringer:
		return typedStringToJSON(v.TypedString())
	default:
		return v.String()
	}
}

func isKnownTimeType(dataType quad.IRI) bool {
	for _, iri := range quad.KnownTimeTypes {
		if iri == dataType {
			return true
		}
	}
	return false
}

func typedStringToJSON(v quad.TypedString) interface{} {
	if AutoConvertTypedString && quad.HasStringConversion(v.Type) && !isKnownTimeType(v.Type) {
		return v.Native()
	}
	return map[string]interface{}{
		"@value": string(v.Value),
		"@type":  string(v.Type),
	}
}

var stringDataType = voc.FullIRI(xsd.String)

func toValue(t ld.Node) quad.Value {
	switch t := t.(type) {
	case *ld.IRI:
		return quad.IRI(t.Value)
	case *ld.BlankNode:
		return quad.BNode(t.Attribute)
	case *ld.Literal:
		if t.Language != "" {
			return quad.LangString{
				Value: quad.String(t.Value),
				Lang:  t.Language,
			}
		} else if t.Datatype != "" && t.Datatype != stringDataType {
			ts := quad.TypedString{
				Value: quad.String(t.Value),
				Type:  quad.IRI(t.Datatype),
			}
			if AutoConvertTypedString {
				if v, err := ts.ParseValue(); err == nil {
					return v
				}
			}
			return ts
		}
		return quad.String(t.Value)
	default:
		panic(fmt.Errorf("unexpected term type: %T", t))
	}
}
