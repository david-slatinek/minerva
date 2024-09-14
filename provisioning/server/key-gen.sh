#!/usr/bin/env bash

function help {
  echo "Generate key pair"
  echo ""
  echo "Usage:"
  echo "./key-gen.sh [-f path] [-c comment] [-n passphrase]"
  echo "-f [path] - path and filename where to store keys. Required."
  echo "-c [comment] - comment. Default: '\$(whoami)'."
  echo "-n [passphrase] - passphrase to encrypt key. Default: ''."
  echo ""
  echo "Examples:"
  echo "./key-gen.sh -f /home/keys/key.pem"
}

comment="$(whoami)"
passphrase=""

while getopts "f:c:n:h" opt; do
  case $opt in
    f)
      filepath="$OPTARG"
      ;;
    c)
      comment="$OPTARG"
      ;;
    n)
      passphrase="$OPTARG"
      ;;
    h)
      help
      exit
      ;;
    \?)
      echo "Invalid option: -$OPTARG"
      help
      exit 1
      ;;
    :)
      echo "Option -$OPTARG requires an argument."
      help
      exit 1
      ;;
  esac
done

if [ -z "$filepath" ]; then
  echo "-f flag is required"
  exit 1
fi

ssh-keygen -t ed25519 -Z aes256-gcm@openssh.com -f "$filepath" -N "$passphrase" -C "$comment" -m pem
