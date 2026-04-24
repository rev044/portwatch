// Package config provides configuration loading and validation for portwatch.
//
// A configuration file is a JSON document. Missing fields fall back to
// DefaultConfig values. Example:
//
//	{
//	  "protocol": "tcp",
//	  "interval": 60000000000,
//	  "format": "json",
//	  "output": "stdout",
//	  "ports": [
//	    {"from": 1,    "to": 1024},
//	    {"from": 8080, "to": 8080}
//	  ]
//	}
//
// interval is expressed in nanoseconds (Go's time.Duration JSON encoding).
package config
