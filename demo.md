
# API DOC

## TOC

- [/config/addTask](#/config/addTask)
  - [github.com/mh-cbon/backup/api.TestConfigAddTask](#github.com/mh-cbon/backup/api.TestConfigAddTask)
  - [github.com/mh-cbon/backup/api.TestConfigAddTaskAlreadyExisting](#github.com/mh-cbon/backup/api.TestConfigAddTaskAlreadyExisting)


- [/config/read](#/config/read)
  - [github.com/mh-cbon/backup/api.TestConfigRead](#github.com/mh-cbon/backup/api.TestConfigRead)


- [/config/rmTask](#/config/rmTask)
  - [github.com/mh-cbon/backup/api.TestConfigRmTask](#github.com/mh-cbon/backup/api.TestConfigRmTask)
  - [github.com/mh-cbon/backup/api.TestConfigRmUnknownTask](#github.com/mh-cbon/backup/api.TestConfigRmUnknownTask)


- [/config/tasks](#/config/tasks)
  - [github.com/mh-cbon/backup/api.TestConfigListTasks](#github.com/mh-cbon/backup/api.TestConfigListTasks)


- [/disksinfo](#/disksinfo)
  - [github.com/mh-cbon/backup/api.TestDiskInfo](#github.com/mh-cbon/backup/api.TestDiskInfo)
  - [github.com/mh-cbon/backup/api.TestDiskInfoFailure](#github.com/mh-cbon/backup/api.TestDiskInfoFailure)


- [/tasks/list](#/tasks/list)
  - [github.com/mh-cbon/backup/api.TestTaskList](#github.com/mh-cbon/backup/api.TestTaskList)


- [/tasks/start](#/tasks/start)
  - [github.com/mh-cbon/backup/api.TestListRunningTasks](#github.com/mh-cbon/backup/api.TestListRunningTasks)
  - [github.com/mh-cbon/backup/api.TestStartTaskAlreadyStarted](#github.com/mh-cbon/backup/api.TestStartTaskAlreadyStarted)
  - [github.com/mh-cbon/backup/api.TestStartUnknownTask](#github.com/mh-cbon/backup/api.TestStartUnknownTask)
  - [github.com/mh-cbon/backup/api.TestTaskStart](#github.com/mh-cbon/backup/api.TestTaskStart)


- [/tasks/stop](#/tasks/stop)
  - [github.com/mh-cbon/backup/api.TestStopTaskNotStarted](#github.com/mh-cbon/backup/api.TestStopTaskNotStarted)
  - [github.com/mh-cbon/backup/api.TestStopUnknownTask](#github.com/mh-cbon/backup/api.TestStopUnknownTask)
  - [github.com/mh-cbon/backup/api.TestTaskStop](#github.com/mh-cbon/backup/api.TestTaskStop)




### /config/addTask
AddTask adds a task to the Config


#### github.com/mh-cbon/backup/api.TestConfigAddTask
TestConfigAddTask add a task to the configuration.


__[200] POST__ /config/addTask

##### Request Body
```
{
    "Dest": "Dest",
    "Name": "Name",
    "Source": "Source"
}
```


##### Response Headers

| Key | Value |
| --- | --- |
| Content-Type | [text/plain; charset=utf-8] |


##### Response Body
```
{
    "Data": {
        "Dest": "Dest",
        "Name": "Name",
        "Source": "Source"
    },
    "Errors": null,
    "Res": true
}
```

[TOP](#api-doc)
___________________

#### github.com/mh-cbon/backup/api.TestConfigAddTaskAlreadyExisting
TestConfigAddTask fail to add a task already existing.


__[200] POST__ /config/addTask

##### Request Body
```
{
    "Dest": "Dest",
    "Name": "Name",
    "Source": "Source"
}
```


##### Response Headers

| Key | Value |
| --- | --- |
| Content-Type | [text/plain; charset=utf-8] |


##### Response Body
```
{
    "Errors": {
        "Name": "Task name already exists"
    },
    "Res": false
}
```

[TOP](#api-doc)
___________________


### /config/read
Read loads the Config


#### github.com/mh-cbon/backup/api.TestConfigRead
TestConfigRead read the configuration.


__[200] GET__ /config/read

##### Response Headers

| Key | Value |
| --- | --- |
| Content-Type | [text/plain; charset=utf-8] |


##### Response Body
```
{
    "Tasks": [
        {
            "Dest": "Dest",
            "Name": "Name",
            "Source": "Source"
        }
    ]
}
```

[TOP](#api-doc)
___________________


### /config/rmTask
RmTask removes a task from the Config


#### github.com/mh-cbon/backup/api.TestConfigRmTask
TestConfigRmTask remove a task from the config.


__[200] POST__ /config/rmTask

##### Request Body
```
{
    "Dest": "Dest",
    "Name": "Name",
    "Source": "Source"
}
```


##### Response Headers

| Key | Value |
| --- | --- |
| Content-Type | [text/plain; charset=utf-8] |


##### Response Body
```
{
    "Errors": null,
    "Res": true
}
```

[TOP](#api-doc)
___________________

#### github.com/mh-cbon/backup/api.TestConfigRmUnknownTask
TestConfigRmUnknownTask fail to remove unknown task.


__[200] POST__ /config/rmTask

##### Request Body
```
{
    "Dest": "Dest",
    "Name": "Name",
    "Source": "Source"
}
```


##### Response Headers

| Key | Value |
| --- | --- |
| Content-Type | [text/plain; charset=utf-8] |


##### Response Body
```
{
    "Errors": {
        "Name": "Task not found"
    },
    "Res": false
}
```

[TOP](#api-doc)
___________________


### /config/tasks
ListTasks lists configured tasks


#### github.com/mh-cbon/backup/api.TestConfigListTasks
TestConfigListTasks read configured tasks.


__[200] GET__ /config/tasks

##### Response Headers

| Key | Value |
| --- | --- |
| Content-Type | [text/plain; charset=utf-8] |


##### Response Body
```
{
    "Data": [
        {
            "Dest": "Dest",
            "Name": "Name",
            "Source": "Source"
        }
    ],
    "Errors": null,
    "Res": true
}
```

[TOP](#api-doc)
___________________


### /disksinfo


#### github.com/mh-cbon/backup/api.TestDiskInfo
TestDiskInfo lists available partitions on the computer.


__[200] GET__ /disksinfo

##### Response Headers

| Key | Value |
| --- | --- |
| Content-Type | [text/plain; charset=utf-8] |


##### Response Body
```
{
    "Data": [
        {
            "IsRemovable": false,
            "Label": "fake",
            "MountPath": "",
            "Path": "",
            "Size": "",
            "SpaceLeft": ""
        }
    ],
    "Errors": null,
    "Res": true
}
```

[TOP](#api-doc)
___________________

#### github.com/mh-cbon/backup/api.TestDiskInfoFailure
TestDiskInfoFailure fail to read the computer partition list.


__[200] GET__ /disksinfo

##### Response Headers

| Key | Value |
| --- | --- |
| Content-Type | [text/plain; charset=utf-8] |


##### Response Body
```
{
    "Errors": "DiskInfo loading failed Disk info was not loaded ",
    "Res": false
}
```

[TOP](#api-doc)
___________________


### /tasks/list


#### github.com/mh-cbon/backup/api.TestTaskList
TestTaskList returns a list of active tasks.


__[200] GET__ /tasks/list

##### Response Headers

| Key | Value |
| --- | --- |
| Content-Type | [text/plain; charset=utf-8] |


##### Response Body
```
{
    "Data": [
        {
            "Dest": "",
            "ETA": "",
            "EndDate": "",
            "Failure": null,
            "Messages": null,
            "Name": "xxxx",
            "Percent": "",
            "Source": "",
            "StartDate": "",
            "Started": false
        }
    ],
    "Errors": null,
    "Res": true
}
```

[TOP](#api-doc)
___________________


### /tasks/start


#### github.com/mh-cbon/backup/api.TestListRunningTasks
TestListRunningTasks liasts running tasks.


__[200] POST__ /tasks/start

##### Request Body
```
{
    "Dest": "Dest",
    "Name": "Name",
    "Source": "Source"
}
```


##### Response Headers

| Key | Value |
| --- | --- |
| Content-Type | [text/plain; charset=utf-8] |


##### Response Body
```
{
    "Data": {
        "Dest": "Dest",
        "ETA": "",
        "EndDate": "",
        "Failure": null,
        "Messages": null,
        "Name": "Name",
        "Percent": "",
        "Source": "Source",
        "StartDate": "",
        "Started": true
    },
    "Errors": null,
    "Res": true
}
```

[TOP](#api-doc)
___________________

#### github.com/mh-cbon/backup/api.TestStartTaskAlreadyStarted
TestStartTaskAlreadyStarted fails to start the same task twice.


__[200] POST__ /tasks/start

##### Request Body
```
{
    "Dest": "Dest",
    "Name": "Name",
    "Source": "Source"
}
```


##### Response Headers

| Key | Value |
| --- | --- |
| Content-Type | [text/plain; charset=utf-8] |


##### Response Body
```
{
    "Errors": "Task was not started: Task not started \"Name\"",
    "Res": false
}
```

[TOP](#api-doc)
___________________

#### github.com/mh-cbon/backup/api.TestStartUnknownTask
TestStartUnknownTask fails to start an unknown task.


__[200] POST__ /tasks/start

##### Request Body
```
{
    "Dest": "",
    "Name": "nop nop",
    "Source": ""
}
```


##### Response Headers

| Key | Value |
| --- | --- |
| Content-Type | [text/plain; charset=utf-8] |


##### Response Body
```
{
    "Errors": {
        "Name": "Task not found"
    },
    "Res": false
}
```

[TOP](#api-doc)
___________________

#### github.com/mh-cbon/backup/api.TestTaskStart
TestTaskStart starts a task.


__[200] POST__ /tasks/start

##### Request Body
```
{
    "Dest": "Dest",
    "Name": "Name",
    "Source": "Source"
}
```


##### Response Headers

| Key | Value |
| --- | --- |
| Content-Type | [text/plain; charset=utf-8] |


##### Response Body
```
{
    "Data": {
        "Dest": "Dest",
        "ETA": "",
        "EndDate": "",
        "Failure": null,
        "Messages": null,
        "Name": "Name",
        "Percent": "",
        "Source": "Source",
        "StartDate": "",
        "Started": true
    },
    "Errors": null,
    "Res": true
}
```

[TOP](#api-doc)
___________________


### /tasks/stop


#### github.com/mh-cbon/backup/api.TestStopTaskNotStarted
TestStopTaskNotStarted fails to stop a task not yet started.


__[200] POST__ /tasks/stop

##### Request Body
```
{
    "Dest": "Dest",
    "Name": "Name",
    "Source": "Source"
}
```


##### Response Headers

| Key | Value |
| --- | --- |
| Content-Type | [text/plain; charset=utf-8] |


##### Response Body
```
{
    "Errors": "Task was not found: Task not stopped \"Name\"",
    "Res": false
}
```

[TOP](#api-doc)
___________________

#### github.com/mh-cbon/backup/api.TestStopUnknownTask
TestStopUnknownTask fails to stop unknown task.


__[200] POST__ /tasks/stop

##### Request Body
```
{
    "Dest": "Dest",
    "Name": "Name",
    "Source": "Source"
}
```


##### Response Headers

| Key | Value |
| --- | --- |
| Content-Type | [text/plain; charset=utf-8] |


##### Response Body
```
{
    "Errors": "Task was not found: Task not stopped \"Name\"",
    "Res": false
}
```

[TOP](#api-doc)
___________________

#### github.com/mh-cbon/backup/api.TestTaskStop
TestTaskStop stop a task.


__[200] POST__ /tasks/stop

##### Request Body
```
{
    "Dest": "Dest",
    "Name": "Name",
    "Source": "Source"
}
```


##### Response Headers

| Key | Value |
| --- | --- |
| Content-Type | [text/plain; charset=utf-8] |


##### Response Body
```
{
    "Data": {
        "Dest": "Dest",
        "ETA": "",
        "EndDate": "",
        "Failure": null,
        "Messages": null,
        "Name": "Name",
        "Percent": "",
        "Source": "Source",
        "StartDate": "",
        "Started": false
    },
    "Errors": null,
    "Res": true
}
```

[TOP](#api-doc)
___________________
