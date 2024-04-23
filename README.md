### PROG2005 Assignment 2

Main task of  this group assignment, is to develop a REST web application in Golang that provides the client with the ability to configure information dashboards that are dynamically populated when requested. The dashboard configurations are saved in your service in a persistent way, and populated based on external services. It will also include a simple notification service that can listen to specific events. The application will be dockerized and deployed using an IaaS system.

The services you will be using for this purpose are:


REST Countries API (instance hosted for this course)

Endpoint: http://129.241.150.113:8080/v3.1

Documentation: http://129.241.150.113:8080/




Open-Meteo APIs (hosted externally, hence please be responsible)

Documentation: https://open-meteo.com/en/features#available-apis

Hint: Take some time to explore the APIs and see which one you can use to solve this assignment.



Currency API

Endpoint: http://129.241.150.113:9090/currency/

Documentation: http://129.241.150.113:9090/



Main tasks:

_NOTE: All examples are shown via local deployment._


## 1. Dashboard configuration 

The initial endpoint (.../dashboard/v1/registrations/) focuses on the management of dashboard configurations that can later be used via the dashboards endpoint. Through different methods dashboard can be registered(POST), retrieved(GET), updated(PUT) and deleted(DELETE).


**Request (POST):**

For given endpoint input http://localhost:8080/dashboard/v1/registrations, using method POST, we are populating dashbord configuration. Example:

{
   "country": "Norway",
   "isoCode": "NO",
   "features": {
                  "temperature": true,
                  "precipitation": true,
                  "capital": true,
                  "coordinates": true,
                  "population": true,
                  "area": true,
                  "targetCurrencies": ["EUR","USD"]
               }
}

As response we get ID of dashboard registration and time when registration is created:

{
    "id": "1",
    "time": "2024-04-22 01:39:53"
}

**Request (GET) all registered dashboard:**

Enables return all registered configurations including IDs and timestamps of last change.

For given endpoint input http://localhost:8080/dashboard/v1/registrations, method GET, results are:

[
    {
        "id": "1",
        "country": "Norway",
        "isoCode": "NO",
        "features": {
            "temperature": true,
            "precipitation": false,
            "capital": false,
            "coordinates": false,
            "population": true,
            "area": true,
            "targetCurrencies": [
                "EUR",
                "USD",
                "NOK"
            ]
        },
        "lastchange": "2024-04-22 01:50:43"
    },
    {
        "id": "2",
        "country": "Finland",
        "isoCode": "FI",
        "features": {
            "temperature": true,
            "precipitation": true,
            "capital": true,
            "coordinates": true,
            "population": true,
            "area": true,
            "targetCurrencies": [
                "EUR",
                "USD"
            ]
        },
        "lastchange": "2024-04-22 01:57:42"
    }
]


**Request (GET) specific registered dashboard:**

Enables retrieval of a specific registered dashboard configuration.

For given endpoint input http://localhost:8080/dashboard/v1/registrations/1, method GET, result is:

{
    "id": "1",
    "country": "Norway",
    "isoCode": "NO",
    "features": {
        "temperature": true,
        "precipitation": true,
        "capital": true,
        "coordinates": true,
        "population": true,
        "area": true,
        "targetCurrencies": [
            "EUR",
            "USD"
        ]
    },
    "lastchange": "2024-04-22 01:39:53"
}

**Request (PUT):**

Enables updating of individual configuration identified by its ID. 

For given endpoint http://localhost:8080/dashboard/v1/registrations/1, method PUT, we can change parameters previously defined through registration. Example:

{
    "country": "Norway",
    "isoCode": "NO",
    "features": {
                  "temperature": true,
                  "precipitation": false,
                  "capital": false,
                  "coordinates": false,
                  "population": true,
                  "area": true,
                  "targetCurrencies": ["EUR", "USD"]
                 }
}

Response:

	Configuration with ID 1 is updated. Dashboard registry updated...


Now, for given endpoint input http://localhost:8080/dashboard/v1/registrations/1, method GET, result is:

{
    "id": "1",
    "country": "Norway",
    "isoCode": "NO",
    "features": {
        "temperature": true,
        "precipitation": false,
        "capital": false,
        "coordinates": false,
        "population": true,
        "area": true,
        "targetCurrencies": [
            "EUR",
            "USD",
            "NOK"
        ]
    },
    "lastchange": "2024-04-22 01:50:43"
}

**Request (DELETE):**

Enables deletion of an individual configuration identified by its ID. This update should lead to a deletion of the configuration on the server.

For given endpoint input http://localhost:8080/dashboard/v1/registrations/1, method DELETE, response is:

	Configuration with ID 1 is removed. Dashboard registry updated...


## 2. Retrieve populated dashboard

Enables retrieving the populated dashboards.

For given endpoint input http://localhost:8080/dashboard/v1/dashboards/1, method GET, result is:

{
    "country": "Norway",
    "isoCode": "NO",
    "features": {
        "temperature": 2.4,
        "precipitation": 0,
        "capital": "Oslo",
        "coordinates": {
            "latitude": 62,
            "langitude": 10
        },
        "population": 5379475,
        "area": 323802,
        "targetcurrencies": {
            "EUR": 0.085114,
            "USD": 0.090739
        }
    },
    "lastchange": "2024-04-22 02:07:01"
}

_NOTE: For this case we have populated all fields with true, showcase example;_


## 3. Managing webhooks for event notifications

IMPORTANT: In order to check webhook functionality when dasboard is created, invoked, changed or deleted, webhook notification MUST be registered first!

As an additional feature, users can register webhooks that are triggered by the service based on specified events, specifically if a new configuration is created, changed or deleted. Users can also register for invocation events, i.e., when a dashboard for a given country is invoked. Users can register multiple webhooks. The registrations should survive a service restart (i.e., be persistently stored).

**Request (POST)**

Enables registering notification, webhook for specific client (in our case we have used webhook.site for generating different URLs to simulate clients), country and invocation type on different events:

_REGISTER_ - webhook is invoked if a new configuration is registered

_CHANGE_ - webhook is invoked if configuration is modified

_DELETE_ - webhook is invoked if configuration is deleted

_INVOKE_ - webhook is invoked if dashboard is retrieved (i.e., populated with values)



For given endpoint input http://localhost:8080/dashboard/v1/notifications, using method POST, we are populating dashbord configuration. Example:

{
    "url": "https://webhook.site/2a2c3e89-8039-4de7-aed9-bf6411ee7b31",
    "country": "NO",
    "event": "REGISTER"
}

Response (Randomly generated 13 character long string): 

{
    "id": "UgzyYGjuameEs"
}

**Request (DELETE)**

Enables deleting specific webhook.

For given endpoint input http://localhost:8080/dashboard/v1/notifications/UgzyYGjuameEs, using method DELETE, response is:

	Notification with ID ' UgzyYGjuameEs ' is removed. Notification registry updated...


**Request (GET) specific registered notification:**

Enables retrieval of a specific registered notification, webhook.

For given endpoint input http://localhost:8080/dashboard/v1/notifications/UgzyYGjuameEs, using method GET, result is:

{
    "id": "UgzyYGjuameEs",
    "url": "https://webhook.site/40401d38-8fbb-49f6-9a04-d7ac025341ab",
    "country": "NO",
    "event": "REGISTER"
}

**Request (GET) all registered notifications:**

Enables retrieval of all registered notifications, webhooks. In this case all events for specific client.

For given endpoint input http://localhost:8080/dashboard/v1/notifications, using method GET, results are:

[
    {
    "id": "UgzyYGjuameEs",
    "url": "https://webhook.site/40401d38-8fbb-49f6-9a04-d7ac025341ab",
    "country": "NO",
    "event": "REGISTER"
    },
    {
    "id": "UgzyYGjuameEs",
    "url": "https://webhook.site/40401d38-8fbb-49f6-9a04-d7ac025341ab",
    "country": "NO",
    "event": "INVOKE"
    },
    {
    "id": "UgzyYGjuameEs",
    "url": "https://webhook.site/40401d38-8fbb-49f6-9a04-d7ac025341ab",
    "country": "NO",
    "event": "CHANGE"
    },
    {
    "id": "UgzyYGjuameEs",
    "url": "https://webhook.site/40401d38-8fbb-49f6-9a04-d7ac025341ab",
    "country": "NO",
    "event": "DELETE"
    }
]

_NOTE: Webhook Invocation (upon trigger)_

When a webhook is triggered, it should send information as follows. Where multiple webhooks are triggered, the information should be sent separately (i.e., one notification per triggered webhook). 

On webhook.site for given url, in this case https://webhook.site/40401d38-8fbb-49f6-9a04-d7ac025341ab (if nottifications with events are registered first), should automatically appear notifications, accordingly to evnet which triger them. Example:

    {
    "id": "ekNapaWcoBkJQ",
    "country": "NO",
    "event": "REGISTER",
    "time": "2024-04-22 02:07:01"
    }

In addition all notifications, with informations about client URL, event upon are triggerd, country for which they are triggered and time when they are created are stored in Firebase/Firestore database system (https://console.firebase.google.com/u/0/project/assignment2-8c8dd/firestore/databases/-default-/data/~2FNotifications~2F1TU6TP508SjAUndl0Ze6, collection:Notifications) in matter that upon server restart (localy and OpenStack/Docker deployment) notifications are not lost. They are retrieved with each server start/restart. 

## 4. Monitoring service availability

The status interface indicates the availability of all individual services this service depends on. 

For given endpoint input http://localhost:8080/dashboard/v1/status, result is:

{
    "countries_api": 200,
    "meteo_api": 200,
    "currency_api": 200,
    "webhooks": 4,
    "version": "v1",
    "uptime": 114.3694919
}


## Deployment
The service is to be deployed on an IaaS solution OpenStack using Docker. URL of deployed service is 10.212.171.38:80.
Accordingly to that we will have five resource root paths:

_Default handler info page:_ 10.212.171.38:80/

_Registration of dashboard:_ 10.212.171.38:80/dashboard/v1/registrations/

_Invoking populated dashboard:_ 10.212.171.38:80/dashboard/v1/dashboards/

_Registration of webhooks:_ 10.212.171.38:80/dashboard/v1/notifications/

_Status o all individual services:_ 10.212.171.38:80/dashboard/v1/status/

All endpoints will have same functionality as those above, deployed localy.


## Authors and acknowledgment
Code written by Nemanja Tosic


