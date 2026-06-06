#!/usr/bin/env bash

ROOT=""

for arg in "$@"; do
  case $arg in
    --root=*)
      ROOT="${arg#*=}"
      ;;
  esac
done

./build/go-train --root="$ROOT"