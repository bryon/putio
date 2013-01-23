#Golang library for Put.io Api

This project aims to be a [Go][1] library for the current [Put.io][2] [V2 api][3].  This is my first go release so there could be bugs.  Please let me know asap and I will correct if you find any.  

Usage
------
To begin you should always call putio.NewPutio with the appropriate values.  This will make a call to put.io to get the appropriate oauth token.  It is up to you to set up and retrieve you usercode to make this call.  Please read the "Obtain an access token" section of api documentation for more information.

Once you have retrieved your putio token all calls follow the api documention in name.  For example to get the files/list you would call putio.FilesList.  This will return a struct, json, and error so you can decide how you want to handle the returned information.

You'll notice that strings and intergers in the structures have a custom type.  This is done to handle the null values that often appear in put.io's returned json.  Go does not handle this as it views them as malformed json and will throw and error.

Please see the example for further usage.

Copyright
---------
Copyright (c) 2013 Bryon Mayo See LICENSE for details.

[1]: http://golang.org
[2]: http://put.io
[3]: https://api.put.io/v2/docs/