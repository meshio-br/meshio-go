#!/usr/bin/env python3
"""Gera os tipos Go a partir de api/openapi.public.json e api/events/.

Uso (na raiz do repositório):

    python3 internal/gen/genmodels.py && gofmt -w models_gen.go params_gen.go events_gen.go

Ferramenta de desenvolvimento: o código gerado é versionado.
"""
import json
import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
SPEC = json.load(open(os.path.join(ROOT, "api", "openapi.public.json"), encoding="utf-8"))
SCHEMAS = SPEC["components"]["schemas"]

INITIALISMS = {"id": "ID", "url": "URL", "api": "API", "http": "HTTP", "csv": "CSV", "ip": "IP", "utc": "UTC", "sla": "SLA", "waba": "WABA", "json": "JSON"}

HEADER = "// Código gerado por internal/gen/genmodels.py a partir do contrato. NÃO EDITE.\n\npackage meshio\n\n"


def camel(name):
    parts = re.split(r"[_\-.\s]+", name)
    out = []
    for p in parts:
        if not p:
            continue
        low = p.lower()
        if low in INITIALISMS:
            out.append(INITIALISMS[low])
        elif re.match(r"^\d", p):
            out.append(p)
        else:
            out.append(p[0].upper() + p[1:])
    s = "".join(out)
    if re.match(r"^\d", s):
        s = "N" + s
    return s


RENAMES = {"ValidationError": "ValidationErrorItem", "EventOut": "Event"}


def type_name(schema_name):
    if schema_name in RENAMES:
        return RENAMES[schema_name]
    if schema_name.endswith("Out"):
        return schema_name[:-3]
    if schema_name.endswith("In"):
        return schema_name[:-2] + "Params"
    return schema_name


def comment(text, indent=""):
    if not text:
        return ""
    lines = []
    for line in str(text).strip().splitlines():
        lines.append(indent + "// " + line.rstrip())
    return "\n".join(lines) + "\n"


def split_nullable(schema):
    """Retorna (schema_sem_null, nullable)."""
    if "anyOf" in schema:
        opts = [o for o in schema["anyOf"] if o.get("type") != "null"]
        nullable = len(opts) != len(schema["anyOf"])
        if len(opts) == 1:
            merged = dict(opts[0])
            for k in ("description", "title"):
                if k in schema and k not in merged:
                    merged[k] = schema[k]
            return merged, nullable
        return {"anyOf": opts}, nullable
    t = schema.get("type")
    if isinstance(t, list):
        non = [x for x in t if x != "null"]
        merged = dict(schema)
        merged["type"] = non[0] if len(non) == 1 else non
        return merged, "null" in t
    return schema, False


class Emitter:
    def __init__(self, resolver=None):
        self.out = []
        self.extra = []  # tipos aninhados
        self.resolver = resolver

    def go_type(self, schema, ctx_name):
        """Tipo Go do schema (sem considerar nullable)."""
        if "$ref" in schema:
            return self.resolver(schema["$ref"], ctx_name)
        if "anyOf" in schema or isinstance(schema.get("type"), list):
            return "any"
        t = schema.get("type")
        if t == "string" or "const" in schema and t is None:
            if schema.get("format") == "date-time":
                return "time.Time"
            return "string"
        if t == "integer":
            return "int64"
        if t == "number":
            return "float64"
        if t == "boolean":
            return "bool"
        if t == "array":
            return "[]" + self.go_field_type(schema.get("items", {}), ctx_name + "Item", required=True)
        if t == "object" or "properties" in schema or "allOf" in schema:
            ap = schema.get("additionalProperties")
            if "properties" in schema or "allOf" in schema:
                self.extra.append((ctx_name, schema))
                return ctx_name
            if isinstance(ap, dict):
                return "map[string]" + self.go_field_type(ap, ctx_name + "Value", required=True)
            return "map[string]any"
        return "any"

    def go_field_type(self, schema, ctx_name, required):
        inner, nullable = split_nullable(schema)
        gt = self.go_type(inner, ctx_name)
        if (nullable or not required) and not gt.startswith(("[]", "map[", "any")):
            if gt == "string" and not nullable:
                return gt
            return "*" + gt
        return gt


def struct_fields(em, schema, owner, is_input):
    lines = []
    props, required = collect_props(schema, em)
    for pname, pschema in props.items():
        req = pname in required
        fname = camel(pname)
        inner, nullable = split_nullable(pschema)
        if is_input:
            gt = em.go_type(inner, owner + fname)
            if req and not nullable:
                ftype, tag = gt, pname
            elif gt in ("string",) and not (req and nullable):
                ftype, tag = gt, pname + ",omitempty"
            elif gt.startswith(("[]", "map[")) or gt == "any":
                ftype, tag = gt, pname + ("" if req else ",omitempty")
            else:
                ftype, tag = "*" + gt, pname + ("" if req else ",omitempty")
        else:
            ftype = em.go_field_type(pschema, owner + fname, required=True if not nullable else False)
            if nullable and not ftype.startswith(("*", "[]", "map[", "any")):
                ftype = "*" + ftype
            tag = pname if req else pname + ",omitempty"
        desc = pschema.get("description") or inner.get("description")
        enum = inner.get("enum")
        if enum:
            vals = ", ".join('"%s"' % e for e in enum if e is not None)
            desc = (desc + "\n" if desc else "") + "Valores: " + vals + "."
        if "const" in inner:
            desc = (desc + "\n" if desc else "") + 'Valor: "%s".' % inner["const"]
        lines.append(comment(desc, "\t") + '\t%s %s `json:"%s"`' % (fname, ftype, tag))
    return lines


def collect_props(schema, em):
    props = {}
    required = set()
    for sub in schema.get("allOf", []):
        sub = em.deref(sub) if "$ref" in sub else sub
        p, r = collect_props(sub, em)
        for k, v in p.items():
            if k in props and "$ref" not in v and "type" not in v and "const" in v:
                v = dict(props[k], **v)
            props[k] = v
        required |= r
    for k, v in schema.get("properties", {}).items():
        props[k] = v
    required |= set(schema.get("required", []))
    return props, required


def emit_struct(em, name, schema, is_input, doc=None):
    fields = "\n".join(struct_fields(em, schema, name, is_input))
    em.out.append(comment(doc or schema.get("description") or ("%s corresponde ao schema do contrato." % name)) + "type %s struct {\n%s\n}\n" % (name, fields))
    while em.extra:
        n, s = em.extra.pop(0)
        emit_struct(em, n, s, is_input)


# ---------------------------------------------------------------- contrato

def contract_models():
    em = Emitter(lambda ref, ctx: type_name(ref.split("/")[-1]))
    em.deref = lambda s: SCHEMAS[s["$ref"].split("/")[-1]]
    for name in sorted(SCHEMAS):
        sch = SCHEMAS[name]
        is_input = name.endswith("In")
        emit_struct(em, type_name(name), sch, is_input, doc="%s corresponde ao schema %s do contrato." % (type_name(name), name))
    body = "\n".join(em.out)
    return HEADER + ('import "time"\n\n' if "time.Time" in body else "") + body


def params_models():
    out = [HEADER, 'import (\n"net/url"\n"strconv"\n"time"\n)\n\n']
    for path, item in SPEC["paths"].items():
        for method, op in item.items():
            qs = [p for p in op.get("parameters", []) if p.get("in") == "query"]
            if not qs:
                continue
            name = camel(op["operationId"]) + "Params"
            fields, enc = [], []
            for p in qs:
                inner, _ = split_nullable(p["schema"])
                fname = camel(p["name"])
                t = inner.get("type")
                desc = p.get("description") or inner.get("description")
                if inner.get("enum"):
                    desc = (desc + "\n" if desc else "") + "Valores: " + ", ".join('"%s"' % e for e in inner["enum"]) + "."
                if inner.get("pattern"):
                    vals = re.findall(r"[a-z_]+", inner["pattern"])
                    desc = (desc + "\n" if desc else "") + "Valores: " + ", ".join('"%s"' % v for v in vals) + "."
                if p.get("required"):
                    desc = (desc + "\n" if desc else "") + "Obrigatório."
                if t == "integer":
                    fields.append(comment(desc, "\t") + "\t%s int" % fname)
                    enc.append('\tif p.%s != 0 {\n\t\tq.Set("%s", strconv.Itoa(p.%s))\n\t}' % (fname, p["name"], fname))
                elif inner.get("format") == "date-time":
                    fields.append(comment(desc, "\t") + "\t%s time.Time" % fname)
                    enc.append('\tif !p.%s.IsZero() {\n\t\tq.Set("%s", p.%s.UTC().Format(time.RFC3339))\n\t}' % (fname, p["name"], fname))
                else:
                    fields.append(comment(desc, "\t") + "\t%s string" % fname)
                    enc.append('\tif p.%s != "" {\n\t\tq.Set("%s", p.%s)\n\t}' % (fname, p["name"], fname))
            out.append("// %s são os filtros de consulta de %s %s (%s).\ntype %s struct {\n%s\n}\n\n" % (name, method.upper(), path, op["operationId"], name, "\n".join(fields)))
            out.append("func (p *%s) query() url.Values {\n\tq := url.Values{}\n\tif p == nil {\n\t\treturn q\n\t}\n%s\n\treturn q\n}\n\n" % (name, "\n".join(enc)))
            has_cursor = any(x["name"] == "cursor" for x in qs)
            if has_cursor:
                out.append("func (p *%s) withCursor(c string) *%s {\n\tcp := %s{}\n\tif p != nil {\n\t\tcp = *p\n\t}\n\tcp.Cursor = c\n\treturn &cp\n}\n\n" % (name, name, name))
                out.append("func (p *%s) cursorValue() string {\n\tif p == nil {\n\t\treturn \"\"\n\t}\n\treturn p.Cursor\n}\n\n" % name)
    return "".join(out)


# ---------------------------------------------------------------- eventos

EV_DIR = os.path.join(ROOT, "api", "events")


def load_event(fname):
    return json.load(open(os.path.join(EV_DIR, "data", fname), encoding="utf-8"))


def events_models():
    files = sorted(f for f in os.listdir(os.path.join(EV_DIR, "data")) if f.endswith(".schema.json"))
    em = Emitter()
    em.resolver = lambda ref, ctx: (_ for _ in ()).throw(Exception("ref inesperado " + ref))
    em.deref = lambda s: load_event(s["$ref"])
    types = []
    for f in files:
        if ".common." in f:
            continue
        etype = f[: -len(".schema.json")]
        name = camel(etype) + "Data"
        sch = load_event(f)
        doc = "%s é o campo data do evento %s.\n%s" % (name, etype, sch.get("description", ""))
        emit_struct(em, name, sch, False, doc=doc)
        types.append((etype, name))
    out = [HEADER, None]
    out.append("// Tipos de evento conhecidos. Novos tipos podem surgir; trate os desconhecidos sem falhar.\nconst (\n")
    for etype, name in types:
        out.append('\tEventType%s = "%s"\n' % (camel(etype), etype))
    out.append(")\n\n")
    out.append("\n".join(em.out))
    out.append("\n// decodeEventData decodifica data no tipo do evento; ok=false para tipos desconhecidos.\nfunc decodeEventData(eventType string, raw json.RawMessage) (v any, ok bool, err error) {\n\tswitch eventType {\n")
    for etype, name in types:
        out.append('\tcase "%s":\n\t\tvar d %s\n\t\terr = json.Unmarshal(raw, &d)\n\t\treturn &d, true, err\n' % (etype, name))
    out.append("\t}\n\treturn nil, false, nil\n}\n")
    if out[1] is None:
        body = "".join(out[2:])
        out[1] = 'import (\n"encoding/json"\n%s)\n\n' % ('"time"\n' if "time.Time" in body else "")
    return "".join(out)


def main():
    for fname, content in (("models_gen.go", contract_models()), ("params_gen.go", params_models()), ("events_gen.go", events_models())):
        with open(os.path.join(ROOT, fname), "w", encoding="utf-8") as fh:
            fh.write(content)
    print("ok", file=sys.stderr)


if __name__ == "__main__":
    main()
