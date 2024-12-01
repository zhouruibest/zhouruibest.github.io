hostnamectl --static set-hostname dev


# 安装配置git


## 安装

```bash
apt update
apt upgrade -y
apt install curl -y
curl -s https://packages.gitlab.com/install/repositories/gitlab/gitlab-ee/script.deb.sh | sudo bash
sudo apt install git -y

```

## 配置

git config --global user.name "zhourui"

git config --global user.email "1572729241@qq.com"

ssh-keygen -t rsa -C "1572729241@qq.com"
