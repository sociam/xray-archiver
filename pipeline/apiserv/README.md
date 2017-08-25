# XRay - App Observatory API

A RESTful API for interacting with the xray App Observatory Database.The App Observatory is a live collection of Mobile Application Metadata, it contains information relating to the packages used in an application, the permissions that an app might request and most importantly, what a specific application is doing with your personal data.

Through this API you can access information held in the App Observatory. An Example of an application that utilises this API can be seen by browising to the [Ethical Data Initiative (EDI)](https://edi.sociam.org) site.

# How to Use
The API allows you to build queries using the url and a series of parameters. All parameters are optional, should any be ommitted default values are used in their place.

Whilst default values are used, many are actually default to be empty strings. This means that queries are not restricted by default. Details of this can be found down the page.

## Parameters
There are a series of parameters that can be used. Each of these can be chained together to produce a single query to the xray database.

```
isFull: bool
title: string
developer: string
genre: string
appId: string
```
### Full Details
The ``` isFull ``` parameter is used to toggle between retrieving all appdata that is present for an app or just the stub. This should be set to ``` True ``` or ``` False ```.

**Example**
```
/api/apps?isFull=True
```

### Title
The ``` title ``` parameter is used to search for applications based on the title field.

When searching for applications based on title, some apps may share some of the text that has been entered. As a result of this, the API will return data that relates to any app that contains the string that was provided

**Example**
```
/api/apps?title=DinnerAtNoon
```

You can also chain title multiple times to search for apps that match either string

**Example**
```
/api/apps?title=dinner&title=at&title=noon
```

The above example would return any apps that contain either ``` dinner ```, ``` at ``` or ``` noon ```.


### Developer Name

The Developer name parameter will allow you to retrieve app data that relates to specific developers.

As with the title field, developers may the term that you pass to the API. Because of this, the API will retireve all apps that contain the string that you provide.

**Example**
```
/api/apps?developer=zynga
```

You can also chain developers together to bring back app data that contain either of the strings that are provided

**Example**
```
/api/apps?developer=zynga&developer=google
```

This would return back app data relating to apps that have a developer name that contains either: ```zynga``` or ```google```.

This can be used to allow you to compare applications that have been produced varying developers.


### Genre Category

The Genre parameter will allw you to retieve metadata for applications that are part of a perticular genre. A list of all the possible Genres can be found further down the page.

**Example**
```
/api/apps?genre=entertainment
```

You can also chain the parameter with more genres. This would allow you to compare app data found in different genres easily.

**Example**
```
/api/apps?genre=lifestyle&genre=food_and_drink
```

As some genres actually share words, putting part of a genre will return meta data for apps, whos genre's contain that string.

**Example**
```
/api/apps?genre=game
```

Requesting meta data where the genre contains the string ```game``` will actually bring back apps that belong to one of the game genres.

**The genres included**
```
 GAME_ADVENTURE
 GAME_BOARD
 GAME_SPORTS
 GAME_CASUAL
 GAME_TRIVIA
 GAME_ROLE_PLAYING
 GAME_ACTION
 GAME_PUZZLE
 GAME_CARD
 GAME_SIMULATION
 GAME_ARCADE
 GAME_STRATEGY
 GAME_WORD
 GAME_RACING
 GAME_EDUCATIONAL
```
Chaining together 15 ```genre=``` is also possible, however this would be much easier to do should you want all metadata for games in the database

#### Application ID

### Example Query

**Example URL**
```
/api/apps?isFull=false&developer=Zynga&title=frozen
```

The above when used with Curl might return JSON similar to:
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

### Genre Types

*TODO: format better / create wiki page for these*
```
 PRODUCTIVITY
 SHOPPING
 FINANCE
 PARENTING
 MUSIC_AND_AUDIO
 PERSONALIZATION
 BUSINESS
 GAME_CARD
 GAME_SIMULATION
 GAME_ARCADE
 ENTERTAINMENT
 GAME_RACING
 LIFESTYLE
 GAME_ADVENTURE
 VIDEO_PLAYERS
 GAME_BOARD
 TOOLS
 GAME_SPORTS
 GAME_CASUAL
 GAME_TRIVIA
 COMMUNICATION
 LIBRARIES_AND_DEMO
 FOOD_AND_DRINK
 BOOKS_AND_REFERENCE
 HOUSE_AND_HOME
 GAME_ROLE_PLAYING
 ART_AND_DESIGN
 HEALTH_AND_FITNESS
 EDUCATION
 TRAVEL_AND_LOCAL
 GAME_ACTION
 PHOTOGRAPHY
 GAME_PUZZLE
 GAME_STRATEGY
 GAME_WORD
 SOCIAL
 NEWS_AND_MAGAZINES
 GAME_EDUCATIONAL
```
