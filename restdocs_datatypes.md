
Data Types
----------

Job
^^^
The Job object represents a cron job.

=================== ======================================= ===================================================================================================================  =======================================
Key                 Type                                    Description                                                                                                          Default *
=================== ======================================= ===================================================================================================================  =======================================
jobId               int                                     Job identifier (read only; ignored during job creation or update)                                                    (auto-assigned)
enabled             boolean                                 Whether the job is enabled (i.e. being executed) or not                                                              ``false``
title               string                                  Job title                                                                                                            (empty)
saveResponses       boolean                                 Whether to save job response header/body or not                                                                      ``false``
url                 string                                  Job URL                                                                                                              (mandatory)
lastStatus          :ref:`JobStatus`                        Last execution status (read only)                                                                                    ``0`` (Unknown / not executed yet)
lastDuration        int                                     Last execution duration in milliseconds (read only)                                                                  \-
lastExecution       int                                     Unix timestamp of last execution (in seconds; read only)                                                             \-
nextExecution       int                                     Unix timestamp of predicted next execution (in seconds), ``null`` if no prediction available (read only)             \-
type                :ref:`JobType`                          Job type (read only)                                                                                                 ``0`` (Default job)
requestTimeout      int                                     Job timeout in seconds                                                                                               ``-1`` (i.e. use default timeout)
redirectSuccess     boolean                                 Whether to treat 3xx HTTP redirect status codes as success or not                                                    ``false``
folderId            int                                     The identifier of the folder this job resides in                                                                     ``0`` (root folder)
schedule            :ref:`JobSchedule`                      Job schedule                                                                                                         ``{}``
requestMethod       :ref:`RequestMethod`                    HTTP request method                                                                                                  ``0`` (GET)
=================== ======================================= ===================================================================================================================  =======================================

`* Value when field is omitted while creating a job.`

DetailedJob
^^^^^^^^^^^
The DetailedJob object represents a cron job with detailed settings. It consists of all members of the
:ref:`Job` object **plus** the following additional fields.

=================== ======================================= ======================================
Key                 Type                                    Description
=================== ======================================= ======================================
auth                :ref:`JobAuth`                          HTTP authentication settings
notification        :ref:`JobNotificationSettings`          Notification settings
extendedData        :ref:`JobExtendedData`                  Extended request data
=================== ======================================= ======================================

JobAuth
^^^^^^^
The JobAuth object represents HTTP (basic) authentication settings for a job.

=================== ======================================= ====================================================== ===========
Key                 Type                                    Description                                            Default *
=================== ======================================= ====================================================== ===========
enable              boolean                                 Whether to enable HTTP basic authentication or not.    ``false``
user                string                                  HTTP basic auth username                               (empty)
password            string                                  HTTP basic auth password                               (empty)
=================== ======================================= ====================================================== ===========

`* Value when field is omitted while creating a job.`

JobNotificationSettings
^^^^^^^^^^^^^^^^^^^^^^^
The JobNotificationSettings specifies notification settings for a job.

=================== ======================================= ======================================================================================= ===========
Key                 Type                                    Description                                                                             Default *
=================== ======================================= ======================================================================================= ===========
onFailure           boolean                                 Whether to send a notification on job failure or not.                                   ``false``
onSuccess           boolean                                 Whether to send a notification when the job succeeds after a prior failure or not.      ``false``
onDisable           boolean                                 Whether to send a notification when the job has been disabled automatically or not.     ``false``
=================== ======================================= ======================================================================================= ===========

`* Value when field is omitted while creating a job.`

JobExtendedData
^^^^^^^^^^^^^^^
The JobExtendedData holds extended request data for a job.

=================== ======================================= ======================================= ===========
Key                 Type                                    Description                             Default *
=================== ======================================= ======================================= ===========
headers             dictionary                              Request headers (key-value dictionary)  ``{}``
body                string                                  Request body data                       (empty)
=================== ======================================= ======================================= ===========

`* Value when field is omitted while creating a job.`

JobStatus
^^^^^^^^^
=================== =========================================================
Value               Description
=================== =========================================================
0                   Unknown / not executed yet
1                   OK
2                   Failed (DNS error)
3                   Failed (could not connect to host)
4                   Failed (HTTP error)
5                   Failed (timeout)
6                   Failed (too much response data)
7                   Failed (invalid URL)
8                   Failed (internal errors)
9                   Failed (unknown reason)
=================== =========================================================

JobType
^^^^^^^
=================== =========================================================
Value               Description
=================== =========================================================
0                   Default job
1                   Monitoring job (used in a status monitor)
=================== =========================================================

JobSchedule
^^^^^^^^^^^
The JobSchedule object represents the execution schedule of a job.

=================== ======================================= ============================================================================================================================================================= =======================================
Key                 Type                                    Description                                                                                                                                                   Default *
=================== ======================================= ============================================================================================================================================================= =======================================
timezone            string                                  Schedule time zone (see `here <https://www.php.net/manual/timezones.php>`_ for a list of supported values)                                                    ``UTC``
expiresAt           int                                     Date/time (in job's time zone) after which the job expires, i.e. after which it is not scheduled anymore (format: `YYYYMMDDhhmmss`, `0` = does not expire)    ``0``
hours               array of int                            Hours in which to execute the job (0-23; `[-1]` = every hour)                                                                                                 ``[]``
mdays               array of int                            Days of month in which to execute the job (1-31; `[-1]` = every day of month)                                                                                 ``[]``
minutes             array of int                            Minutes in which to execute the job (0-59; `[-1]` = every minute)                                                                                             ``[]``
months              array of int                            Months in which to execute the job (1-12; `[-1]` = every month)                                                                                               ``[]``
wdays               array of int                            Days of week in which to execute the job (0=Sunday - 6=Saturday; `[-1]` = every day of week)                                                                  ``[]``
=================== ======================================= ============================================================================================================================================================= =======================================

`* Value when field is omitted while creating a job.`

RequestMethod
^^^^^^^^^^^^^
=================== =========================================================
Value               Description
=================== =========================================================
0                   GET
1                   POST
2                   OPTIONS
3                   HEAD
4                   PUT
5                   DELETE
6                   TRACE
7                   CONNECT
8                   PATCH
=================== =========================================================

HistoryItem
^^^^^^^^^^^
The HistoryItem object represents a job history log entry corresponding to one execution of the job.

=================== ======================================= ======================================
Key                 Type                                    Description
=================== ======================================= ======================================
jobId               int                                     Identifier of the associated cron job
identifier          string                                  Identifier of the history item
date                int                                     Unix timestamp (in seconds) of the actual execution
datePlanned         int                                     Unix timestamp (in seconds) of the planned/ideal execution
jitter              int                                     Scheduling jitter in milliseconds
url                 string                                  Job URL at time of execution
duration            int                                     Actual job duration in milliseconds
status              :ref:`JobStatus`                        Status of execution
statusText          string                                  Detailed job status Description
httpStatus          int                                     HTTP status code returned by the host, if any
headers             string or ``null``                      Raw response headers returned by the host (``null`` if unavailable)
body                string or ``null``                      Raw response body returned by the host (``null`` if unavailable)
stats               :ref:`HistoryItemStats`                 Additional timing information for this request
=================== ======================================= ======================================

HistoryItemStats
^^^^^^^^^^^^^^^^
The HistoryItemStats object contains additional timing information for a job execution history item.

=================== ======================================= ======================================
Key                 Type                                    Description
=================== ======================================= ======================================
nameLookup          int                                     Time from transfer start until name lookups completed (in microseconds)
connect             int                                     Time from transfer start until socket connect completed (in microseconds)
appConnect          int                                     Time from transfer start until SSL handshake completed (n microseconds) - ``0`` if not using SSL
preTransfer         int                                     Time from transfer start until beginning of data transfer (in microseconds)
startTransfer       int                                     Time from transfer start until the first response byte is received (in microseconds)
total               int                                     Total transfer time (in microseconds)
=================== ======================================= ======================================