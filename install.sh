VERSION="v0.1.1"

cd /tmp

wget https://github.com/cristianoliveira/ergo/releases/download/$VERSION/ergo-$VERSION--linux-amd64.tar.gz
tar -xf ergo-$VERSION-x86_64-unknown-linux-gnu.tar.gz
sudo cp ergo /usr/local/bin

echo "Application was installed in /usr/local/bin. To uninstall just do rm /usr/local/bin/ergo"
