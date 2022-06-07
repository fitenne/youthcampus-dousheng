# 抖声项目

## 团队介绍

- 队名：**going小分队**
- 队号：**1257836400**
- 队员：**张康寒、陈奇廷、孙世铭、杨可、姚天野、程攻**



## 项目介绍
项目缩略图依赖ffmpeg，文件存储在工作目录下的public文件夹

- 技术栈：**gin gorm mysql**

- 目录结构：
```
DOUSHENG
│  config.yaml      //配置文件
│  dousheng.sql     //建表语句
│
│
├─cmd
│  └─dousheng
│          main.go
│          router.go
│
├─internal
│  ├─common     //公用部分    
│  │  │
│  │  ├─jwt     //token相关
│  │  │
│  │  ├─mid     //中间件
│  │  │
│  │  └─settings        //配置文件相关
│  │
│  ├─controller     //控制层
│  │      comment.go
│  │      common.go
│  │      favorite.go
│  │      feed.go
│  │      publish.go
│  │      relation.go
│  │      user.go
│  │      user_test.go
│  │
│  ├─repository     //数据访问层
│  │      comment.go
│  │      db_provider.go
│  │      favorite.go
│  │      follow.go
│  │      user.go
│  │      video.go
│  │
│  └─service        //服务层
│      │  comment.go
│      │  favorite.go
│      │  feed.go
│      │  follow.go
│      │  user.go
│      │  video.go
│      
│     
│
├─pkg
│  └─model      //模型存储
│          comment.go
│          comment_test.go
│          favorite.go
│          favorite_test.go
│          user.go
│          video.go
│          video_test.go
│
└─public        //存储数据
```




