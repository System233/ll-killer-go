#!/bin/bash
set -e
cp -arfLT "$1" $PREFIX/share 
rm -rf "$1"
ln -sfvn $PREFIX/share "$1"