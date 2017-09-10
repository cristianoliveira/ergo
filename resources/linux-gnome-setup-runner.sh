#!/bin/sh

finish() {
  # these should probably revert to original settings
  gsettings set org.gnome.system.proxy mode "none"
  gsettings set org.gnome.system.proxy autoconfig-url ''
}
trap finish EXIT

gsettings set org.gnome.system.proxy mode 'auto'
gsettings set org.gnome.system.proxy autoconfig-url http://127.0.0.1:2000/proxy.pac

ergo run
