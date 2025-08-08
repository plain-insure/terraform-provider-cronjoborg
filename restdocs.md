REST API
========

Introduction
------------
cron-job.org provides an API which enables users to programatically create, update, delete and view cron jobs.
The API provides a REST-like interface with API-key based authorization and JSON request and response payload.
Given this, the API should be easily usable from virtually any programming language.

Limitations
-----------
In order to prevent abuse of our service, we enforce a daily usage limit. By default, this limit is 100 requests
per day, but can be increased upon request. For sustaining members, a higher limit of 5,000 requests per day applies.

Apart from the daily request limit, individual rate limits might apply depending on the specific API call made.
Those limits are mentioned in the documentation of the specific API method.


Requests
--------

Endpoint
^^^^^^^^
The API endpoint is reachable via HTTPS at::

    https://api.cron-job.org/

Authentication
^^^^^^^^^^^^^^
All requests to the API must be authenticated via an API key supplied in the ``Authorization`` bearer token header.
API keys can be retrieved/generated in the `cron-job.org Console <https://console.cron-job.org>`_ at "Settings".
An example header could look like::

    Authorization: Bearer zaX78aqKJuIH4l4RX6njoqADn77MQNJJ

Access via a specific API key might be restricted to certain IP addresses, if configured in the Console. In this
case, requests from non-allowlisted IP addresses will be rejected with an HTTP error code of ``403``.

.. warning::
    API keys are secrets, just like a password. They allow access to your cron-job.org account via the API
    and should always be treated confidentially. We highly recommend to also enable the IP address restriction
    whenever possible.

Content Type
^^^^^^^^^^^^
In case a request requires a payload, it must be JSON-encoded and the ``Content-Type`` header must be set to
``application/json``. In case the ``Content-Type`` header is missing or contains a different value, the request
payload will be ignored.

HTTP Status Codes
^^^^^^^^^^^^^^^^^
The following status codes can be returned by the API:

================    =======================================================
Status code         Description
================    =======================================================
200                 OK: Request succeeded
400                 Bad request: Invalid request / invalid input data
401                 Unauthorized: Invalid API key
403                 Forbidden: API key cannot be used from this origin
404                 Not found: The requested resource could not be found
409                 Conflict, e.g. because a resource already exists
429                 API key quota, resource quota or rate limit exceeded
500                 Internal server error
================    =======================================================

Examples
^^^^^^^^
Example request via curl:

.. code-block:: console

    curl -X PATCH \
         -H 'Content-Type: application/json' \
         -H 'Authorization: Bearer zaX78aqKJuIH4l4RX6njoqADn77MQNJJ' \
         -d '{"job":{"enabled":true}}' \
         https://api.cron-job.org/jobs/12345

Example request via Python:

.. code-block:: python

    import json
    import requests

    ENDPOINT = 'https://api.cron-job.org'

    headers = {
        'Authorization': 'Bearer zaX78aqKJuIH4l4RX6njoqADn77MQNJJ',
        'Content-Type': 'application/json'
    }
    payload = {
        'job': {
            'enabled': True
        }
    }

    result = requests.patch(ENDPOINT + '/jobs/12345', headers=headers, data=json.dumps(payload))
    print(result.json())


API Methods
-----------

Listing Cron Jobs
^^^^^^^^^^^^^^^^^
List all jobs in this account::

    GET /jobs

**Input Object**

None.

**Output Object**

=================== ======================================= ======================================
Key                 Type                                    Description
=================== ======================================= ======================================
jobs                array of :ref:`Job`                     List of jobs present in the account
someFailed          boolean                                 ``true`` in case some jobs could not be retrieved because of internal errors and the list might be incomplete, otherwise ``false``
=================== ======================================= ======================================

**curl Example**

.. code-block:: console

    curl -X GET \
         -H 'Content-Type: application/json' \
         -H 'Authorization: Bearer zaX78aqKJuIH4l4RX6njoqADn77MQNJJ' \
         https://api.cron-job.org/jobs

**Response Example**

.. code-block:: json

    {
        "jobs": [
            {
                "jobId": 12345,
                "enabled": true,
                "title": "Example Job",
                "saveResponses": false,
                "url": "https:\/\/example.com\/",
                "lastStatus": 0,
                "lastDuration": 0,
                "lastExecution": 0,
                "nextExecution": 1640187240,
                "type": 0,
                "requestTimeout": 300,
                "redirectSuccess": false,
                "folderId": 0,
                "schedule": {
                    "timezone": "Europe/Berlin",
                    "expiresAt": 0,
                    "hours": [
                        -1
                    ],
                    "mdays": [
                        -1
                    ],
                    "minutes": [
                        0,
                        15,
                        30,
                        45
                    ],
                    "months": [
                        -1
                    ],
                    "wdays": [
                        -1
                    ]
                },
                "requestMethod": 0
            }
        ],
        "someFailed": false
    }


**Rate Limit**

Max. 5 requests per second.


Retrieving Cron Job Details
^^^^^^^^^^^^^^^^^^^^^^^^^^^
Retrieve detailed information for a specific cron job identified by its `jobId`::

    GET /jobs/<jobId>

**Input Object**

None.

**Output Object**

=================== ======================================= ======================================
Key                 Type                                    Description
=================== ======================================= ======================================
jobDetails          array of :ref:`DetailedJob`             Job details
=================== ======================================= ======================================

**curl Example**

.. code-block:: console

    curl -X GET \
         -H 'Content-Type: application/json' \
         -H 'Authorization: Bearer zaX78aqKJuIH4l4RX6njoqADn77MQNJJ' \
         https://api.cron-job.org/jobs/12345

**Response Example**

.. code-block:: json

    {
        "jobDetails": {
            "jobId": 12345,
            "enabled": true,
            "title": "Example Job",
            "saveResponses": false,
            "url": "https:\/\/example.com\/",
            "lastStatus": 0,
            "lastDuration": 0,
            "lastExecution": 0,
            "nextExecution": 1640189160,
            "auth": {
                "enable": false,
                "user": "",
                "password": ""
            },
            "notification": {
                "onFailure": false,
                "onSuccess": false,
                "onDisable": false
            },
            "extendedData": {
                "headers": {
                    "X-Foo": "Bar"
                },
                "body": "Hello World!"
            },
            "type": 0,
            "requestTimeout": 300,
            "redirectSuccess": false,
            "folderId": 0,
            "schedule": {
                "timezone": "Europe/Berlin",
                "expiresAt": 0,
                "hours": [
                    -1
                ],
                "mdays": [
                    -1
                ],
                "minutes": [
                    0,
                    15,
                    30,
                    45
                ],
                "months": [
                    -1
                ],
                "wdays": [
                    -1
                ]
            },
            "requestMethod": 0
        }
    }


**Rate Limit**

Max. 5 requests per second.


Creating a Cron Job
^^^^^^^^^^^^^^^^^^^
Creating a new cron job::

    PUT /jobs


**Input Object**

=================== ======================================= ======================================
Key                 Type                                    Description
=================== ======================================= ======================================
job                 :ref:`DetailedJob`                      Job (only the ``url`` field is mandatory)
=================== ======================================= ======================================

**Output Object**

=================== ======================================= ======================================
Key                 Type                                    Description
=================== ======================================= ======================================
jobId               int                                     Identifier of the created job
=================== ======================================= ======================================

**curl Example**

.. code-block:: console

    curl -X PUT \
         -H 'Content-Type: application/json' \
         -H 'Authorization: Bearer zaX78aqKJuIH4l4RX6njoqADn77MQNJJ' \
         -d '{"job":{"url":"https://example.com","enabled":true,"saveResponses":true,"schedule":{"timezone":"Europe/Berlin","expiresAt":0,"hours":[-1],"mdays":[-1],"minutes":[-1],"months":[-1],"wdays":[-1]}}}' \
         https://api.cron-job.org/jobs

**Response Example**

.. code-block:: json

    {
        "jobId": 12345
    }


**Rate Limit**

Max. 1 request per second and 5 requests per minute.


Updating a Cron Job
^^^^^^^^^^^^^^^^^^^
Updating a cron job identified by its `jobId`::

    PATCH /jobs/<jobId>

**Input Object**

=================== ======================================= ======================================
Key                 Type                                    Description
=================== ======================================= ======================================
job                 :ref:`DetailedJob`                      Job delta (only include changed fields - unchanged fields can be left out)
=================== ======================================= ======================================

**Output Object**

Empty object.

**curl Example**

.. code-block:: console

    curl -X PATCH \
         -H 'Content-Type: application/json' \
         -H 'Authorization: Bearer zaX78aqKJuIH4l4RX6njoqADn77MQNJJ' \
         -d '{"job":{"enabled":true}}' \
         https://api.cron-job.org/jobs/12345

**Response Example**

.. code-block:: json

    {}


**Rate Limit**

Max. 5 requests per second.



Deleting a Cron Job
^^^^^^^^^^^^^^^^^^^
Deleting a cron job identified by its `jobId`::

    DELETE /jobs/<jobId>

**Input Object**

None.

**Output Object**

Empty object.

**curl Example**

.. code-block:: console

    curl -X DELETE \
         -H 'Content-Type: application/json' \
         -H 'Authorization: Bearer zaX78aqKJuIH4l4RX6njoqADn77MQNJJ' \
         https://api.cron-job.org/jobs/12345

**Response Example**

.. code-block:: json

    {}


**Rate Limit**

Max. 5 requests per second.


Retrieving the Job Execution History
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
Retrieve the execution history for a specific cron job identified by its `jobId`::

    GET /jobs/<jobId>/history

**Input Object**

None.

**Output Object**

=================== ======================================= ======================================
Key                 Type                                    Description
=================== ======================================= ======================================
history             array of :ref:`HistoryItem`             The last execution history items
predictions         array of int                            Unix timestamps (in seconds) of the predicted next executions (up to 3)
=================== ======================================= ======================================

Please note that the `headers` and `body` fields of the `HistoryItem` objects will not be populated.
In order to retrieve headers and body, see :ref:`Retrieving Job Execution History Item Details`.

**curl Example**

.. code-block:: console

    curl -X GET \
         -H 'Content-Type: application/json' \
         -H 'Authorization: Bearer zaX78aqKJuIH4l4RX6njoqADn77MQNJJ' \
         https://api.cron-job.org/jobs/12345/history

**Response Example**

.. code-block:: json

    {
        "history": [
            {
                "jobLogId": 4946,
                "jobId": 12345,
                "identifier": "12345-22-11-4946",
                "date": 1640189711,
                "datePlanned": 1640189700,
                "jitter": 11257,
                "url": "http:\/\/example.com\/",
                "duration": 239,
                "status": 1,
                "statusText": "OK",
                "httpStatus": 200,
                "headers": null,
                "body": null,
                "stats": {
                    "nameLookup": 1003,
                    "connect": 85516,
                    "appConnect": 0,
                    "preTransfer": 85548,
                    "startTransfer": 238112,
                    "total": 238129
                }
            }
        ],
        "predictions": [
            1640190600,
            1640191500,
            1640192400
        ]
    }


**Rate Limit**

Max. 5 requests per second.


Retrieving Job Execution History Item Details
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
Retrieve details for a specific history item identified by its `identifier` for a specific cron job identified by its `jobId`::

    GET /jobs/<jobId>/history/<identifier>

**Input Object**

None.

**Output Object**

=================== ======================================= ======================================
Key                 Type                                    Description
=================== ======================================= ======================================
jobHistoryDetails   :ref:`HistoryItem`                      History item
=================== ======================================= ======================================

**curl Example**

.. code-block:: console

    curl -X GET \
         -H 'Content-Type: application/json' \
         -H 'Authorization: Bearer zaX78aqKJuIH4l4RX6njoqADn77MQNJJ' \
         https://api.cron-job.org/jobs/12345/history/12345-22-11-4946

**Response Example**

.. code-block:: json

    {
        "jobHistoryDetails": {
            "jobLogId": 4946,
            "jobId": 12345,
            "identifier": "12345-22-11-4946",
            "date": 1640189711,
            "datePlanned": 1640189700,
            "jitter": 11257,
            "url": "http:\/\/example.com\/",
            "duration": 239,
            "status": 1,
            "statusText": "OK",
            "httpStatus": 200,
            "headers": "Accept-Ranges: bytes\r\nCache-Control: max-age=604800\r\nContent-Type: text\/html; charset=UTF-8...\r\n\r\n",
            "body": "<!doctype html>\n<html>\n<head>\n    <title>Example Domain<\/title>...\n",
            "stats": {
                "nameLookup": 1003,
                "connect": 85516,
                "appConnect": 0,
                "preTransfer": 85548,
                "startTransfer": 238112,
                "total": 238129
            }
        }
    }


**Rate Limit**

Max. 5 requests per second