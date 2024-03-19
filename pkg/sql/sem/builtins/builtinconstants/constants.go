// Copyright 2022 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package builtinconstants

import "time"

// SequenceNameArg represents the name of sequence (string) arguments in
// builtin functions.
// Namely, it exists to classify overloads of functions which in postgres
// only take `REGCLASS`, but in cockroach db take both `REGCLASS` and
// `STRING` because we need to be backwards compatible.
const SequenceNameArg = "sequence_name"

// DefaultFollowerReadDuration represents the default time span back from the
// statement time which we wish to be recent and old enough for a follower read.
// Such a default will be returned if we do *not* have an enterprise license
// on a CCL distribution, which may not result in reading from the nearest replica.
const DefaultFollowerReadDuration = -4800 * time.Millisecond

// MaxAllocatedStringSize represents the maximum allowed string length
// in various string related builtin function.
const MaxAllocatedStringSize = 128 * 1024 * 1024

// ErrInsufficientArgsFmtString represents illegal or unknown argument(s) to
// builtin functions.
const ErrInsufficientArgsFmtString = "unknown signature: %s()"

// The following constants are used to categorize builtin functions
// for documentation.
const (
	CategoryArray               = "Array"
	CategoryComparison          = "Comparison"
	CategoryCompatibility       = "Compatibility"
	CategoryCast                = "Cast"
	CategoryCrypto              = "Cryptographic"
	CategoryDateAndTime         = "Date and time"
	CategoryEnum                = "Enum"
	CategoryFullTextSearch      = "Full Text Search"
	CategoryGenerator           = "Set-returning"
	CategoryTrigram             = "Trigrams"
	CategoryFuzzyStringMatching = "Fuzzy String Matching"
	CategoryIDGeneration        = "ID generation"
	CategoryJSON                = "JSONB"
	CategoryMultiRegion         = "Multi-region"
	CategoryMultiTenancy        = "Multi-tenancy"
	CategorySequences           = "Sequence"
	CategorySpatial             = "Spatial"
	CategoryString              = "String and byte"
	CategorySystemInfo          = "System info"
	CategorySystemRepair        = "System repair"
	CategoryStreamIngestion     = "Stream Ingestion"
	CategoryTesting             = "Testing"
	CategoryMigrations          = "Migrations"
)

const (
	// GatewayRegionBuiltinName is the name for the builtin that returns the gateway
	// region of the current node.
	GatewayRegionBuiltinName = "gateway_region"
	// DefaultToDatabasePrimaryRegionBuiltinName is the name for the builtin that
	// takes in a region and returns it if it is a valid region on the database.
	// Otherwise, it returns the primary region.
	DefaultToDatabasePrimaryRegionBuiltinName = "default_to_database_primary_region"
	// RehomeRowBuiltinName is the name for the builtin that rehomes a row to the
	// user's gateway region, defaulting to the database primary region.
	RehomeRowBuiltinName = "rehome_row"
	// CreateSchemaTelemetryJobBuiltinName is the name for the builtin that
	// creates a job that logs SQL schema telemetry.
	CreateSchemaTelemetryJobBuiltinName = "crdb_internal.create_sql_schema_telemetry_job"
)

// A unique int generated by GenerateUniqueInt is a 64-bit integer with
// the following format:
//
//	[1 leading zero bit][48 bits for timestamp][15 bits for nodeID]
const (
	// UniqueIntLeadingZeroBits is the number of leading zero bits in a unique
	// int generated by GenerateUniqueInt.
	UniqueIntLeadingZeroBits = 1

	// UniqueIntTimestampBits is the number of bits in the timestamp segment
	// in a unique int generated by GenerateUniqueInt.
	UniqueIntTimestampBits = 48

	// UniqueIntNodeIDBits is the number of bits in the node ID segment
	// in a unique int generated by GenerateUniqueInt.
	UniqueIntNodeIDBits = 15

	// UniqueIntNodeIDMask is a bitmask for the node ID in a unique int
	// generated by GenerateUniqueInt.
	UniqueIntNodeIDMask = 1<<UniqueIntNodeIDBits - 1

	// UniqueIntTimestampMask is a bitmask for the timestamp in a unique int
	// generated by GenerateUniqueInt.
	UniqueIntTimestampMask = (1<<UniqueIntTimestampBits - 1) << UniqueIntNodeIDBits
)
