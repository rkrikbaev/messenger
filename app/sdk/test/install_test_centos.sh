#!/bin/bash

sudo cp -a test/*.pem /usr/share/pki/ca-trust-source/anchors/

sudo update-ca-trust extract
