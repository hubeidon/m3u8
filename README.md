# M3u8下载

## 下载地址
[m3u8](http://git.kaidon.cn/attachments/f80462f5-20d6-41fb-81cc-2a5d0a5bcc19)

## 使用方法
在配置文件 conf.yaml 中添加本地路径或者网络路径
```yaml
address:
   - 
    # 网络路径或者m3u8本地文件
    path: https://hw-vod.cdn.huya.com/1048585/1279520613919/23092785/51f1dae7421423e68e800eff52338eb4.m3u8?hyvid=514993415&hyauid=1279520613919&hyroomid=1279520613919&hyratio=4000&hyscence=vod&appid=66&domainid=25&srckey=NjZfMjVfNTE0OTkzNDE1&bitrate=4044&client=115&definition=yuanhua&pid=1279520613919&scene=vod&vid=514993415&u=1685195357&t=100&sv=2211141506
    # [非必填] 当m3u8文件内没有host时,在地址前添加prefix 
    prefix: https://hw-vod.cdn.huya.com/1048585/1279520613919/23092785
    # [非必填] 保存文件名称(默认从网络地址中提取名称 51f1dae7421423e68e800eff52338eb4.m3u8)
    # 最终文件地址 = dir + fname + ext
    fname: huya
```

添加完成后直接运行m3u8,无需指定配置文件路径.

```shell
./m3u8
```