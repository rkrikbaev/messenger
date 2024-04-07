#!/bin/bash

sudo cp -a production/*.pem /usr/share/pki/ca-trust-source/anchors/

sudo update-ca-trust extract
