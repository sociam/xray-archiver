# RESTful api for accessing the app archive

A RESTful API fetching. See it working on [EDI](https://edi.sociam.org)

Dependencies:

postgresdb

To run locally:

```
go install 
$GOPATH/bin/apiserv
```

# How to Use
The API allows you to build queries using the url and a series of parameters

### Params include:
```
isFull: bool
title: string
developer: string
genre: string
appId: string
```
future updates need to include:

```
permission: string
package: string
hosts: string
```

Example URL:
```/api/apps/?isFull=false/developer=Zynga/title=frozen```

The above when used with Curl might return a json object similar to:
```
[
     {
          "Title":"Frozen Free Fall",
          "App":"com.disney.frozensaga_goo"
     },
     {
          "Title":"Candy Crush Jelly Saga",
          "App":"com.king.candycrushjellysaga"
     },
     {
          "Title":"Word Search",
          "App":"com.asgardsoft.words"
     }
]
```
