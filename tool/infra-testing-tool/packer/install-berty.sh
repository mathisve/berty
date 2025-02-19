#!/bin/bash

# update system
sudo yum update -y
sudo yum install gcc gcc-c++ make htop wget git gcc -y

# install golang
echo "installing golang"
wget -q https://golang.org/dl/go1.16.4.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.16.4.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version

# install berty
echo "installing berty"
cd ~/ || exit
git clone https://github.com/berty/berty.git
cd berty/go || exit
make install -j 4
sudo mv ~/go/bin/* /usr/local/bin
berty version

# daemon
echo "installing daemon"
cd ~/ || exit
wget -q https://bertytest.s3.eu-west-3.amazonaws.com/infra.zip
unzip -q infra.zip
cd infra-testing-tool/daemon || exit
go build -o infra-daemon .
sudo mv infra-daemon /usr/local/bin

# add daemon to systemd init
sudo mv infra-daemon.service /etc/systemd/system/
sudo systemctl enable infra-daemon.service

echo "creating aws config"
# create config

sudo mkdir /home/ec2-user/.aws
sudo touch /home/ec2-user/.aws/config
sudo echo "[default]" | sudo tee /home/ec2-user/.aws/config > /dev/null
sudo echo "region = $(ec2-metadata --availability-zone | sed 's/placement: \(.*\).$/\1/')" | sudo tee -a /home/ec2-user/.aws/config > /dev/null

cd /home/ec2-user/infra-testing-tool/packer || exit
cat set-ip-vars.sh | sudo tee -a /home/ec2-user/.bashrc > /dev/null

mkdir ~/logs
