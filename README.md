# Quad formats for Go

![Tests](https://github.com/cayleygraph/quad/actions/workflows/tests.yml/badge.svg)

This library provides encoding and decoding support for NQuad/NTriple-compatible formats.

## Supported formats

| ID            | Name         | Read | Write | Ext           |
|---------------|--------------|------|-------|---------------|
| `nquads`      | NQuads       | +    | +     | `.nq`, `.nt`  |
| `jsonld`      | JSON-LD      | +    | +     | `.jsonld`     |
| `graphviz`    | DOT/Graphviz | -    | +     | `.gv`, `.dot` |
| `gml`         | GML          | -    | +     | `.gml`        |
| `graphml`     | GraphML      | -    | +     | `.graphml`    |
| `pquads`      | ProtoQuads   | +    | +     | `.pq`         |
| `json`        | JSON         | +    | +     | `.json`       |
| `json-stream` | JSON Stream  | +    | +     | -             |

## Community

* Slack: [cayleygraph.slack.com](https://cayleygraph.slack.com) -- Invite [here](https://cayley-slackin.herokuapp.com/)
