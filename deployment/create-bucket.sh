#!/bin/bash

mc alias set minio http://minio:9000 ROOTNAME CHANGEME123
mc mb minio/avatar
mc anonymous set download minio/avatar
