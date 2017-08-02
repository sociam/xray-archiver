# RESTful api for accessing the app archive

A RESTful API fetching. See it working on [EDI](https://edi.sociam.org)

Dependencies:

postgresdb

To run locally:

```
go install 
$GOPATH/bin/apiserv

```


##Show App Version

Returns the latest version information of a single app

* **Method URL**

```http
GET /api/apps/:appId/:version
```  
*  **URL Params**

   **Required:**
 
   `appId=[integer]`

    **Optional:**
 
   `version=[string]`

* **Data Params**
  None

* **Success Response:**
  * **Code:** 200 <br />
    **Content:** `{ appid : {{types.AppVersion}}}`
 
* **Error Response:**
  TODO: but but its labourful

* **Sample Call:**

curl /api/apps/500282 -H 'Accept: application/json'

curl /api/apps/500282/playen_US12312 -H 'Accept: application/json'

curl /api/apps/500282/play/en_US/12312 -H 'Accept: application/json'


##Show App Selection

Returns the a selection of app meta data given a start pos and a count

* **Method URL**

```http
GET /api/apps/:amount/:startPos
```  
*  **URL Params**

   **Required:**
 
   `amount=[integer]`

   `startPos=[integer]`

* **Data Params**
  None

* **Success Response:**
  * **Code:** 200 <br />
    **Content:** `{ appid : {{types.AppVersion}}}`
 
* **Error Response:**
  TODO: but but its labourful

* **Sample Call:**

curl /api/apps/10/2 -H 'Accept: application/json'




TODO: need to select by package name and get back more than version metadata...
use previous to find this


##Show App 

Returns a app and all metadata

* **Method URL**

```http
GET /api/apps/:packagename
```  
*  **URL Params**

   **Required:**
 
   `packagename=[string]`

* **Data Params**
  None

* **Success Response:**
  * **Code:** 200 <br />
    **Content:** `{ appid : {{types.storeInfo}}}`
 
* **Error Response:**
  TODO: but but its labourful

* **Sample Call:**

curl /api/apps/com.google.app -H 'Accept: application/json'








/api/apps/<appid>/<version>


 /api/apps/<appid>/<version string>

/api/apps/<amount>/<startPos>

##companies

# url format
/api/companies/<companyId>
/api/companies/<amount>/<startPos>

WISHFUlapi 

/api/apps/categroy/ -> apps
/api/apps/trending/
/api/apps/developer/

