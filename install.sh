VERSION=`curl https://raw.githubusercontent.com/cristianoliveira/ergo/master/.version`

cd /tmp/

echo "Dowloading ergo-$VERSION from repo..."
wget https://github.com/cristianoliveira/ergo/releases/download/$VERSION/ergo-$VERSION-linux-amd64.tar.gz
tar -xf ergo-$VERSION-linux-amd64.tar.gz

echo "It is going to copy the binary ./ergo into /usr/local/bin and may need sudo."
sudo cp ergo /usr/local/bin

echo "Application was successfully installed."
echo "For uninstalling execute: rm /usr/local/bin/ergo"
