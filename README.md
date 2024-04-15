# TRAFFIC INFORMATION

Retrieves traffic information of a given origin and destination `latitude,longitude` coordinates.

## USAGE

This program requires google map api key inorder to utilize google map functionalities.

Export the key before using the program by following the sample below (for unix based OS)

```shell
export GOOGLE_MAP_API_KEY=***your api key***
```

Clone the repository and changed directory to the project root directory, then execute
the following command to use the program:
```shell
go run main.go -origin latitude,longitude -destination latitude,longitude
```
