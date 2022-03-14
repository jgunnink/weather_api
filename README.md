# Weather API

This project is designed to retrieve weather information for Sydney when a HTTP call to the service is made.

The expected response will be in JSON format in the form:

```json
{
  "wind_speed": 20,
  "temperature_degrees": 29
}
```

## Usage

To consume the service, first build the Dockerfile like so:

```
docker build -t weatherapi .
```

Then, run the service with the following:

```
docker run -p 8080:8080 -e WEATHERSTACK_KEY="your_key_here" weatherapi
```

After the service is spinning, you can `curl` it with:

```
http://localhost:8080/v1/weather?city=sydney
```
