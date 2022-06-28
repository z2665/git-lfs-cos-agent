# git-lfs-cos-agent
## Principle
see [custom-transfers](https://github.com/git-lfs/git-lfs/blob/main/docs/custom-transfers.md)
## usage
```
#setup conf file
touch conf.toml
echo 'secretID = "xxxx"' >> conf.toml
echo 'secretKey = "xxxx"' >> conf.toml
echo 'bucketName = "xxx-appid"' >> conf.toml
echo 'region = "ap-chengdu"' >> conf.toml
echo 'tmpdir = "./tmp/"' >> conf.toml

#set your git project
git config lfs.standalonetransferagent "cos-lfs"
git config lfs.customtransfer.cos-lfs.path pathto/this/agent
git config lfs.customtransfer.cos-lfs.args pathto/this/conf
```
## thanks to
[git-lfs-rsync-agent](https://github.com/aleb/git-lfs-rsync-agent)