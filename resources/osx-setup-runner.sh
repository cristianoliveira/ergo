#!/bin/sh

finish() {
  for s in `networksetup -listallnetworkservices | sed -e '1d'`; do
    sudo networksetup -setautoproxyurl "$s" ""
  done
}
trap finish EXIT

for s in `networksetup -listallnetworkservices | sed -e '1d'`; do
  sudo networksetup -setautoproxyurl "$s" "http://localhost:2000/proxy.pac"
done

ergo run
