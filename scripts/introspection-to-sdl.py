#!/usr/bin/env python3
"""
Convert a GraphQL introspection JSON response (from stdin) to SDL format (stdout).

Usage:
    curl -s -X POST <endpoint> -H "Content-Type: application/json" \
        -d '{"query":"<introspection query>"}' | python3 scripts/introspection-to-sdl.py
"""

import json
import sys

BUILTIN_SCALARS = {"Boolean", "String", "Int", "Float", "ID"}
SKIP_PREFIX = ("__",)


def type_ref(t):
    if t is None:
        return "String"
    if t["kind"] == "NON_NULL":
        return type_ref(t["ofType"]) + "!"
    if t["kind"] == "LIST":
        return "[" + type_ref(t["ofType"]) + "]"
    return t["name"]


def field_args(args):
    if not args:
        return ""
    parts = []
    for a in args:
        default = f" = {a['defaultValue']}" if a.get("defaultValue") else ""
        parts.append(f"{a['name']}: {type_ref(a['type'])}{default}")
    return "(" + ", ".join(parts) + ")"


def skip(name):
    return any(name.startswith(p) for p in SKIP_PREFIX)


data = json.load(sys.stdin)
schema = data["data"]["__schema"]
types_by_name = {t["name"]: t for t in schema["types"]}

lines = []


def emit(s=""):
    lines.append(s)


# Custom scalars
for t in schema["types"]:
    if skip(t["name"]) or t["kind"] != "SCALAR" or t["name"] in BUILTIN_SCALARS:
        continue
    emit(f"scalar {t['name']}")
    emit()

# Enums
for t in schema["types"]:
    if skip(t["name"]) or t["kind"] != "ENUM":
        continue
    if t.get("description"):
        emit('"""')
        emit(t["description"])
        emit('"""')
    emit(f"enum {t['name']} {{")
    for v in t["enumValues"] or []:
        emit(f"  {v['name']}")
    emit("}")
    emit()

# Input objects
for t in schema["types"]:
    if skip(t["name"]) or t["kind"] != "INPUT_OBJECT":
        continue
    emit(f"input {t['name']} {{")
    for f in t["inputFields"] or []:
        default = f" = {f['defaultValue']}" if f.get("defaultValue") else ""
        emit(f"  {f['name']}: {type_ref(f['type'])}{default}")
    emit("}")
    emit()

# Interfaces
for t in schema["types"]:
    if skip(t["name"]) or t["kind"] != "INTERFACE":
        continue
    emit(f"interface {t['name']} {{")
    for f in t["fields"] or []:
        emit(f"  {f['name']}{field_args(f['args'])}: {type_ref(f['type'])}")
    emit("}")
    emit()

# Unions
for t in schema["types"]:
    if skip(t["name"]) or t["kind"] != "UNION":
        continue
    members = " | ".join(p["name"] for p in t["possibleTypes"] or [])
    emit(f"union {t['name']} = {members}")
    emit()

# Object types (non-root)
root_names = {schema["queryType"]["name"]}
if schema.get("mutationType"):
    root_names.add(schema["mutationType"]["name"])
if schema.get("subscriptionType"):
    root_names.add(schema["subscriptionType"]["name"])

for t in schema["types"]:
    if skip(t["name"]) or t["kind"] != "OBJECT" or t["name"] in root_names:
        continue
    interfaces = " & ".join(i["name"] for i in t["interfaces"] or [])
    iface_str = f" implements {interfaces}" if interfaces else ""
    emit(f"type {t['name']}{iface_str} {{")
    for f in t["fields"] or []:
        emit(f"  {f['name']}{field_args(f['args'])}: {type_ref(f['type'])}")
    emit("}")
    emit()

# Query type
qt = types_by_name[schema["queryType"]["name"]]
emit("type Query {")
for f in qt["fields"] or []:
    emit(f"  {f['name']}{field_args(f['args'])}: {type_ref(f['type'])}")
emit("}")
emit()

# Mutation type
if schema.get("mutationType"):
    mt = types_by_name[schema["mutationType"]["name"]]
    emit("type Mutation {")
    for f in mt["fields"] or []:
        emit(f"  {f['name']}{field_args(f['args'])}: {type_ref(f['type'])}")
    emit("}")
    emit()

# Subscription type
if schema.get("subscriptionType"):
    st = types_by_name[schema["subscriptionType"]["name"]]
    emit("type Subscription {")
    for f in st["fields"] or []:
        emit(f"  {f['name']}{field_args(f['args'])}: {type_ref(f['type'])}")
    emit("}")
    emit()

print("\n".join(lines))
