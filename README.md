# geolocation
Geolocation REST API, obtain IP address geolocation data
Location is extracted by looking up user IP in a Geo DB (from Maxmind).

## Instructions

1. Make sure you have Go installed ([download](https://golang.org/dl/)).
2. Make sure you update the MaxMind license key in the configuration with your key.
3. Then you can either run ``make run`` to run locally, or ``make api-up`` to run in a docker container.
4. To build an executable binary file, run `make build` in a command line terminal. This will build and create `geolocation` file.

## Configuration:
- There are two configuration files by default: `dev` and `localdocker`, but we can add as many config files per each environment, for example `prod`. To use a specific config file before running the application, specify `SERVER_MODE` environment variable. for example `SERVER_MODE=dev`, `SERVER_MODE=prod` or `SERVER_MODE=your_config_file`. This will tell the application to look for `dev.yaml`, `prod.yaml` or `your_config_file.yaml` file and reads the configurations from that file. If you don't specify a value for `SERVER_MODE` environment, `dev` will be selected by default.
- Configuration files are located under ``configs`` folder.
- Format of the configuration file is YAML.

#### GEO_CITY_DB:
The file name of GEO DB. (The latest version will be downloaded from https://dev.maxmind.com/geoip/geoip2/geolite2/ each time you run the api, and will be scheduled to download it every day. It needs you to sign up and create a License Key)
It is located under the ``geodb`` folder.

## How to run

#### Build and run locally:
- Build the application (read the instructions section)
- Run this command in the terminal:

```
SERVER_MODE=dev ./geolocation
```

Output:
```
11:22AM INF pkg/config/config.go:13 > Configs are being initialized
11:22AM INF pkg/maxmind/maxmind.go:55 > Downloading maxmind db...
11:22AM INF cmd/api/main.go:115 > starting api server addr=:3000

```
#### Run in a docker container:
- Build and run api (`make api-up`)
  Docker will run and, then you can use the webservice with port 3000 (configurable in the configs files).

### command line arguments:
- `loglevel` : log level can be a value from `trace`, `debug`, `info`, `warn`, `error`
  example `make run -loglevel=info`

### Consuming the webservice

- To check if the webservice is running, call this utl in a browser (http://127.0.0.1:3000). Note: Make sure the port 3000 is already free in the hosting environment.

- ![FHGEO home page](https://raw.githubusercontent.com/hojabri/geolocation/main/static/first_page.png) This is the swagger documentation of the webservice.
  In this page you can also download the swagger file.
- ### Usage

GET http://localhost:3000/api/v1/geo/1.2.3.4

Response:
```json
{
  "City": {
    "GeoNameID": 0,
    "Names": null
  },
  "Continent": {
    "Code": "OC",
    "GeoNameID": 6255151,
    "Names": {
      "de": "Ozeanien",
      "en": "Oceania",
      "es": "Oceanía",
      "fr": "Océanie",
      "ja": "オセアニア",
      "pt-BR": "Oceania",
      "ru": "Океания",
      "zh-CN": "大洋洲"
    }
  },
  "Country": {
    "GeoNameID": 2077456,
    "IsInEuropeanUnion": false,
    "IsoCode": "AU",
    "Names": {
      "de": "Australien",
      "en": "Australia",
      "es": "Australia",
      "fr": "Australie",
      "ja": "オーストラリア",
      "pt-BR": "Austrália",
      "ru": "Австралия",
      "zh-CN": "澳大利亚"
    }
  },
  "Location": {
    "AccuracyRadius": 1000,
    "Latitude": -33.494,
    "Longitude": 143.2104,
    "MetroCode": 0,
    "TimeZone": "Australia/Sydney"
  },
  "Postal": {
    "Code": ""
  },
  "RegisteredCountry": {
    "GeoNameID": 2077456,
    "IsInEuropeanUnion": false,
    "IsoCode": "AU",
    "Names": {
      "de": "Australien",
      "en": "Australia",
      "es": "Australia",
      "fr": "Australie",
      "ja": "オーストラリア",
      "pt-BR": "Austrália",
      "ru": "Австралия",
      "zh-CN": "澳大利亚"
    }
  },
  "RepresentedCountry": {
    "GeoNameID": 0,
    "IsInEuropeanUnion": false,
    "IsoCode": "",
    "Names": null,
    "Type": ""
  },
  "Subdivisions": null,
  "Traits": {
    "IsAnonymousProxy": false,
    "IsSatelliteProvider": false
  }
}
```

- geoip-golang (https://github.com/oschwald/geoip2-golang)

To read from geoip db file.

- viper (https://github.com/spf13/viper)

Viper is a complete configuration solution for Go applications including 12-Factor apps. It is designed to work within an application, and can handle all types of configuration needs and formats.

_GeoLite2:
This product includes GeoLite2 data created by MaxMind, available from
https://www.maxmind.com_