[![Build Status](https://travis-ci.org/coderlindacheng/lzw.svg?branch=master)](https://travis-ci.org/coderlindacheng/lzw)

# 先来个源码拉取的办法(行内人士用的)

    go get github.com/coderlindacheng/lzw

    下面是完美的分割线
---

# 使用说明

## 先来一波下载连接

* [Windows版本](https://github.com/coderlindacheng/lzw/raw/master/lzw.exe)这个是64位版本 
* [Mac版本](https://github.com/coderlindacheng/lzw/raw/master/lzw) 这个是64位版本
* [Linux版本](https://github.com/coderlindacheng/lzw/raw/master/lzw_linux64) 这个是64位版本
* [评分标准.xlsx](https://github.com/coderlindacheng/lzw/raw/master/%E8%AF%84%E5%88%86%E6%A0%87%E5%87%86.xlsx) 这个是样板
* [原始表.xlsx](https://github.com/coderlindacheng/lzw/raw/master/%E5%8E%9F%E5%A7%8B%E8%A1%A8.xlsx) 这个是样板
* [分数表.xlsx](https://github.com/coderlindacheng/lzw/raw/master/%E5%88%86%E6%95%B0%E8%A1%A8.xlsx) 这个是样板
    
## 这回真的是使用说明了

### 程序执行说明

1. 首先下载对应版本的程序
2. 创建一个叫 评分标准.xlsx 的文件
3. 创建一个叫 原始表.xlsx 的文件
4. 创建一个文件夹 把 1,2,3 提到的文件放到同一个文件夹里面,然后执行程序,就可以了
5. 程序会自动把 原始表.xlsx 的输入,根据 评分标准.xlsx 的规则输出 分数表.xlsx

### 所有表的统一格式
    
1.第一行是表头,不能填数值,而且每个表都有固定的格式
2.评分标准.xlsx 的作用是把原始表的输入
    
### 评分标准.xlsx 应该怎么填
    
1. 先来个样板[样板](https://github.com/coderlindacheng/lzw/raw/master/%E8%AF%84%E5%88%86%E6%A0%87%E5%87%86.xlsx)

2. 第一列的表头必须是"分值"两个字, 下面的值顾名思义就是分数,一定要按照降序(就是上面大下面小)去填,程序会按照这个分数间隔去这算结果的分数,同一个分数就别填两行了,别那么调皮

3. 第二列开始的表头都是以"#"号来分割,程序可以从而读取参数, 格式是 "项目名称#性别(就是"男"或者"女"没有"人妖"哦)#计量单位#数值的排列顺序(默认是升序,就是上面小下面大)",项目名称会对应 原始表.xlsx 的项目名称,性别会对应 原始表.xlsx 的性别,计量单位是用来把人手输入的格式转换成程序计算的需要的格式的,数值的排列顺序也是和程序计算相关

4. 第二列开始就是具体的项目数值,要对应第一列的分数来填数值,程序会以这个为依据输出结果

5. 计量单位现在只对"时间"这个单位做了区别对待,其他可以不填,有默认处理

6. 其实只有项目名称和性别是必须的,其他都有默认处理

### 原始表.xlsx 应该怎么填

1. 也是先来个[样板](https://github.com/coderlindacheng/lzw/raw/master/%E5%8E%9F%E5%A7%8B%E8%A1%A8.xlsx)
2. 反正第一行肯定是表头,需要特别注意的是"性别"这个表头是必须有的,而且必需要出现项目名称表头之前
3. 项目名称表头,就是这个表头名称要对应 评分标准.xlsx 表头格式里面的 项目名称(即#号前的一个名字)
4. 程序会根据项目名称表头,项目名称对应的人手输入的格式和性别表头,以及它的值(也就是男或者女),来输出 分数表.xlsx

### 特别设计的**时间**格式
    
1. 格式就是 分"秒'毫秒x100
2. 毫秒x100 其实就是零点几秒的意思
3. 分钟的值不能超过60,秒的值不能超过60,毫秒x100的值不能超过10
    
### 需要注意的
    
程序的输入第几个表第几行第几列都是从0开始的




