# mgorm

## Getting Start

``` Bash
$ go get -u github.com/siskinc/mgorm
```

## Example
[example file](./example)

### Set Default Infomation
> In a general way, we will only use one mongodb database, so we can set default infomation, some like mongodb host, database name, database auth infomation(username and password) and the time out second. So if you set these infomation, this package will generate default Session and Database Object.

Well, you can write the code by this way:
``` Golang
err := mgorm.DefaultMgoInfo(
    "127.0.0.1:27017",
    "database_name",
    "username",
    "password",
    30, //  CollectionTimeoutSecond
)
```
Cool, we havd set the default infomation of mongodb.

Now, we should study how to get a MongoDBClient for our Business Logic Layer.

In the above, we have set the default infomation of mongodb, we could just one parameter can get our MongoDBClient Object.
``` Golang
client := mgorm.DefaultMongoDBClient("name")
```
Well, we get the client object, let's curd! see the file [example1.go](./example/example1.go)