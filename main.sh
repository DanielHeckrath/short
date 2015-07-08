#!/bin/bash

GRPC=""
DEBUG=""
HTTP=""

if [ -z "$GRPC_ADDR" ]; then
	GRPC = "-grpc.addr=$GRPC_ADDR"
fi

if [ -z "$HTTP_ADDR" ]; then
	HTTP = "-http.addr=$HTTP_ADDR"
fi

if [ -z "$DEBUG_ADDR" ]; then
	DEBUG = "-debug.addr=$DEBUG_ADDR"
fi

short $GRPC $HTTP $DEBUG $*
