#Minimalistic Video Transcoding

This is an attempt to implement the minimalistic transcoding in golang, The job is submitted to the api endpoint where the POST parameters
contain the input file and desired output format, where the api returns the job_id is returned for calback and the job is processed 
asynchronously through a distributed worker, the callback will return the output file when job is successfully completed, I have attempted
to transcode the videos to HLS in this implementation

## Technologies Used

[ Programming Language - go-1.11.1 ](https://dl.google.com/go/go1.11.1.darwin-amd64.pkg)
[ Distributed Task Queue - RabbitMQ ](http://www.rabbitmq.com/install-homebrew.html)
[ Jobs State - MySQL ](https://dev.mysql.com/get/Downloads/MySQL-8.0/mysql-8.0.12-macos10.13-x86_64.dmg)
[ ffmpeg - To Transocde ](https://evermeet.cx/ffmpeg/ffmpeg-4.0.2.dmg)
[ Machinery - Asynchronous Processing ](https://github.com/RichardKnop/machinery)
[ Gorilla Mux - HTTP Router ](https://github.com/gorilla/mux)
[ gorm - Struct to Table Mapper ](https://github.com/jinzhu/gorm)

## API Usage

### Submitting the Encoding Job
```bash
bseenu@C02Q942ZG8WP-lm ~/Downloads> curl -v -X POST -d '{ "format": "hls", "input": "/Users/bseenu/Downloads/test.mp4" }' 'http://localhost:8000/api/v1/jobs'
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8000 (#0)
> POST /api/v1/jobs HTTP/1.1
> Host: localhost:8000
> User-Agent: curl/7.54.0
> Accept: */*
> Content-Length: 64
> Content-Type: application/x-www-form-urlencoded
>
* upload completely sent off: 64 out of 64 bytes
< HTTP/1.1 201 Created
< Date: Mon, 08 Oct 2018 02:35:15 GMT
< Content-Length: 49
< Content-Type: text/plain; charset=utf-8
<
{"JobID":"6e398e59-3f54-4d8a-88cc-7b1192ff341d"}
* Connection #0 to host localhost left intact
```

### Querying the job - Callback

```bash
bseenu@C02Q942ZG8WP-lm ~/Downloads> curl -v 'http://localhost:8000/api/v1/jobs/6e398e59-3f54-4d8a-88cc-7b1192ff341d'
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8000 (#0)
> GET /api/v1/jobs/6e398e59-3f54-4d8a-88cc-7b1192ff341d HTTP/1.1
> Host: localhost:8000
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Mon, 08 Oct 2018 02:35:28 GMT
< Content-Length: 189
< Content-Type: text/plain; charset=utf-8
<
{"job_id":"6e398e59-3f54-4d8a-88cc-7b1192ff341d","created":"2018-10-08T02:35:16Z","completed":null,"format":"hls","input":"/Users/bseenu/Downloads/test.mp4","output":"","status":"PENDING"}
* Connection #0 to host localhost left intact


bseenu@C02Q942ZG8WP-lm ~/Downloads> curl -v 'http://localhost:8000/api/v1/jobs/6e398e59-3f54-4d8a-88cc-7b1192ff341d'
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8000 (#0)
> GET /api/v1/jobs/6e398e59-3f54-4d8a-88cc-7b1192ff341d HTTP/1.1
> Host: localhost:8000
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Mon, 08 Oct 2018 02:35:37 GMT
< Content-Length: 189
< Content-Type: text/plain; charset=utf-8
<
{"job_id":"6e398e59-3f54-4d8a-88cc-7b1192ff341d","created":"2018-10-08T02:35:16Z","completed":null,"format":"hls","input":"/Users/bseenu/Downloads/test.mp4","output":"","status":"PENDING"}
* Connection #0 to host localhost left intact
bseenu@C02Q942ZG8WP-lm ~/Downloads> curl -v 'http://localhost:8000/api/v1/jobs/6e398e59-3f54-4d8a-88cc-7b1192ff341d'
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8000 (#0)
> GET /api/v1/jobs/6e398e59-3f54-4d8a-88cc-7b1192ff341d HTTP/1.1
> Host: localhost:8000
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Mon, 08 Oct 2018 02:35:39 GMT
< Content-Length: 284
< Content-Type: text/plain; charset=utf-8
<
{"job_id":"6e398e59-3f54-4d8a-88cc-7b1192ff341d","created":"2018-10-08T02:35:16Z","completed":"2018-10-08T02:35:38Z","format":"hls","input":"/Users/bseenu/Downloads/test.mp4","output":"/Users/bseenu/Downloads/6e398e59-3f54-4d8a-88cc-7b1192ff341d/playlist.m3u8","status":"SUCCESSFUL"}
* Connection #0 to host localhost left intact
```

### Media playlist and segment files

```bash
bseenu@C02Q942ZG8WP-lm ~/Downloads> ls -ltr 6e398e59-3f54-4d8a-88cc-7b1192ff341d/
total 63696
-rw-r--r--  1 bseenu  staff  1505880 Oct  7 19:35 0.ts
-rw-r--r--  1 bseenu  staff  1217112 Oct  7 19:35 1.ts
-rw-r--r--  1 bseenu  staff  1120668 Oct  7 19:35 2.ts
-rw-r--r--  1 bseenu  staff   822312 Oct  7 19:35 3.ts
-rw-r--r--  1 bseenu  staff  1076864 Oct  7 19:35 4.ts
-rw-r--r--  1 bseenu  staff  1386876 Oct  7 19:35 5.ts
-rw-r--r--  1 bseenu  staff   457968 Oct  7 19:35 6.ts
-rw-r--r--  1 bseenu  staff  2389292 Oct  7 19:35 7.ts
-rw-r--r--  1 bseenu  staff  1033812 Oct  7 19:35 8.ts
-rw-r--r--  1 bseenu  staff  1344388 Oct  7 19:35 9.ts
-rw-r--r--  1 bseenu  staff  1401352 Oct  7 19:35 10.ts
-rw-r--r--  1 bseenu  staff   955228 Oct  7 19:35 11.ts
-rw-r--r--  1 bseenu  staff  3271764 Oct  7 19:35 12.ts
-rw-r--r--  1 bseenu  staff   695788 Oct  7 19:35 13.ts
-rw-r--r--  1 bseenu  staff  1051672 Oct  7 19:35 14.ts
-rw-r--r--  1 bseenu  staff  1649136 Oct  7 19:35 15.ts
-rw-r--r--  1 bseenu  staff    88172 Oct  7 19:35 16.ts
-rw-r--r--  1 bseenu  staff  1901244 Oct  7 19:35 17.ts
-rw-r--r--  1 bseenu  staff   496508 Oct  7 19:35 18.ts
-rw-r--r--  1 bseenu  staff   941504 Oct  7 19:35 19.ts
-rw-r--r--  1 bseenu  staff   666272 Oct  7 19:35 20.ts
-rw-r--r--  1 bseenu  staff   394048 Oct  7 19:35 21.ts
-rw-r--r--  1 bseenu  staff   750496 Oct  7 19:35 22.ts
-rw-r--r--  1 bseenu  staff   873448 Oct  7 19:35 23.ts
-rw-r--r--  1 bseenu  staff   364908 Oct  7 19:35 24.ts
-rw-r--r--  1 bseenu  staff   730004 Oct  7 19:35 25.ts
-rw-r--r--  1 bseenu  staff      725 Oct  7 19:35 playlist.m3u8
```

#### Playing the hls 

```bash
bseenu@C02Q942ZG8WP-lm ~/Downloads> ffplay 6e398e59-3f54-4d8a-88cc-7b1192ff341d/playlist.m3u8
```


