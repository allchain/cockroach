// Copyright 2017 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package types

import (
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/lib/pq/oid"
)

var (
	// Oid is the type of an OID. Can be compared with ==.
	Oid = TOid{oid.T_oid}
	// RegClass is the type of an regclass OID variant. Can be compared with ==.
	RegClass = TOid{oid.T_regclass}
	// RegNamespace is the type of an regnamespace OID variant. Can be compared with ==.
	RegNamespace = TOid{oid.T_regnamespace}
	// RegProc is the type of an regproc OID variant. Can be compared with ==.
	RegProc = TOid{oid.T_regproc}
	// RegProcedure is the type of an regprocedure OID variant. Can be compared with ==.
	RegProcedure = TOid{oid.T_regprocedure}
	// RegType is the type of an regtype OID variant. Can be compared with ==.
	RegType = TOid{oid.T_regtype}

	// Name is a type-alias for String with a different OID. Can be
	// compared with ==.
	Name = WrapTypeWithOid(String, oid.T_name)
	// IntVector is a type-alias for an IntArray with a different OID. Can
	// be compared with ==.
	IntVector = WrapTypeWithOid(TArray{Int}, oid.T_int2vector)
	// OidVector is a type-alias for an OidArray with a different OID. Can
	// be compared with ==.
	OidVector = WrapTypeWithOid(TArray{Oid}, oid.T_oidvector)
	// NameArray is the type family of a DArray containing the Name alias type.
	// Can be compared with ==.
	NameArray T = TArray{Name}
)

var (
	// Unexported wrapper types. These exist for Postgres type compatibility.
	typeInt2   = WrapTypeWithOid(Int, oid.T_int2)
	typeInt4   = WrapTypeWithOid(Int, oid.T_int4)
	typeFloat4 = WrapTypeWithOid(Float, oid.T_float4)
	typeBit    = WrapTypeWithOid(BitArray, oid.T_bit)

	// typeVarChar is the "standard SQL" string type of varying length.
	//
	// It is reported as VARCHAR in SHOW CREATE and "character varying" in
	// introspection for compatibility with PostgreSQL.
	//
	// It has no default maximum length but can be associated with one in the
	// syntax.
	typeVarChar = WrapTypeWithOid(String, oid.T_varchar)

	// typeBpChar is the "standard SQL" string type of fixed length, where "bp"
	// stands for "blank padded".
	//
	// It is reported as CHAR in SHOW CREATE and "character" in introspection for
	// compatibility with PostgreSQL.
	//
	// Its default maximum with is 1. It always has a maximum width.
	typeBpChar = WrapTypeWithOid(String, oid.T_bpchar)

	// typeQChar is a special PostgreSQL-only type supported for compatibility.
	// It behaves like VARCHAR, its maximum width cannot be modified, and has a
	// peculiar name in the syntax and introspection.
	//
	// It is reported as "char" (with double quotes included) in SHOW CREATE and
	// "char" in introspection for compatibility with PostgreSQL.
	typeQChar = WrapTypeWithOid(String, oid.T_char)
)

var semanticTypeToOid = map[SemanticType]oid.Oid{
	BOOL:           oid.T_bool,
	INT:            oid.T_int8,
	FLOAT:          oid.T_float8,
	DECIMAL:        oid.T_numeric,
	DATE:           oid.T_date,
	TIMESTAMP:      oid.T_timestamp,
	INTERVAL:       oid.T_interval,
	STRING:         oid.T_text,
	BYTES:          oid.T_bytea,
	TIMESTAMPTZ:    oid.T_timestamptz,
	COLLATEDSTRING: oid.T_text,
	OID:            oid.T_oid,
	NULL:           oid.T_unknown,
	UUID:           oid.T_uuid,
	ARRAY:          oid.T_anyarray,
	INET:           oid.T_inet,
	TIME:           oid.T_time,
	JSON:           oid.T_jsonb,
	TUPLE:          oid.T_record,
	BIT:            oid.T_bit,
	ANY:            oid.T_any,
}

// OidToType maps Postgres object IDs to CockroachDB types.  We export the map
// instead of a method so that other packages can iterate over the map directly.
// Note that additional elements for the array Oid types are added in init().
var OidToType = map[oid.Oid]T{
	oid.T_anyelement:   Any,
	oid.T_bit:          typeBit,
	oid.T_bool:         Bool,
	oid.T_bpchar:       typeBpChar,
	oid.T_bytea:        Bytes,
	oid.T_char:         typeQChar,
	oid.T_date:         Date,
	oid.T_float4:       typeFloat4,
	oid.T_float8:       Float,
	oid.T_int2:         typeInt2,
	oid.T_int2vector:   IntVector,
	oid.T_int4:         typeInt4,
	oid.T_int8:         Int,
	oid.T_inet:         INet,
	oid.T_interval:     Interval,
	oid.T_jsonb:        Jsonb,
	oid.T_name:         Name,
	oid.T_numeric:      Decimal,
	oid.T_oid:          Oid,
	oid.T_oidvector:    OidVector,
	oid.T_record:       EmptyTuple,
	oid.T_regclass:     RegClass,
	oid.T_regnamespace: RegNamespace,
	oid.T_regproc:      RegProc,
	oid.T_regprocedure: RegProcedure,
	oid.T_regtype:      RegType,
	oid.T_text:         String,
	oid.T_time:         Time,
	oid.T_timestamp:    Timestamp,
	oid.T_timestamptz:  TimestampTZ,
	oid.T_uuid:         Uuid,
	oid.T_varbit:       BitArray,
	oid.T_varchar:      typeVarChar,
}

// oidToArrayOid maps scalar type Oids to their corresponding array type Oid.
var oidToArrayOid = map[oid.Oid]oid.Oid{
	oid.T_anyelement:   oid.T_anyarray,
	oid.T_bit:          oid.T__bit,
	oid.T_bool:         oid.T__bool,
	oid.T_bpchar:       oid.T__bpchar,
	oid.T_bytea:        oid.T__bytea,
	oid.T_char:         oid.T__char,
	oid.T_date:         oid.T__date,
	oid.T_float4:       oid.T__float4,
	oid.T_float8:       oid.T__float8,
	oid.T_inet:         oid.T__inet,
	oid.T_int2:         oid.T__int2,
	oid.T_int2vector:   oid.T__int2vector,
	oid.T_int4:         oid.T__int4,
	oid.T_int8:         oid.T__int8,
	oid.T_interval:     oid.T__interval,
	oid.T_jsonb:        oid.T__jsonb,
	oid.T_name:         oid.T__name,
	oid.T_numeric:      oid.T__numeric,
	oid.T_oid:          oid.T__oid,
	oid.T_oidvector:    oid.T__oidvector,
	oid.T_record:       oid.T__record,
	oid.T_regclass:     oid.T__regclass,
	oid.T_regnamespace: oid.T__regnamespace,
	oid.T_regproc:      oid.T__regproc,
	oid.T_regprocedure: oid.T__regprocedure,
	oid.T_regtype:      oid.T__regtype,
	oid.T_text:         oid.T__text,
	oid.T_time:         oid.T__time,
	oid.T_timestamp:    oid.T__timestamp,
	oid.T_timestamptz:  oid.T__timestamptz,
	oid.T_uuid:         oid.T__uuid,
	oid.T_varbit:       oid.T__varbit,
	oid.T_varchar:      oid.T__varchar,
}

// ArrayOids is a set of all oids which correspond to an array type.
var ArrayOids = map[oid.Oid]struct{}{}

func init() {
	if len(oidToArrayOid) != len(oidToArrayOid) {
		panic("missing some mapping from array element OID to array OID")
	}
	for o, ao := range oidToArrayOid {
		ArrayOids[ao] = struct{}{}
		OidToType[ao] = TArray{OidToType[o]}
	}
}

// TOid represents an alias to the Int type with a different Postgres OID.
type TOid struct {
	oidType oid.Oid
}

func (TOid) SemanticType() SemanticType { return OID }

func (t TOid) String() string { return t.SQLName() }

// Equivalent implements the T interface.
func (t TOid) Equivalent(other T) bool { return isTypeOrAny(other.SemanticType(), OID) }

// Oid implements the T interface.
func (t TOid) Oid() oid.Oid { return t.oidType }

// SQLName implements the T interface.
func (t TOid) SQLName() string {
	switch t.oidType {
	case oid.T_oid:
		return "oid"
	case oid.T_regclass:
		return "regclass"
	case oid.T_regnamespace:
		return "regnamespace"
	case oid.T_regproc:
		return "regproc"
	case oid.T_regprocedure:
		return "regprocedure"
	case oid.T_regtype:
		return "regtype"
	default:
		panic(pgerror.NewAssertionErrorf("unexpected oidType: %v", log.Safe(t.oidType)))
	}
}

// IsAmbiguous implements the T interface.
func (TOid) IsAmbiguous() bool { return false }

// TOidWrapper is a T implementation which is a wrapper around a T, allowing
// custom Oid values to be attached to the T. The T is used by DOidWrapper
// to permit type aliasing with custom Oids without needing to create new typing
// rules or define new Datum types.
type TOidWrapper struct {
	T
	oid oid.Oid
}

var customOidNames = map[oid.Oid]string{
	oid.T_name: "name",
}

func (t TOidWrapper) String() string {
	// Allow custom type names for specific Oids, but default to wrapped String.
	if s, ok := customOidNames[t.oid]; ok {
		return s
	}
	return t.T.String()
}

// Oid implements the T interface.
func (t TOidWrapper) Oid() oid.Oid { return t.oid }

// WrapTypeWithOid wraps a T with a custom Oid.
func WrapTypeWithOid(t T, oid oid.Oid) T {
	switch v := t.(type) {
	case tUnknown, tAny, TOidWrapper:
		panic(pgerror.NewAssertionErrorf("cannot wrap %T with an Oid", v))
	}
	return TOidWrapper{
		T:   t,
		oid: oid,
	}
}

// UnwrapType returns the base T type for a provided type, stripping
// a *TOidWrapper if present. This is useful for cases like type switches,
// where type aliases should be ignored.
func UnwrapType(t T) T {
	if w, ok := t.(TOidWrapper); ok {
		return w.T
	}
	return t
}
