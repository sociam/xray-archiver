# XRay - App Observatory API

A RESTful API for interacting with the xray App Observatory Database.The App Observatory is a live collection of Mobile Application Metadata, it contains information relating to the packages used in an application, the permissions that an app might request and most importantly, what a specific application is doing with your personal data.

Through this API you can access information held in the App Observatory. An Example of an application that utilises this API can be seen by browising to the [Ethical Data Initiative (EDI)](https://edi.sociam.org) site.

# How to Use
The API allows you to build queries using the url and a series of parameters. All parameters are optional, should any be ommitted default values are used in their place.

Whilst default values are used, many are actually default to be empty strings. This means that queries are not restricted by default. Details of this can be found down the page.

### Parameters
There are a series of parameters that can be used. Each of these can be chained together to produce a single query to the xray database.

```
isFull: bool
title: string
developer: string
genre: string
appId: string
```
#### Full Details
The ``` isFull ``` parameter is used to toggle between retrieving all appdata that is present for an app or just the stub. This should be set to ``` True ``` or ``` False ```.

Example
```
/api/apps/?isFull=True
```

#### Title
The ``` title ``` parameter is used to search for applications based on the title field.

When searching for applications based on title, some apps may share some of the text that has been entered. As a result of this, the API will return data that relates to any app that contains the string that was provided

Example
```
/api/apps/?title=DinnerAtNoon
```

You can also chain title multiple times to search for apps that match either string

Example
```
/api/apps/?title=dinner&title=at&title=noon
```

The above example would return any apps that contain either ``` dinner ```, ``` at ``` or ``` noon ```.


#### Developer Name

#### Genre Category

#### Application ID


### Example Query

Example URL:
```
/api/apps/?isFull=false&developer=Zynga&title=frozen
```

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

## Dependencies
postgresql
Golang

To run locally:

```
go install 
$GOPATH/bin/apiserv
```


### Future updates 
At the moment, only the base information is included with the API. 
Additional updates are required so that the API will include:
```
permission: string
package: string
hosts: string
```
